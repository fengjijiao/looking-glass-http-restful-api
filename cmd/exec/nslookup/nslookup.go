package nslookup

import (
	"fmt"
	"log"
	"looking-glass-http-restful-api/cmd"
	"runtime"
)

func Run(opChan chan string, host string) {
	var c string
	switch runtime.GOOS {
	case "linux":
		c = fmt.Sprintf("nslookup %s", host)
		break
	case "windows":
		c = fmt.Sprintf("nslookup %s", host)
		break
	default:
		log.Fatal("not support this platform!")
	}

	cmd.Run(opChan, c)
}
