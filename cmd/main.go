package main

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/mannyfres11/cfupdater/utils/api"
	"github.com/mannyfres11/cfupdater/utils/network"
)

var DOMAIN = os.Getenv("DOMAIN")
var API_TOKEN = os.Getenv("CF_API_TOKEN")

func main() {

	ctx := context.Background()

	externalIP := strings.Trim(network.Netter.GetPublicIP(), "\n")

	cfIP, recordID := api.CFSetter.GetDNSRecords(ctx, DOMAIN)

	zoneID := api.CFSetter.GetZoneID(DOMAIN)
	token := api.CFSetter.VerifyToken(ctx)

	if token == "active" {
		if cfIP != externalIP {
			fmt.Printf("IP does not match DNS record. Cloudflare IP is %s, expected %s\n", cfIP, externalIP)
			fmt.Println("Now updating DNS record...")

			api.CFSetter.UpdateDNSRecord(ctx, zoneID, recordID, externalIP)

			fmt.Println("Updated record")
		} else {
			fmt.Printf("IP matches DNS record. Cloudflare IP f%s, expected %s\n", cfIP, externalIP)
		}
	} else {
		fmt.Println("Token not active.")
	}
}
