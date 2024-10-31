package api

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

type CFAuth struct {
	Token string
}

func NewCFAuth(token string) (CloudflareAPI, error) {

	api, err := cloudflare.NewWithAPIToken(token)
	if err != nil {
		return CloudflareAPI{}, fmt.Errorf("Failed to get api: %v", err)
	}

	return CloudflareAPI{
		api: api,
	}, nil
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

	switch token.Status {
	case "active":
		return token.Status
	case "disabled", "expired":
		return fmt.Sprintf("Your token is either disabled or expired.")
	default:
		return fmt.Sprintf("This is the token status: %s", token.Status)
	}
}

func (c CloudflareAPI) GetDNSRecordIP(ctx context.Context, domain string) (string, string, error) {

	zoneID := c.GetZoneID(domain)

	param := cloudflare.ListDNSRecordsParams{
		Name: domain,
	}

	records, _, err := c.api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), param)
	if err != nil {
		return "", "", fmt.Errorf("Error getting DNS record IP: %v", err)
	}

	for _, r := range records {
		if r.Name == domain {
			return r.Content, r.ID, nil
		}
	}

	return "", "", fmt.Errorf("No matching DNS record found in domain: %s", domain)
}

func (c CloudflareAPI) UpdateDNSRecord(ctx context.Context, zoneID, externalIP, recordID string) error {

	params := cloudflare.UpdateDNSRecordParams{
		ID:      recordID,
		Content: externalIP,
	}

	if _, err := c.api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), params); err != nil {
		return fmt.Errorf("Error updating DNS Record: %v", err)
	}

	SendHook("Clodflare IP has changed.")

	return nil
}
