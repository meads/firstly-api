package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
)

type FirstlyServer struct {
	router *gin.Engine
}

// NewFirstlyAPI creates a new Http Server and sets up routing.
func NewFirstlyServer(store db.Store, router *gin.Engine) *FirstlyServer {
	s := &FirstlyServer{
		router: router,
	}

	s.router.GET("/firstly/login/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	s.router.GET("/firstly/images/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	s.router.GET("/image/", s.ListImagesHandler(store))
	s.router.POST("/image/", s.CreateImageHandler(store))
	s.router.DELETE("/image/:id/", s.DeleteImageHandler(store))

	return s
}

// Start runs the Http server on the supplied address.
func (server *FirstlyServer) Start(address string) error {
	return server.router.Run(address)
}

func (server *FirstlyServer) LoadHTMLTemplates() {
	server.router.LoadHTMLGlob("www/*.html")
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
