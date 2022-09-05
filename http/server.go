package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	api "github.com/meads/firstly-api/api"
)

type FirstlyServer struct {
	api    *api.ImageAPI
	router *gin.Engine
}

// NewFirstlyAPI creates a new Http Server and sets up routing.
func NewFirstlyServer(api *api.ImageAPI, router *gin.Engine) *FirstlyServer {
	s := &FirstlyServer{
		api:    api,
		router: router,
	}

	s.router.GET("/firstly/login/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	s.router.GET("/firstly/images/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	s.router.GET("/image/", s.ListImagesHandler)
	s.router.POST("/image/", s.CreateImageHandler)
	s.router.DELETE("/image/:id/", s.DeleteImageHandler)

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
