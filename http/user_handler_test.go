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
	db "github.com/meads/firstly-api/db"
	"github.com/meads/firstly-api/security"
	"go.uber.org/mock/gomock"
)

func passClaimsMiddleware(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
	tokenString := "mocktoken"
	// os.Setenv("SECRET", "test")
	r.Header.Add("Authorization", "Bearer mocktoken")

	userClaims := &security.UserClaims{Type: "access"}
	tokener.EXPECT().VerifyToken(tokenString).Return(userClaims, nil)
}

func TestUserHandler(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		method            string
		name              string
		responseCode      int
		route             string
		isList            bool
		setupExpectations func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier)
	}{
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			method:       http.MethodPost,
			name:         "register handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/register/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				req := registerUserRequest{Username: "newuser", Password: "message"}
				querier.EXPECT().UserExists(gomock.Any(), req.Username).Return(false, nil)

				hashedPassword := "generated_hash"

				hasher.EXPECT().HashPassword(req.Password).Return(hashedPassword, nil)

				dbUser := db.User{ID: 1, Username: "newuser"}

				querier.EXPECT().CreateUser(
					gomock.Any(), db.CreateUserParams{Username: req.Username, Password: hashedPassword},
				).Return(dbUser, nil)

				accessTokenClaims, _ := security.NewUserClaims(dbUser.ID, dbUser.Username, "access", 5*time.Minute)
				accessTokenString := "mockaccesstoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "access", 5*time.Minute).
					Return(accessTokenString, accessTokenClaims, nil)

				refreshTokenClaims, _ := security.NewUserClaims(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour)
				refreshTokenString := "mockrefreshtoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour).
					Return(refreshTokenString, refreshTokenClaims, nil)

				querier.EXPECT().CreateSession(gomock.Any(), db.CreateSessionParams{
					ID:           refreshTokenClaims.RegisteredClaims.ID,
					Username:     dbUser.Username,
					RefreshToken: refreshTokenString,
					IsRevoked:    false,
					ExpiresAt:    refreshTokenClaims.RegisteredClaims.ExpiresAt.Time,
				})
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			method:       http.MethodPost,
			name:         "register handler responds with Status Code 500 when error generating refresh token",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				req := registerUserRequest{Username: "newuser", Password: "message"}
				querier.EXPECT().UserExists(gomock.Any(), req.Username).Return(false, nil)

				hashedPassword := "generated_hash"
				hasher.EXPECT().HashPassword(req.Password).Return(hashedPassword, nil)

				dbUser := db.User{ID: 1, Username: "newuser"}
				querier.EXPECT().CreateUser(
					gomock.Any(), db.CreateUserParams{Username: req.Username, Password: hashedPassword},
				).Return(dbUser, nil)

				accessTokenClaims, _ := security.NewUserClaims(dbUser.ID, dbUser.Username, "access", 5*time.Minute)
				accessTokenString := "mockaccesstoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "access", 5*time.Minute).
					Return(accessTokenString, accessTokenClaims, nil)

				refreshTokenClaims, _ := security.NewUserClaims(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour)
				refreshTokenString := "mockrefreshtoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour).
					Return(refreshTokenString, refreshTokenClaims, errors.New("error generating refresh token"))

				querier.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			method:       http.MethodPost,
			name:         "register handler responds with Status Code 500 when error creating session",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				req := registerUserRequest{Username: "newuser", Password: "message"}
				querier.EXPECT().UserExists(gomock.Any(), req.Username).Return(false, nil)

				hashedPassword := "generated_hash"
				hasher.EXPECT().HashPassword(req.Password).Return(hashedPassword, nil)

				dbUser := db.User{ID: 1, Username: "newuser"}
				querier.EXPECT().CreateUser(
					gomock.Any(), db.CreateUserParams{Username: req.Username, Password: hashedPassword},
				).Return(dbUser, nil)

				accessTokenClaims, _ := security.NewUserClaims(dbUser.ID, dbUser.Username, "access", 5*time.Minute)
				accessTokenString := "mockaccesstoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "access", 5*time.Minute).
					Return(accessTokenString, accessTokenClaims, nil)

				refreshTokenClaims, _ := security.NewUserClaims(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour)
				refreshTokenString := "mockrefreshtoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour).
					Return(refreshTokenString, refreshTokenClaims, nil)

				querier.EXPECT().CreateSession(gomock.Any(), db.CreateSessionParams{
					ID:           refreshTokenClaims.RegisteredClaims.ID,
					Username:     dbUser.Username,
					RefreshToken: refreshTokenString,
					IsRevoked:    false,
					ExpiresAt:    refreshTokenClaims.RegisteredClaims.ExpiresAt.Time,
				}).Return(db.Session{}, errors.New("error creating session"))
			},
		},
		{
			name:         "register handler responds with Status Code 400 given invalid params are supplied",
			body:         bytes.NewBufferString("{\"username\":\"\",\"password\":\"\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusBadRequest,
			route:        "/register/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
			},
		},
		{
			name:         "register handler responds with Status Code 500 given there is an error hashing the password",
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				req := registerUserRequest{Username: "newuser", Password: "message"}
				querier.EXPECT().UserExists(gomock.Any(), req.Username).Return(false, nil)
				hasher.EXPECT().HashPassword(req.Password).Return("", errors.New("error hashing password"))
			},
		},
		{
			name:         "register handler responds with Status Code 500 given there is some server error with user exists",
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				req := registerUserRequest{Username: "newuser", Password: "message"}
				querier.EXPECT().UserExists(gomock.Any(), req.Username).Return(false, errors.New("oops"))
			},
		},
		{
			name:         "register handler responds with Status Code 500 given there is some server error before create",
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				req := registerUserRequest{Username: "newuser", Password: "message"}
				querier.EXPECT().UserExists(gomock.Any(), req.Username).Return(false, nil)
				hasher.EXPECT().HashPassword("message").Return("generated_hash", nil)
				querier.EXPECT().CreateUser(
					gomock.Any(),
					db.CreateUserParams{Username: "newuser", Password: "generated_hash"}).
					Return(db.User{}, errors.New("oops"))
			},
		},
		{
			name:         "register handler responds with Status Code 400 given a user already exists with username x",
			body:         bytes.NewBufferString("{\"username\":\"existinguser\",\"password\":\"valid\"}"),
			method:       http.MethodPost,
			responseCode: http.StatusBadRequest,
			route:        "/register/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				req := registerUserRequest{Username: "existinguser", Password: "valid"}
				querier.EXPECT().UserExists(gomock.Any(), req.Username).Return(true, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			method:       http.MethodPost,
			name:         "register handler responds with Status Code 500 when get five minute expiration token returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				req := registerUserRequest{Username: "newuser", Password: "message"}
				querier.EXPECT().UserExists(gomock.Any(), req.Username).Return(false, nil)
				hashedPassword := "generated_hash"
				hasher.EXPECT().HashPassword(req.Password).Return(hashedPassword, nil)
				user := db.User{ID: 1, Username: "newuser"}
				querier.EXPECT().CreateUser(
					gomock.Any(), db.CreateUserParams{Username: req.Username, Password: hashedPassword},
				).Return(user, nil)
				tokenString := "mocktoken"
				tokener.EXPECT().GenerateToken(user.ID, user.Username, "access", 5*time.Minute).
					Return(tokenString, &security.UserClaims{}, errors.New("token create error"))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 200 given valid request",
			method:       http.MethodDelete,
			responseCode: http.StatusOK,
			route:        "/users/1/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				querier.EXPECT().DeleteUser(gomock.Any(), int64(1)).Return(nil)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 400 given param id not supplied",
			method:       http.MethodDelete,
			responseCode: http.StatusBadRequest,
			route:        "/users//",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 400 given param id is not a valid integer",
			method:       http.MethodDelete,
			responseCode: http.StatusBadRequest,
			route:        "/users/invalid/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "delete handler responds with Status Code 500 given there is a server error",
			method:       http.MethodDelete,
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				querier.EXPECT().DeleteUser(gomock.Any(), int64(1)).Return(errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 400 given limit param is invalid int",
			method:       http.MethodGet,
			responseCode: http.StatusBadRequest,
			route:        "/users/?limit=invalid",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 400 given offset param is invalid int",
			method:       http.MethodGet,
			responseCode: http.StatusBadRequest,
			route:        "/users/?offset=invalid",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 500 given there is a server error",
			method:       http.MethodGet,
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				params := db.ListUsersParams{Limit: 50, Offset: 0}
				querier.EXPECT().ListUsers(gomock.Any(), params).Return([]db.User{}, errors.New("oops."))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			name:         "list handler responds with Status Code 200 given a valid request",
			method:       http.MethodGet,
			responseCode: http.StatusOK,
			route:        "/users/",
			isList:       true,
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				params := db.ListUsersParams{Limit: 50, Offset: 0}
				querier.EXPECT().ListUsers(gomock.Any(), params).Return([]db.User{
					{ID: 69, Username: "foo", CreatedAt: sql.NullTime{}},
				}, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"new\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/users/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				req := updateUserRequest{
					ID:              1,
					Username:        "user",
					CurrentPassword: "current",
					NewPassword:     "new",
				}
				user := db.User{
					ID:       1,
					Username: "user",
					Password: "currenthashed",
				}
				hasher.EXPECT().HashPassword(req.CurrentPassword).Return("currenthashed", nil)
				querier.EXPECT().GetUser(gomock.Any(), req.ID).Return(user, nil)
				hasher.EXPECT().HashPassword(req.NewPassword).Return("newhashed", nil)
				params := db.UpdateUserParams{ID: int64(1), Password: "newhashed"}
				querier.EXPECT().UpdateUser(gomock.Any(), params).Return(nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69,\"username\":\"user\",\"wrong\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 400 when invalid data supplied",
			responseCode: http.StatusBadRequest,
			route:        "/users/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 404 when record not found",
			responseCode: http.StatusNotFound,
			route:        "/users/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				req := updateUserRequest{
					ID:              1,
					Username:        "user",
					CurrentPassword: "current",
					NewPassword:     "newpass",
				}
				hashedPassword := "currenthashed"
				hasher.EXPECT().HashPassword(req.CurrentPassword).Return(hashedPassword, nil)
				querier.EXPECT().GetUser(gomock.Any(), req.ID).Return(db.User{}, sql.ErrNoRows)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"newpass\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on get before update",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				req := updateUserRequest{
					ID:              1,
					Username:        "user",
					CurrentPassword: "current",
					NewPassword:     "newpass",
				}
				hashedPassword := "currenthashed"
				hasher.EXPECT().HashPassword(req.CurrentPassword).Return(hashedPassword, nil)
				querier.EXPECT().GetUser(gomock.Any(), req.ID).Return(db.User{}, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"new\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on hashing request current password",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				req := updateUserRequest{
					ID:              1,
					Username:        "user",
					CurrentPassword: "current",
					NewPassword:     "new",
				}
				hasher.EXPECT().HashPassword(req.CurrentPassword).Return("", errors.New("error hashing current password"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"new\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on update",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				req := updateUserRequest{
					ID:              1,
					Username:        "user",
					CurrentPassword: "current",
					NewPassword:     "new",
				}
				user := db.User{
					ID:       1,
					Username: "user",
					Password: "currenthashed",
				}
				hasher.EXPECT().HashPassword(req.CurrentPassword).Return("currenthashed", nil)
				querier.EXPECT().GetUser(gomock.Any(), req.ID).Return(user, nil)
				hasher.EXPECT().HashPassword(req.NewPassword).Return("newhashed", nil)
				params := db.UpdateUserParams{ID: int64(1), Password: "newhashed"}
				querier.EXPECT().UpdateUser(gomock.Any(), params).Return(errors.New("db server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"new\"}"),
			method:       http.MethodPatch,
			name:         "update handler responds with Status Code 500 when server error on get user",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				req := updateUserRequest{
					ID:              1,
					Username:        "user",
					CurrentPassword: "current",
					NewPassword:     "new",
				}
				hasher.EXPECT().HashPassword(req.CurrentPassword).Return("currenthashed", nil)
				querier.EXPECT().GetUser(gomock.Any(), req.ID).Return(db.User{}, errors.New("connection refused"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			router := gin.Default()
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockQuerier := db.NewMockQuerier(ctrl)
			mockHasher := security.NewMockHasher(ctrl)
			mockTokener := security.NewMockTokener(ctrl)

			NewFirstlyServer(mockTokener, mockHasher, router, mockQuerier)
			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(test.method, test.route, test.body)
			test.setupExpectations(request, mockTokener, mockHasher, mockQuerier)
			router.ServeHTTP(responseRecorder, request)

			result := responseRecorder.Result()
			defer result.Body.Close()

			// Assert
			assert.Equal(t, test.responseCode, result.StatusCode)

			if !test.isList {
				response := db.User{}

				if result.Body != http.NoBody {
					if err := json.NewDecoder(result.Body).Decode(&response); err != nil && !errors.Is(err, io.EOF) {
						t.Errorf("Error decoding response body: %v", err)
						t.Log()
						t.Log(responseRecorder.Body)
						t.Log()
					}
				}
			} else {
				response := []db.User{}
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
