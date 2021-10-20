package tcpreceiver

import "math"

type WrappingInt32 uint32

// 序列号：第一个SYN的序列号为isn，后续字节从isn+1开始；FIN也占一个序列号。32位，超出最大从0开始，相当于mod2^32
// 绝对序列号：第一个SYN的序列号为0, 后续字节从1开始，FIN也占一个绝对序列号。64位，可以到2^64，被认为是无穷多的，因为这样多的字节按照100Gb/s的速率传输需要50年。
// 流索引：不计算SYN和FIN，从传输字节流的第一个字节索引0开始计数。同样也是64位，不考虑SYN，FIN

// 从绝对序列号->序列号
func Wrap(n uint64, isn WrappingInt32) WrappingInt32 {
	return WrappingInt32(n) + isn
}

// 从序列号->绝对序列号
func UnWrap(n, isn WrappingInt32, checkpoint uint64) uint64 {
	absoluteSeqno := uint64(n - isn)
	for math.Abs(float64(absoluteSeqno-checkpoint)) > 0x10000000 {
		absoluteSeqno += 0x100000000
	}
	return absoluteSeqno
}
