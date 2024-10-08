package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	"github.com/golang/mock/gomock"
	"github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
)

func TestAccountHandler(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		method            string
		name              string
		responseCode      int
		route             string
		isList            bool
		setupExpectations func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore)
	}{
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"phrase\":\"message\"}"),
			method:       http.MethodPost,
			name:         "create handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				store.EXPECT().GetAccountByUsername(gomock.Any(), "newuser").Return(db.Account{Deleted: true}, nil)
				os.Setenv("SECRET", "test")
				hash := []byte("generated_hash")
				hasher.EXPECT().GenerateSalt().Return("salt").Times(1)
				hasher.EXPECT().GeneratePasswordHash([]byte("message"), "salt").Return(hash, nil)
				store.EXPECT().CreateAccount(
					gomock.Any(),
					db.CreateAccountParams{Username: "newuser", Phrase: hash, Salt: "salt"},
				).Return(db.Account{Username: "newuser"}, nil)
				tokenString := "mocktoken"
				expirationTime := time.Now().Add(5 * time.Minute)
				claimer.EXPECT().GetFiveMinuteExpirationToken("newuser").Return(tokenString, expirationTime, nil)

			},
		},
		{
			name:         "create handler responds with Status Code 400 given invalid params are supplied",
			body:         bytes.NewBufferString("{\"username\":\"\",\"phrase\":\"\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusBadRequest,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
			},
		},
		{
			name:         "create handler responds with Status Code 500 given there is an error hashing the phrase, no env var set",
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"phrase\":\"message\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusInternalServerError,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				store.EXPECT().GetAccountByUsername(gomock.Any(), "newuser").Return(db.Account{ID: 0}, nil)
				hasher.EXPECT().GenerateSalt().Return("somesalt")
				os.Unsetenv("SECRET")
				hasher.EXPECT().GeneratePasswordHash(gomock.Any(), "somesalt").Return(
					nil,
					errors.New("secret not set"),
				)
			},
		},
		{
			name:         "create handler responds with Status Code 500 given there is some server error with get account",
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"phrase\":\"message\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusInternalServerError,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				store.EXPECT().GetAccountByUsername(gomock.Any(), "newuser").Return(db.Account{}, errors.New("oops"))
			},
		},
		{
			name:         "create handler responds with Status Code 500 given there is some server error before create",
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"phrase\":\"message\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusInternalServerError,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				store.EXPECT().GetAccountByUsername(gomock.Any(), "newuser").Return(db.Account{ID: 0}, nil)
				os.Setenv("SECRET", "test")
				hash := []byte("generated_hash")
				hasher.EXPECT().GenerateSalt().Return("salt").Times(1)
				hasher.EXPECT().GeneratePasswordHash([]byte("message"), "salt").Return(hash, nil)
				store.EXPECT().CreateAccount(
					gomock.Any(),
					db.CreateAccountParams{Username: "newuser", Phrase: []byte("generated_hash"), Salt: "salt"}).
					Return(db.Account{}, errors.New("oops"))
			},
		},
		{
			name:         "create handler responds with Status Code 400 given a user already exists with username x",
			body:         bytes.NewBufferString("{\"username\":\"invalid\",\"phrase\":\"valid\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusBadRequest,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				store.EXPECT().GetAccountByUsername(gomock.Any(), "invalid").
					Return(db.Account{ID: 1}, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"phrase\":\"message\"}"),
			method:       http.MethodPost,
			name:         "create handler responds with Status Code 500 when get five minute expiration token returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				store.EXPECT().GetAccountByUsername(gomock.Any(), "newuser").Return(db.Account{Deleted: true}, nil)
				os.Setenv("SECRET", "test")
				hash := []byte("generated_hash")
				hasher.EXPECT().GenerateSalt().Return("salt").Times(1)
				hasher.EXPECT().GeneratePasswordHash([]byte("message"), "salt").Return(hash, nil)
				store.EXPECT().CreateAccount(
					gomock.Any(),
					db.CreateAccountParams{Username: "newuser", Phrase: hash, Salt: "salt"},
				).Return(db.Account{Username: "newuser"}, nil)
				tokenString := ""
				expirationTime := time.Time{}
				claimer.EXPECT().GetFiveMinuteExpirationToken("newuser").Return(tokenString, expirationTime, errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 200 given valid request",
			method:       http.MethodDelete,
			responseCode: http.StatusOK,
			route:        "/account/69/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				store.EXPECT().DeleteAccount(gomock.Any(), int64(69)).Return(nil)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 400 given param id not supplied",
			method:       http.MethodDelete,
			responseCode: http.StatusBadRequest,
			route:        "/account//",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 400 given param id is not a valid integer",
			method:       http.MethodDelete,
			responseCode: http.StatusBadRequest,
			route:        "/account/invalid/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 500 given there is a server error",
			method:       http.MethodDelete,
			responseCode: http.StatusInternalServerError,
			route:        "/account/69/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				store.EXPECT().DeleteAccount(gomock.Any(), int64(69)).Return(errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 400 given limit param is invalid int",
			method:       http.MethodGet,
			responseCode: http.StatusBadRequest,
			route:        "/account/?limit=invalid",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 400 given offset param is invalid int",
			method:       http.MethodGet,
			responseCode: http.StatusBadRequest,
			route:        "/account/?offset=invalid",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 500 given there is a server error",
			method:       http.MethodGet,
			responseCode: http.StatusInternalServerError,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				params := db.ListAccountsParams{Limit: 50, Offset: 0}
				store.EXPECT().ListAccounts(gomock.Any(), params).Return([]db.ListAccountsRow{}, errors.New("oops."))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 200 given a valid request",
			method:       http.MethodGet,
			responseCode: http.StatusOK,
			route:        "/account/",
			isList:       true,
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				params := db.ListAccountsParams{Limit: 50, Offset: 0}
				store.EXPECT().ListAccounts(gomock.Any(), params).Return([]db.ListAccountsRow{
					{ID: 69, Username: "foo", Created: "", Deleted: false},
				}, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69,\"username\":\"user\",\"phrase\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				store.EXPECT().AccountExists(gomock.Any(), int64(69)).Return(true, nil)
				store.EXPECT().GetAccount(gomock.Any(), int64(69)).Return(db.Account{ID: 69, Phrase: []byte("newpass"), Salt: "salt the snail"}, nil)
				hasher.EXPECT().GenerateSalt().Times(0)
				hasher.EXPECT().GeneratePasswordHash([]byte("newpass"), "salt the snail").Return([]byte("newhash"), nil)
				params := db.UpdateAccountParams{ID: int64(69), Phrase: []byte("newhash")}
				store.EXPECT().UpdateAccount(gomock.Any(), params).Return(nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69,\"username\":\"user\",\"wrong\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 400 when invalid data supplied",
			responseCode: http.StatusBadRequest,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":68,\"username\":\"user\",\"phrase\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 404 when record not found",
			responseCode: http.StatusNotFound,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				store.EXPECT().AccountExists(gomock.Any(), int64(68)).Return(false, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69,\"username\":\"user\",\"phrase\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on get before update",
			responseCode: http.StatusInternalServerError,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				store.EXPECT().AccountExists(gomock.Any(), int64(69)).Return(true, nil)
				store.EXPECT().GetAccount(gomock.Any(), int64(69)).
					Return(db.Account{}, errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69,\"username\":\"user\",\"phrase\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on update hash",
			responseCode: http.StatusInternalServerError,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				store.EXPECT().AccountExists(gomock.Any(), int64(69)).Return(true, nil)
				store.EXPECT().GetAccount(gomock.Any(), int64(69)).Return(
					db.Account{ID: 69, Salt: "somesalt", Phrase: []byte("newpass")}, nil)
				hasher.EXPECT().GenerateSalt().Times(0)
				os.Unsetenv("SECRET")
				hasher.EXPECT().GeneratePasswordHash(gomock.Any(), "somesalt").Return(
					nil,
					errors.New("secret not set"),
				)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69,\"username\":\"user\",\"phrase\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on update",
			responseCode: http.StatusInternalServerError,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				store.EXPECT().AccountExists(gomock.Any(), int64(69)).Return(true, nil)
				store.EXPECT().GetAccount(gomock.Any(), int64(69)).Return(db.Account{ID: 69, Phrase: []byte("newpass"), Salt: "salt the snail"}, nil)
				hasher.EXPECT().GenerateSalt().Times(0)
				hasher.EXPECT().GeneratePasswordHash([]byte("newpass"), "salt the snail").Return([]byte("newhash"), nil)
				params := db.UpdateAccountParams{ID: 69, Phrase: []byte("newhash")}
				store.EXPECT().UpdateAccount(gomock.Any(), params).Return(errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69,\"username\":\"user\",\"phrase\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on update account exists",
			responseCode: http.StatusInternalServerError,
			route:        "/account/",
			setupExpectations: func(r *http.Request, claimer *security.MockClaimer, hasher *security.MockHasher, store *db.MockStore) {
				passClaimsMiddleware(r, claimer, hasher, store)
				store.EXPECT().AccountExists(gomock.Any(), int64(69)).Return(false, errors.New("oops"))
				store.EXPECT().GetAccount(gomock.Any(), gomock.Any()).Times(0)
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

			NewFirstlyServer(mockClaimer, mockHasher, router, mockStore)
			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(test.method, test.route, test.body)
			test.setupExpectations(request, mockClaimer, mockHasher, mockStore)
			router.ServeHTTP(responseRecorder, request)

			result := responseRecorder.Result()
			defer result.Body.Close()

			// Assert
			assert.Equal(t, test.responseCode, result.StatusCode)

			if !test.isList {
				response := db.Account{}

				if result.Body != http.NoBody {
					if err := json.NewDecoder(result.Body).Decode(&response); err != nil && !errors.Is(err, io.EOF) {
						t.Errorf("Error decoding response body: %v", err)
						t.Log()
						t.Log(responseRecorder.Body)
						t.Log()
					}
				}
			} else {
				response := []db.Account{}
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
