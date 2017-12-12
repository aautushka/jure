package main

import (
    "net"
    "log"
    "bufio"
    "strings"
    "strconv"
    "hash/fnv"
    "runtime"
)

type object struct {
    str string
    arr []*object
    dict *map[string]*object
}

func (o *object) set(value string) {
    o.str = value
}

func (o *object) get() string {
    return o.str
}

type shard struct {
    storage map[string]*object
}

type cache struct {
    shards []*shard
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

type exists_command struct {
    *command_data
}

type append_command struct {
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

func (c *exists_command) apply(cc *cache) {
    if cc.exists(c.args[0]) {
        c.c <- "$1\r\n"
    } else {
        c.c <- "$0\r\n"
    }
}

func (c *append_command) apply(cc *cache) {
    new_length := cc.append(c.args[0], c.args[1])
    c.c <- ":" + strconv.Itoa(new_length) + "\r\n"
}

func (c *shard) get(key string) string {
    if val, ok := c.storage[key]; ok {
        return val.get()
    }
    return ""
}

func (c *shard) set(key string, value string) {
    if val, ok := c.storage[key]; ok {
        val.set(value)
    } else {
        c.storage[key] = &object{ str : value };
    }
}

func (c *shard) del(key string) bool {
    delete(c.storage, key)
    return true
}

func (c *shard) exists(key string) bool {
    if _, ok := c.storage[key]; ok {
        return true
    } else {
        return false
    }
}

func (c *shard) append(key string, value string) int {
    prev, _ := c.storage[key]
    newval := prev.get() + value
    c.storage[key].set(newval)
    return len(newval)
}

func (c *cache) set(key string, value string) {
    c.get_shard(key).set(key, value)
}

func (c* cache) get(key string) string {
    return c.get_shard(key).get(key)
}

func (c* cache) del(key string) bool {
    return c.get_shard(key).del(key)
}

func (c* cache) exists(key string) bool {
    return c.get_shard(key).exists(key)
}

func (c* cache) append(key string, value string) int {
    return c.get_shard(key).append(key, value)
}

func hash(s string) uint32 {
    h := fnv.New32a()
    h.Write([]byte(s))
    return h.Sum32()
}

func (c* cache) get_shard(key string) *shard {
    num := hash(key) % uint32(len(c.shards))
    return c.shards[num]
}

func make_shard() *shard {
    return &shard{storage: make(map[string]*object)}
}

func make_cache_n(num_shards int) *cache {
    shards := make([]*shard, 0, num_shards)
    for i := 0; i < 10; i++ {
        s := make_shard()
        shards = append(shards, s)
    }
    return &cache{shards: shards}
}

func make_cache() *cache {
    num_cpu := runtime.NumCPU()
    return make_cache_n(num_cpu + 1)
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
    cache := make_cache()
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
        case "exists":
            dat.args = tokens[1:]
            cmd := &exists_command{}
            cmd.command_data = dat
            return cmd
        case "append":
            dat.args = tokens[1:]
            cmd := &append_command{}
            cmd.command_data = dat
            return cmd
        }
    }

    return nil
}


