package main

import (
    "errors"
    "fmt"
    "log"
    "io"
    "os"
    "strings"
)

const inputFilePath = "messages.txt"

func main() {
    f,err := os.Open(inputFilePath)
    if err != nil {
        log.Fatalf("Could not open %s: %s\n", inputFilePath, err)
    }

    defer f.Close()

    fmt.Printf("Reading data from %s\n", inputFilePath)
    fmt.Println("=====================================")

    currentLine := ""

    for {
        buf := make([]byte, 8, 8)

        chunk,err := f.Read(buf)
        if err != nil {

            if errors.Is(err, io.EOF) {
                break
            }

            fmt.Printf("error: %s\n", err.Error())
            break
        }

        parts := strings.Split(string(buf[:chunk]), "\n")

        currentLine += parts[0]

        if len(parts) == 1 {
            continue
        }

        fmt.Printf("read: %s\n", currentLine)
        for i := 1; i < len(parts)-1; i++ {
            fmt.Printf("read: %s\n", parts[i])
        }

        currentLine = parts[len(parts)-1]
    }
}
