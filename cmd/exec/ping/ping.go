package ping

import (
	"fmt"
	"log"
	"looking-glass-http-restful-api/cmd"
	"runtime"
)

func Run(opChan chan string, host string, ifName string, protocolVersion string) {
	var c string
	switch runtime.GOOS {
	case "linux":
		if ifName == "" {
			c = fmt.Sprintf("ping -%s -c 9 %s", protocolVersion, host)
		} else {
			c = fmt.Sprintf("ping -%s -I %s -c 9 %s", protocolVersion, host, ifName)
		}
		break
	case "windows":
		c = fmt.Sprintf("ping -%s %s", protocolVersion, host)
		break
	default:
		log.Fatal("not support this platform!")
	}

	cmd.Run(opChan, c)
}
