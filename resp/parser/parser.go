package parser

import (
	"bufio"
	"errors"
	"godis/interface/resp"
	"io"
	"strconv"
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
	go parser0(reader, ch)
	return ch
}
func parser0(reader io.Reader, ch chan *Payload) {

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
