package inMemByteStream

// 在本机上提供一个可读/写的字节流
// startPtr == endPtr时，存在两种情况：缓冲区满，或缓冲区空。判断的依据是，若capacity == 0，则缓冲区满；否则空。
type ByteStream struct {
	buffer      []byte
	startPtr    int  // 当前读开始位置
	endPtr      int  // 当前写开始位置
	capacity    int  // 当前剩余可写容量
	bufferSize  int  // 缓冲区总大小
	writeDone   bool // 写结束标识
	err         bool // 出错标识
	bytesWriten int  // 已写字节总数
	bytesRead   int  // 已读字节总数
}

func NewByteStream(capacity int) *ByteStream {
	if capacity <= 0 {
		panic("capacity should be larger than 0")
	}
	return &ByteStream{
		buffer:      make([]byte, capacity),
		startPtr:    0,
		endPtr:      0,
		capacity:    capacity,
		bufferSize:  capacity,
		writeDone:   false,
		err:         false,
		bytesWriten: 0,
		bytesRead:   0,
	}
}

func (byteStream *ByteStream) Write(data string) int {
	dataBytes := []byte(data)
	dataLength := len(dataBytes)
	curCap := byteStream.capacity
	for _, content := range dataBytes {
		if byteStream.capacity == 0 {
			break
		}
		byteStream.buffer[byteStream.endPtr] = content
		byteStream.endPtr = (byteStream.endPtr + 1) % byteStream.bufferSize
		byteStream.capacity--
	}

	if dataLength <= curCap {
		byteStream.bytesWriten += dataLength
		return dataLength
	} else {
		byteStream.bytesWriten += curCap
		return curCap
	}
}

func (byteStream *ByteStream) Remaining_capacity() int {
	return byteStream.capacity
}

// 考虑到缓冲区满时也可能调用这个函数，在buffer中写入结束标志可能会覆盖数据，因此使用一个标志位来标识写的完成。
// 读时，若 bufferSize - capacity == 0, 且写完成标识位为真，则表示已经读取完了所有文件。
func (byteStream *ByteStream) End_input() {
	byteStream.writeDone = true
}

func (byteStream *ByteStream) Set_error() {
	byteStream.err = true
}

// 如果peek的长度超过了已写的长度，调整为peek所有写入的数据
func (byteStream *ByteStream) Peek_output(len int) string {
	var res string
	if len > byteStream.bufferSize-byteStream.capacity {
		len = byteStream.bufferSize - byteStream.capacity
	}
	ptr := byteStream.startPtr
	for i := 0; i < len; i++ {
		res += string(byteStream.buffer[ptr])
		ptr++
	}
	return res
}

// 如果pop的长度超过了已写的长度，调整为pop掉所有的写入数据
// 考虑直接丢掉的数据不算进已读字节数中，这里不对已读字节总数做出调整
func (byteStream *ByteStream) Pop_output(len int) {
	if len > byteStream.bufferSize-byteStream.capacity {
		len = byteStream.bufferSize - byteStream.capacity
	}
	byteStream.startPtr = (byteStream.startPtr + len) % byteStream.bufferSize
	byteStream.capacity += len
}

func (byteStream *ByteStream) Read(len int) string {
	if len > byteStream.bufferSize-byteStream.capacity {
		len = byteStream.bufferSize - byteStream.capacity
	}
	byteStream.bytesRead += len
	res := byteStream.Peek_output(len)
	byteStream.Pop_output(len)
	return res
}

func (byteStream *ByteStream) Input_ended() bool {
	return byteStream.capacity == byteStream.bufferSize
}

func (byteStream *ByteStream) Eof() bool {
	if byteStream.writeDone && byteStream.capacity == byteStream.bufferSize {
		return true
	}
	return false
}

func (byteStream *ByteStream) Error() bool {
	return byteStream.err
}

func (byteStream *ByteStream) Buffer_size() int {
	return byteStream.bufferSize
}

func (byteStream *ByteStream) Buffer_empty() bool {
	return byteStream.capacity == byteStream.bufferSize
}

func (byteStream *ByteStream) Bytes_written() int {
	return byteStream.bytesWriten
}

func (byteStream *ByteStream) Bytes_read() int {
	return byteStream.bytesRead
}

// func main() {
// 	byteStream := NewByteStream(5)
// 	fmt.Println("byteStream capacity = ", byteStream.capacity, ", expected: 5")
// 	writeLen := byteStream.write("Hel")
// 	fmt.Println("writing length = ", writeLen, ", expected:3 ")
// 	fmt.Println("byteStream capacity = ", byteStream.capacity, ", expected: 2")
// 	writeLen = byteStream.write("loo")
// 	fmt.Println("writing length = ", writeLen, ", expected: 2")
// 	fmt.Println("byteStream capacity = ", byteStream.capacity, ", expected: 0")
// 	fmt.Println("buffer content: ")
// 	for _, content := range byteStream.buffer {
// 		fmt.Print(string(content))
// 	}
// 	fmt.Println()
// 	fmt.Println("call remainning_capacity, result = ", byteStream.remaining_capacity(), ", expected: 0")
// 	byteStream.end_input()
// 	fmt.Println("after end input, write_done = ", byteStream.writeDone)
// 	fmt.Println("peek out put, len = 3: ", byteStream.peek_output(3), ", expected: Hel")
// 	fmt.Println("peek out put, len = 6: ", byteStream.peek_output(6), ", expected: Hello")
// 	byteStream.pop_output(2)
// 	fmt.Println("after pop_output(2), start pointer = ", byteStream.startPtr, ", expected: 2")
// 	fmt.Println("peek out put, len = 2: ", byteStream.peek_output(2), ", expected: ll")
// 	fmt.Println("read 2: ", byteStream.read(2), ", then startPtr = ", byteStream.startPtr, ", expected: 4")
// 	fmt.Println("Now capcity = ", byteStream.capacity, ", expected: 4")
// 	fmt.Println("error() = ", byteStream.error(), ", expectd: false")
// 	byteStream.set_error()
// 	fmt.Println("after set error, error() = ", byteStream.error(), ", expected: true")
// 	fmt.Println("eof() = ", byteStream.eof(), ", expected: false")
// 	fmt.Println("input_ended() = ", byteStream.input_ended(), ", expected: false")
// 	fmt.Println("buffer_size() = ", byteStream.buffer_size(), ", expected: 5")
// 	fmt.Println("buffer_empty() = ", byteStream.buffer_empty(), ", expected: false")
// 	fmt.Println("bytes_written() = ", byteStream.bytes_written(), ", expected: 5")
// 	fmt.Println("bytes_read() = ", byteStream.bytes_read(), ", expected: 2")
// }
