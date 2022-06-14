package main

import (
	"flag"
	"github.com/gin-gonic/gin"
	"log"
)

var (
	listen string
)

func init() {
	flag.StringVar(&listen, "l", ":8080", "service listen host&port")
	flag.Parse()
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	r := Router()
	log.Fatal(r.Run(listen))
}
