package gapi

import (
	"context"
	"fmt"
	"strings"

	"github.com/alrasyidin/simplebank-go/token"
	"google.golang.org/grpc/metadata"
)

const (
	authorizationHeader = "authorization"
	authorizationType   = "bearer"
)

func (server *Server) authorizeUser(ctx context.Context) (*token.PayloadJWT, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, fmt.Errorf("missing metadata")
	}

	values := md.Get(authorizationHeader)
	if len(values) == 0 {
		return nil, fmt.Errorf("missing authorization header")
	}
	authHeader := values[0]
	fields := strings.Fields(authHeader)
	if len(fields) < 2 {
		return nil, fmt.Errorf("invalid authorization header format")
	}

	if strings.ToLower(fields[0]) != authorizationType {
		return nil, fmt.Errorf("authorization type not supported")
	}

	accessToken := fields[1]
	payload, err := server.tokenGenerator.VerifyToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("invalid access token")
	}

	return payload, nil
}
