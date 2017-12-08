package main

import (
    "net"
    "log"
    "bufio"
    "strings"
    "strconv"
)

type cache struct {
    storage map[string]string
}

type command interface {
    apply(c *cache)
}

type command_data struct {
    args []string
    c chan string
}

type get_command struct {
    *command_data
}

type set_command struct {
    *command_data
}

type del_command struct {
    *command_data
}

func (c *get_command) apply(cc *cache) {
    value := cc.get(c.args[0])
    if value != "" {
        c.c <- string("$" + strconv.Itoa(len(value)) + "\r\n" + value + "\r\n")
    } else {
        c.c <- "$-1\r\n"
    }
}

func (c *set_command) apply(cc *cache) {
    cc.set(c.args[0], c.args[1])
    c.c <- "+OK\r\n"
}

func (c *del_command) apply(cc *cache) {
    if cc.del(c.args[0]) {
        c.c <- ":1\r\n"
    } else {
        c.c <- ":0\r\n"
    }
}

func (c *cache) get(key string) string {
    if val, ok := c.storage[key]; ok {
        return val
    }
    return ""
}

func (c * cache) set(key string, value string) {
    log.Println("set ", key)
    c.storage[key] = value

    val, ok := c.storage[key]
    log.Println("set ", val, ok)
}

func (c * cache) del(key string) bool {
    delete(c.storage, key)
    return true
}

func main() {
    endpoint := "localhost:7777"
    listener, err := net.Listen("tcp", endpoint)
    if err != nil {
        log.Fatal(err)
    }
    defer listener.Close()

    db_control := make(chan command)
    go manage_cache(db_control)

    log.Print("Listening incoming connections on ", endpoint)
    accept_connections(listener, db_control)
}

func manage_cache(c chan command) {
    cache := &cache{storage: make(map[string]string)}
    for {
        cmd := <-c
        cmd.apply(cache)
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
        request, err := reader.ReadString('\n')
        request = strings.Trim(request, "\r\n")
        if err != nil {
            break
        }

        dat := &command_data{c: feedback}
        cmd := parse_command(request, dat)
        if cmd != nil {
            c <- cmd
            response := <-feedback
            conn.Write([]byte(response))
        } else {
            log.Print("Unexpected command ", request)
        }

    }
}

func parse_command(line string, dat *command_data) command {
    tokens := strings.Split(line, " ")
    if len(tokens) > 0 {
        switch tokens[0] {
        case "get":
            dat.args = tokens[1:]
            cmd := &get_command{}
            cmd.command_data = dat
            return cmd
        case "set":
            dat.args = tokens[1:]
            cmd := &set_command{}
            cmd.command_data = dat
            return cmd
        case "del":
            dat.args = tokens[1:]
            cmd := &del_command{}
            cmd.command_data = dat
            return cmd
        }
    }

    return nil
}


