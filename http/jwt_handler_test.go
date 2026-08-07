package http

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	db "github.com/meads/firstly-api/db/sqlc"
	"github.com/meads/firstly-api/internal/security"
	"go.uber.org/mock/gomock"
)

func TestJWTSignInHandler(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		method            string
		name              string
		responseCode      int
		route             string
		setupExpectations func(store *db.MockQuerier, hasher *security.MockHasher, tokener *security.MockTokener, rr *httptest.ResponseRecorder, r *http.Request)
	}{
		{
			body:         bytes.NewBufferString("{\"password\":\"blah\",\"invalid\":\"test\"}"),
			method:       http.MethodPost,
			name:         "login returns status code bad request when invalid json supplied",
			responseCode: http.StatusBadRequest,
			route:        "/login/",
			setupExpectations: func(store *db.MockQuerier, hasher *security.MockHasher, tokener *security.MockTokener, rr *httptest.ResponseRecorder, r *http.Request) {
				store.EXPECT().GetUserByUsername(gomock.Any(), gomock.Any()).Times(0)
				// hasher.EXPECT().IsValidPassword(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			body:         bytes.NewBufferString("{\"password\":\"valid\",\"username\":\"invalid\"}"),
			method:       http.MethodPost,
			name:         "login returns status code 401 when username supplied was not found",
			responseCode: http.StatusUnauthorized,
			route:        "/login/",
			setupExpectations: func(store *db.MockQuerier, hasher *security.MockHasher, tokener *security.MockTokener, rr *httptest.ResponseRecorder, r *http.Request) {
				store.EXPECT().GetUserByUsername(gomock.Any(), "invalid").Return(db.User{}, sql.ErrNoRows)
				hasher.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			body:         bytes.NewBufferString("{\"password\":\"valid\",\"username\":\"invalid\"}"),
			method:       http.MethodPost,
			name:         "login returns status code 500 when error querying username",
			responseCode: http.StatusInternalServerError,
			route:        "/login/",
			setupExpectations: func(store *db.MockQuerier, hasher *security.MockHasher, tokener *security.MockTokener, rr *httptest.ResponseRecorder, r *http.Request) {
				store.EXPECT().GetUserByUsername(gomock.Any(), "invalid").Return(db.User{}, errors.New("server error"))
				hasher.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			body:         bytes.NewBufferString("{\"password\":\"invalid\",\"username\":\"valid\"}"),
			method:       http.MethodPost,
			name:         "login returns status code unauthorized when call fails to validate password",
			responseCode: http.StatusUnauthorized,
			route:        "/login/",
			setupExpectations: func(store *db.MockQuerier, hasher *security.MockHasher, tokener *security.MockTokener, rr *httptest.ResponseRecorder, r *http.Request) {
				expectedUser := db.User{
					ID:       1,
					Password: "invalids-hash",
					Username: "valid",
				}
				store.EXPECT().GetUserByUsername(gomock.Any(), expectedUser.Username).Return(expectedUser, nil)
				hasher.EXPECT().ComparePassword(expectedUser.Password, "invalid").
					Return(errors.New("invalid username or password"))
				tokener.EXPECT().GenerateToken(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			body:         bytes.NewBufferString("{\"password\":\"valid\",\"username\":\"invalid\"}"),
			method:       http.MethodPost,
			name:         "login returns status code unauthorized when invalid username supplied",
			responseCode: http.StatusUnauthorized,
			route:        "/login/",
			setupExpectations: func(store *db.MockQuerier, hasher *security.MockHasher, tokener *security.MockTokener, rr *httptest.ResponseRecorder, r *http.Request) {
				expectedUser := db.User{}
				store.EXPECT().GetUserByUsername(gomock.Any(), "invalid").Return(expectedUser, sql.ErrNoRows)
				hasher.EXPECT().ComparePassword(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		// {
		// 	body:         bytes.NewBufferString("{\"password\":\"valid\",\"username\":\"valid\"}"),
		// 	method:       http.MethodPost,
		// 	name:         "login handler given valid credentials sets the Set-Cookie header with valid jwt claims token",
		// 	responseCode: http.StatusOK,
		// 	route:        "/login/",
		// 	setupExpectations: func(store *db.MockQuerier, hasher *security.MockHasher, tokener *security.MockTokener, rr *httptest.ResponseRecorder, r *http.Request) {
		// 		expectedAccount := db.Account{
		// 			Username: "valid",
		// 			Password: []byte("valid"),
		// 			Salt:     "salt",
		// 		}
		// 		store.EXPECT().GetAccountByUsername(gomock.Any(), expectedAccount.Username).
		// 			Return(expectedAccount, nil)
		// 		hasher.EXPECT().IsValidPassword(expectedAccount.Password, expectedAccount.Salt, "valid").
		// 			Return(true, nil)

		// 		// Create the JWT claims, which includes the username and expiry time
		// 		tokenString := "mocktoken"
		// 		expirationTime := time.Now().Add(5 * time.Minute)
		// 		claimer.EXPECT().GenerateToken(expectedAccount.Username).Return(tokenString, nil)

		// 		// Finally, we set the client cookie for "token" as the JWT we just generated
		// 		// we also set an expiry time which is the same as the token itself
		// 		r.AddCookie(&http.Cookie{
		// 			Name:    "token",
		// 			Value:   tokenString,
		// 			Expires: expirationTime,
		// 		})
		// 	},
		// },
		// {
		// 	body:         bytes.NewBufferString("{\"password\":\"valid\",\"username\":\"valid\"}"),
		// 	method:       http.MethodPost,
		// 	name:         "login handler given an error with call to GetFiveMinuteExpirationToken is encountered, 500 status code is the response",
		// 	responseCode: http.StatusInternalServerError,
		// 	route:        "/login/",
		// 	setupExpectations: func(store *db.MockQuerier, hasher *security.MockHasher, tokener *security.MockTokener, rr *httptest.ResponseRecorder, r *http.Request) {
		// 		expectedAccount := db.Account{
		// 			Username: "valid",
		// 			Password: []byte("valid"),
		// 			Salt:     "salt",
		// 		}
		// 		store.EXPECT().GetAccountByUsername(gomock.Any(), expectedAccount.Username).
		// 			Return(expectedAccount, nil)
		// 		hasher.EXPECT().IsValidPassword(expectedAccount.Password, expectedAccount.Salt, "valid").
		// 			Return(true, nil)

		// 		// Create the JWT claims, which includes the username and expiry time
		// 		tokenString := "mocktoken"
		// 		expirationTime := time.Now().Add(5 * time.Minute)
		// 		claimer.EXPECT().GenerateToken(expectedAccount.Username).Return("", errors.New("oops"))

		// 		// Finally, we set the client cookie for "token" as the JWT we just generated
		// 		// we also set an expiry time which is the same as the token itself
		// 		r.AddCookie(&http.Cookie{
		// 			Name:    "token",
		// 			Value:   tokenString,
		// 			Expires: expirationTime,
		// 		})
		// 	},
		// },
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := gin.New()
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockQuerier := db.NewMockQuerier(ctrl)
			mockHasher := security.NewMockHasher(ctrl)
			mockTokener := security.NewMockTokener(ctrl)

			NewFirstlyServer(mockTokener, mockHasher, router, mockQuerier)
			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(test.method, test.route, test.body)
			test.setupExpectations(mockQuerier, mockHasher, mockTokener, responseRecorder, request)

			// Act
			router.ServeHTTP(responseRecorder, request)

			result := responseRecorder.Result()
			defer result.Body.Close()

			// Assert
			assert.Equal(t, test.responseCode, result.StatusCode)
		})
	}
}
