package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/meads/firstly-api/db"
)

func TestImageHandler(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
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
			method:       http.MethodPost,
			name:         "should respond with Status Code 200 when valid data supplied",
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
		{
			name: "should respond with Status Code 400 given no name is supplied",
			body: func() *bytes.Buffer {
				return bytes.NewBufferString("{\"data\":\"\"}")
			}(),
			method:       http.MethodPost,
			responseCode: http.StatusBadRequest,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
			},
		},
		{
			name: "should respond with Status Code 500 given there is some server error",
			body: func() *bytes.Buffer {
				return bytes.NewBufferString("{\"data\":\"server error\"}")
			}(),
			method:       http.MethodPost,
			responseCode: http.StatusInternalServerError,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
				store.EXPECT().Create(gomock.Any(), "server error").Return(db.Image{}, errors.New("oops"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
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

			// Assert
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
