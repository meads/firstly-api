package http

import (
	"log"
	"os"
	"strings"

	cors "github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
)

type FirstlyServer struct {
	claimer security.Claimer
	hasher  security.Hasher
	router  *gin.Engine
	store   db.Querier
}

var firstly = &FirstlyServer{}

// NewFirstlyAPI creates a new Http Server and sets up routing.
func NewFirstlyServer(claimer security.Claimer, hasher security.Hasher, router *gin.Engine, store db.Querier) *FirstlyServer {
	firstly.claimer = claimer
	firstly.hasher = hasher
	firstly.router = router
	firstly.store = store

	// https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies
	firstly.router.SetTrustedProxies(nil)

	// https://github.com/gin-gonic/gin/issues/3681
	origins := os.Getenv("ALLOW_ORIGINS")
	origins = strings.TrimSpace(origins)
	var ALLOW_ORIGINS []string
	if origins == "" {
		log.Println("ALLOW_ORIGINS not set in config!")
		os.Exit(1)
	} else if strings.Contains(origins, ",") {
		ALLOW_ORIGINS = strings.Split(origins, ",")
	} else {
		ALLOW_ORIGINS = []string{origins}
	}

	config := cors.DefaultConfig()
	config.AllowOrigins = ALLOW_ORIGINS
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"}
	config.ExposeHeaders = []string{"Authorization"}
	firstly.router.Use(cors.New(config))

	firstly.router.POST("/signin/", signinHandler)
	// firstly.router.GET("/welcome/", welcomeHandler)
	firstly.router.POST("/refresh/", refreshHandler)

	firstly.router.POST("/account/", createAccountHandler)
	firstly.router.GET("/account/", claimsMiddleware(listAccountsHandler))
	firstly.router.PATCH("/account/", claimsMiddleware(updateAccountHandler))
	firstly.router.DELETE("/account/:id/", claimsMiddleware(deleteAccountHandler))

	firstly.router.GET("/protected/", claimsMiddleware(listProtectedHandler))
	// firstly.router.GET("/image/", claimsMiddleware(listImagesHandler))
	// firstly.router.GET("/image/", claimsMiddleware(listImagesHandler))
	// firstly.router.POST("/image/", claimsMiddleware(createImageHandler(store)))
	// firstly.router.DELETE("/image/:id/", claimsMiddleware(deleteImageHandler(store)))
	// firstly.router.PATCH("/image/", claimsMiddleware(updateImageHandler(store)))

	return firstly
}

// Start runs the Http server on the supplied address.
func (server *FirstlyServer) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
