package cluster

import (
	"godis/interface/resp"
)

func execSelect(cluster *ClusterDatabase, c resp.Connection, cmdAndArgs [][]byte) resp.Reply {
	return cluster.db.Exec(c, cmdAndArgs)
}
