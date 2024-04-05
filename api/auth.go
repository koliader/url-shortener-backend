package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	db "github.com/koliadervyanko/url-shortener-backend.git/db/sqlc"
	"github.com/koliadervyanko/url-shortener-backend.git/util"
)

const authError = "incorrect login or password"

type registerUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}
type authRes struct {
	Token string `json:"token" binding:"required"`
}

func (s *Server) registerUser(ctx *gin.Context) {
	var req registerUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("error hashing password %v", err)))
		return
	}
	arg := db.CreateUserParams{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Color:    util.RandomColor(),
	}
	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.ErrUniqueViolation.Code {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("user with this email or username is created")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	token, err := s.tokenMaker.CreateToken(user.Email, user.Username, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := authRes{
		Token: token,
	}
	ctx.JSON(http.StatusOK, res)
}

type loginUserReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(err.Error() == pgx.ErrNoRows.Error())
		if err.Error() == pgx.ErrNoRows.Error() {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf(authError)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(user.Password, req.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf(authError)))
		return
	}

	token, err := s.tokenMaker.CreateToken(user.Email, user.Username, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := authRes{
		Token: token,
	}
	ctx.JSON(http.StatusOK, res)

}
