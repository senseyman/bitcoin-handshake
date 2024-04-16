package service

import (
	"encoding/binary"
	"io"

	"github.com/senseyman/bitcoin-handshake/model"
)

type DecodeService struct {
}

func NewDecodeService() *DecodeService {
	return &DecodeService{}
}

func (s *DecodeService) DecodeElements(r io.Reader, elements ...any) error {
	for _, element := range elements {
		err := s.decodeElement(r, element)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *DecodeService) decodeElement(r io.Reader, element any) error {
	switch e := element.(type) {
	case *int32:
		rv, err := s.uint32(r, littleEndian)
		if err != nil {
			return err
		}
		*e = int32(rv)
		return nil

	case *uint32:
		rv, err := s.uint32(r, littleEndian)
		if err != nil {
			return err
		}
		*e = rv
		return nil

	case *int64:
		rv, err := s.uint64(r, littleEndian)
		if err != nil {
			return err
		}
		*e = int64(rv)
		return nil

	case *uint64:
		rv, err := s.uint64(r, littleEndian)
		if err != nil {
			return err
		}
		*e = rv
		return nil

	case *bool:
		rv, err := s.uint8(r)
		if err != nil {
			return err
		}
		if rv == 0x00 {
			*e = false
		} else {
			*e = true
		}
		return nil

	// Message header checksum.
	case *[4]byte:
		_, err := io.ReadFull(r, e[:])
		if err != nil {
			return err
		}
		return nil

	// Message header command.
	case *[model.CommandSize]uint8:
		_, err := io.ReadFull(r, e[:])
		if err != nil {
			return err
		}
		return nil

	case *model.NetAddress:
		na, err := s.decodeNetAddress(r)
		if err != nil {
			return err
		}

		*e = na

		return nil
	case *string:
		_, err := io.ReadFull(r, []byte(*e))
		return err
	}

	return binary.Read(r, littleEndian, element)
}

func (s *DecodeService) decodeNetAddress(r io.Reader) (model.NetAddress, error) {
	var (
		services uint64
		ip       [16]byte
		port     uint16
	)

	buf := make([]byte, 8)

	if _, err := io.ReadFull(r, buf); err != nil {
		return model.NetAddress{}, err
	}
	services = littleEndian.Uint64(buf)

	if _, err := io.ReadFull(r, ip[:]); err != nil {
		return model.NetAddress{}, err
	}

	if _, err := io.ReadFull(r, buf[:2]); err != nil {
		return model.NetAddress{}, err
	}
	port = bigEndian.Uint16(buf[:2])

	return model.NetAddress{
		Services: services,
		IP:       ip[:],
		Port:     port,
	}, nil
}

func (s *DecodeService) uint64(r io.Reader, byteOrder binary.ByteOrder) (uint64, error) {
	buf := make([]byte, 8)

	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	rv := byteOrder.Uint64(buf)

	return rv, nil
}

func (s *DecodeService) uint32(r io.Reader, byteOrder binary.ByteOrder) (uint32, error) {
	buf := make([]byte, 4)

	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	rv := byteOrder.Uint32(buf)

	return rv, nil
}

func (s *DecodeService) uint8(r io.Reader) (uint8, error) {
	buf := make([]byte, 1)

	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	rv := buf[0]

	return rv, nil
}
