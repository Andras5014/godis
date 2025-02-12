package database

import (
	"godis/datastruct/dict"
	"godis/interface/database"
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

func (db *DB) GetEntity(key string) (*database.DataEntity, bool) {
	val, ok := db.data.Get(key)
	if !ok {
		return nil, false
	}
	entity, _ := val.(*database.DataEntity)
	return entity, true
}

func (db *DB) PutEntity(key string, entity *database.DataEntity) int {
	return db.data.Put(key, entity)
}
func (db *DB) PutIfExists(key string, entity *database.DataEntity) int {
	return db.data.PutIfExists(key, entity)
}
func (db *DB) PutIfAbsent(key string, entity *database.DataEntity) int {
	return db.data.PutIfAbsent(key, entity)
}
func (db *DB) Remove(key string) {
	db.data.Remove(key)
}

func (db *DB) Removes(keys ...string) (deleted int) {
	deleted = 0
	for _, key := range keys {
		_, exists := db.data.Get(key)
		if exists {
			db.Remove(key)
			deleted++
		}
	}
	return deleted
}

// Flush clean database
func (db *DB) Flush() {
	db.data.Clear()

}
