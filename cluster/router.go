package cluster

import "godis/interface/resp"

func newRouter() map[string]CmdFunc {
	routerMap := make(map[string]CmdFunc)
	routerMap["exists"] = defaultFunc
	routerMap["set"] = defaultFunc
	routerMap["setnx"] = defaultFunc
	routerMap["get"] = defaultFunc
	routerMap["getset"] = defaultFunc
	routerMap["ping"] = ping
	routerMap["del"] = Del
	routerMap["flushdb"] = FlushDB
	return routerMap
}
func defaultFunc(cluster *ClusterDatabase, conn resp.Connection, cmdArgs [][]byte) resp.Reply {
	key := string(cmdArgs[1])
	node := cluster.peerPicker.PickNode(key)
	return cluster.relay(node, conn, cmdArgs)
}
