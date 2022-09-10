package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
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
			body:         bytes.NewBufferString("{\"data\":\"test\"}"),
			method:       http.MethodPost,
			name:         "create handler responds with Status Code 200 when valid data supplied",
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
			name:         "create handler responds with Status Code 400 given no name is supplied",
			body:         bytes.NewBufferString("{\"data\":\"\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusBadRequest,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
			},
		},
		{
			name:         "create handler responds with Status Code 500 given there is some server error",
			body:         bytes.NewBufferString("{\"data\":\"server error\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusInternalServerError,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
				store.EXPECT().Create(gomock.Any(), "server error").Return(db.Image{}, errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 200 given valid request",
			method:       http.MethodDelete,
			responseCode: http.StatusOK,
			route:        "/image/69/",
			setupExpectations: func(store *db.MockStore) {
				store.EXPECT().Delete(gomock.Any(), int64(69)).Return(nil)
			},
		},
		{
			body:              bytes.NewBufferString(""),
			name:              "delete handler responds with Status Code 400 given param id not supplied",
			method:            http.MethodDelete,
			responseCode:      http.StatusBadRequest,
			route:             "/image//",
			setupExpectations: func(store *db.MockStore) {},
		},
		{
			body:              bytes.NewBufferString(""),
			name:              "delete handler responds with Status Code 400 given param id is not a valid integer",
			method:            http.MethodDelete,
			responseCode:      http.StatusBadRequest,
			route:             "/image/invalid/",
			setupExpectations: func(store *db.MockStore) {},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 500 given there is a server error",
			method:       http.MethodDelete,
			responseCode: http.StatusInternalServerError,
			route:        "/image/69/",
			setupExpectations: func(store *db.MockStore) {
				store.EXPECT().Delete(gomock.Any(), int64(69)).Return(errors.New("oops"))
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

			result := responseRecorder.Result()
			defer result.Body.Close()

			// Assert
			assert.Equal(t, test.responseCode, result.StatusCode)

			response := db.Image{}

			if result.Body != http.NoBody {
				if err := json.NewDecoder(result.Body).Decode(&response); err != nil && !errors.Is(err, io.EOF) {
					t.Errorf("Error decoding response body: %v", err)
					t.Log()
					t.Log(responseRecorder.Body)
					t.Log()
				}
			}
		})
	}
}
