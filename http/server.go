package http

import (
	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
)

type FirstlyServer struct {
	claimer security.Claimer
	hasher  security.Hasher
	router  *gin.Engine
	store   db.Store
}

var firstly = &FirstlyServer{}

// NewFirstlyAPI creates a new Http Server and sets up routing.
func NewFirstlyServer(claimer security.Claimer, hasher security.Hasher, router *gin.Engine, store db.Store) *FirstlyServer {
	firstly.claimer = claimer
	firstly.hasher = hasher
	firstly.router = router
	firstly.store = store

	firstly.router.POST("/signin/", signinHandler)
	firstly.router.GET("/welcome/", welcomeHandler)
	firstly.router.POST("/refresh/", refreshHandler)

	firstly.router.POST("/account/", createAccountHandler)
	firstly.router.GET("/account/", claimsMiddleware(listAccountsHandler))
	firstly.router.PATCH("/account/", claimsMiddleware(updateAccountHandler))
	firstly.router.DELETE("/account/:id/", claimsMiddleware(deleteAccountHandler))

	firstly.router.GET("/image/", claimsMiddleware(listImagesHandler))
	firstly.router.POST("/image/", claimsMiddleware(createImageHandler(store)))
	firstly.router.DELETE("/image/:id/", claimsMiddleware(deleteImageHandler(store)))
	firstly.router.PATCH("/image/", claimsMiddleware(updateImageHandler(store)))

	return firstly
}

// Start runs the Http server on the supplied address.
func (server *FirstlyServer) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
