package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
)

type FirstlyAPI struct {
	store  db.Store
	router *gin.Engine
}

// NewFirstlyAPI creates a new Http Server and sets up routing.
func NewFirstlyAPI(store db.Store, router *gin.Engine, loadGlobs bool) *FirstlyAPI {
	fapi := &FirstlyAPI{
		store:  store,
		router: router,
	}

	if loadGlobs {
		fapi.router.LoadHTMLGlob("www/*.html")
	}

	fapi.router.GET("/firstly/login/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	fapi.router.GET("/firstly/images/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", nil)
	})

	fapi.router.GET("/image/", fapi.ListImagesHandler)
	fapi.router.POST("/image/", fapi.CreateImageHandler)
	fapi.router.DELETE("/image/:id/", fapi.DeleteImageHandler)

	return fapi
}

// Start runs the Http server on the supplied address.
func (fapi *FirstlyAPI) Start(address string) error {
	return fapi.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
