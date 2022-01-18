// package server implements a RFC 1350 tftp server
package server

import (
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	pkt "github.com/maxwellconover/golang-tftp/packet"
)

var Timeout = time.Second * 30
var Retry = time.Second * 3

var ErrTimeout = errors.New("response timed out")
var ErrBlockNum = errors.New("recieved unexpected block number")

type Reader func(filename string) (r io.Reader, err error)
type Writer func(filename string) (r io.Writer, err error)

type Server struct {
	Directory string
	ReadFunc  Reader
	WriteFunc Writer
}

func NewServer(directory string, rf Reader, wf Writer) *Server {
	return &Server{
		Directory: directory,
		ReadFunc:  rf,
		WriteFunc: wf,
	}
}

// TODO: Handle netascii and mail modes - currently only sends files in octet
func (s *Server) HandleClientReq(addr *net.UDPAddr, req pkt.Packet) error {
	reqpkt := req.(*pkt.ReqPacket)

	clientaddr, err := net.ResolveUDPAddr("udp", addr.String())
	if err != nil {
		return err
	}

	//Port 0 requests dynamic system-allocated port for server
	localaddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", localaddr, clientaddr)
	if err != nil {
		return err
	}

	// Send initial acknowledge packet
	ackPkt := pkt.AckPacket{
		Type:        pkt.ACK,
		BlockNumber: uint16(0),
	}
	_, err = conn.Write(ackPkt.Serialize())
	if err != nil {
		return err
	}

	switch reqpkt.TypeCode() {
	case pkt.RRQ:
		err := s.HandleRead(reqpkt, conn)
		if err != nil {
			return err
		}
	case pkt.WRQ:
		err := s.HandleWrite(reqpkt, conn)
		if err != nil {
			return err
		}
	default:
		return pkt.ErrPacketType
	}
	return nil
}

func (s *Server) HandleWrite(wrq *pkt.ReqPacket, conn *net.UDPConn) error {
	filewriter, err := s.WriteFunc(s.Directory + "/" + wrq.Filename)
	if err != nil {
		return err
	}

	// Recieve file transfer
	var b []byte
	curblknum := uint16(1)
	for {
		n, _, err := conn.ReadFromUDP(b)
		if err != nil {
			return err
		}

		pkti, err := pkt.PacketDeserialize(b[:n])
		if err != nil {
			return err
		}
		switch pkti.TypeCode() {
		case pkt.ERROR:
			return s.HandleErrPacket(pkti, conn)
		case pkt.DATA:
		default:
			return pkt.ErrPacketType
		}

		data := pkti.(*pkt.DataPacket)

		if data.BlockNumber == curblknum-1 {
			// ACK was not recieved by client
			ackPkt := pkt.AckPacket{
				Type:        pkt.ACK,
				BlockNumber: curblknum - 1,
			}
			_, err = conn.Write(ackPkt.Serialize())
			if err != nil {
				return err
			}
			continue
		} else if data.BlockNumber != curblknum {
			return errors.New("unexpected block number. stopping file transfer")
		}

		_, err = filewriter.Write(data.Data)
		if err != nil {
			return err
		}

		ackPkt := pkt.AckPacket{
			Type:        pkt.ACK,
			BlockNumber: curblknum,
		}
		_, err = conn.Write(ackPkt.Serialize())
		if err != nil {
			return err
		}

		// On last packet, wait to ensure final ack was recieved
		if len(data.Data) < pkt.BlockLength {
			success := false
			for !success {
				time.Sleep(time.Second * Timeout)
				_, _, err := conn.ReadFromUDP(b)
				if err != nil {
					success = true
				} else {
					if data.BlockNumber == curblknum-1 {
						// ACK was not recieved by client
						ackPkt := pkt.AckPacket{
							Type:        pkt.ACK,
							BlockNumber: curblknum - 1,
						}
						_, err = conn.Write(ackPkt.Serialize())
						if err != nil {
							return err
						}
						continue
					} else if data.BlockNumber != curblknum {
						return errors.New("unexpected block number. stopping file transfer")
					}

					_, err = filewriter.Write(data.Data)
					if err != nil {
						return err
					}

					ackPkt := pkt.AckPacket{
						Type:        pkt.ACK,
						BlockNumber: curblknum,
					}
					_, err = conn.Write(ackPkt.Serialize())
					if err != nil {
						return err
					}
				}
			}

			return nil
		}

		curblknum++
	}

}

