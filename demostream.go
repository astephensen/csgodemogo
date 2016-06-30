package csgodemogo

import (
	"encoding/binary"
	"io"
)

type demoStream struct {
	r   io.ReadSeeker
	pos int
}

func (stream *demoStream) GetVarInt() uint64 {
	var x uint64
	var s uint
	buf := make([]byte, 1)
	for i := 0; ; i++ {
		_, err := stream.r.Read(buf)
		stream.pos++
		if err != nil {
			panic(err)
		}
		if buf[0] < 0x80 {
			if i > 9 || i == 9 && buf[0] > 1 {
				panic("overflow")
			}
			return x | uint64(buf[0])<<s
		}
		x |= uint64(buf[0]&0x7f) << s
		s += 7
	}
}

func (stream *demoStream) GetCurrentOffset() int {
	return stream.pos
}

func (stream *demoStream) GetByte() byte {
	buf := make([]byte, 1)
	n, err := stream.r.Read(buf)
	if err != nil {
		panic(err)
	}
	stream.pos += n
	return buf[0]
}

func (stream *demoStream) GetInt() int32 {
	var x int32
	err := binary.Read(stream.r, binary.LittleEndian, &x)
	if err != nil {
		panic(err)
	}
	stream.pos += 4
	return x
}

func DemoStream(reader io.ReadSeeker) *demoStream {
	stream := demoStream{r: reader, pos: 0}
	return &stream
}

func (stream *demoStream) Read(out []byte) (int, error) {
	n, err := stream.r.Read(out)
	stream.pos += n
	return n, err
}

func (stream *demoStream) Skip(n int64) {
	stream.pos += int(n)
	stream.r.Seek(n, 1)
}
