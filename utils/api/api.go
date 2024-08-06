package api

import (
	"context"
)

type CFSetter interface {
	GetZoneID(domain string) string

	VerifyToken(ctx context.Context) string

	GetDNSRecords(ctx context.Context, domain string) (string, string)

	UpdateDNSRecord(ctx context.Context, domain, zoneID, externalIP string)
}

type Auth interface {
	New(ctx context.Context) CFSetter
}
