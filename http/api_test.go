package http_test

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	db "github.com/meads/firstly-api/db"
	// . "github.com/onsi/ginkgo"
	// . "github.com/onsi/gomega"
)

// Helper function to create a router during testing
func getRouter() *gin.Engine {
	r := gin.Default()
	r.LoadHTMLGlob("www/*")
	return r
}

// Helper function to process a request and test its response
func testHTTPResponse(t *testing.T, r *gin.Engine, req *http.Request, f func(w *httptest.ResponseRecorder) bool) {
	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	if !f(w) {
		t.Fail()
	}
}

func getImageRecord(id int64, data string, deleted sql.NullInt32, created string) *db.Image {
	return &db.Image{
		ID:      id,
		Data:    data,
		Created: created,
		Deleted: deleted,
	}
}

// var _ = Describe("ImageHandlers", func() {
// 	var (
// 		ctrl   *gomock.Controller
// 		store  *mockapi.MockStore
// 		sut    *server.FirstlyAPI
// 		t      GinkgoTInterface
// 		router *gin.Engine
// 	)
// 	BeforeEach(func() {
// 		t = GinkgoT()
// 		ctrl = gomock.NewController(t)
// 		store = mockapi.NewMockStore(ctrl)
// 		router = gin.Default()
// 		sut = server.NewFirstlyAPI(store, router)
// 	})
// 	Describe("should have an image handler with POST method that expects data to create", func() {
// 		requestBody, _ := json.Marshal(map[string]string{
// 			"data": "test",
// 		})
// 		request, err := http.NewRequest(http.MethodPost, "/image/", bytes.NewBuffer(requestBody))
// 		if err != nil {
// 			GinkgoT().Fail()
// 		}

// 		testHTTPResponse(t, router, request, func(w *httptest.ResponseRecorder) bool {
// 			response := w.Body.Bytes()
// 			fmt.Println(string(response))

// 		})
// 	})
// })
