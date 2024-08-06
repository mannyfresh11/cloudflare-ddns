package api

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
)

type CFAuth struct {
	Token string
}

func (a CFAuth) New(_ context.Context) CloudflareAPI {
	api, err := cloudflare.NewWithAPIToken(a.Token)
	if err != nil {
		log.Fatalf("Failed to get api: %v", err)
	}

	c := CloudflareAPI{
		api: api,
	}

	return c
}

type CloudflareAPI struct {
	api *cloudflare.API
}

func (c CloudflareAPI) GetZoneID(domain string) string {

	zoneID, err := c.api.ZoneIDByName(domain)
	if err != nil {
		fmt.Printf("Error getting zone id: %v\n", err)
	}

	return zoneID
}

func (c CloudflareAPI) VerifyToken(ctx context.Context) string {

	token, err := c.api.VerifyAPIToken(ctx)
	if err != nil {
		fmt.Printf("Error verifying token: %v\n", err)
	}

	return token.Status
}

func (c CloudflareAPI) GetDNSRecordIP(ctx context.Context, domain string) (string, string) {

	zoneID := c.GetZoneID(domain)

	param := cloudflare.ListDNSRecordsParams{
		Name: domain,
	}

	records, _, err := c.api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), param)
	if err != nil {
		log.Fatalf("Error getting DNS record IP: %v", err)
	}

	var rec string
	var recID string

	for _, r := range records {
		if r.Name == domain {
			rec = r.Content
			recID = r.ID
		}
	}

	return rec, recID
}

func (c CloudflareAPI) UpdateDNSRecord(ctx context.Context, zoneID, externalIP, recordID string) {

	params := cloudflare.UpdateDNSRecordParams{
		ID:      recordID,
		Content: externalIP,
	}

	if _, err := c.api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), params); err != nil {

		log.Fatalf("Error updating DNS Record: %v", err)
	}
}
