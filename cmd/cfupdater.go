package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/mannyfresh11/cloudflare-ddns/utils/api"
	"github.com/mannyfresh11/cloudflare-ddns/utils/logger"
	"github.com/mannyfresh11/cloudflare-ddns/utils/network"
)

var (
	DOMAIN    = os.Getenv("DOMAIN")
	API_TOKEN = os.Getenv("CF_API_TOKEN")
	INTERVAL  = os.Getenv("INTERVAL")
)

func checkEnv() int {
	if DOMAIN == "" {
		log.Fatal("Domain empty. Please provide value.")
	}

	if API_TOKEN == "" {
		log.Fatal("API token empty. Please provide value.")
	}

	interval, err := strconv.Atoi(INTERVAL)
	if err != nil {
		interval = 60
	}

	return interval
}

func main() {
	logger := logger.NewLogger(slog.Level(0))
	interval := checkEnv()
	ctx := context.Background()

	ticker := time.NewTicker(time.Minute * time.Duration(interval))
	defer ticker.Stop()

	cf, err := api.NewCFAuth(API_TOKEN)
	if err != nil {
		logger.Error("Could not authenticate to CF. Check API token.")
		os.Exit(1)
	}

	fmt.Println("starting application...")
	Run := func() {
		zoneID := cf.GetZoneID(DOMAIN)
		token := cf.VerifyToken(ctx)

		externalIP, err := network.GetPublicIP()
		if err != nil {
			logger.Info("Could not get a public IP.")
			fmt.Println("Could not get a public IP.")
		}

		cfIP, recordID, err := cf.GetDNSRecordIP(ctx, DOMAIN)
		if err != nil {
			logger.Info("Could not get a Cloudflare IP from DNS record.")
			fmt.Println("Could not get a Cloudflare IP from DNS record.")
		}

		if token == "active" {
			if cfIP != externalIP {
				fmt.Printf("IP does not match DNS record. Cloudflare IP is %s, expected %s\n", cfIP, externalIP)

				fmt.Println("Now updating DNS record...")
				err = cf.UpdateDNSRecord(ctx, zoneID, recordID, externalIP)
				if err != nil {
					logger.Error("Could not update DNS record.")
					fmt.Println("Could not update DNS record.")
				}

				fmt.Println("Record updated!")
			} else {
				logger.Error("IP matches DNS record. Cloudflare IP: %s - Expected IP: %s\n", cfIP, externalIP)
				fmt.Printf("IP matches DNS record. Cloudflare IP: %s - Expected IP: %s\n", cfIP, externalIP)
			}
		} else {
			logger.Error("Token not active.")
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
