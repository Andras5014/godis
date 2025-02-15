package cluster

import (
	"godis/interface/resp"
	"godis/resp/reply"
)

func Rename(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	if len(cmdArgs) != 3 {
		return reply.NewErrReply("ERR Wrong number args")
	}
	src := string(cmdArgs[1])
	dest := string(cmdArgs[2])

	srcPeer := cluster.peerPicker.PickNode(src)
	destPeer := cluster.peerPicker.PickNode(dest)
	if srcPeer != destPeer {
		return reply.NewErrReply("ERR rename must within on peer")
	}
	return cluster.relay(srcPeer, conn, cmdArgs)
}
