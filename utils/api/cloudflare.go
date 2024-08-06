package api

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudflare/cloudflare-go"
)

type Cloudflare struct {
	api *cloudflare.API
}

func (c Cloudflare) GetZoneID(domain string) string {

	zoneID, err := c.api.ZoneIDByName(domain)
	if err != nil {
		fmt.Printf("Error getting zone id: %v\n", err)
	}

	return zoneID
}

func (c Cloudflare) VerifyToken(ctx context.Context) string {

	token, err := c.api.VerifyAPIToken(ctx)
	if err != nil {
		fmt.Printf("Error verifying token: %v\n", err)
	}

	return token.Status
}

func (c Cloudflare) GetDNSRecordIP(ctx context.Context, domain string) (string, string) {

	zoneID := GetZoneID(c.api, domain)

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

func (c Cloudflare) UpdateDNSRecord(ctx context.Context, zoneID, externalIP, recordID string) {

	params := cloudflare.UpdateDNSRecordParams{
		ID:      recordID,
		Content: externalIP,
	}

	if _, err := c.api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), params); err != nil {

		log.Fatalf("Error updating DNS Record: %v", err)
	}
}
