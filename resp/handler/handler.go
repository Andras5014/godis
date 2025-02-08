package handler

import (
	"context"
	"errors"
	"godis/interface/database"
	"godis/lib/logger"
	"godis/lib/sync/atomic"
	"godis/resp/connection"
	"godis/resp/parser"
	"godis/resp/reply"
	"io"
	"net"
	"strings"
	"sync"
)

var (
	unknownErrReplyBytes = []byte("-ERR unknown\r\n")
)

type RespHandler struct {
	activeConn sync.Map
	db         database.Database
	closing    atomic.Boolean
}

func NewRespHandler(db database.Database) *RespHandler {
	return &RespHandler{
		db: db,
	}
}

func (r *RespHandler) closeClient(client *connection.Connection) {
	_ = client.Close()
	r.db.AfterClientClose(client)
	r.activeConn.Delete(client)
	logger.Info("connection closed: " + client.RemoteAddr().String())
}
func (r *RespHandler) Handle(ctx context.Context, conn net.Conn) {
	if r.closing.Get() {
		// closing handler refuse new connection
		_ = conn.Close()
	}
	client := connection.NewConnection(conn)
	r.activeConn.Store(client, struct{}{})
	ch := parser.ParseStream(conn)
	for payload := range ch {
		if payload.Err != nil {
			if errors.Is(payload.Err, io.EOF) || errors.Is(payload.Err, io.ErrUnexpectedEOF) || strings.Contains(payload.Err.Error(), "use of closed network connection") {

			}
			errReply := reply.NewErrReply(payload.Err.Error())
			err := client.Write(errReply.ToBytes())
			if err != nil {
				r.closeClient(client)
				logger.Info("connection closed: " + client.RemoteAddr().String())
				return
			}
			continue
		}
		if payload.Data == nil {
			continue
		}
		replyResult, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		result := r.db.Exec(client, replyResult.Args)
		if result != nil {
			_ = client.Write(result.ToBytes())
		} else {
			_ = client.Write(unknownErrReplyBytes)
		}
	}
}

func (r *RespHandler) Close() error {
	logger.Info("handler shutting down...")
	r.closing.Set(true)
	r.activeConn.Range(func(key interface{}, val interface{}) bool {
		client := key.(*connection.Connection)
		_ = client.Close()
		return true
	})
	return nil
}
