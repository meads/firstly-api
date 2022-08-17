package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewServer creates a new Http Server and sets up routing.
func NewServer(store db.Store, router *gin.Engine) *Server {
	s := &Server{
		store:  store,
		router: router,
	}
	s.router.RedirectTrailingSlash = false

	s.router.LoadHTMLGlob("www/*.html")

	s.router.GET("/firstly/login/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	s.router.GET("/firstly/images/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	s.router.GET("/image/", s.listImages)
	s.router.POST("/image/", s.createImage)
	s.router.DELETE("/image/:id/", s.deleteImage)

	return s
}

// Start runs the Http server on the supplied address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
