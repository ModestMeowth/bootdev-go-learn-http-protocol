package main

import (
    "errors"
    "fmt"
    "log"
    "io"
    "net"
    "strings"
)

const inputFilePath = "messages.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {
    c := make(chan string)

    go func() {
        defer f.Close()
        defer close(c)

        current := ""

        for {
            buffer := make([]byte, 8, 8)

            n,err := f.Read(buffer)
            if err != nil {
                if current != "" {
                    c <- current
                }
                if errors.Is(err, io.EOF) {
                    break
                }

                fmt.Printf("error: %s\n", err.Error())
                return
            }

            str := string(buffer[:n])
            parts := strings.Split(str, "\n")

            for i := 0; i < len(parts)-1; i++ {
                c <- fmt.Sprintf("%s%s", current, parts[i])
                current = ""
            }

            current += parts[len(parts)-1]
        }
    }()

    return c
}

func main() {
    l,err := net.Listen("tcp", ":42069")
    if err != nil {
        log.Fatalf("error: %s\n", err)
    }
    defer l.Close()

    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatalf("error: %s\n", err)
        }
        fmt.Println("connection accepted")

        c := getLinesChannel(conn)

        for line := range c {
            fmt.Printf("read: %s\n", line)
        }
    }
}