func (s *Server) HandleRead(rrq *pkt.ReqPacket, conn *net.UDPConn) error {
	filereader, err := s.ReadFunc(s.Directory + "/" + rrq.Filename)
	if err != nil {
		return err
	}

	// Transmit file
	b := make([]byte, pkt.BlockLength)
	curblknum := uint16(1)
	for {
		n, readerr := filereader.Read(b)
		if readerr != nil && err != io.EOF {
			return err
		}

		b = b[:n]

		datapkt := &pkt.DataPacket{
			Data:        b,
			BlockNumber: curblknum,
		}

		err := s.SendData(datapkt, conn)
		if err != nil {
			return err
		}

		if readerr == io.EOF {
			fmt.Println("transfer complete")
			return nil
		}
		curblknum++
	}

}

func (s *Server) ServeRequests(addr string) error {
	udpaddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpaddr)
	if err != nil {
		return err
	}

	for {
		var b []byte
		n, retaddr, err := conn.ReadFromUDP(b)
		if err != nil {
			return err
		}

		b = b[:n]
		packet, err := pkt.PacketDeserialize(b)
		if err != nil {
			continue
		}

		s.HandleClientReq(retaddr, packet)
	}
}

func (s *Server) HandleErrPacket(pkti pkt.Packet, conn *net.UDPConn) error {
	errpkt := pkti.(*pkt.ErrPacket)
	m := errpkt.ErrMsg

	switch errpkt.ErrCode {
	case pkt.ErrNotDefinedTFTP:
		return errors.New("error packet recieved with code 0 Not defined: " + m)
	case pkt.ErrFileNotFoundTFTP:
		return errors.New("error packet recieved with code 1 File not found: " + m)
	case pkt.ErrAccessViolationTFTP:
		return errors.New("error packet recieved with code 2 Access violation: " + m)
	case pkt.ErrDiskFullAllocExceededTFTP:
		return errors.New("error packet recieved with code 3 Disk full or allocation exceeded: " + m)
	case pkt.ErrIllegalOperationTFTP:
		return errors.New("error packet recieved with code 4 Illegal tftp operation: " + m)
	case pkt.ErrUnknownTransferIdTFTP:
		return errors.New("error packet recieved with code 5 Unknown transfer ID: " + m)
	case pkt.ErrFileAlreadyExistsTFTP:
		return errors.New("error packet recieved with code 6 File already exists: " + m)
	case pkt.ErrNoSuchUserTFTP:
		return errors.New("error packet recieved with code 7 No such user: " + m)
	default:
		return pkt.ErrPacketStructure
	}
}

// TODO: Rewrite to handle netascii and mail modes
// TODO: Potentially use goroutines for concurrency
func (s *Server) SendData(data *pkt.DataPacket, conn *net.UDPConn) error {
	_, err := conn.Write(data.Serialize())
	if err != nil {
		return err
	}

	timeout := time.After(Timeout)
	retry := time.After(Retry)
	for {
		select {
		case <-timeout:
			return ErrTimeout
		case <-retry:
			_, err := conn.Write(data.Serialize())
			if err != nil {
				return err
			}
			err = s.RecieveAck(data, conn)
			if err != ErrBlockNum {
				return err
			} else {
				continue
			}
		}

	}

}

func (s *Server) RecieveAck(data *pkt.DataPacket, conn *net.UDPConn) error {
	var b []byte
	n, _, err := conn.ReadFromUDP(b)
	if err != nil {
		return err
	}

	pkti, err := pkt.PacketDeserialize(b[:n])
	if err != nil {
		return err
	}

	switch pkti.TypeCode() {
	case pkt.ERROR:
		return s.HandleErrPacket(pkti, conn)
	case pkt.ACK:
	default:
		return pkt.ErrPacketType
	}

	ack := pkti.(*pkt.AckPacket)

	if ack.BlockNumber != data.BlockNumber {
		return ErrBlockNum
	}
	return nil
}
