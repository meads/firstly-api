package http

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db/sqlc"
	"github.com/meads/firstly-api/security"
)

type FirstlyServer struct {
	tokener security.Tokener
	hasher  security.Hasher
	router  *gin.Engine
	store   db.Querier
}

var firstly = &FirstlyServer{}

// NewFirstlyAPI creates a new Http Server and sets up routing.
func NewFirstlyServer(tokener security.Tokener, hasher security.Hasher, router *gin.Engine, store db.Querier) *FirstlyServer {
	firstly.tokener = tokener
	firstly.hasher = hasher
	firstly.router = router
	firstly.store = store

	// https://github.com/gin-gonic/gin/blob/master/docs/doc.md#dont-trust-all-proxies
	firstly.router.SetTrustedProxies(nil)

	// https://github.com/gin-gonic/gin/issues/3681
	config := cors.DefaultConfig()
	ALLOW_ORIGINS := os.Getenv("ALLOW_ORIGINS")
	if ALLOW_ORIGINS == "" {
		ALLOW_ORIGINS = "http://localhost:3000"
	}
	config.AllowOrigins = []string{ALLOW_ORIGINS}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"}
	config.ExposeHeaders = []string{"Authorization"}
	firstly.router.Use(cors.New(config))

	firstly.router.POST("/register/", registerUserHandler)
	firstly.router.POST("/login/", loginHandler)
	firstly.router.POST("/refresh/", renewAccessTokenHandler)
	firstly.router.POST("/logout/:sessionid", logoutHandler)
	firstly.router.POST("/revoke/:sessionid", revokeSessionHandler)

	// firstly.router.GET("/users/:id", claimsMiddleware(getUserHandler))
	firstly.router.GET("/users/", claimsMiddleware(listUsersHandler))
	firstly.router.PATCH("/users/", claimsMiddleware(patchUserHandler))
	firstly.router.DELETE("/users/:userid/", claimsMiddleware(deleteUserHandler))

	firstly.router.POST("/users/:userid/notes/", claimsMiddleware(createNoteHandler))
	firstly.router.PUT("/users/:userid/notes/", claimsMiddleware(updateNoteHandler))
	firstly.router.GET("/users/:userid/notes/", claimsMiddleware(listNotesHandler))
	firstly.router.DELETE("/users/:userid/notes/:noteid", claimsMiddleware(deleteNoteHandler))

	return firstly
}

// Start runs the Http server on the supplied address.
func (server *FirstlyServer) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
