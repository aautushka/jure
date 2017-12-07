package main

import (
    "fmt"
    "net"
    "log"
)

func main() {
    endpoint := "localhost:7777"
    l, err := net.Listen("tcp", endpoint)
    if err != nil {
        log.Fatal(err)
    }
    defer l.Close()

    log.Print("Listening incoming connections on ", endpoint)
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }

        log.Print("Connection accepted from ", conn.RemoteAddr())
        go handle_connection(conn)
    }
}

func handle_connection(conn net.Conn) {
    defer log.Print("Connection closed ", conn.RemoteAddr())
    defer conn.Close()

    buf := make([]byte, 1024)

    for {
        bytesRead, err := conn.Read(buf)
        if err != nil{
            break
        }
        fmt.Println("Read some bytes")
        conn.Write(buf[0: bytesRead])
    }
}
