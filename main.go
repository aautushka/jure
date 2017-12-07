package main

import (
    "net"
    "log"
    "bufio"
)

type database struct {
}

type command struct {
    command string
    c chan string
}

func (d* database) get(key string) string {
    return key
}

func main() {
    endpoint := "localhost:7777"
    listener, err := net.Listen("tcp", endpoint)
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    db_control := make(chan command)
    go access_database(db_control)

    log.Print("Listening incoming connections on ", endpoint)
    accept_connections(listener, db_control)
}

func access_database(c chan command) {
    database := &database{}
    for {
        cmd := <-c
        cmd.c <- database.get(cmd.command)
    }
}

func accept_connections(listener net.Listener, c chan command) {
    for {
        conn, err := listener.Accept()
        if err != nil {
            log.Fatal(err)
        }

        log.Print("Connection accepted from ", conn.RemoteAddr())
        go handle_connection(conn, c)
    }
}

func handle_connection(conn net.Conn, c chan command) {
    defer log.Print("Connection closed ", conn.RemoteAddr())
    defer conn.Close()

    feedback := make(chan string)
    reader := bufio.NewReader(conn)

    for {
        cmd, err := reader.ReadString('\n')
        if err != nil {
            break
        }
        c <- command{cmd, feedback}
        response := <-feedback
        conn.Write([]byte(response))
    }
}

