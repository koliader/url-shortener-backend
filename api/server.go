package api

import (
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/koliadervyanko/url-shortener-backend.git/db/sqlc"
	"github.com/koliadervyanko/url-shortener-backend.git/token"
	"github.com/koliadervyanko/url-shortener-backend.git/util"
	"golang.org/x/oauth2"
)

type Server struct {
	config      util.Config
	store       db.Store
	router      *gin.Engine
	tokenMaker  token.Maker
	oauthConfig oauth2.Config
}

var (
	githubScopes = []string{"user:email"}
)

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewJWTMaker(config.TokenKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create tokenManager: %v", err)
	}
	oauthConfig := oauth2.Config{
		ClientID:     config.GithubClientId,
		ClientSecret: config.GithubClientSecret,
		Endpoint: oauth2.Endpoint{
			TokenURL: "https://github.com/login/oauth/access_token",
		},
		Scopes:      githubScopes,
		RedirectURL: "https://url-shortener-frontend-git-main-koliadervyankos-projects.vercel.app/",
	}

	server := &Server{store: store, config: config, tokenMaker: tokenMaker, oauthConfig: oauthConfig}

	server.setupRouter()

	return server, nil
}

func (s *Server) setupRouter() {
	router := gin.Default()

	// cors
	c := cors.New(cors.Config{
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodDelete,
		},
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowHeaders:     []string{"Content-Type", "Authorization"},
	})

	router.Use(c)
	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker, s.store))

	// auth
	router.POST("/auth/register", s.registerUser)
	router.POST("/auth/login", s.loginUser)
	router.GET("/auth/github/:code", s.githubAuth)

	// urls
	router.POST("/urls/guest", s.createGuestUrl)
	router.GET("/urls/:code", s.getUrlByCode)
	router.PUT("/urls/:code", s.updateClicks)
	authRoutes.POST("/urls", s.createUrl)
	authRoutes.GET("/urls/myUrls", s.listUrlsByOwner)

	// users
	authRoutes.GET("/users/:username", s.getUserByUsername)

	s.router = router
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

type Success struct {
	Success bool `json:"success"`
}

type Color struct {
	Color string `json:"color"`
}
