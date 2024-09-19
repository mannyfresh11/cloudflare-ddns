package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/mannyfresh11/cloudflare-ddns/utils/api"
	"github.com/mannyfresh11/cloudflare-ddns/utils/network"
)

var (
	DOMAIN    = os.Getenv("DOMAIN")
	API_TOKEN = os.Getenv("CF_API_TOKEN")
	INTERVAL  = os.Getenv("INTERVAL")
)

func auth(ctx context.Context) (api.CloudflareAPI, error) {
	token := api.CFAuth{
		Token: API_TOKEN,
	}

	a, err := api.CFAuth.New(token, ctx)
	if err != nil {
		return api.CloudflareAPI{}, err
	}

	return a, nil
}

func main() {

	ctx := context.Background()

	cf, err := auth(ctx)
	if err != nil {
		fmt.Println("Could not authenticate to CF. Check API token.")
	}

	interval, err := strconv.Atoi(INTERVAL)
	if err != nil {
		interval = 60
	}

	ticker := time.NewTicker(time.Minute * time.Duration(interval))

	Run := func() {

		externalIP, err := network.GetPublicIP()
		if err != nil {
			fmt.Println("Could not get a public IP.")
		}

		cfIP, recordID, err := cf.GetDNSRecordIP(ctx, DOMAIN)
		if err != nil {
			fmt.Println("Could not get a Cloudflare IP from DNS record.")
		}
		zoneID := cf.GetZoneID(DOMAIN)
		token := cf.VerifyToken(ctx)

		if token == "active" {
			if cfIP != externalIP {
				fmt.Printf("IP does not match DNS record. Cloudflare IP is %s, expected %s\n", cfIP, externalIP)

				fmt.Println("Now updating DNS record...")
				err = cf.UpdateDNSRecord(ctx, zoneID, recordID, externalIP)
				if err != nil {
					fmt.Println("Could not update DNS record.")
				}

				fmt.Println("Record updated!")
			} else {
				fmt.Printf("IP matches DNS record. Cloudflare IP: %s - Expected IP: %s\n", cfIP, externalIP)
			}
		} else {
			fmt.Println("Token not active.")
		}
	}

	Run()

	for {
		select {
		case <-ticker.C:
			Run()
		}
	}
}
