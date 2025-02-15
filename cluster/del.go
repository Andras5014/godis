package cluster

import (
	"godis/interface/resp"
	"godis/resp/reply"
)

func Del(cluster *ClusterDatabase, c resp.Connection, args [][]byte) resp.Reply {
	replies := cluster.broadcast(c, args)
	var errReply reply.ErrorReply
	var deleted int64 = 0
	for _, v := range replies {
		if reply.IsErrorReply(v) {
			errReply = v.(reply.ErrorReply)
			break
		}
		intReply, ok := v.(*reply.IntReply)
		if !ok {
			errReply = reply.NewErrReply("error")
		}
		deleted += intReply.Code
	}

	if errReply == nil {
		return reply.NewIntReply(deleted)
	}
	return reply.NewErrReply("error occurs: " + errReply.Error())
}
