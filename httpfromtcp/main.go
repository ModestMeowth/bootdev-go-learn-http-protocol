package main

import (
    "fmt"
    "io"
    "os"
)

func main() {
    f,err := os.Open("./messages.txt")
    if err != nil {
        panic(err)
    }

    defer f.Close()

    buf := make([]byte, 8)

    for {
        chunk,err := f.Read(buf)
        if err != nil && err != io.EOF {
            panic(err)
        }

        if err == io.EOF {
            break
        }

        fmt.Printf("read: %s\n", buf[:chunk])
    }
}
