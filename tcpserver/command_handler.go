package tcpserver

import (
	"context"
	"errors"
	"strings"
	"time"
)

var (
	QuitOpcode = 0
	noopAPI    = 1
	delayAPI   = 2
)

type CmdHandler struct {
	cmdStrMap map[string]int
	rateCntr  *RateLimitController
	monChan   chan int
}

func (c *CmdHandler) InitOpcodeTable() {
	c.cmdStrMap = make(map[string]int)
	c.cmdStrMap["quit"] = QuitOpcode
	c.cmdStrMap["noopAPI"] = noopAPI
	c.cmdStrMap["delayAPI"] = delayAPI
}

func ErrRateLimitReached() error {
	return errors.New("API Rate Limit Reached")
}

func ErrAPIUnavailable() error {
	return errors.New("API unavailable or unreachable")
}

func (c *CmdHandler) ExecuteCommand(cmd string) (int, error) {
	cmd = strings.TrimSuffix(cmd, "\n")
	cmd = strings.TrimSpace(cmd)

	opcode, found := c.cmdStrMap[cmd]

	//setting much lower than 8 seconds deliberately for failure testing
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(3)*time.Second)
	defer cancel()

	if found {
		if c.rateCntr.GetToken() {
			c.monChan <- IncCmdInQue
			err := c.runCommand(opcode, ctx)
			c.monChan <- DecCmdInQue

			if err == nil && opcode != QuitOpcode {
				c.monChan <- IncCmdExecuted
			}
			return opcode, err
		} else {
			return opcode, ErrRateLimitReached()
		}
	}

	return -1, nil
}

func (c *CmdHandler) runCommand(opcode int, ctx context.Context) error {
	errChan := make(chan error)

	switch opcode {
	case 0:
		return nil
	case 1:
		go RunNoop(errChan)
	case 2:
		go RunDelayingNoop(errChan)
	default:
		return nil
	}

	select {
	case <-ctx.Done():
		return ErrAPIUnavailable()
	case err := <-errChan:
		return err
	}

	return nil
}

func StartANewCommandHandler(monChan chan int, rateCntr *RateLimitController) (*CmdHandler, error) {
	var cmdHandler CmdHandler
	cmdHandler.InitOpcodeTable()
	cmdHandler.monChan = monChan
	cmdHandler.rateCntr = rateCntr
	return &cmdHandler, nil
}
