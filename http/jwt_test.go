package http

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	db "github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
)

// func TestJWTRefreshHandler(t *testing.T) {
// 	tests := []struct {
// 		body              *bytes.Buffer
// 		method            string
// 		name              string
// 		responseCode      int
// 		route             string
// 		setupExpectations func()
// 	}{
// 		{
// 			body:         bytes.NewBufferString(""),
// 			method:       http.MethodPost,
// 			name:         "refresh returns refresh_token cookie given valid claims for existing token",
// 			responseCode: http.StatusAccepted,
// 			route:        "/refresh/",
// 			createClaimsToken: func(username string) string {
// 				expirationTime := time.Now().Add(5 * time.Minute)

// 				// Create the JWT claims, which includes the username and expiry time
// 				claims := &Claims{
// 					Username: req.Username,
// 					StandardClaims: jwt.StandardClaims{
// 						// In JWT, the expiry time is expressed as unix milliseconds
// 						ExpiresAt: expirationTime.Unix(),
// 					},
// 				}

// 				// Declare the token with the algorithm used for signing, and the claims
// 				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 				// Create the JWT string
// 				tokenString, err := token.SignedString(jwtKey)
// 			},
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			// Arrange
// 			router := gin.Default()
// 			gin.SetMode(gin.TestMode)

// 			ctrl := gomock.NewController(t)

// 			mockStore := db.NewMockStore(ctrl)
// 			mockHasher := security.NewMockHasher(ctrl)
// 			test.setupExpectations()

// 			NewFirstlyServer(mockStore, mockHasher, router)
// 			responseRecorder := httptest.NewRecorder()

// 			// Act
// 			request := httptest.NewRequest(test.method, test.route, test.body)
// 			// set the cookie on the request?
// 			if len(test.cookie) > 0 {
// 				request.AddCookie(&http.Cookie{Name: "token", Value: test.cookie})
// 				// Copy the Cookie over to a new Request
// 				// request = &http.Request{Header: http.Header{"Cookie": responseRecorder.HeaderMap["Set-Cookie"]}}
// 			}
// 			router.ServeHTTP(responseRecorder, request)

// 			result := responseRecorder.Result()

// 			fmt.Println(result.Header.Get("Cookie"))
// 			defer result.Body.Close()

// 			// Assert
// 			assert.Equal(t, test.responseCode, result.StatusCode)
// 		})
// 	}
// }

func TestJWTSignInHandler(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		method            string
		name              string
		responseCode      int
		route             string
		setupExpectations func(store *db.MockStore, hasher *security.MockHasher, rr *httptest.ResponseRecorder, r *http.Request)
	}{
		{
			body:         bytes.NewBufferString("{\"phrase\":\"blah\",\"invalid\":\"test\"}"),
			method:       http.MethodPost,
			name:         "signin returns status code bad request when invalid json supplied",
			responseCode: http.StatusBadRequest,
			route:        "/signin/",
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, rr *httptest.ResponseRecorder, r *http.Request) {
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
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, rr *httptest.ResponseRecorder, r *http.Request) {
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
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, rr *httptest.ResponseRecorder, r *http.Request) {
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
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, rr *httptest.ResponseRecorder, r *http.Request) {
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
			setupExpectations: func(store *db.MockStore, hasher *security.MockHasher, rr *httptest.ResponseRecorder, r *http.Request) {
				expectedAccount := db.Account{
					Phrase: []byte("valid"),
					Salt:   "salt",
				}
				store.EXPECT().GetAccountByUsername(gomock.Any(), "valid").Return(expectedAccount, nil)
				hasher.EXPECT().IsValidPassword(expectedAccount.Phrase, expectedAccount.Salt, "valid").Return(true, nil)
				// Create the JWT claims, which includes the username and expiry time

				expirationTime := time.Now().Add(5 * time.Minute)

				claims := &Claims{
					Username: "valid",
					StandardClaims: jwt.StandardClaims{
						// In JWT, the expiry time is expressed as unix milliseconds
						ExpiresAt: expirationTime.Unix(),
					},
				}

				// Declare the token with the algorithm used for signing, and the claims
				token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				// Create the JWT string
				tokenString, _ := token.SignedString(jwtKey)

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

			NewFirstlyServer(mockStore, mockHasher, router)
			responseRecorder := httptest.NewRecorder()

			request := httptest.NewRequest(test.method, test.route, test.body)
			test.setupExpectations(mockStore, mockHasher, responseRecorder, request)

			// Act
			router.ServeHTTP(responseRecorder, request)

			result := responseRecorder.Result()
			defer result.Body.Close()

			// Assert
			assert.Equal(t, test.responseCode, result.StatusCode)
		})
	}
}
