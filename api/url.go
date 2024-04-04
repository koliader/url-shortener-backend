package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	db "github.com/koliadervyanko/url-shortener-backend.git/db/sqlc"
	"github.com/koliadervyanko/url-shortener-backend.git/token"
	"github.com/koliadervyanko/url-shortener-backend.git/util"
)

type convertedUrl struct {
	Url  string         `json:"url" binding:"required,http_url"`
	Code string         `json:"code" binding:"required"`
	User *convertedUser `json:"user" binding:"required"`
}

func (s *Server) convertUrl(ctx *gin.Context, url db.Url) (*convertedUrl, error) {
	var convertedUser *convertedUser
	if url.Owner != nil {
		user, err := s.store.GetUserByEmail(ctx, *url.Owner)
		if err != nil {
			return nil, err
		}
		res := s.convertUser(user)
		convertedUser = &res
	}
	converted := convertedUrl{
		Url:  url.Url,
		Code: url.Code,
		User: convertedUser,
	}
	return &converted, nil
}

type createUrlReq struct {
	Url string `json:"url" binding:"required,http_url"`
}

func (s *Server) createGuestUrl(ctx *gin.Context) {
	var req createUrlReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateUrlParams{
		Url:   req.Url,
		Code:  util.RandomString(5),
		Owner: nil,
	}
	url, err := s.store.CreateUrl(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, url)
}

func (s *Server) createUrl(ctx *gin.Context) {
	var req createUrlReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	arg := db.CreateUrlParams{
		Url:   req.Url,
		Code:  util.RandomString(5),
		Owner: &authPayload.Email,
	}
	url, err := s.store.CreateUrl(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	converted, err := s.convertUrl(ctx, url)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, converted)
}

type getUrlByCodeReq struct {
	Code string `uri:"code" binding:"required"`
}

func (s *Server) getUrlByCode(ctx *gin.Context) {
	var req getUrlByCodeReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	url, err := s.store.GetUrlByCode(ctx, req.Code)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("url not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	converted, err := s.convertUrl(ctx, url)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("owner not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, converted)
}

func (s *Server) listUrlsByOwner(ctx *gin.Context) {
	var converted []convertedUrl
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	urls, err := s.store.ListUrlsByUser(ctx, &authPayload.Email)
	if err != nil {
		if err == pgx.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user not found")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	for _, url := range urls {
		convertedUrl, err := s.convertUrl(ctx, url)
		if err != nil {
			if err == pgx.ErrNoRows {
				ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf("user not found")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		converted = append(converted, *convertedUrl)
	}
	ctx.JSON(http.StatusOK, converted)
}
