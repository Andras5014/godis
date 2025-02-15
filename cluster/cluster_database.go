package cluster

import (
	"context"
	pool "github.com/jolestar/go-commons-pool/v2"
	"godis/config"
	"godis/database"
	db "godis/interface/database"
	"godis/interface/resp"
	"godis/lib/consistenthash"
	"godis/lib/logger"
	"godis/resp/reply"
	"strings"
)

type CmdFunc func(cluster *ClusterDatabase, conn resp.Connection, args [][]byte) resp.Reply

var router = newRouter()

type ClusterDatabase struct {
	self           string
	nodes          []string
	peerPicker     *consistenthash.NodeMap
	peerConnection map[string]*pool.ObjectPool
	db             db.Database
}

func NewClusterDatabase() *ClusterDatabase {
	cluster := &ClusterDatabase{
		self:           config.Properties.Self,
		db:             database.NewStandaloneDatabase(),
		peerPicker:     consistenthash.NewNodeMap(nil),
		peerConnection: make(map[string]*pool.ObjectPool),
	}
	nodes := make([]string, 0, len(config.Properties.Peers)+1)
	for _, node := range config.Properties.Peers {
		nodes = append(nodes, node)
	}
	nodes = append(nodes, cluster.self)
	cluster.peerPicker.AddNode(nodes...)
	ctx := context.Background()
	for _, peer := range config.Properties.Peers {
		cluster.peerConnection[peer] = pool.NewObjectPoolWithDefaultConfig(ctx, &connectionFactory{Peer: peer})
	}
	cluster.nodes = nodes
	return cluster
}
func (c *ClusterDatabase) Exec(conn resp.Connection, args [][]byte) resp.Reply {
	var result resp.Reply
	defer func() {
		if err := recover(); err != nil {
			logger.Error(err)
			result = &reply.UnknownErrReply{}
		}
	}()
	cmdName := strings.ToLower(string(args[0]))
	cmdFunc, ok := router[cmdName]
	if !ok {
		result = reply.NewErrReply("ERR unknown command '" + cmdName + "'")
		return result
	}
	result = cmdFunc(c, conn, args)
	return result
}

func (c *ClusterDatabase) Close() {
	c.db.Close()
}

func (c *ClusterDatabase) AfterClientClose(conn resp.Connection) {
	c.db.AfterClientClose(conn)
}
