// packet package provides serialization/deserialization of TFTP packets
package packet

import (
	"encoding/binary"
)

// Type codes
const (
	RRQ   = uint16(iota + 1) // 1
	WRQ                      //2
	ACK                      //3
	ERROR                    //4
	OACK                     //5
)

type Packet interface {
	// Return packet type
	GetType() uint16

	// Serialize packet
	ToBytes() []byte
}

// RRQ and WRQ Packet types
type ReqPacket struct {
	Type     uint16
	Filename string
	Mode     string
	Length   int
}

func (p *ReqPacket) GetType() uint16 {
	return p.Type
}

func (p *ReqPacket) ToBytes() []byte {
	var b []byte

	//opcode is Type code in big endian
	opcode := make([]byte, 2)
	binary.BigEndian.PutUint16(opcode, p.GetType())

	// *-------------RRQ/WRQ Header format-------------*
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
