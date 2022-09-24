package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
)

type FirstlyServer struct {
	router *gin.Engine
}

// NewFirstlyAPI creates a new Http Server and sets up routing.
func NewFirstlyServer(store db.Store, hasher security.Hasher, router *gin.Engine) *FirstlyServer {
	s := &FirstlyServer{
		router: router,
	}

	s.router.Static("/assets", "./www/assets")

	s.router.GET("/firstly/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	s.router.GET("/firstly/images/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "images.html", nil)
	})

	s.router.GET("/firstly/sign-up/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sign-up.html", nil)
	})

	s.router.POST("/account/", s.CreateAccountHandler(store, hasher))
	s.router.GET("/account/", s.ListAccountsHandler(store))
	s.router.PATCH("/account/", s.UpdateAccountHandler(store, hasher))
	s.router.DELETE("/account/:id/", s.DeleteAccountHandler(store))

	s.router.POST("/login/", s.LoginAccountHandler(store, hasher))

	s.router.GET("/image/", s.ListImagesHandler(store))
	s.router.POST("/image/", s.CreateImageHandler(store))
	s.router.DELETE("/image/:id/", s.DeleteImageHandler(store))
	s.router.PATCH("/image/", s.UpdateImageHandler(store))

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
