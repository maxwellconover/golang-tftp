// packet package provides serialization/deserialization of TFTP packets
package packet

import (
	"testing"
)

func TestRrqSerializeDeserialize(t *testing.T) {
	testfname := "testfile.txt"
	testmode := "octet"
	rrq := ReqPacket{
		TypeCode: 01,
		Filename: testfname,
		Mode:     testmode,
	}
	rrqbytes := rrq.Serialize()
	rrqds, err := PacketDeserialize(rrqbytes)
	if err != nil {
		t.Fatal(err)
	}
	if rrqds.TypeCode != RRQ {
		t.Fatal("incorrect type - RRQ packet code should be %v not %v", RRQ, rrqds.TypeCode)
	}
	if rrqds.Filename != testfname {
		t.Fatal("incorrect filename - Filename should be %v not %v", testfname, rrqds.Filename)
	}
	if rrqds.Mode != testmode {
		t.Fatal("incorrect mode - Mode should be %v not %v", testmode, rrqds.Mode)
	}
}

func TestWrqSerializeDeserialize(t *testing.T) {
	testfname := "testfile.txt"
	testmode := "octet"
	wrq := ReqPacket{
		TypeCode: 02,
		Filename: testfname,
		Mode:     testmode,
	}
	wrqbytes := wrq.Serialize()
	wrqds, err := PacketDeserialize(wrqbytes)
	if err != nil {
		t.Fatal(err)
	}
	if wrqds.TypeCode != WRQ {
		t.Fatal("incorrect type - WRQ packet code should be %v not %v", RRQ, wrqds.TypeCode)
	}
	if wrqds.Filename != testfname {
		t.Fatal("incorrect filename - Filename should be %v not %v", testfname, wrqds.Filename)
	}
	if wrqds.Mode != testmode {
		t.Fatal("incorrect mode - Mode should be %v not %v", testmode, wrqds.Mode)
	}
}
