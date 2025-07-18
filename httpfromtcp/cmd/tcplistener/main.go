package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"

    "github.com/ModestMeowth/bootdev-go-learn-http-protocol/internal/request/request"
)

const port = "127.0.0.1:42069"

func main() {
	l, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("error listening for TCP traffic: %s\n", err.Error())
	}
	defer l.Close()

	fmt.Printf("Listening for TCP traffic on %s\n", port)
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatalf("error: %s\n", err.Error())
		}
		fmt.Println("Accepted connection from", conn.RemoteAddr())

		c := getLinesChannel(conn)

		for line := range c {
			fmt.Println(line)
		}
		fmt.Println("Connection to ", conn.RemoteAddr(), "closed")
	}
}

func getLinesChannel(f io.ReadCloser) <-chan string {
	lines := make(chan string)
	go func() {
		defer f.Close()
		defer close(lines)

		current := ""
		for {
			b := make([]byte, 8, 8)
			n, err := f.Read(b)
			if err != nil {
				if current != "" {
					lines <- current
				}
				if errors.Is(err, io.EOF) {
					break
				}
				fmt.Printf("error: %s\n", err.Error())
				return
			}
            parts := strings.Split(string(b[:n]), "\n")
			for i := 0; i < len(parts)-1; i++ {
				lines <- current + parts[i]
				current = ""
			}
			current += parts[len(parts)-1]
		}
	}()
	return lines
}
