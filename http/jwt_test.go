package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	db "github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
)

func TestJWTHandler(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		method            string
		name              string
		responseCode      int
		route             string
		setupExpectations func(store *db.MockStore, hasher *security.MockHasher)
	}{
		{
			body:         bytes.NewBufferString("{\"phrase\":\"blah\",\"invalid\":\"test\"}"),
			method:       http.MethodPost,
			name:         "signin returns status code bad request when invalid json supplied",
			responseCode: http.StatusBadRequest,
			route:        "/signin/",
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher) {
				store.EXPECT().GetAccountByUsername(gomock.Any(), gomock.Any()).Times(0)
				hasher.EXPECT().IsValidPassword(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			body:         bytes.NewBufferString("{\"phrase\":\"valid\",\"username\":\"invalid\"}"),
			method:       http.MethodPost,
			name:         "signin returns status code bad request when invalid username supplied",
			responseCode: http.StatusBadRequest,
			route:        "/signin/",
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher) {
				store.EXPECT().GetAccountByUsername(gomock.Any(), "invalid").Return(db.Account{}, errors.New("oops"))
				hasher.EXPECT().IsValidPassword(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			body:         bytes.NewBufferString("{\"phrase\":\"valid\",\"username\":\"valid\"}"),
			method:       http.MethodPost,
			name:         "signin returns status code unauthorized when call fails to validate phrase",
			responseCode: http.StatusInternalServerError,
			route:        "/signin/",
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher) {
				expectedAccount := db.Account{
					Phrase: []byte("invalid"),
					Salt:   "salt",
				}
				store.EXPECT().GetAccountByUsername(gomock.Any(), "valid").Return(expectedAccount, nil)
				hasher.EXPECT().IsValidPassword(expectedAccount.Phrase, expectedAccount.Salt, "valid").Return(
					false,
					errors.New("secret not set"),
				)
			},
		},
		{
			body:         bytes.NewBufferString("{\"phrase\":\"invalid\",\"username\":\"valid\"}"),
			method:       http.MethodPost,
			name:         "signin returns status code unauthorized when invalid phrase supplied",
			responseCode: http.StatusUnauthorized,
			route:        "/signin/",
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher) {
				expectedAccount := db.Account{
					Phrase: []byte("valid"),
					Salt:   "salt",
				}
				store.EXPECT().GetAccountByUsername(gomock.Any(), "valid").Return(expectedAccount, nil)
				hasher.EXPECT().IsValidPassword(expectedAccount.Phrase, expectedAccount.Salt, "invalid").Return(false, nil)
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
			test.setupExpectations(mockStore, mockHasher)

			NewFirstlyServer(mockStore, mockHasher, router)
			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(test.method, test.route, test.body)
			router.ServeHTTP(responseRecorder, request)

			result := responseRecorder.Result()
			defer result.Body.Close()

			// Assert
			assert.Equal(t, test.responseCode, result.StatusCode)

			// response := db.Image{}

			// if result.Body != http.NoBody {
			// 	if err := json.NewDecoder(result.Body).Decode(&response); err != nil && !errors.Is(err, io.EOF) {
			// 		t.Errorf("Error decoding response body: %v", err)
			// 		t.Log()
			// 		t.Log(responseRecorder.Body)
			// 		t.Log()
			// 	}
			// }
		})
	}
}
