package http

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	db "github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
)

func TestJWTSignInHandler(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		method            string
		name              string
		responseCode      int
		route             string
		setupExpectations func(store *db.MockStore, hasher *security.MockHasher, claimer *security.MockClaimer, rr *httptest.ResponseRecorder, r *http.Request)
	}{
		{
			body:         bytes.NewBufferString("{\"phrase\":\"blah\",\"invalid\":\"test\"}"),
			method:       http.MethodPost,
			name:         "signin returns status code bad request when invalid json supplied",
			responseCode: http.StatusBadRequest,
			route:        "/signin/",
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, claimer *security.MockClaimer, rr *httptest.ResponseRecorder, r *http.Request) {
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
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, claimer *security.MockClaimer, rr *httptest.ResponseRecorder, r *http.Request) {
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
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, claimer *security.MockClaimer, rr *httptest.ResponseRecorder, r *http.Request) {
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
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, claimer *security.MockClaimer, rr *httptest.ResponseRecorder, r *http.Request) {
				expectedAccount := db.Account{
					Phrase: []byte("valid"),
					Salt:   "salt",
				}
				store.EXPECT().GetAccountByUsername(gomock.Any(), "valid").Return(expectedAccount, nil)
				hasher.EXPECT().IsValidPassword(expectedAccount.Phrase, expectedAccount.Salt, "invalid").Return(false, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"phrase\":\"valid\",\"username\":\"valid\"}"),
			method:       http.MethodPost,
			name:         "signin handler given valid credentials sets the Set-Cookie header with valid jwt claims token",
			responseCode: http.StatusOK,
			route:        "/signin/",
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, claimer *security.MockClaimer, rr *httptest.ResponseRecorder, r *http.Request) {
				expectedAccount := db.Account{
					Username: "valid",
					Phrase:   []byte("valid"),
					Salt:     "salt",
				}
				store.EXPECT().GetAccountByUsername(gomock.Any(), expectedAccount.Username).
					Return(expectedAccount, nil)
				hasher.EXPECT().IsValidPassword(expectedAccount.Phrase, expectedAccount.Salt, "valid").
					Return(true, nil)

				// Create the JWT claims, which includes the username and expiry time
				tokenString := "mocktoken"
				expirationTime := time.Now().Add(5 * time.Minute)
				claimer.EXPECT().GetFiveMinuteExpirationToken(expectedAccount.Username).
					Return(tokenString, expirationTime, nil)

				// Finally, we set the client cookie for "token" as the JWT we just generated
				// we also set an expiry time which is the same as the token itself
				r.AddCookie(&http.Cookie{
					Name:    "token",
					Value:   tokenString,
					Expires: expirationTime,
				})
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

			NewFirstlyServer(mockStore, mockHasher, mockClaimer, router)
			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(test.method, test.route, test.body)
			test.setupExpectations(mockStore, mockHasher, mockClaimer, responseRecorder, request)

			// Act
			router.ServeHTTP(responseRecorder, request)

			result := responseRecorder.Result()
			defer result.Body.Close()

			// Assert
			assert.Equal(t, test.responseCode, result.StatusCode)
		})
	}
}

func TestJWTWelcomeHandler(t *testing.T) {
	tests := []struct {
		name              string
		responseCode      int
		route             string
		expectedBody      string
		setupExpectations func(claimer security.MockClaimer, r *http.Request, rr *httptest.ResponseRecorder)
	}{
		{
			name:         "welcome handler given no token cookie will respond with status unauthorized",
			responseCode: http.StatusUnauthorized,
			route:        "/welcome/",
			expectedBody: "",
			setupExpectations: func(claimer security.MockClaimer, r *http.Request, rr *httptest.ResponseRecorder) {
				cookies := r.Cookies()
				if len(cookies) > 0 {
					t.Fatal("expected no cookies in request scenario")
				}
			},
		},
		// {
		// 	name:         "welcome handler given valid token cookie will respond with status ok",
		// 	responseCode: http.StatusOK,
		// 	route:        "/welcome/",
		// 	expectedBody: "Welcome valid!",
		// 	setupExpectations: func(claimer security.MockClaimer, r *http.Request, rr *httptest.ResponseRecorder) {
		// 		tokenString := "mocktoken"
		// 		claimsValidator := security.ClaimsValidator{}
		// 		claimToken := &security.ClaimToken{}
		// 		claimer.EXPECT().ParseWithClaims(tokenString, claimsValidator, func(*jwt.Token) (interface{}, error) {
		// 			return gomock.Any(), nil
		// 		}).Return(
		// 			claimToken, nil,
		// 		)
		// 		// Finally, we set the client cookie for "token" as the JWT we just generated
		// 		// we also set an expiry time which is the same as the token itself
		// 		r.AddCookie(&http.Cookie{
		// 			Name:  "token",
		// 			Value: tokenString,
		// 		})
		// 	},
		// },
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

			NewFirstlyServer(mockStore, mockHasher, mockClaimer, router)
			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(http.MethodGet, test.route, nil)
			test.setupExpectations(*mockClaimer, request, responseRecorder)

			// Act
			router.ServeHTTP(responseRecorder, request)

			result := responseRecorder.Result()
			defer result.Body.Close()

			// Assert
			assert.Equal(t, test.responseCode, result.StatusCode)

			responseBody, _ := io.ReadAll(result.Body)
			assert.Equal(t, string(responseBody), test.expectedBody)
		})
	}
}
