package http_test

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	mockapi "github.com/meads/firstly-api/db/mock"
	db "github.com/meads/firstly-api/db/sqlc"
	http "github.com/meads/firstly-api/http"
	. "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

var _ = Describe("Image Http Handler", func() {
	var (
		ctrl   *gomock.Controller
		store  *mockapi.MockStore
		sut    http.Server
		router *gin.Engine
	)
	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		store = mockapi.NewMockStore(ctrl)
		router = gin.Default()
		sut = *http.NewServer(store, router)
	})
	Describe("should create an image", func() {
		BeforeEach(func() {
			sut.Start(gomock.Any().String())
		})
		store.EXPECT().CreateImage(gomock.Any(), "foo").Return(db.Image{})
	})
})
