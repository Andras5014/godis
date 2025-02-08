package database

import "godis/interface/resp"

type CmdLine = [][]byte
type Database interface {
	Exec(conn resp.Connection, args [][]byte) resp.Reply
	Close()
	AfterClientClose(c resp.Connection)
}

type DataEntity struct {
	Data interface{}
}
