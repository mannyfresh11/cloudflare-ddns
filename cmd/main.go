package main

import (
	"context"
	"fmt"
	"os"

	"github.com/mannyfresh11/cloudflare-ddns/utils/api"
	"github.com/mannyfresh11/cloudflare-ddns/utils/network"
)

var DOMAIN = os.Getenv("DOMAIN")
var API_TOKEN = os.Getenv("CF_API_TOKEN")

func InitStart(ctx context.Context) api.CloudflareAPI {
	auth := api.CFAuth{
		Token: API_TOKEN,
	}

	a := api.Auth.New(auth, ctx)

	return a
}

func main() {

	ctx := context.Background()

	a := InitStart(ctx)

	externalIP := network.GetPublicIP()

	cfIP, recordID := a.GetDNSRecordIP(ctx, DOMAIN)
	zoneID := a.GetZoneID(DOMAIN)
	token := a.VerifyToken(ctx)

	if token == "active" {
		if cfIP != externalIP {
			fmt.Printf("IP does not match DNS record. Cloudflare IP is %s, expected %s\n", cfIP, externalIP)
			fmt.Println("Now updating DNS record...")
			fmt.Printf("This is the record ID: %s and this is the zoneID: %s", recordID, zoneID)

			a.UpdateDNSRecord(ctx, zoneID, recordID, externalIP)

			fmt.Println("Updated record")
		} else {
			fmt.Printf("IP matches DNS record. Cloudflare IP %s, expected %s\n", cfIP, externalIP)
		}
	} else {
		fmt.Println("Token not active.")
	}
}
