package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx"
	db "github.com/koliadervyanko/url-shortener-backend.git/db/sqlc"
	"github.com/koliadervyanko/url-shortener-backend.git/util"
)

const authError = "incorrect login or password"

var emptyUser db.User = db.User{}

type registerUserReq struct {
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
		Password: &hashedPassword,
		Color:    util.RandomColor(),
	}
	user, err := s.store.CreateUser(ctx, arg)
	if err != nil {
		if db.ErrorCode(err) == db.ErrUniqueViolation.Code {
			ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("user with this username is created")))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	token, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
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
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}

func (s *Server) loginUser(ctx *gin.Context) {
	var req loginUserReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := s.store.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if err.Error() == pgx.ErrNoRows.Error() {
			ctx.JSON(http.StatusNotFound, errorResponse(fmt.Errorf(authError)))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	// check if github auth user trying to login with simple auth
	if user != emptyUser && user.Password == nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(fmt.Errorf("login with github")))
		return
	}

	token, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := authRes{
		Token: token,
	}
	ctx.JSON(http.StatusOK, res)

}

type getGithubAccessTokenReq struct {
	Code string `uri:"code" binding:"required"`
}
type githubUserResponse struct {
	Login string `json:"login"`
	ID    int    `json:"id"`
	Name  string `json:"name"`
}

func (s *Server) githubAuth(ctx *gin.Context) {
	var req getGithubAccessTokenReq
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	//* get github token using oauthConfig
	token, err := s.oauthConfig.Exchange(ctx, req.Code)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to exchange token %v", err)))
		return
	}
	httpClient := s.oauthConfig.Client(ctx, token)
	resp, err := httpClient.Get("https://api.github.com/user")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to fetch user data")))
		return
	}
	defer resp.Body.Close()

	var user githubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(fmt.Errorf("failed to decode user data")))
		return
	}
	dbUser, err := s.store.GetUserByUsername(ctx, user.Login)
	if err != nil {
		if err.Error() != pgx.ErrNoRows.Error() {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return

		}
		// if err is equal to not not found
		arg := db.CreateUserParams{
			Username: user.Login,
			Password: nil,
			Color:    util.RandomColor(),
		}
		user, err := s.store.CreateUser(ctx, arg)
		if err != nil {
			if db.ErrorCode(err) == db.UniqueViolation {
				ctx.JSON(http.StatusBadRequest, errorResponse(fmt.Errorf("user with this data was already created")))
				return
			}
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		jwtToken, err := s.tokenMaker.CreateToken(user.Username, s.config.AccessTokenDuration)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}
		res := authRes{
			Token: jwtToken,
		}
		ctx.JSON(http.StatusOK, res)
		return
	}
	jwtToken, err := s.tokenMaker.CreateToken(dbUser.Username, s.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := authRes{
		Token: jwtToken,
	}
	ctx.JSON(http.StatusOK, res)
}
