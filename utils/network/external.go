package network

import (
	"io"
	"net/http"
	"strings"
)

func GetPublicIP() (string, error) {

	url := "https://one.one.one.one/cdn-cgi/trace"

	res, err := http.Get(url)
	if err != nil {
		return "", err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	bodyStripped := strings.Split(string(body), "\n")

	for _, v := range bodyStripped {
		if strings.HasPrefix(v, "ip=") {

			ipaddr := strings.TrimPrefix(v, "ip=")

			return ipaddr, nil
		}
	}

	return "", nil
}
