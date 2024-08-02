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
	cfIP, recordID := getDNSRecordIP(ctx, api)

	if token == "active" {
		if cfIP != externalIP {
			fmt.Printf("IP does not match DNS record. Cloudflare IP is %s, expected %s\n", cfIP, externalIP)
			fmt.Println("Now updating DNS record...")
			fmt.Printf("This is the zoneID: %s\n This is the RecordID: %s\n", zoneID, recordID)

			// updateDNSRecord(ctx, api, zoneID, recordID, externalIP)

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

func getZoneID(api *cloudflare.API, domain string) string {

	zoneID, err := api.ZoneIDByName(domain)
	if err != nil {
		fmt.Printf("Error getting zone id: %v\n", err)
	}

	return zoneID
}

func getDNSRecordIP(ctx context.Context, api *cloudflare.API) (string, string) {

	zoneID := getZoneID(api, DOMAIN)

	param := cloudflare.ListDNSRecordsParams{
		Name: DOMAIN,
	}

	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), param)
	if err != nil {
		log.Fatalf("Error getting DNS record IP: %v", err)
	}

	var rec string
	var recID string

	for _, r := range records {
		if r.Name == DOMAIN {
			rec = r.Content
			recID = r.ID
		}
	}

	return rec, recID
}

func verifyToken(ctx context.Context, api *cloudflare.API) string {

	token, err := api.VerifyAPIToken(ctx)
	if err != nil {
		fmt.Printf("Error verifying token: %v\n", err)
	}

	return token.Status
}

func updateDNSRecord(ctx context.Context, api *cloudflare.API, zoneID, externalIP, recordID string) {

	params := cloudflare.UpdateDNSRecordParams{
		ID:      recordID,
		Content: externalIP,
	}

	if _, err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), params); err != nil {

		log.Fatalf("Error updating DNS Record: %v", err)
	}
}
