package git

import (
	"fmt"
	"io"
)

type pktStatus int

const (
	pktInvalid pktStatus = iota
	pktNromal
	pktFlush
	pktDelim
	pktResponseEnd
)

func (ps pktStatus) String() string {
	switch ps {
	case pktInvalid:
		return "invalid"
	case pktNromal:
		return "normal"
	case pktFlush:
		return "flush"
	case pktDelim:
		return "delim"
	case pktResponseEnd:
		return "resp-end"
	default:
		return "<unknown>"
	}
}

type gitWireReader struct {
	sizeBuf [4]byte
	reader  io.Reader

	remainder int
}

func (r *gitWireReader) readSize() (int, pktStatus, error) {
	return readPktSize(r.reader, r.sizeBuf[:4])
}

func (r *gitWireReader) SideBandReader() (*SideBandReader, error) {
	if r.remainder != 0 {
		return nil, fmt.Errorf("invalid last packet read not finished")
	}

	return &SideBandReader{
		sizeBuf: r.sizeBuf[:],
		reader:  r.reader,
	}, nil
}

func (r *gitWireReader) ReadFlush() error {
	_, status, err := r.readSize()
	if err != nil {
		return err
	}

	if status != pktFlush {
		return fmt.Errorf("not flush packet: %q", status)
	}

	return nil
}

func (r *gitWireReader) ReadPkt() ([]byte, error) {
	if r.remainder > 0 {
		return nil, fmt.Errorf("invalid last packet read not finished")
	}

readPkt:
	size, status, err := r.readSize()
	if err != nil {
		return nil, err
	}

	if status != pktNromal {
		// flush pkt
		goto readPkt
	}

	buf := make([]byte, size)
	_, err = io.ReadFull(r.reader, buf)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *gitWireReader) Read(p []byte) (int, error) {
	if r.remainder > 0 {
		// not finished reading last line
		maxRead := len(p)
		if maxRead > r.remainder {
			maxRead = r.remainder
		}

		n, err := io.ReadFull(r.reader, p[:maxRead])
		r.remainder -= n

		return n, err
	}

readSize:
	size, status, err := r.readSize()
	if err != nil {
		return 0, err
	}

	if status != pktNromal {
		// flush pkt
		goto readSize
	}

	maxRead := len(p)
	if maxRead > size {
		maxRead = size
	}

	n, err := io.ReadFull(r.reader, p[:maxRead])
	r.remainder = size - n
	if err != nil {
		return n, err
	}

	return n, nil
}
