// packet package provides serialization/deserialization of TFTP packets
package packet

import (
	"bytes"
	"testing"
)

// *-------------RRQ/WRQ Packet Format-------------*
//
//  2 bytes     string    1 byte     string   1 byte
//  ------------------------------------------------
// | Opcode |  Filename  |   0  |    Mode    |   0  |
//  ------------------------------------------------
func TestRrqSerializeDeserialize(t *testing.T) {
	testfname := "testfile.txt"
	testmode := "octet"
	rrq := ReqPacket{
		Type:     uint16(1),
		Filename: testfname,
		Mode:     testmode,
	}
	rrqbytes := rrq.Serialize()
	rrqi, err := PacketDeserialize(rrqbytes)
	if err != nil {
		t.Fatal(err)
	}
	rrqpkt := rrqi.(*ReqPacket) //Retrieve RRQ type from Packet interface
	if rrqpkt.TypeCode() != RRQ {
		t.Fatalf("incorrect type - RRQ packet code should be %v not %v", RRQ, rrqpkt.TypeCode())
	}
	if rrqpkt.Filename != testfname {
		t.Fatalf("incorrect filename - Filename should be %v not %v", testfname, rrqpkt.Filename)
	}
	if rrqpkt.Mode != testmode {
		t.Fatalf("incorrect mode - Mode should be %v not %v", testmode, rrqpkt.Mode)
	}
}

// *-------------RRQ/WRQ Packet Format-------------*
//
//  2 bytes     string    1 byte     string   1 byte
//  ------------------------------------------------
// | Opcode |  Filename  |   0  |    Mode    |   0  |
//  ------------------------------------------------
func TestWrqSerializeDeserialize(t *testing.T) {
	testfname := "testfile.txt"
	testmode := "octet"
	wrq := ReqPacket{
		Type:     uint16(2),
		Filename: testfname,
		Mode:     testmode,
	}
	wrqbytes := wrq.Serialize()
	wrqi, err := PacketDeserialize(wrqbytes)
	if err != nil {
		t.Fatal(err)
	}
	wrqpkt := wrqi.(*ReqPacket)
	if wrqpkt.TypeCode() != WRQ {
		t.Fatalf("incorrect type - WRQ packet code should be %v not %v", WRQ, wrqpkt.TypeCode())
	}
	if wrqpkt.Filename != testfname {
		t.Fatalf("incorrect filename - Filename should be %v not %v", testfname, wrqpkt.Filename)
	}
	if wrqpkt.Mode != testmode {
		t.Fatalf("incorrect mode - Mode should be %v not %v", testmode, wrqpkt.Mode)
	}
}

// *--------DATA Packet Format--------*
//
//  2 bytes     2 bytes      n bytes
//  ----------------------------------
// | Opcode |   Block #  |   Data     |
//  ----------------------------------
func TestDataSerializeDeserialize(t *testing.T) {
	testbnum := uint16(48)
	testdata := []byte("Testing a Data packet")
	data := DataPacket{
		Type:        uint16(3),
		Data:        testdata,
		BlockNumber: testbnum,
	}
	databytes := data.Serialize()
	datai, err := PacketDeserialize(databytes)
	if err != nil {
		t.Fatal(err)
	}
	datapkt := datai.(*DataPacket)
	if datapkt.TypeCode() != DATA {
		t.Fatalf("incorrect type - DATA packet code should be %v not %v", DATA, datapkt.TypeCode())
	}
	if !bytes.Equal(datapkt.Data, testdata) {
		t.Fatalf("incorrect data - Data value should be %v not %v", testdata, datapkt.Data)
	}
	if datapkt.BlockNumber != testbnum {
		t.Fatalf("incorrect blocknumber - Block Number should be %v not %v", testbnum, datapkt.BlockNumber)
	}
}

// *--ACK Packet Format--*
//
//  2 bytes     2 bytes
//  ---------------------
// | Opcode |   Block #  |
//  ---------------------
func TestAckSerializeDeserialize(t *testing.T) {
	testbnum := uint16(14)
	ack := AckPacket{
		Type:        uint16(4),
		BlockNumber: testbnum,
	}
	ackbytes := ack.Serialize()
	acki, err := PacketDeserialize(ackbytes)
	if err != nil {
		t.Fatal(err)
	}
	ackpkt := acki.(*AckPacket)
	if ackpkt.TypeCode() != ACK {
		t.Fatalf("incorrect type - ACK packet code should be %v not %v", ACK, ackpkt.TypeCode())
	}
	if ackpkt.BlockNumber != testbnum {
		t.Fatalf("incorrect blocknumber - Block Number should be %v not %v", testbnum, ackpkt.BlockNumber)
	}
}

// *-----------ERROR Packet Format-----------*
//
//  2 bytes     2 bytes      string    1 byte
//  -----------------------------------------
// | Opcode |  ErrorCode |   ErrMsg   |   0  |
//  -----------------------------------------
func TestErrorSerializeDeserialize(t *testing.T) {
	testerrcode := ErrIllegalOperationTFTP
	testerrmsg := "illegal tftp operation"
	errorvals := ErrPacket{
		Type:    uint16(5),
		ErrCode: testerrcode,
		ErrMsg:  testerrmsg,
	}
	errorbytes := errorvals.Serialize()
	errori, err := PacketDeserialize(errorbytes)
	if err != nil {
		t.Fatal(err)
	}
	errorpkt := errori.(*ErrPacket)
	if errorpkt.TypeCode() != ERROR {
		t.Fatalf("incorrect type - ERROR packet code should be %v not %v", ERROR, errorpkt.TypeCode())
	}
	if errorpkt.ErrCode != testerrcode {
		t.Fatalf("incorrect error code - Error code should be %v not %v", testerrcode, errorpkt.ErrCode)
	}
	if errorpkt.ErrMsg != testerrmsg {
		t.Fatalf("incorrect error message - Error message should be %v not %v", testerrmsg, errorpkt.ErrMsg)
	}
}
