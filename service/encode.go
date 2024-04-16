package service

import (
	"encoding/binary"
	"io"
	"math"

	"github.com/senseyman/bitcoin-handshake/model"
)

type EncodeService struct {
}

func NewEncodeService() *EncodeService {
	return &EncodeService{}
}

func (s *EncodeService) EncodeVersionMessage(w io.Writer, msg model.VersionMessage) error {
	err := s.EncodeElements(w, msg.Version, msg.Services, msg.Timestamp)
	if err != nil {
		return err
	}

	err = s.encodeNetAddress(w, &msg.AddrRecv)
	if err != nil {
		return err
	}

	err = s.encodeNetAddress(w, &msg.AddrFrom)
	if err != nil {
		return err
	}

	err = s.encodeElement(w, msg.Nonce)
	if err != nil {
		return err
	}

	err = s.encodeVarString(w, msg.UserAgent)
	if err != nil {
		return err
	}

	err = s.encodeElement(w, msg.StartHeight)
	if err != nil {
		return err
	}

	return nil
}

func (s *EncodeService) EncodeElements(w io.Writer, elements ...any) error {
	for _, element := range elements {
		err := s.encodeElement(w, element)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *EncodeService) encodeElement(w io.Writer, element any) error {
	switch e := element.(type) {
	case int32:
		err := s.putUint32(w, littleEndian, uint32(e))
		if err != nil {
			return err
		}
		return nil

	case uint32:
		err := s.putUint32(w, littleEndian, e)
		if err != nil {
			return err
		}
		return nil

	case int64:
		err := s.putUint64(w, littleEndian, uint64(e))
		if err != nil {
			return err
		}
		return nil

	case uint64:
		err := s.putUint64(w, littleEndian, e)
		if err != nil {
			return err
		}
		return nil

	// Message header checksum.
	case [4]byte:
		_, err := w.Write(e[:])
		if err != nil {
			return err
		}
		return nil

	// Message header command.
	case [model.CommandSize]uint8:
		err := s.putBytes(w, e[:])
		if err != nil {
			return err
		}
		return nil
	}

	return binary.Write(w, littleEndian, element)
}

func (s *EncodeService) encodeVarStringBuf(w io.Writer, str string, buf []byte) error {
	err := s.encodeVarIntBuf(w, uint64(len(str)), buf)
	if err != nil {
		return err
	}

	_, err = w.Write([]byte(str))
	return err
}

func (s *EncodeService) encodeVarString(w io.Writer, str string) error {
	buf := make([]byte, 8)

	err := s.encodeVarStringBuf(w, str, buf)
	return err
}

func (s *EncodeService) encodeNetAddress(w io.Writer, na *model.NetAddress) error {
	buf := make([]byte, 8)
	err := s.encodeNetAddressBuf(w, na, buf)

	return err
}

func (s *EncodeService) encodeNetAddressBuf(w io.Writer, na *model.NetAddress, buf []byte) error {
	littleEndian.PutUint64(buf, na.Services)
	if _, err := w.Write(buf); err != nil {
		return err
	}

	// Ensure to always write 16 bytes even if the ip is nil.
	var ip [16]byte
	if na.IP != nil {
		copy(ip[:], na.IP.To16())
	}
	if _, err := w.Write(ip[:]); err != nil {
		return err
	}

	// Sigh.  Bitcoin protocol mixes little and big endian.
	bigEndian.PutUint16(buf[:2], na.Port)
	_, err := w.Write(buf[:2])

	return err
}

func (s *EncodeService) encodeVarIntBuf(w io.Writer, val uint64, buf []byte) error {
	switch {
	case val < 0xfd:
		buf[0] = uint8(val)
		_, err := w.Write(buf[:1])
		return err

	case val <= math.MaxUint16:
		buf[0] = 0xfd
		littleEndian.PutUint16(buf[1:3], uint16(val))
		_, err := w.Write(buf[:3])
		return err

	case val <= math.MaxUint32:
		buf[0] = 0xfe
		littleEndian.PutUint32(buf[1:5], uint32(val))
		_, err := w.Write(buf[:5])
		return err

	default:
		buf[0] = 0xff
		if _, err := w.Write(buf[:1]); err != nil {
			return err
		}

		littleEndian.PutUint64(buf, val)
		_, err := w.Write(buf)
		return err
	}
}

func (s *EncodeService) putUint32(w io.Writer, byteOrder binary.ByteOrder, val uint32) error {
	buf := make([]byte, 4)

	byteOrder.PutUint32(buf, val)
	_, err := w.Write(buf)

	return err
}

func (s *EncodeService) putUint64(w io.Writer, byteOrder binary.ByteOrder, val uint64) error {
	buf := make([]byte, 8)

	byteOrder.PutUint64(buf, val)
	_, err := w.Write(buf)

	return err
}

func (s *EncodeService) putBytes(w io.Writer, e []byte) error {
	_, err := w.Write(e)
	return err
}
