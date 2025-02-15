package aof

import (
	"errors"
	"godis/config"
	"godis/interface/database"
	"godis/lib/logger"
	"godis/lib/utils"
	"godis/resp/connection"
	"godis/resp/parser"
	"godis/resp/reply"
	"io"
	"os"
	"strconv"
)

const (
	aofQueueSize = 1 << 16
)

type CmdLine = [][]byte
type payload struct {
	cmdLine CmdLine
	dbIndex int
}
type AofHandler struct {
	database    database.Database
	aofFile     *os.File
	aofFileName string
	currentDB   int
	aofChan     chan *payload
}

func NewAofHandler(database database.Database) (*AofHandler, error) {
	aofFileName := config.Properties.AppendFilename
	aofFile, err := os.OpenFile(aofFileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	aofHandler := &AofHandler{
		database:    database,
		aofFile:     aofFile,
		aofFileName: aofFileName,
		currentDB:   0,
		aofChan:     make(chan *payload, aofQueueSize),
	}
	aofHandler.LoadAof()
	go aofHandler.handleAof()
	return aofHandler, nil
}
func (h *AofHandler) AddToAof(dbIndex int, cmdLine CmdLine) {
	if config.Properties.AppendOnly && h.aofChan != nil {
		h.aofChan <- &payload{
			cmdLine: cmdLine,
			dbIndex: dbIndex,
		}
	}
}
func (h *AofHandler) handleAof() {
	for p := range h.aofChan {

		if p.dbIndex != h.currentDB {
			// select db
			data := reply.NewMultiBulkReply(utils.ToCmdLine("SELECT", strconv.Itoa(p.dbIndex))).ToBytes()
			_, err := h.aofFile.Write(data)
			if err != nil {
				logger.Warn(err)
				continue // skip this command
			}
			h.currentDB = p.dbIndex
		}
		data := reply.NewMultiBulkReply(p.cmdLine).ToBytes()
		_, err := h.aofFile.Write(data)
		if err != nil {
			logger.Warn(err)
		}
	}
}

func (h *AofHandler) LoadAof() {
	aofFile, err := os.Open(h.aofFileName)
	if err != nil {
		logger.Error(err)
		return
	}
	defer func() {
		if err := aofFile.Close(); err != nil {
			logger.Error(err)
		}
	}()
	payloadCh := parser.ParseStream(aofFile)
	fakeConn := &connection.Connection{}
	for payload := range payloadCh {
		if payload.Err != nil {
			if errors.Is(payload.Err, io.EOF) {
				break
			}
		}
		if payload.Data == nil {
			logger.Error("empty payload")
			continue
		}
		r, ok := payload.Data.(*reply.MultiBulkReply)
		if !ok {
			logger.Error("require multi bulk reply")
			continue
		}
		rep := h.database.Exec(fakeConn, r.Args)
		if reply.IsErrorReply(rep) {
			logger.Error(rep)
		}
	}
}
