package main

import (
	"log"
	"net"
	"time"
)

const THREADS = 10
const TIMEOUT = 5 * time.Second

func handleConnection(connChan chan net.Conn) {
	for {
		select {
		case conn := <-connChan:
			buffer := make([]byte, 1024)

			select {
			default:
				err := conn.SetDeadline(time.Now().Add(TIMEOUT))
				if err != nil {
					log.Fatal(err)
				}

				_, err = conn.Read(buffer)
				if err != nil {
					log.Println("Read error:", err)
					conn.Write([]byte("HTTP/1.1 500 Internal Server Error\r\n\r\nFailed to read request\r\n"))
					conn.Close()
					continue
				}

				workDone := make(chan bool)
				go func() {
					// Theoretical work here on request
					// time.Sleep(10*time.Second)
					workDone <- true
				}()

				select {
				case <-workDone:
					conn.Write([]byte("HTTP/1.1 200 OK\r\n\r\nHello World!\r\n"))
					conn.Close()
				case <-time.After(TIMEOUT):
					log.Printf("Timeout: workload took longer than %s\n", TIMEOUT)
					conn.Close()
				}
			}
		}
	}
}

func main() {
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	ch := make(chan net.Conn)
	defer close(ch)

	for range THREADS {
		go handleConnection(ch)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		ch <- conn
	}
}
