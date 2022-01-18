package main

import (
	"flag"
	"io"
	"os"

	"github.com/maxwellconover/golang-tftp/server"
)

func reader(path string) (r io.Reader, err error) {
	r, err = os.Open(path)
	return
}

func writer(path string) (w io.Writer, err error) {
	w, err = os.Create(path)
	return
}

func main() {
	const tftpport = "69"

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dir := flag.String("dir", wd, "provide a directory tftp server will serve files from")
	address := flag.String("address", "localhost", "provide an address tftp server will listen on")
	flag.Parse()

	s := server.NewServer(*dir, reader, writer)
	panic(s.ServeRequests(*address + ":" + tftpport))
}
