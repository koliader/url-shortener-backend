package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/koliadervyanko/url-shortener-backend.git/db/sqlc"
	"github.com/koliadervyanko/url-shortener-backend.git/token"
)

const userNotFound = "user not found"

type convertedUser struct {
	Username string `json:"username" binding:"required"`
	Color    string `json:"color" binding:"required"`
}

func (s *Server) convertUser(user db.User) convertedUser {
	converted := convertedUser{
		Username: user.Username,
		Color:    user.Color,
	}
	return converted
}

type getUserByUsernameReq struct {
	Username string `uri:"username" binding:"required"`
}

func (s *Server) getUserByUsername(ctx *gin.Context) {
	var req getUserByUsernameReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if authPayload.Username != req.Username {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("you can't get data of another user")))
		return
	}
	user, err := s.store.GetUserByUsername(ctx, req.Username)
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
