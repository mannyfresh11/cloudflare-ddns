package main

import (
	"context"
	"fmt"
	"log"
	"os"
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
	zoneID := GetZoneID(api, DOMAIN)
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
