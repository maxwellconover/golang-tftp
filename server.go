// package server implements a RFC 1350 tftp server
package server

type Server struct {
	dir string
	readfunc
}

func NewServer()
