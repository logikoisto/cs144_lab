package tcphelpers

import (
	tcpr "cs144/tcpreceiver"
)

type TCPHeader struct {
	length int                // 报头长度，无option为20字节
	sport  uint16             // 源端口
	dport  uint16             // 目标端口
	seqno  tcpr.WrappingInt32 // sequence number
	ackno  tcpr.WrappingInt32 // acknowledge number
	doff   uint8              // data offset
	utg    bool               // urgent flag
	ack    bool               // ack flag
	psh    bool               // push flag
	rst    bool               // rst flag
	syn    bool               // syn flag
	fin    bool               // fin flag
	win    uint16             // window size
	cksum  uint16             // check sum
	uptr   uint16             // urgentptr
}
