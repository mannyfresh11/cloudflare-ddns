package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/cloudflare/cloudflare-go"
)

var DOMAIN = os.Getenv("DOMAIN")
var API_TOKEN = os.Getenv("CF_API_TOKEN")

func main() {

	ctx := context.Background()

	api, err := cloudflare.NewWithAPIToken(API_TOKEN)
	if err != nil {
		log.Fatalf("Failed to get api: %v", err)
	}

	externalIP := strings.Trim(getPublicAddr(), "\n")
	zoneID := getZoneID(api, DOMAIN)
	token := verifyToken(ctx, api)
	cfIP, recordID := getDNSRecordIP(ctx, api, DOMAIN)

	if token == "active" {
		if cfIP != externalIP {
			fmt.Printf("IP does not match DNS record. Cloudflare IP is %s, expected %s\n", cfIP, externalIP)
			fmt.Println("Now updating DNS record...")

			updateDNSRecord(ctx, api, zoneID, recordID, externalIP)

			fmt.Println("Updated record")
		} else {
			fmt.Printf("IP matches DNS record. Cloudflare IP is %s, expected %s\n", cfIP, externalIP)
		}
	} else {
		fmt.Println("Token not active.")
	}
}

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
