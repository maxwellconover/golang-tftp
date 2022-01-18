// packet package provides serialization/deserialization of TFTP packets
package packet

import (
	"bytes"
	"encoding/binary"
	"errors"
)

var ErrPacketType = errors.New("packet has invalid type code")
var ErrPacketStructure = errors.New("packet structure is invalid")

// Type codes
const (
	RRQ   = uint16(iota + 1) //1
	WRQ                      //2
	DATA                     //3
	ACK                      //4
	ERROR                    //5
)

//Error Codes for ERROR packet
const (
	ErrNotDefinedTFTP            = uint16(iota) //0
	ErrFileNotFoundTFTP                         //1
	ErrAccessViolationTFTP                      //2
	ErrDiskFullAllocExceededTFTP                //3
	ErrIllegalOperationTFTP                     //4
	ErrUnknownTransferIdTFTP                    //5
	ErrFileAlreadyExistsTFTP                    //6
	ErrNoSuchUserTFTP                           //7
)

const (
	BlockLength    = 512
	DatagramLength = 516
)

type Packet interface {
	TypeCode() uint16
	Serialize() []byte
}

// RRQ and WRQ Packet types
type ReqPacket struct {
	Type     uint16
	Filename string
	Mode     string
}

func (p *ReqPacket) TypeCode() uint16 {
	return p.Type
}

func (p *ReqPacket) Serialize() []byte {
	var b []byte

	//opcode is Type code in big endian
	opcode := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, p.TypeCode())

	// *-------------RRQ/WRQ Packet Format-------------*
	//
	//  2 bytes     string    1 byte     string   1 byte
	//  ------------------------------------------------
	// | Opcode |  Filename  |   0  |    Mode    |   0  |
	//  ------------------------------------------------
	b = append(b, opcode...)
	b = append(b, p.Filename...)
	b = append(b, 0)
	b = append(b, p.Mode...)
	b = append(b, 0)

	return b
}

type DataPacket struct {
	Type        uint16
	BlockNumber uint16
	Data        []byte
}

func (p *DataPacket) TypeCode() uint16 {
	return p.Type
}

func (p *DataPacket) Serialize() []byte {
	var b []byte

	//opcode and bnum are Type code and BlockNumber in big endian
	opcode := make([]byte, 2)
	bnum := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, p.TypeCode())
	binary.BigEndian.PutUint16(bnum, p.BlockNumber)

	// *--------DATA Packet Format--------*
	//
	//  2 bytes     2 bytes      n bytes
	//  ----------------------------------
	// | Opcode |   Block #  |   Data     |
	//  ----------------------------------
	b = append(b, opcode...)
	b = append(b, bnum...)
	b = append(b, p.Data...)

	return b
}

type AckPacket struct {
	Type        uint16
	BlockNumber uint16
}

func (p *AckPacket) TypeCode() uint16 {
	return p.Type
}

func (p *AckPacket) Serialize() []byte {
	var b []byte

	//opcode and BNum are Type code and BlockNumber in big endian
	opcode := make([]byte, 2)
	bnum := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, p.TypeCode())
	binary.BigEndian.PutUint16(bnum, p.BlockNumber)

	// *--ACK Packet Format--*
	//
	//  2 bytes     2 bytes
	//  ---------------------
	// | Opcode |   Block #  |
	//  ---------------------
	b = append(b, opcode...)
	b = append(b, bnum...)

	return b
}

type ErrPacket struct {
	Type    uint16
	ErrCode uint16
	ErrMsg  string
}

func (p *ErrPacket) TypeCode() uint16 {
	return p.Type
}

func (p *ErrPacket) Serialize() []byte {
	var b []byte

	//opcode and errorcode are Type code and ErrCode in big endian
	opcode := make([]byte, 2)
	errorcode := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, p.TypeCode())
	binary.BigEndian.PutUint16(errorcode, p.ErrCode)

	// *-----------ERROR Packet Format-----------*
	//
	//  2 bytes     2 bytes      string    1 byte
	//  -----------------------------------------
	// | Opcode |  ErrorCode |   ErrMsg   |   0  |
	//  -----------------------------------------
	b = append(b, opcode...)
	b = append(b, errorcode...)
	b = append(b, p.ErrMsg...)
	b = append(b, 0)

	return b
}

func PacketDeserialize(b []byte) (Packet, error) {
	if len(b) < 4 {
		return nil, ErrPacketStructure
	}
	opcode := binary.BigEndian.Uint16(b[0:2])

	switch opcode {
	case RRQ, WRQ:
		vals := bytes.Split(b[2:], []byte{0})
		if len(vals) < 2 {
			return nil, ErrPacketStructure
		}
		return &ReqPacket{
			Type:     opcode,
			Filename: string(vals[0]),
			Mode:     string(vals[1]),
		}, nil
	case ACK:
		blocknum := binary.BigEndian.Uint16(b[2:4])
		return &AckPacket{
			Type:        opcode,
			BlockNumber: blocknum,
		}, nil
	case DATA:
		blocknum := binary.BigEndian.Uint16(b[2:4])
		return &DataPacket{
			Type:        opcode,
			BlockNumber: blocknum,
			Data:        b[4:],
		}, nil
	case ERROR:
		if len(b) < 5 {
			return nil, ErrPacketStructure
		}
		errcode := binary.BigEndian.Uint16(b[2:4])
		return &ErrPacket{
			Type:    opcode,
			ErrCode: errcode,
			ErrMsg:  string(b[4 : len(b)-1]),
		}, nil
	default:
		return nil, ErrPacketType
	}
}
