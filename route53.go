package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/route53"
	"github.com/aws/aws-sdk-go-v2/service/route53/types"
)

func createRecord(ctx context.Context, name, value, hostedZoneName string, ttl int64) error {
	sugar.Debugw("getting hosted zone", "dns_name", hostedZoneName)
	hostedZones, err := client.ListHostedZonesByName(ctx, &route53.ListHostedZonesByNameInput{
		DNSName: pointer(hostedZoneName),
	})

	if err != nil {
		return fmt.Errorf("error getting hosted zone: %w", err)
	}

	if len(hostedZones.HostedZones) <= 0 {
		return fmt.Errorf("no hosted zones found")
	}

	var hostedZoneId string
	for _, hostedZone := range hostedZones.HostedZones {
		if hostedZone.Name != nil && *hostedZone.Name == fmt.Sprintf("%s.", hostedZoneName) {
			if hostedZone.Id != nil {
				hostedZoneId = *hostedZone.Id
				break
			}
			return fmt.Errorf("invalid hosted zone id")
		}
	}

	if len(hostedZoneId) <= 0 {
		return fmt.Errorf("empty hosted zone id")
	}

	sugar.Debugw(
		"creating record in hosted zone",
		"hosted_zone_id", hostedZoneId,
		"name", name,
		"value", value,
		"type", "A",
	)
	_, err = client.ChangeResourceRecordSets(ctx, &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: pointer(hostedZoneId),
		ChangeBatch: &types.ChangeBatch{
			Changes: []types.Change{
				{
					Action: "UPSERT",
					ResourceRecordSet: &types.ResourceRecordSet{
						Name: pointer(name),
						Type: "A",
						TTL:  pointer(ttl),
						ResourceRecords: []types.ResourceRecord{
							{
								Value: pointer(value),
							},
						},
					},
				},
			},
		},
	})

	if err != nil {
		return fmt.Errorf("error creating record: %w", err)
	}

	return nil
}
