package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/koliadervyanko/url-shortener-backend.git/db/sqlc"
)

const userNotFound = "user not found"

type convertedUser struct {
	Email    string `json:"email" binding:"required,email"`
	Username string `json:"username" binding:"required"`
	Color    string `json:"color" binding:"required"`
}

func (s *Server) convertUser(user db.User) convertedUser {
	converted := convertedUser{
		Email:    user.Email,
		Username: user.Username,
		Color:    user.Color,
	}
	return converted
}

type getUserByEmailReq struct {
	Email string `uri:"email" binding:"required,email"`
}

func (s *Server) getUserByEmail(ctx *gin.Context) {
	var req getUserByEmailReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := s.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf(userNotFound)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	converted := s.convertUser(user)
	ctx.JSON(http.StatusOK, converted)
}
