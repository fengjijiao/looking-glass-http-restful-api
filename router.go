package main

import (
	"embed"
	"github.com/gin-gonic/gin"
	"log"
	"looking-glass-http-restful-api/cmd"
	"looking-glass-http-restful-api/cmd/exec/nslookup"
	"looking-glass-http-restful-api/cmd/exec/ping"
	"looking-glass-http-restful-api/cmd/exec/route"
	"looking-glass-http-restful-api/cmd/exec/traceroute"
	"net/http"
	"path"
	"strings"
)

//go:embed static/*
var staticFS embed.FS

func Router() *gin.Engine {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		r, err := staticFS.ReadFile("static/index.html")
		if err != nil {
			log.Println(err)
			return
		}
		_, err = c.Writer.Write(r)
		if err != nil {
			log.Println(err)
			return
		}
	})
	//r.StaticFS("/static", http.FS(staticFS))
	r.GET("/:folder/*filepath", func(c *gin.Context) {
		requestUri := strings.TrimPrefix(path.Join("static", c.Request.URL.Path), "/")
		r, err := staticFS.ReadFile(requestUri)
		if err != nil {
			log.Println(err)
			return
		}
		c.Writer.Header().Set("Content-Type", getContentType(requestUri))
		_, err = c.Writer.Write(r)
		if err != nil {
			log.Println(err)
			return
		}
	})
	r.GET("/hello", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hello, world!"})
	})
	apiRouteGroup := r.Group("/api")
	apiRouteGroup.GET("/ping/:host", func(c *gin.Context) {
		host := c.Param("host")
		ifName := c.GetHeader("interface")
		if ifName == "" {
			ifName = ""
		}
		protocolVersion := c.GetHeader("protocol")
		if protocolVersion == "" {
			protocolVersion = "4"
		}
		opChan := make(chan string)
		go ping.Run(opChan, host, ifName, protocolVersion)
		c.Stream(cmd.Stream(opChan, c))
	})
	apiRouteGroup.GET("/traceroute/:host", func(c *gin.Context) {
		host := c.Param("host")
		ifName := c.GetHeader("interface")
		if ifName == "" {
			ifName = ""
		}
		protocolVersion := c.GetHeader("protocol")
		if protocolVersion == "" {
			protocolVersion = "4"
		}
		opChan := make(chan string)
		go traceroute.Run(opChan, host, ifName, protocolVersion)
		c.Stream(cmd.Stream(opChan, c))
	})
	apiRouteGroup.GET("/nslookup/:host", func(c *gin.Context) {
		host := c.Param("host")
		protocolVersion := c.GetHeader("protocol")
		if protocolVersion == "" {
			protocolVersion = "4"
		}
		opChan := make(chan string)
		go nslookup.Run(opChan, host)
		c.Stream(cmd.Stream(opChan, c))
	})
	apiRouteGroup.GET("/route", func(c *gin.Context) {
		protocolVersion := c.GetHeader("protocol")
		if protocolVersion == "" {
			protocolVersion = "4"
		}
		opChan := make(chan string)
		go route.Run(opChan, protocolVersion)
		c.Stream(cmd.Stream(opChan, c))
	})
	return r
}

var contentTypes = make(map[string]string)

func init() {
	contentTypes["html"] = "text/html"
	contentTypes["css"] = "text/css"
	contentTypes["js"] = "text/javascript"
	contentTypes["eot"] = "application/vnd.ms-fontobject"
	contentTypes["svg"] = "image/svg+xml"
	contentTypes["ttf"] = "font/ttf"
	contentTypes["woff"] = "font/woff"
}

func getContentType(uri string) string {
	for k, v := range contentTypes {
		if strings.HasSuffix(uri, k) {
			return v
		}
	}
	return ""
}
