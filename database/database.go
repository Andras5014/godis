package database

import (
	"godis/config"
	"godis/interface/resp"
	"godis/lib/logger"
	"godis/resp/reply"
	"strconv"
	"strings"
)

type Database struct {
	dbSet []*DB
}

func NewDatabase() *Database {
	if config.Properties.Databases == 0 {
		config.Properties.Databases = 16
	}
	dbSet := make([]*DB, 0, config.Properties.Databases)
	for i := 0; i < config.Properties.Databases; i++ {
		dbSet = append(dbSet, NewDB())
		dbSet[i].index = i
	}
	return &Database{
		dbSet: dbSet,
	}
}
func (d *Database) Exec(conn resp.Connection, args [][]byte) resp.Reply {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	if cmdName == "select" {
		if len(args) != 2 {
			return reply.NewArgsErrReply("select")
		}
		return execSelect(conn, d, args[1:])
	}
	dbIndex := conn.GetDBIndex()
	return d.dbSet[dbIndex].Exec(conn, args)

}

func (d *Database) Close() {
	//TODO implement me
	panic("implement me")
}

func (d *Database) AfterClientClose(c resp.Connection) {
	//TODO implement me
	panic("implement me")
}

func execSelect(c resp.Connection, database *Database, args [][]byte) resp.Reply {
	dbIndex, err := strconv.Atoi(string(args[0]))
	if err != nil {
		return reply.NewErrReply("ERR invalid DB index")
	}
	if dbIndex >= len(database.dbSet) {
		return reply.NewErrReply("ERR DB index is out of range")
	}
	c.SelectDB(dbIndex)
	return reply.NewOkReply()
}
