package database

import "strings"

var cmdTable = make(map[string]*command)

type command struct {
	executor ExecutorFunc
	arity    int //参数数量
}

func RegisterCommand(name string, executor ExecutorFunc, arity int) {
	cmdTable[strings.ToLower(name)] = &command{
		executor: executor,
		arity:    arity,
	}
}
