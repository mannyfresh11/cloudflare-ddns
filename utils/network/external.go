package network

import (
	"bytes"
	"errors"
	"log"
	"os/exec"
)

func getPublicAddr() string {

	var buffer bytes.Buffer

	defer buffer.Reset()

	command := `curl -sf4 https://one.one.one.one/cdn-cgi/trace | grep 'ip' | tr -d 'ip='`

	cmd := exec.Command("bash", "-c", command)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}

	cmd.Stdout = &buffer

	if err := cmd.Run(); err != nil {
		log.Fatalf("Error runing cmd: %v\n", err)
	}

	return buffer.String()
}
