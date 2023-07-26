package gapi

import (
	"context"
	"log"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcUserAgentGatewayHeader = "grpcgateway-user-agent"
	grpcXForwardedForHeader    = "x-forwarded-for"
	userAgentHeader            = "user-agent"
)

type Metadata struct {
	ClientIP  string
	UserAgent string
}

func extractMetadata(ctx context.Context) *Metadata {
	md := new(Metadata)
	if metadata, ok := metadata.FromIncomingContext(ctx); ok {
		log.Printf("md: %+v", metadata)
		if userAgents := metadata.Get(grpcUserAgentGatewayHeader); len(userAgents) > 0 {
			md.UserAgent = userAgents[0]
		}
		if userAgents := metadata.Get(userAgentHeader); len(userAgents) > 0 {
			md.UserAgent = userAgents[0]
		}
		if clientIPs := metadata.Get(grpcXForwardedForHeader); len(clientIPs) > 0 {
			md.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		md.ClientIP = p.Addr.String()
	}
	return md
}
