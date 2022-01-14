// packet package provides serialization/deserialization of TFTP packets
package packet

import (
	"encoding/binary"
)

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

type Packet interface {
	// Serialize packet
	ToBytes() []byte
}

// RRQ and WRQ Packet types
type ReqPacket struct {
	TypeCode uint16
	Filename string
	Mode     string
	Length   int
}

func (p *ReqPacket) ToBytes() []byte {
	var b []byte

	//opcode is Type code in big endian
	opcode := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, p.TypeCode)

	// *-------------RRQ/WRQ Header Format-------------*
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
	TypeCode    uint16
	Data        []byte
	BlockNumber uint16
}

func (p *DataPacket) ToBytes() []byte {
	var b []byte

	//opcode and bnum are Type code and BlockNumber in big endian
	opcode := make([]byte, 2)
	bnum := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, p.TypeCode)
	binary.BigEndian.PutUint16(bnum, p.BlockNumber)

	// *--------DATA Header Format--------*
	//
	//  2 bytes     2 bytes      n bytes
	//  ----------------------------------
	// | Opcode |   Block #  |   Data     |
	//  ----------------------------------
	b = append(b, opcode...)
	b = append(b, byte(p.BlockNumber))
	b = append(b, p.Data...)

	return b
}

type AckPacket struct {
	TypeCode    uint16
	BlockNumber uint16
}

func (p *AckPacket) ToBytes() []byte {
	var b []byte

	//opcode and BNum are Type code and BlockNumber in big endian
	opcode := make([]byte, 2)
	bnum := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, p.TypeCode)
	binary.BigEndian.PutUint16(bnum, p.BlockNumber)

	// *--ACK Header Format--*
	//
	//  2 bytes     2 bytes
	//  ---------------------
	// | Opcode |   Block #  |
	//  ---------------------
	b = append(b, opcode...)
	b = append(b, byte(p.BlockNumber))

	return b
}

type ErrPacket struct {
	TypeCode uint16
	ErrCode  uint16
	ErrMsg   string
}

func (p *ErrPacket) ToBytes() []byte {
	var b []byte

	//opcode and errorcode are Type code and ErrCode in big endian
	opcode := make([]byte, 2)
	errorcode := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, p.TypeCode)
	binary.BigEndian.PutUint16(errorcode, p.ErrCode)

	// *-----------ERROR Header Format-----------*
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

//TODO

// func PacketDeserialize(b []byte) (Packet, error) {
// 	var ErrPacketInvalid = errors.New("packet structure is invalid")

// 	if len(b) < 2 {
// 		return nil, ErrPacketInvalid
// 	}

// 	opcode := binary.BigEndian.Uint16(b[0:2])
// 	switch opcode {
// 	case RRQ, WRQ:
// 		vals := bytes.Split(b[2:], []byte{0})
// 		if len(vals) < 2 {
// 			return nil, ErrPacketInvalid
// 		}
// 		return &ReqPacket{
// 			TypeCode: opcode,
// 			Filename: string(vals[0]),
// 			Mode:     string(vals[1]),
// 		}, nil
// 	}
// }
