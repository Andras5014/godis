package reply

var pongBytes = []byte("+PONG\r\n")
var okBytes = []byte("+OK\r\n")
var nullBulkBytes = []byte("$-1\r\n")      //null
var emptyMultiBulkBytes = []byte("*0\r\n") //空数组
var noBytes = []byte("")

type PongReply struct {
}

func NewPongReply() *PongReply {
	return &PongReply{}
}
func (p *PongReply) ToBytes() []byte {
	return pongBytes
}

type OkReply struct {
}

func NewOkReply() *OkReply {
	return &OkReply{}
}
func (o *OkReply) ToBytes() []byte {
	return okBytes
}

type NullBulkReply struct {
}

func NewNullBulkReply() *NullBulkReply {
	return &NullBulkReply{}
}
func (n *NullBulkReply) ToBytes() []byte {
	return nullBulkBytes
}

type EmptyMultiBulkReply struct {
}

func NewEmptyMultiBulkReply() *EmptyMultiBulkReply {
	return &EmptyMultiBulkReply{}
}
func (e *EmptyMultiBulkReply) ToBytes() []byte {
	return emptyMultiBulkBytes
}

type NoReply struct {
}

func NewNoReply() *NoReply {
	return &NoReply{}
}
func (n *NoReply) ToBytes() []byte {
	return noBytes
}
