package gapi

import (
	"context"
	"errors"

	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/alrasyidin/simplebank-go/pb"
	"github.com/alrasyidin/simplebank-go/util"
	"github.com/alrasyidin/simplebank-go/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (server *Server) LoginUser(ctx context.Context, param *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	violations := validateLoginUserRequest(param)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	user, err := server.store.GetUser(ctx, param.Username)
	if err != nil {
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed get user: %w", err)
	}

	err = util.CheckPassword(param.Password, user.HashedPassword)

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "incorrect password")
	}

	accessToken, accessPayload, err := server.tokenGenerator.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create access token")
	}
	refreshToken, refreshPayload, err := server.tokenGenerator.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create refresh token")
	}

	mtdt := extractMetadata(ctx)
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshPayload.Payload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    mtdt.UserAgent,
		ClientIp:     mtdt.ClientIP,
		ExpiresAt:    refreshPayload.ExpiresAt.Time,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create session")
	}

	resp := &pb.LoginUserResponse{
		User:                  convertUser(user),
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessPayload.ExpiresAt.Time),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshPayload.ExpiresAt.Time),
	}

	return resp, nil
}

func validateLoginUserRequest(param *pb.LoginUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(param.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validation.ValidatePassword(param.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}

	return violations
}
