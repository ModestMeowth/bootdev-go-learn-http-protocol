package main

import (
    "fmt"
    "log"
    "io"
    "os"
    "strings"
)

const inputFilePath = "messages.txt"

func getLinesChannel(f io.ReadCloser) <-chan string {
    c := make(chan string)
    var current string

    go func() {
        for {
            buffer := make([]byte, 8, 8)

            n,err := f.Read(buffer)
            if err != nil {
                if current != "" {
                    c <- current
                    current = ""
                }
                break
            }

            str := string(buffer[:n])
            parts := strings.Split(str, "\n")

            for i := 0; i < len(parts)-1; i++ {
                c <- current + parts[i]
                current = ""
            }

            current += parts[len(parts)-1]
        }
        close(c)
    }()

    return c
}

func main() {
    f,err := os.Open(inputFilePath)
    if err != nil {
        log.Fatalf("Could not open %s: %s\n", inputFilePath, err)
    }
    defer f.Close()


    fmt.Printf("Reading data from %s\n", inputFilePath)
    fmt.Println("=====================================")

    c := getLinesChannel(f)

    for line := range c {
        fmt.Printf("read: %s\n", line)
    }
}
