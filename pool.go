package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"net/http"
	"time"
)

var pool *Pool

func init() {
	pool = NewPool()
}

type Pool struct {
	redisPool chan redis.Conn
}

func NewPool() *Pool {
	p := &Pool{
		redisPool: make(chan redis.Conn, 20),
	}
	for i := 0; i < 20; i++ {
		c, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			log.Println(err)
			i--
			continue
		}
		p.redisPool <- c
	}
	return p
}

func (p *Pool) Get() redis.Conn {
	return <-p.redisPool //如果没有新的连接，则等待
}

func (p *Pool) Close(c redis.Conn) {
	p.redisPool <- c
}

func (p *Pool) ClosePool() {
	for c := range p.redisPool {
		c.Close()
	}
}

func redisServer(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	c := pool.Get()
	if ok, err := redis.Bool(c.Do("LPUSH", "bughunter", "bughunter")); ok {
	} else {
		log.Println(err)
	}
	msg := fmt.Sprintf("用时:%s", time.Now().Sub(startTime))
	io.WriteString(w, msg+"\n\n")
	pool.Close(c)
}

func main() {
	http.HandleFunc("/", redisServer)
	http.ListenAndServe(":9527", nil)
}
