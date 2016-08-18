package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"
)

// 重写生成连接池方法
func newPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:   80,
		MaxActive: 12000, // max number of connections
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}

// 生成连接池
var pool = newPool()

func redisServer(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	// 从连接池里面获得一个连接
	c := pool.Get()
	// 连接完关闭，其实没有关闭，是放回池里，也就是队列里面，等待下一个重用
	defer c.Close()
	dbkey := "netgame:info"
	if ok, err := redis.Bool(c.Do("LPUSH", dbkey, "yangzetao")); ok {
	} else {
		log.Print(err)
	}
	msg := fmt.Sprintf("用时：%s", time.Now().Sub(startTime))
	io.WriteString(w, msg+"\n\n")
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	http.HandleFunc("/", redisServer)
	http.ListenAndServe(":9527", nil)
}
