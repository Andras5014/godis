package database

import (
	"godis/datastruct/dict"
	"godis/interface/resp"
	"godis/resp/reply"
	"strings"
)

type ExecutorFunc func(db *DB, args [][]byte) resp.Reply
type CmdLine = [][]byte
type DB struct {
	index int
	data  dict.Dict
}

func NewDB() *DB {
	return &DB{
		data: dict.NewSyncDict(),
	}
}

func (db *DB) Exec(conn resp.Connection, cmdLine CmdLine) resp.Reply {
	cmdName := strings.ToLower(string(cmdLine[0]))
	cmd, ok := cmdTable[cmdName]
	if !ok {
		return reply.NewErrReply("ERR unknown command '" + cmdName + "'")
	}
	if !validateArity(cmd.arity, cmdLine) {
		return reply.NewArgsErrReply(cmdName)
	}
	fun := cmd.executor
	return fun(db, cmdLine[1:])
}

func validateArity(arity int, cmdArgs CmdLine) bool {

	return true
}
