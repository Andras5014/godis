package database

import (
	"godis/interface/resp"
	"godis/resp/reply"
)

func init() {
	RegisterCommand("PING", Ping, 1)
}
func Ping(db *DB, args [][]byte) resp.Reply {
	return reply.NewPongReply()
}
