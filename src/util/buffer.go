package util

type Buffer struct {
	storage string
	offset  int
}

type BufferList struct {
	 buffers []Buffer
}

func NewBuffer() *Buffer {
	return &Buffer{}
}

func NewBufferByStr(str string) *Buffer {
	return &Buffer{
		storage: str,
		offset:  0,
	}
}

func (buf *Buffer) Str() string {
	return buf.storage[buf.offset:]
}

func (buf *Buffer) At(n int) uint8 {
	return buf.Str()[n]
}

func (buf *Buffer) Size() int {
	return len(buf.Str())
}

func (buf *Buffer) Copy() string {
	return buf.Str()
}

func (buf *Buffer) RemovePrefix(n int) {
	if n > buf.Size() {
		panic("out of range in RemovePrefix")
	}

	buf.offset += n
}

func NewBufferList() *BufferList {
	return &BufferList{}
}

func NewBufferListByBuffer(buffer *Buffer) *BufferList {
	bufList := BufferList{}
	bufList.buffers = append(bufList.buffers, *buffer)
	return &bufList
}

func NewBufferListByString(str string) *BufferList {
	buffer := NewBufferByStr(str)
	return NewBufferListByBuffer(buffer)
}

func (bufList *BufferList) Buffers() []Buffer {
	return bufList.buffers
}

func (bufList *BufferList) Append(other *BufferList) {
	for _, buf := range other.buffers {
		bufList.buffers = append(bufList.buffers, buf)
	}
}