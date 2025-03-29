package handler

import (
	"strconv"

	"github.com/R-I-S-H-A-B-H-S-I-N-G-H/keyvDB/internal/db"
	"github.com/tidwall/redcon"
)

type CmdHandler struct{}

// CommandType represents supported database commands.
type CommandType string

const (
	PING = "PING"
	GET  = "GET"
	SET  = "SET"
)

func NewCmdHandler() *CmdHandler {
	return &CmdHandler{}
}

var store = db.NewStore(db.NewDb())

func (h *CmdHandler) HandleDbCommands(cmd redcon.Command) []byte {
	parsedCmd := h.parseCommand(cmd)
	if parsedCmd[0] == PING {
		return h.ping()
	}

	if parsedCmd[0] == GET {
		return h.Get(parsedCmd)
	}

	if parsedCmd[0] == SET {
		return h.Set(parsedCmd)

	}

	return redcon.AppendError(nil, "ERR unknown command")
}

func (h *CmdHandler) parseCommand(cmd redcon.Command) []string {
	commands := []string{}
	for i := 0; i < len(cmd.Args); i++ {
		commands = append(commands, string(cmd.Args[i]))
	}
	return commands
}

func (h *CmdHandler) ping() []byte {
	return redcon.AppendString(nil, "PONG")
}

func (h *CmdHandler) Get(cmd []string) []byte {
	if len(cmd) < 2 {
		return redcon.AppendError(nil, "ERR wrong number of arguments for 'get' command")
	}
	key := cmd[1]
	val, err := store.Get(key)
	if err != nil {
		return redcon.AppendNull(nil)
	}
	return redcon.AppendString(nil, val)
}

func (h *CmdHandler) Set(cmd []string) []byte {
	if len(cmd) < 4 {
		return redcon.AppendError(nil, "ERR wrong number of arguments for 'set' command")
	}
	key := cmd[1]
	val := cmd[2]
	exp, err := strconv.ParseInt(cmd[3], 10, 64)
	if err != nil {
		return redcon.AppendError(nil, "ERR invalid expiration")
	}
	store.Set(key, val, exp)
	return redcon.AppendOK(nil)
}
