package gapi

import (
	"context"
	"log"

	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/alrasyidin/simplebank-go/pb"
	"github.com/alrasyidin/simplebank-go/validation"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) VerifyEmail(ctx context.Context, param *pb.VerifyEmailRequest) (*pb.VerifyEmailResponse, error) {
	violations := validateVerifyEmailRequest(param)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	arg := db.VerifyEmailTxParam{
		EmailId:    param.EmailId,
		SecretCode: param.SecretCode,
	}
	result, err := server.store.VerifyEmailTx(ctx, arg)
	if err != nil {
		log.Println(err)
		return nil, status.Errorf(codes.Internal, "failed to verify email user")
	}

	response := &pb.VerifyEmailResponse{
		IsVerified: result.User.IsEmailActivated,
	}
	return response, nil
}

func validateVerifyEmailRequest(param *pb.VerifyEmailRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateEmailID(param.GetEmailId()); err != nil {
		violations = append(violations, fieldViolation("email_id", err))
	}
	if err := validation.ValidateSecretCode(param.GetSecretCode()); err != nil {
		violations = append(violations, fieldViolation("secret_code", err))
	}

	return violations
}
