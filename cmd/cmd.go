package cmd

import (
	"bufio"
	"bytes"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"os/exec"
	"runtime"
	"sync"
)

func Run(opChan chan string, cmd string) {
	var c *exec.Cmd
	switch runtime.GOOS {
	case "linux":
		c = exec.Command("bash", "-c", cmd)
		break
	case "windows":
		c = exec.Command("cmd", "/c", cmd)
		break
	default:
		log.Fatal("not support this platform!")
	}
	stdout, err := c.StdoutPipe()
	if err != nil {
		log.Print(err)
		return
	}
	defer func(stdout io.ReadCloser) {
		err := stdout.Close()
		if err != nil {
			log.Print(err)
		}
	}(stdout)
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		reader := bufio.NewReader(stdout)
		b := make([]byte, 28)
		for {
			n, err := reader.Read(b)
			if err != nil || err == io.EOF {
				break
			}
			opChan <- string(b[:n])
		}
		close(opChan)
	}()
	err = c.Start()
	wg.Wait()
	if err != nil {
		log.Print(err)
		return
	}
}

func Stream(opChan <-chan string, c *gin.Context) func(w io.Writer) bool {
	return func(w io.Writer) bool {
		output, ok := <-opChan
		if !ok {
			return false
		}
		outputBytes := bytes.NewBufferString(output)
		_, err := c.Writer.Write(outputBytes.Bytes())
		if err != nil {
			return false
		}
		return true
	}
}
