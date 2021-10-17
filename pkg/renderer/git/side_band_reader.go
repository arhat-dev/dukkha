package git

import (
	"fmt"
	"io"
)

type sideBandType byte

const (
	sideBandPrimary     sideBandType = 0x01
	sideBandSecondary   sideBandType = 0x02
	sideBandRemoteError sideBandType = 0x03
)

type SideBandReader struct {
	sizeBuf []byte

	reader io.Reader

	remainder int
}

func (sbr *SideBandReader) readSizeAndBand() (int, pktStatus, error) {
	// (in git srouce code) when send_sideband called with band < 0
	// there is no band type after the size, so here we read as pkt-size
	// to avoid over reading
	return readPktSize(sbr.reader, sbr.sizeBuf)
}

func (sbr *SideBandReader) Read(p []byte) (int, error) {
	if sbr.remainder > 0 {
		// not finished reading last line
		maxRead := len(p)
		if maxRead > sbr.remainder {
			maxRead = sbr.remainder
		}

		n, err := io.ReadFull(sbr.reader, p[:maxRead])
		sbr.remainder -= n

		if sbr.remainder == 0 {
			p[maxRead-1] = 0x00
		}

		return n, err
	}

readSize:
	size, status, err := sbr.readSizeAndBand()
	if err != nil {
		return 0, err
	}

	if status != pktNromal {
		goto readSize
	}

	maxRead := len(p)
	if maxRead == 0 {
		// for connectivity testing
		return sbr.reader.Read(p)
	}

	// read one byte first to determine whether it's band type
	backupFirstByte := p[0]
	_, err = sbr.reader.Read(p[:1])
	if err != nil {
		p[0] = backupFirstByte
		return 0, err
	}

	band := sideBandType(p[0] & 0xff)
	switch band {
	case sideBandPrimary:
		// TODO: meant for terminal message, the side band ends after this pkt is read
	case sideBandSecondary:
		// TODO: continue (usual case)
	case sideBandRemoteError:
		// TODO: print error?
	default:
		p[0] = backupFirstByte
		// protocol error
		return 0, fmt.Errorf("invalid side band type %q", band)
	}

	size--
	if maxRead > size {
		maxRead = size
	}

	n, err := io.ReadFull(sbr.reader, p[:maxRead])
	sbr.remainder = size - n
	if err != nil {
		return 0, err
	}

	if sbr.remainder == 0 {
		p[maxRead-1] = 0x00
	}

	return n, nil
}
