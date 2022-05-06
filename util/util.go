package util

import (
	"io"
	"log"
	"net"
)

func ReadAll(conn net.Conn) ([]byte, error) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err == io.EOF {
				log.Println("IO.EOF")
				return nil, err
			}
			log.Println("Read", err)
			return nil, err
		}
		return buffer[:n], nil
	}
}

//copy
func Copy(src, dst io.ReadWriter) {
	buff := make([]byte, 0xffff)
	for {
		n, err := src.Read(buff)
		if err != nil {
			return
		}
		b := buff[:n]
		//write out result
		n, err = dst.Write(b)
		if err != nil {
			return
		}
	}
}
