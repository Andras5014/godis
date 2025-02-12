package database

import (
	"godis/interface/resp"
	"godis/lib/wildcard"
	"godis/resp/reply"
)

func execDel(db *DB, args [][]byte) resp.Reply {
	keys := make([]string, len(args))
	for i, arg := range args {
		keys[i] = string(arg)
	}
	deleted := db.Removes(keys...)
	return reply.NewIntReply(int64(deleted))
}

func execExists(db *DB, args [][]byte) resp.Reply {
	result := int64(0)
	for _, arg := range args {
		key := string(arg)
		_, exists := db.GetEntity(key)
		if exists {
			result++
		}
	}
	return reply.NewIntReply(result)
}

func execFlushDB(db *DB, args [][]byte) resp.Reply {
	db.Flush()
	return &reply.OkReply{}
}

func execType(db *DB, args [][]byte) resp.Reply {
	key := string(args[0])
	entity, exists := db.GetEntity(key)
	if !exists {
		return reply.NewStatusReply("none")
	}
	switch entity.Data.(type) {
	case []byte:
		return reply.NewStatusReply("string")
	}
	//todo list hash set
	return &reply.UnknownErrReply{}
}

func execRename(db *DB, args [][]byte) resp.Reply {
	if len(args) != 2 {
		return reply.NewErrReply("ERR wrong number of arguments for 'rename' command")
	}
	src := string(args[0])
	dest := string(args[1])

	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.NewErrReply("no such key")
	}
	db.PutEntity(dest, entity)
	db.Remove(src)
	return &reply.OkReply{}
}

func execRenameNx(db *DB, args [][]byte) resp.Reply {
	src := string(args[0])
	dest := string(args[1])

	_, ok := db.GetEntity(dest)
	if ok {
		return reply.NewIntReply(0)
	}

	entity, ok := db.GetEntity(src)
	if !ok {
		return reply.NewErrReply("no such key")
	}
	db.Removes(src, dest) // clean src and dest with their ttl
	db.PutEntity(dest, entity)
	return reply.NewIntReply(1)
}
func execKeys(db *DB, args [][]byte) resp.Reply {
	pattern := wildcard.CompilePattern(string(args[0]))
	result := make([][]byte, 0)
	db.data.ForEach(func(key string, val interface{}) bool {
		if pattern.IsMatch(key) {
			result = append(result, []byte(key))
		}
		return true
	})
	return reply.NewMultiBulkReply(result)
}

func init() {
	RegisterCommand("Del", execDel, -2)
	RegisterCommand("Exists", execExists, -2)
	RegisterCommand("Keys", execKeys, 2)
	RegisterCommand("FlushDB", execFlushDB, -1)
	RegisterCommand("Type", execType, 2)
	RegisterCommand("Rename", execRename, 3)
	RegisterCommand("RenameNx", execRenameNx, 3)
}
