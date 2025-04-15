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
	PING              = "PING"
	GET               = "GET"
	SET               = "SET"
	GEOADD            = "GEOADD"
	GEODIST           = "GEODIST"
	GEORADIUS         = "GEORADIUS"
	GEORADIUSBYMEMBER = "GEORADIUSBYMEMBER"
	GEOHASH           = "GEOHASH"
	GEOPOS            = "GEOPOS"
	GEOSEARCH         = "GEOSEARCH"
	GEOSEARCHSTORE    = "GEOSEARCHSTORE"
)

func NewCmdHandler() *CmdHandler {
	return &CmdHandler{}
}

var store = db.NewStore(db.NewDb())
var geoDb = db.NewGeoDB()

func (h *CmdHandler) HandleDbCommands(cmd redcon.Command) []byte {
	parsedCmd := h.parseCommand(cmd)
	SUFF_CMD := CommandType(parsedCmd[0])

	if SUFF_CMD == PING {
		return h.ping()
	}

	if SUFF_CMD == GET {
		return h.Get(parsedCmd)
	}

	if SUFF_CMD == SET {
		return h.Set(parsedCmd)
	}

	if SUFF_CMD == GEOADD {
		return h.GeoAdd(parsedCmd)
	}

	if SUFF_CMD == GEODIST {
		return h.GEODIST(parsedCmd)
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

func (h *CmdHandler) GeoAdd(cmd []string) []byte {
	if len(cmd) < 5 {
		return redcon.AppendError(nil, "ERR wrong number of arguments for 'geoadd' command")
	}

	key := cmd[1]

	lat, err := strconv.ParseFloat(cmd[2], 64)
	if err != nil {
		return redcon.AppendError(nil, "ERR invalid latitude")
	}

	lon, err := strconv.ParseFloat(cmd[3], 64)
	if err != nil {
		return redcon.AppendError(nil, "ERR invalid longitude")
	}

	member := cmd[4]
	geoDb.GeoAdd(key, lon, lat, member)
	return redcon.AppendOK(nil)
}

func (h *CmdHandler) GEODIST(cmd []string) []byte {
	if len(cmd) < 4 {
		return redcon.AppendError(nil, "ERR wrong number of arguments for 'geodist' command")
	}
	key := cmd[1]
	member1 := cmd[2]
	member2 := cmd[3]
	dist := geoDb.GeoDist(key, member1, member2)
	return redcon.AppendBulkFloat(nil, dist)
}
