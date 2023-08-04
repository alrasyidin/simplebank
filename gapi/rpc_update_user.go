package gapi

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/alrasyidin/simplebank-go/pb"
	"github.com/alrasyidin/simplebank-go/util"
	"github.com/alrasyidin/simplebank-go/validation"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) UpdateUser(ctx context.Context, param *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	payload, err := server.authorizeUser(ctx)
	if err != nil {
		return nil, unauthenticatedError(err)
	}

	violations := validateUpdateUserRequest(param)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	if payload.Username != param.Username {
		return nil, status.Errorf(codes.PermissionDenied, "cannot update other user's info")
	}

	data := db.UpdateUserUsingCaseSecondParams{
		Username: param.GetUsername(),
		Email: pgtype.Text{
			String: param.GetEmail(),
			Valid:  param.Email != nil,
		},
		FullName: pgtype.Text{
			String: param.GetFullName(),
			Valid:  param.FullName != nil,
		},
	}

	if param.Password != nil {
		hashedPassword, err := util.HashPassword(param.GetPassword())
		if err != nil {
			return nil, status.Errorf(codes.Internal, "failed to hash the password %s", hashedPassword)
		}

		data.HashedPassword = pgtype.Text{
			String: hashedPassword,
			Valid:  true,
		}

		data.PasswordChangedAt = pgtype.Timestamptz{
			Time:  time.Now(),
			Valid: true,
		}
	}

	user, err := server.store.UpdateUserUsingCaseSecond(ctx, data)
	if err != nil {
		log.Println(err)
		if errors.Is(err, db.ErrRecordNotFound) {
			return nil, status.Errorf(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, "failed to update user")
	}

	response := &pb.UpdateUserResponse{
		User: convertUser(user),
	}
	return response, nil
}

func validateUpdateUserRequest(param *pb.UpdateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(param.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	fmt.Printf("param %+v", param)
	if param.Password != nil {
		if err := validation.ValidatePassword(param.GetPassword()); err != nil {
			violations = append(violations, fieldViolation("password", err))
		}
	}
	if param.FullName != nil {
		if err := validation.ValidateFullname(param.GetFullName()); err != nil {
			violations = append(violations, fieldViolation("full_name", err))
		}
	}
	if param.Email != nil {
		if err := validation.ValidateEmail(param.GetEmail()); err != nil {
			violations = append(violations, fieldViolation("email", err))
		}
	}

	return violations
}
