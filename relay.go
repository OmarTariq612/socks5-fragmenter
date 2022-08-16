package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

type Relay struct {
	bindAddr   string
	serverAddr string
}

func NewRelay(bindAddr, serverAddr string) *Relay {
	return &Relay{bindAddr: bindAddr, serverAddr: serverAddr}
}

const timeoutDuration = 5 * time.Second

func (r *Relay) Serve() error {
	listener, err := net.Listen("tcp", r.bindAddr)
	if err != nil {
		return err
	}
	defer listener.Close()
	log.Printf("Serving on %v\n", r.bindAddr)
	log.Printf("the provided server address is %v\n", r.serverAddr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return fmt.Errorf("could not accept connections")
		}
		go func() {
			defer conn.Close()
			serverConn, err := net.DialTimeout("tcp", r.serverAddr, timeoutDuration)
			if err != nil {
				log.Printf("could not establish server connection, %v", err)
				return
			}
			defer serverConn.Close()
			err = handleConnection(conn, serverConn)
			if err != nil {
				log.Printf("connection failed, %v", err)
			}
		}()
	}
}

func handleConnection(clientConn, serverConn io.ReadWriter) error {
	var buf [2]byte
	_, err := io.ReadFull(clientConn, buf[:])
	if err != nil {
		return fmt.Errorf("could not read handshake header (socks_version + n_methods) from the client")
	}

	_, err = serverConn.Write(buf[:])
	if err != nil {
		return fmt.Errorf("could not write handshake header (socks_version + n_methods) to the server")
	}

	methods := make([]byte, buf[1])
	_, err = io.ReadFull(clientConn, methods)
	if err != nil {
		return fmt.Errorf("could not read methods from the client")
	}

	_, err = serverConn.Write(methods)
	if err != nil {
		return fmt.Errorf("could not write methods to the server")
	}

	errc := make(chan error, 2)
	go func() {
		_, err := io.Copy(clientConn, serverConn)
		if err != nil {
			err = fmt.Errorf("could not copy from server to client, %v", err)
		}
		errc <- err
	}()

	go func() {
		_, err := io.Copy(serverConn, clientConn)
		if err != nil {
			err = fmt.Errorf("could not copy from client to server, %v", err)
		}
		errc <- err
	}()

	return <-errc
}
