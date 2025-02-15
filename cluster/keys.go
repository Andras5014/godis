package cluster

import (
	"godis/interface/resp"
	"godis/resp/reply"
)

func FlushDB(cluster *ClusterDatabase, c resp.Connection, args [][]byte) resp.Reply {
	replies := cluster.broadcast(c, args)
	var errReply reply.ErrorReply
	for _, v := range replies {
		if reply.IsErrorReply(v) {
			errReply = v.(reply.ErrorReply)
			break
		}
	}
	if errReply == nil {
		return &reply.OkReply{}
	}
	return reply.NewErrReply("error occurs: " + errReply.Error())
}
