package parser

import (
	"bufio"
	"errors"
	"godis/interface/resp"
	"godis/lib/logger"
	"godis/resp/reply"
	"io"
	"runtime/debug"
	"strconv"
	"strings"
)

type Payload struct {
	Data resp.Reply
	Err  error
}

type readState struct {
	readingMultiLine  bool
	expectedArgsCount int
	msgType           byte
	args              [][]byte
	bulkLen           int64
}

func (r *readState) finished() bool {
	return r.expectedArgsCount > 0 && len(r.args) == r.expectedArgsCount
}

func ParseStream(reader io.Reader) <-chan *Payload {
	ch := make(chan *Payload)
	go parse0(reader, ch)
	return ch
}
func parse0(reader io.Reader, ch chan *Payload) {
	defer func() {
		if err := recover(); err != nil {
			logger.Error(string(debug.Stack()))
		}
	}()
	bufReader := bufio.NewReader(reader)
	var (
		msg   []byte
		state readState
		err   error
	)
	for {
		var ioErr bool
		msg, ioErr, err = readLine(bufReader, &state)
		if err != nil {
			ch <- &Payload{Err: err}
			if ioErr {
				close(ch)
				return
			}
			state = readState{}
			continue
		}
		// 多行读取
		if !state.readingMultiLine {
			if msg[0] == '*' {
				if err := parseMultiBulkHeader(msg, &state); err != nil {
					ch <- &Payload{Err: err}
					state = readState{}
					continue
				}
				if state.expectedArgsCount == 0 {
					ch <- &Payload{Data: reply.NewEmptyMultiBulkReply()}
					state = readState{}
					continue
				}
			} else if msg[0] == '$' {
				if err := parseBulkHeader(msg, &state); err != nil {
					ch <- &Payload{Err: err}
					state = readState{}
					continue
				}
				// $-1\r\n
				if state.bulkLen == -1 {
					ch <- &Payload{Data: reply.NewNullBulkReply()}
					state = readState{}
					continue
				}
			} else {
				lineReply, err := parserSingleLineReply(msg)
				ch <- &Payload{
					Data: lineReply,
					Err:  err,
				}
				state = readState{}
				continue

			}
		} else {
			if err := readBody(msg, &state); err != nil {
				ch <- &Payload{Err: err}
				state = readState{}
				continue
			}
			if state.finished() {
				var result resp.Reply
				if state.msgType == '*' {
					result = reply.NewMultiBulkReply(state.args)
				} else if state.msgType == '$' {
					result = reply.NewBulkReply(state.args[0])
				}
				ch <- &Payload{
					Data: result,
					Err:  nil,
				}
				state = readState{}
			}
		}
	}
}

// bool判定是否io错误
func readLine(bufReader *bufio.Reader, state *readState) ([]byte, bool, error) {
	var (
		msg []byte
		err error
	)
	if state.bulkLen == 0 { // \r\n切分
		msg, err = bufReader.ReadBytes('\n')
		if err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' {
			return nil, false, errors.New("protocol error" + string(msg))
		}
	} else { // 读取到了$数字，读取数字长度
		//msg, err = bufReader.Peek(int(state.bulkLen) + 2)
		msg = make([]byte, int(state.bulkLen)+2)
		if _, err = io.ReadFull(bufReader, msg); err != nil {
			return nil, true, err
		}
		if len(msg) == 0 || msg[len(msg)-2] != '\r' || msg[len(msg)-1] != '\n' {
			return nil, false, errors.New("protocol error" + string(msg))
		}
		state.bulkLen = 0
	}
	return msg, false, nil
}

// *3\r\n$3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n
func parseMultiBulkHeader(msg []byte, state *readState) error {
	var (
		expectedLine uint64
		err          error
	)
	expectedLine, err = strconv.ParseUint(string(msg[1:len(msg)-2]), 10, 32)
	if err != nil {
		return errors.New("protocol error" + string(msg))
	}
	if expectedLine == 0 {
		state.expectedArgsCount = 0
		return nil
	} else if expectedLine > 0 {
		state.msgType = msg[0]
		state.expectedArgsCount = int(expectedLine)
		state.readingMultiLine = true
		state.args = make([][]byte, 0, expectedLine)
		return nil
	} else {
		return errors.New("protocol error" + string(msg))
	}
}

func parseBulkHeader(msg []byte, state *readState) error {
	var err error
	state.bulkLen, err = strconv.ParseInt(string(msg[1:len(msg)-2]), 10, 64)
	if err != nil {
		return errors.New("protocol error: " + string(msg))
	}
	if state.bulkLen == -1 { // null bulk
		return nil
	} else if state.bulkLen > 0 {
		state.msgType = msg[0]
		state.readingMultiLine = true
		state.expectedArgsCount = 1
		state.args = make([][]byte, 0, 1)
		return nil
	} else {
		return errors.New("protocol error: " + string(msg))
	}
}

// +OK\r\n -err\r\n :5\r\n
func parserSingleLineReply(msg []byte) (resp.Reply, error) {
	var result resp.Reply
	str := strings.TrimSuffix(string(msg), "\r\n")
	switch msg[0] {
	case '+':
		result = reply.NewStatusReply(str[1:])
	case '-':
		result = reply.NewErrReply(str[1:])
	case ':':
		val, err := strconv.ParseInt(str[1:], 10, 64)
		if err != nil {
			return nil, errors.New("protocol error:" + string(msg))
		}
		result = reply.NewIntReply(val)
	}
	return result, nil
}

// $3\r\nSET\r\n$3\r\nKEY\r\n$5\r\nVALUE\r\n
// PING\r\n
func readBody(msg []byte, state *readState) error {
	line := msg[0 : len(msg)-2]
	var err error
	if line[0] == '$' {
		state.bulkLen, err = strconv.ParseInt(string(line[1:]), 10, 64)
		if err != nil {
			return errors.New("protocol error:" + string(msg))
		}
		//$0\r\n
		if state.bulkLen <= 0 {
			state.args = append(state.args, []byte{})
			state.bulkLen = 0
		}
	} else {
		state.args = append(state.args, line)
	}
	return nil
}
