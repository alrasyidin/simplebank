package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"time"

	db "github.com/alrasyidin/simplebank-go/db/sqlc"
	"github.com/alrasyidin/simplebank-go/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createUserParams struct {
	Username string `json:"username" binding:"required,alphanum"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

type createUserResponse struct {
	Username          string
	Email             string
	FullName          string
	CreatedAt         time.Time
	PasswordChangedAt time.Time
}

func newUserResponse(user db.User) createUserResponse {
	return createUserResponse{
		Username:          user.Username,
		Email:             user.Email,
		FullName:          user.FullName,
		CreatedAt:         user.CreatedAt,
		PasswordChangedAt: user.PasswordChangedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var param createUserParams

	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(param.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	data := db.CreateUserParams{
		Username:       param.Username,
		Email:          param.Email,
		HashedPassword: hashedPassword,
		FullName:       param.FullName,
	}

	user, err := server.store.CreateUser(ctx, data)
	if err != nil {
		log.Println(err)
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	response := newUserResponse(user)

	ctx.JSON(http.StatusOK, response)
}

type loginUserParams struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}
type loginUserResponse struct {
	AccessToken string  `json:"access_token"`
	User        db.User `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var param loginUserParams
	if err := ctx.ShouldBindJSON(&param); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, param.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(param.Password, user.HashedPassword)

	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	token, err := server.tokenGenerator.CreateToken(user.Username, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	resp := loginUserResponse{
		AccessToken: token,
		User:        user,
	}

	ctx.JSON(http.StatusOK, resp)
}
