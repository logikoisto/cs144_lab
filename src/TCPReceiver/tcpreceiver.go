package tcpreceiver

import (
	bs "cs144/inMemByteStream"
	sr "cs144/streamreassembler"
)

// capacity: 最大的缓冲区大小
type TCPreceiver struct {
	streamreassembler sr.StreamReassembler
	capacity          int
	isn               WrappingInt32
}

func NewTCPreceiver(cap int) *TCPreceiver {
	return &TCPreceiver{
		streamreassembler: *sr.NewStreamReassembler(cap),
		capacity:          cap,
	}
}

func (recv *TCPreceiver) Unassembled_bytes() int {
	return recv.streamreassembler.Unassembled_bytes()
}

func (recv *TCPreceiver) Stream_out() *bs.ByteStream {
	return recv.streamreassembler.Stream_out()
}
 
// func (recv *TCPreceiver) Ackno() WrappingInt32 {
// 	return 0
// }
