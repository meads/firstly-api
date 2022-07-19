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
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()
	router.LoadHTMLGlob("www/*.html")

	router.GET("/app/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	router.GET("/app/images", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	router.GET("/app/image", server.listImages)
	router.DELETE("/app/image/:id", server.deleteImage)
	router.POST("/app/image", server.createImage)
	// router.GET("/config/:id", server.getConfig)
	server.router = router
	return server
}

// Start runs the Http server on the supplied address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
