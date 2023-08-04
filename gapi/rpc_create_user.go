package gapi

import (
	"context"
	"errors"
	"time"

	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/alrasyidin/simplebank-go/pb"
	"github.com/alrasyidin/simplebank-go/util"
	"github.com/alrasyidin/simplebank-go/validation"
	"github.com/alrasyidin/simplebank-go/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/rs/zerolog/log"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (server *Server) CreateUser(ctx context.Context, param *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	violations := validateCreateUserRequest(param)
	if violations != nil {
		return nil, invalidArgumentError(violations)
	}

	hashedPassword, err := util.HashPassword(param.GetPassword())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash the password %s", hashedPassword)
	}

	data := db.CreateUserParams{
		Username:       param.GetUsername(),
		Email:          param.GetEmail(),
		HashedPassword: hashedPassword,
		FullName:       param.GetFullName(),
	}

	arg := db.CreateUserTxParam{
		CreateUserParams: data,
		AfterCreate: func(user db.User) error {
			taskPayload := &worker.PayloadSendVerifyEmail{
				Username: user.Username,
			}
			opts := []asynq.Option{
				asynq.MaxRetry(10),
				asynq.ProcessIn(time.Second * 15),
				asynq.Queue(worker.QueueCritical),
			}
			return server.taskDistributor.DistributeTaskSendVerifyEmail(ctx, taskPayload, opts...)
		},
	}

	result, err := server.store.CreateUserTx(ctx, arg)
	if err != nil {
		log.Print(err)
		var pqErr *pgconn.PgError
		if errors.As(err, &pqErr) {
			switch pqErr.Code {
			case db.ErrUniqueViolations:
				return nil, status.Errorf(codes.AlreadyExists, "username already exist")
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user")
	}

	response := &pb.CreateUserResponse{
		User: convertUser(result.User),
	}
	return response, nil
}

func validateCreateUserRequest(param *pb.CreateUserRequest) (violations []*errdetails.BadRequest_FieldViolation) {
	if err := validation.ValidateUsername(param.GetUsername()); err != nil {
		violations = append(violations, fieldViolation("username", err))
	}
	if err := validation.ValidatePassword(param.GetPassword()); err != nil {
		violations = append(violations, fieldViolation("password", err))
	}
	if err := validation.ValidateFullname(param.GetFullName()); err != nil {
		violations = append(violations, fieldViolation("full_name", err))
	}
	if err := validation.ValidateEmail(param.GetEmail()); err != nil {
		violations = append(violations, fieldViolation("email", err))
	}

	return violations
}
