package network

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

// func GetPublicIP() string {
//
// 	var buffer bytes.Buffer
//
// 	defer buffer.Reset()
//
// 	command := `curl -sf4 https://one.one.one.one/cdn-cgi/trace | grep 'ip' | tr -d 'ip='`
//
// 	cmd := exec.Command("bash", "-c", command)
// 	if errors.Is(cmd.Err, exec.ErrDot) {
// 		cmd.Err = nil
// 	}
//
// 	cmd.Stdout = &buffer
//
// 	if err := cmd.Run(); err != nil {
// 		log.Fatalf("Error runing cmd: %v\n", err)
// 	}
//
// 	ip := strings.Trim(buffer.String(), "\n")
//
// 	return ip
// }

func GetPublicIP() string {

	url := "https://one.one.one.one/cdn-cgi/trace"

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("Cannot reach url: %s", url)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("Cannot read response")
	}

	ipS := string(body)

	ip := strings.Split(ipS, "\n")

	for _, v := range ip {
		if strings.HasPrefix(v, "ip=") {
			return strings.TrimPrefix(v, "ip=")
		}
	}

	return ""
}
