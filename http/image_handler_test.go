package http

import (
	"bytes"
	"database/sql"
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
	"github.com/meads/firstly-api/security"
)

func TestImageHandler(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		method            string
		name              string
		responseCode      int
		route             string
		isList            bool
		setupExpectations func(store *db.MockStore)
	}{
		{
			body:         bytes.NewBufferString("{\"data\":\"test\"}"),
			method:       http.MethodPost,
			name:         "create handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
				store.EXPECT().CreateImage(gomock.Any(), "test").Return(
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
				store.EXPECT().CreateImage(gomock.Any(), "server error").Return(db.Image{}, errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 200 given valid request",
			method:       http.MethodDelete,
			responseCode: http.StatusOK,
			route:        "/image/69/",
			setupExpectations: func(store *db.MockStore) {
				store.EXPECT().DeleteImage(gomock.Any(), int64(69)).Return(nil)
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
				store.EXPECT().DeleteImage(gomock.Any(), int64(69)).Return(errors.New("oops"))
			},
		},
		{
			body:              bytes.NewBufferString(""),
			name:              "list handler responds with Status Code 400 given limit param is invalid int",
			method:            http.MethodGet,
			responseCode:      http.StatusBadRequest,
			route:             "/image/?limit=invalid",
			setupExpectations: func(store *db.MockStore) {},
		},
		{
			body:              bytes.NewBufferString(""),
			name:              "list handler responds with Status Code 400 given offset param is invalid int",
			method:            http.MethodGet,
			responseCode:      http.StatusBadRequest,
			route:             "/image/?offset=invalid",
			setupExpectations: func(store *db.MockStore) {},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 500 given there is a server error",
			method:       http.MethodGet,
			responseCode: http.StatusInternalServerError,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
				params := db.ListImagesParams{Limit: 50, Offset: 0}
				store.EXPECT().ListImages(gomock.Any(), params).Return([]db.Image{}, errors.New("oops."))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 200 given a valid request",
			method:       http.MethodGet,
			responseCode: http.StatusOK,
			route:        "/image/",
			isList:       true,
			setupExpectations: func(store *db.MockStore) {
				params := db.ListImagesParams{Limit: 50, Offset: 0}
				store.EXPECT().ListImages(gomock.Any(), params).Return([]db.Image{
					{
						ID:      69,
						Data:    "foo",
						Created: "",
						Deleted: false,
					},
				}, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69, \"memo\": \"memo test\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
				params := db.UpdateImageParams{ID: int64(69), Memo: "memo test"}
				store.EXPECT().GetImage(gomock.Any(), params.ID).Return(db.Image{
					ID:   int64(69),
					Memo: "",
				}, nil)
				store.EXPECT().UpdateImage(gomock.Any(), params).Return(nil)
			},
		},
		{
			body:              bytes.NewBufferString("{\"id\":69}"),
			method:            http.MethodPatch,
			name:              "update handler responds with Status Code 400 when invalid data supplied",
			responseCode:      http.StatusBadRequest,
			route:             "/image/",
			setupExpectations: func(store *db.MockStore) {},
		},
		{
			body:         bytes.NewBufferString("{\"id\":68, \"memo\":\"memo test\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 404 when record not found",
			responseCode: http.StatusNotFound,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
				params := db.UpdateImageParams{ID: int64(68), Memo: "memo test"}
				store.EXPECT().GetImage(gomock.Any(), params.ID).Return(db.Image{}, sql.ErrNoRows)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":68, \"memo\":\"memo test\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on get before update",
			responseCode: http.StatusInternalServerError,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
				params := db.UpdateImageParams{ID: int64(68), Memo: "memo test"}
				store.EXPECT().GetImage(gomock.Any(), params.ID).Return(db.Image{}, errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69, \"memo\": \"memo test\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on update",
			responseCode: http.StatusInternalServerError,
			route:        "/image/",
			setupExpectations: func(store *db.MockStore) {
				params := db.UpdateImageParams{ID: int64(69), Memo: "memo test"}
				store.EXPECT().GetImage(gomock.Any(), params.ID).Return(db.Image{
					ID:   int64(69),
					Memo: "",
				}, nil)
				store.EXPECT().UpdateImage(gomock.Any(), params).Return(errors.New("oops"))
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
			mockHasher := security.NewMockHasher(ctrl)
			mockClaimer := security.NewMockClaimer(ctrl)

			test.setupExpectations(mockStore)

			NewFirstlyServer(mockClaimer, mockHasher, router, mockStore)
			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(test.method, test.route, test.body)
			router.ServeHTTP(responseRecorder, request)

			result := responseRecorder.Result()
			defer result.Body.Close()

			// Assert
			assert.Equal(t, test.responseCode, result.StatusCode)

			if !test.isList {
				response := db.Image{}

				if result.Body != http.NoBody {
					if err := json.NewDecoder(result.Body).Decode(&response); err != nil && !errors.Is(err, io.EOF) {
						t.Errorf("Error decoding response body: %v", err)
						t.Log()
						t.Log(responseRecorder.Body)
						t.Log()
					}
				}
			} else {
				response := []db.Image{}
				if result.Body != http.NoBody {
					if err := json.NewDecoder(result.Body).Decode(&response); err != nil && !errors.Is(err, io.EOF) {
						t.Errorf("Error decoding response body: %v", err)
						t.Log()
						t.Log(responseRecorder.Body)
						t.Log()
					}
				}
			}
		})
	}
}
