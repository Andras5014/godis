package database

import (
	"godis/interface/resp"
	"godis/resp/reply"
)

type EchoDatabase struct {
}

func NewEchoDatabase() *EchoDatabase {
	return &EchoDatabase{}
}
func (e *EchoDatabase) Exec(conn resp.Connection, args [][]byte) resp.Reply {
	return reply.NewMultiBulkReply(args)
}

func (e *EchoDatabase) Close() {
	//TODO implement me
	panic("implement me")
}

func (e *EchoDatabase) AfterClientClose(c resp.Connection) {
	//TODO implement me
	panic("implement me")
}
