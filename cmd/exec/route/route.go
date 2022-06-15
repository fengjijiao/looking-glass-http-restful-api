package route

import (
	"fmt"
	"log"
	"looking-glass-http-restful-api/cmd"
	"runtime"
)

func Run(opChan chan string, protocolVersion string) {
	var c string
	switch runtime.GOOS {
	case "linux":
		c = fmt.Sprintf("ip -%s route list", protocolVersion)
		break
	case "windows":
		if protocolVersion == "4" {
			c = "netstat -rn -p IP"
		} else {
			c = "netstat -rn -p IPv6"
		}
		break
	default:
		log.Fatal("not support this platform!")
	}
	cmd.Run(opChan, c)
}
