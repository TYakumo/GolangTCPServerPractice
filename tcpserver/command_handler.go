package tcpserver

import (
	"errors"
	"strings"
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

func (c *CmdHandler) ExecuteCommand(cmd string) (int, error) {
	cmd = strings.TrimSuffix(cmd, "\n")
	cmd = strings.TrimSpace(cmd)

	opcode, found := c.cmdStrMap[cmd]

	if found {
		if c.rateCntr.GetToken() {
			c.monChan <- IncCmdInQue
			err := c.runCommand(opcode)
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

func (c *CmdHandler) runCommand(opcode int) error {
	switch opcode {
	case 0:
		return nil
	case 1:
		return RunNoop()
	case 2:
		return RunDelayingNoop()
	default:
		return nil
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
