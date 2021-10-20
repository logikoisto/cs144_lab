package main

import (
	sr "cs144/streamreassembler"
	rcv "cs144/tcpreceiver"
	"cs144/util"
	"fmt"
)

func test_streamReassembler() {
	reassembler := sr.NewStreamReassembler(10)
	output := reassembler.Stream_out()
	reassembler.Push_substring("llo World", 2, false)
	reassembler.Push_substring("Hell", 0, false)
	firstRead := output.Read(8)
	fmt.Println(firstRead)
	fmt.Println("expected: Hello Wo")
	secondRead := output.Read(5)
	fmt.Println(secondRead)
	fmt.Println("expected: rl")
	reassembler.Push_substring("tea", 10, false)
	thirdRead := output.Read(4)
	fmt.Println(thirdRead)
	fmt.Println("expected: tea")
}

func test_wrappingInt32() {
	var seqno rcv.WrappingInt32 = 0x100000000 - 2
	var absseqno uint64 = 0
	var isn rcv.WrappingInt32 = 0x100000000 - 2
	fmt.Println("wrap, convert absseqno -> seqno", rcv.Wrap(absseqno, isn), ", expected: ", seqno)
	fmt.Println("unwrap, convert seqno -> abssequno", rcv.UnWrap(seqno, isn, 0), ", expected: ", absseqno)
	thirdSeqno := seqno + 3
	fmt.Println("thirdSeqno = ", thirdSeqno, ", expected: 1")
	fmt.Println("unwrap thirdSeqno, ", rcv.UnWrap(thirdSeqno, isn, 0), ", expected: ", 3)
}

func testBuffer() {
	buf := util.NewBufferByStr("Hello")
	fmt.Println(buf.Str(), ", expected: Hello")
	fmt.Println(buf.Size(), ", expected: 5")
	buf.RemovePrefix(3)
	fmt.Println(buf.Str(), ", expected: lo")
	fmt.Println(buf.Size(), ", expected: 2")
}

func main() {
	test_streamReassembler()
	test_wrappingInt32()
	fmt.Println("----------")
	testBuffer()
}
