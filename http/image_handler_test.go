package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/meads/firstly-api/db"
)

func TestImageAPI(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		callback          func(w *httptest.ResponseRecorder) bool
		method            string
		name              string
		responseCode      int
		route             string
		setupExpectations func(store *db.MockStore)
	}{
		{
			body: func() *bytes.Buffer {
				return bytes.NewBufferString("{\"data\":\"test\"}")
			}(),
			callback: func(w *httptest.ResponseRecorder) bool {
				return true
			},
			method:       http.MethodPost,
			name:         "should create image record when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
				store.EXPECT().Create(gomock.Any(), "test").Return(
					db.Image{
						ID:      1,
						Data:    "test",
						Created: time.Now().String(),
						Deleted: false,
					}, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			router := gin.Default()
			gin.SetMode(gin.TestMode)

			ctrl := gomock.NewController(t)
			mockStore := db.NewMockStore(ctrl)
			test.setupExpectations(mockStore)

			NewFirstlyServer(mockStore, router)
			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(test.method, test.route, test.body)
			router.ServeHTTP(responseRecorder, request)

			assert.Equal(t, test.responseCode, responseRecorder.Code)

			response := db.Image{}

			t.Log()
			t.Log(responseRecorder.Body)
			t.Log()

			if err := json.NewDecoder(responseRecorder.Body).Decode(&response); err != nil {
				t.Errorf("Error decoding response body: %v", err)
			}

		})
	}
}
