package main

import (
	"fmt"
)

func main() {
	redisPool := make(chan int, 12)
	redisPool <- 4
	fmt.Println(len(redisPool))
	redisPool <- 3
	fmt.Println(len(redisPool))
	<-redisPool
	<-redisPool
	fmt.Println(len(redisPool))
}
