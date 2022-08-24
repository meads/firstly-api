package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db/sqlc"
)

type FirstlyAPI struct {
	store  db.Store
	router *gin.Engine
}

// NewFirstlyAPI creates a new Http Server and sets up routing.
func NewFirstlyAPI(store db.Store, router *gin.Engine) *FirstlyAPI {
	api := &FirstlyAPI{
		store:  store,
		router: router,
	}

	api.router.LoadHTMLGlob("www/*.html")

	api.router.GET("/firstly/login/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	api.router.GET("/firstly/images/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	api.router.GET("/image/", api.listImages)
	api.router.POST("/image/", api.createImage)
	api.router.DELETE("/image/:id/", api.deleteImage)

	return api
}

// Start runs the Http server on the supplied address.
func (api *FirstlyAPI) Start(address string) error {
	return api.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
