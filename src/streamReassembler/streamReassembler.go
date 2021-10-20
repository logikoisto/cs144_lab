package streamreassembler

import (
	bs "cs144/inMemByteStream"
)

// 无序地读入完整字节流中的各个子片段（可能有重叠，但一定完整覆盖了整个字节流），一边读入一边使得片段有序，
// 通过ByteStream提供输出流，读入的数据若能够写入输出，则尽快写入输出，否则缓存在unassembledBuffer中
// 如何尽快写入？由于unread的部分处于output输出流中，unread的部分加上unassambled的部分不会超过capacity,
// 因此可以在unassembledBuffer的start指针指向的位置写入数据时即将可以被写入输出流的部分写到输出流中

// capacity: 从（完整字节流的）第一个未读出的字节开始计算，超出capacity的部分在被接收到时直接丢弃。控制内存使用。
// output: 提供有序的输出
// unassembledBuffer: 没有按顺序到来的字节片段首先被存放在这个buffer中, 最大可能的长度是capacity
// unassembledBufferHasCont: 用一个简单的方式来记录buffer上特定位置有没有记录完整字节流上的某些信息，即一个bool数组，在
// 		复制字节时将这一位置为true，当这一位的数据被写入到byteStream中时，将这一位置为false
// unassembledBufferStartPtr: 由于unassembledBuffer的头位置时常要被写到byteStream中，对这个数组时常做复制操作
// 		效率太低，所以设置一个StartPtr指向这个buffer初始可用的位置。
// unassembledLocationInWholeStream指向完整的字节流上
// currentLargerThanCapIndex: 用于记录可以被写入到buffer里的完整字节流里“最远的位置”，大于等于此位置的信息由于capacity的
// 		限制如果被提前收到会直接丢弃
type StreamReassembler struct {
	capacity                         int
	output                           bs.ByteStream
	unassembledBuffer                []byte
	unassembledBufferHasCont         []bool
	unassembledBufferStartPtr        int
	unassembledLocationInWholeStream int
	currentLargerThanCapIndex        int
	unAssembledBytes                 int
}

func NewStreamReassembler(cap int) *StreamReassembler {
	newStreamReassembler := StreamReassembler{
		capacity:                         cap,
		output:                           *bs.NewByteStream(cap),
		unassembledBuffer:                make([]byte, cap),
		unassembledBufferHasCont:         make([]bool, cap),
		unassembledBufferStartPtr:        0,
		unassembledLocationInWholeStream: 0,
		currentLargerThanCapIndex:        cap,
		unAssembledBytes:                 0,
	}
	return &newStreamReassembler
}

// data: 要写入的子串；index: 在完整字节流中的位置；eof：这个子串的结束是不是整个流的结束
func (streamReassembler *StreamReassembler) Push_substring(data string, index int, eof bool) {
	dataBytes := []byte(data)
	writingLength := len(dataBytes)
	// 如果待写入的串已经存在于输出流的有序字节流中，则不用写入
	if index+writingLength < streamReassembler.unassembledLocationInWholeStream {
		return
	}
	// 最远可写入的字节位置实际上由byteStream的输出当前读了多少决定的，在此处更新最远可写入字节的位置
	bytesRead := streamReassembler.output.Bytes_read()
	streamReassembler.currentLargerThanCapIndex = bytesRead + streamReassembler.capacity
	// dataBytes 完整/或在超出capacity时丢失 写入到unassembledBuffer中
	for i, content := range dataBytes {
		contentInWholeStreamLoc := index + i
		writeInBufferLocBeforeMod := contentInWholeStreamLoc - streamReassembler.unassembledLocationInWholeStream +
			streamReassembler.unassembledBufferStartPtr
		if writeInBufferLocBeforeMod >= streamReassembler.currentLargerThanCapIndex {
			break
		}
		writeInBufferLoc := writeInBufferLocBeforeMod % streamReassembler.capacity
		streamReassembler.unassembledBuffer[writeInBufferLoc] = content
		streamReassembler.unAssembledBytes++
		streamReassembler.unassembledBufferHasCont[writeInBufferLoc] = true
	}

	// 尝试将已经有序的字节写入到输出流当中
	streamReassembler.tryWriteToOutput()
}

func (streamReassembler *StreamReassembler) tryWriteToOutput() {
	// 从指向buffer实际初始位置的指针开始，为真则将该位置的信息写入到output中，并将该位置置false。
	// 同时修正指向初始位置的指针
	var writeToOutput string
	for streamReassembler.unassembledBufferHasCont[streamReassembler.unassembledBufferStartPtr] {
		writeToOutput += string(streamReassembler.unassembledBuffer[streamReassembler.unassembledBufferStartPtr])
		streamReassembler.unassembledBufferHasCont[streamReassembler.unassembledBufferStartPtr] = false
		streamReassembler.unassembledBufferStartPtr = (streamReassembler.unassembledBufferStartPtr + 1) % streamReassembler.capacity
		streamReassembler.unAssembledBytes--
	}
	streamReassembler.output.Write(writeToOutput)
}

func (streamReassembler *StreamReassembler) Stream_out() *bs.ByteStream {
	return &streamReassembler.output
}

func (streamReassembler *StreamReassembler) Unassembled_bytes() int {
	return streamReassembler.unAssembledBytes
}

// 返回储存无序字节的buffer是否为空（不考虑已经写到output的部分）
func (streamReassembler *StreamReassembler) Empty() bool {
	return streamReassembler.unAssembledBytes == 0
}
