package http

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	db "github.com/meads/firstly-api/db/sqlc"
	"github.com/meads/firstly-api/security"
	"go.uber.org/mock/gomock"
)

func passClaimsMiddleware(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
	tokenString := "mocktoken"
	r.Header.Add("Authorization", "Bearer mocktoken")

	userClaims := &security.UserClaims{Type: "access"}
	tokener.EXPECT().VerifyToken(tokenString).Return(userClaims, nil)
}

func TestRegisterUserHandler_Post(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              registerUserResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier)
	}{
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/register/",
			want: registerUserResponse{
				SessionID:             "uuid",
				AccessToken:           "mockaccesstoken",
				RefreshToken:          "mockrefreshtoken",
				AccessTokenExpiresAt:  time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
				RefreshTokenExpiresAt: time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC),
				Username:              "newuser",
				UserID:                1,
			},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				requestUsername, requestPassword := "newuser", "message"
				querier.EXPECT().UsernameExists(gomock.Any(), requestUsername).Return(false, nil)

				hashedPassword := "generated_hash"
				hasher.EXPECT().HashPassword(requestPassword).Return(hashedPassword, nil)

				dbUser := db.User{ID: 1, Username: requestUsername}
				querier.EXPECT().CreateUser(gomock.Any(), db.CreateUserParams{
					Username: requestUsername, Password: hashedPassword},
				).Return(dbUser, nil)

				accessTokenClaims := &security.UserClaims{
					ID:       123,
					Username: requestUsername,
					Type:     "access",
					RegisteredClaims: &jwt.RegisteredClaims{
						ID:        "uuidstring",
						Subject:   requestUsername,
						IssuedAt:  jwt.NewNumericDate(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)),
						ExpiresAt: jwt.NewNumericDate(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)),
					},
				}
				accessTokenString := "mockaccesstoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "access", 15*time.Minute).
					Return(accessTokenString, accessTokenClaims, nil)

				refreshTokenClaims := &security.UserClaims{
					ID:       1234,
					Username: requestUsername,
					Type:     "refresh",
					RegisteredClaims: &jwt.RegisteredClaims{
						ID:        "uuidstring",
						Subject:   requestUsername,
						IssuedAt:  jwt.NewNumericDate(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)),
						ExpiresAt: jwt.NewNumericDate(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)),
					},
				}

				refreshTokenString := "mockrefreshtoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour).
					Return(refreshTokenString, refreshTokenClaims, nil)

				querier.EXPECT().CreateSession(gomock.Any(), db.CreateSessionParams{
					ID:           refreshTokenClaims.RegisteredClaims.ID,
					UserID:       dbUser.ID,
					RefreshToken: refreshTokenString,
					IsRevoked:    false,
					ExpiresAt:    refreshTokenClaims.RegisteredClaims.ExpiresAt.Time,
				}).Return(db.Session{ID: "uuid"}, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 500 when error generating refresh token",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			want:         registerUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				requestUsername, requestPassword := "newuser", "message"
				querier.EXPECT().UsernameExists(gomock.Any(), requestUsername).Return(false, nil)

				hashedPassword := "generated_hash"
				hasher.EXPECT().HashPassword(requestPassword).Return(hashedPassword, nil)

				dbUser := db.User{ID: 1, Username: requestUsername}
				querier.EXPECT().CreateUser(gomock.Any(),
					db.CreateUserParams{Username: requestUsername, Password: hashedPassword},
				).Return(dbUser, nil)

				accessTokenClaims := &security.UserClaims{
					ID:       dbUser.ID,
					Username: dbUser.Username,
					Type:     "access",
					RegisteredClaims: &jwt.RegisteredClaims{
						ID:        "uuidstring",
						Subject:   dbUser.Username,
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
					},
				}
				accessTokenString := "mockaccesstoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "access", 15*time.Minute).
					Return(accessTokenString, accessTokenClaims, nil)

				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour).
					Return("", nil, errors.New("error generating refresh token"))

				querier.EXPECT().CreateSession(gomock.Any(), gomock.Any()).Times(0)
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 500 when error creating session",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			want:         registerUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				requestUsername, requestPassword := "newuser", "message"
				querier.EXPECT().UsernameExists(gomock.Any(), requestUsername).Return(false, nil)

				hashedPassword := "generated_hash"
				hasher.EXPECT().HashPassword(requestPassword).Return(hashedPassword, nil)

				dbUser := db.User{ID: 1, Username: requestUsername}
				querier.EXPECT().CreateUser(
					gomock.Any(), db.CreateUserParams{Username: requestUsername, Password: hashedPassword},
				).Return(dbUser, nil)

				accessTokenClaims := &security.UserClaims{
					ID:       dbUser.ID,
					Username: dbUser.Username,
					Type:     "access",
					RegisteredClaims: &jwt.RegisteredClaims{
						ID:        "uuidstring",
						Subject:   dbUser.Username,
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
					},
				}
				accessTokenString := "mockaccesstoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "access", 15*time.Minute).
					Return(accessTokenString, accessTokenClaims, nil)

				refreshTokenClaims := &security.UserClaims{
					ID:       1234,
					Username: requestUsername,
					Type:     "refresh",
					RegisteredClaims: &jwt.RegisteredClaims{
						ID:        "uuidstring",
						Subject:   requestUsername,
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
					},
				}
				refreshTokenString := "mockrefreshtoken"
				tokener.EXPECT().GenerateToken(dbUser.ID, dbUser.Username, "refresh", 24*time.Hour).
					Return(refreshTokenString, refreshTokenClaims, nil)

				querier.EXPECT().CreateSession(gomock.Any(), db.CreateSessionParams{
					ID:           refreshTokenClaims.RegisteredClaims.ID,
					UserID:       dbUser.ID,
					RefreshToken: refreshTokenString,
					IsRevoked:    false,
					ExpiresAt:    refreshTokenClaims.RegisteredClaims.ExpiresAt.Time,
				}).Return(db.Session{}, errors.New("error creating session"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"\",\"password\":\"\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 400 given invalid params are supplied",
			responseCode: http.StatusBadRequest,
			route:        "/register/",
			want:         registerUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 500 given there is an error hashing the password",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			want:         registerUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				requestUsername, requestPassword := "newuser", "message"
				querier.EXPECT().UsernameExists(gomock.Any(), requestUsername).Return(false, nil)
				hasher.EXPECT().HashPassword(requestPassword).Return("", errors.New("error hashing password"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 500 given there is some server error with user exists",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			want:         registerUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				querier.EXPECT().UsernameExists(gomock.Any(), "newuser").Return(false, errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 500 given there is some server error before create",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			want:         registerUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				requestUsername, requestPassword := "newuser", "message"
				querier.EXPECT().UsernameExists(gomock.Any(), requestUsername).Return(false, nil)
				hasher.EXPECT().HashPassword(requestPassword).Return("generated_hash", nil)
				querier.EXPECT().CreateUser(
					gomock.Any(),
					db.CreateUserParams{Username: requestUsername, Password: "generated_hash"}).
					Return(db.User{}, errors.New("oops"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"existinguser\",\"password\":\"valid\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 400 given a user already exists with username x",
			responseCode: http.StatusBadRequest,
			route:        "/register/",
			want:         registerUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				querier.EXPECT().UsernameExists(gomock.Any(), "existinguser").Return(true, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "register handler responds with Status Code 500 when generate token returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/register/",
			want:         registerUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				requestUsername, requestPassword := "newuser", "message"
				querier.EXPECT().UsernameExists(gomock.Any(), requestUsername).Return(false, nil)
				hashedPassword := "generated_hash"
				hasher.EXPECT().HashPassword(requestPassword).Return(hashedPassword, nil)
				user := db.User{ID: 1, Username: "newuser"}
				querier.EXPECT().CreateUser(
					gomock.Any(), db.CreateUserParams{Username: requestUsername, Password: hashedPassword},
				).Return(user, nil)
				tokenString := "mocktoken"
				tokener.EXPECT().GenerateToken(user.ID, user.Username, "access", 15*time.Minute).
					Return(tokenString, &security.UserClaims{}, errors.New("token create error"))
			},
		},
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

			// Act
			request := httptest.NewRequest(http.MethodPost, test.route, test.body)
			test.setupExpectations(request, mockTokener, mockHasher, mockQuerier)
			router.ServeHTTP(responseRecorder, request)

			// result := responseRecorder.Result()
			// defer result.Body.Close()

			// Assert that the response recorder http.Response Status code matches the test expectations
			// assert.Equal(t, test.responseCode, result.StatusCode)
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got registerUserResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if got != test.want {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              DeleteUserResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier)
	}{
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "delete handler responds with Status Code 200 given valid request",
			responseCode: http.StatusOK,
			route:        "/users/1/",
			want:         DeleteUserResponse{Message: "Resource successfully deleted"},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				querier.EXPECT().DeleteUser(gomock.Any(), int64(1)).Return(nil)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "delete handler responds with Status Code 400 given param id not supplied",
			responseCode: http.StatusBadRequest,
			route:        "/users//",
			want:         DeleteUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "delete handler responds with Status Code 400 given param id is not a valid integer",
			responseCode: http.StatusBadRequest,
			route:        "/users/invalid/",
			want:         DeleteUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "delete handler responds with Status Code 500 given there is a server error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/",
			want:         DeleteUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				querier.EXPECT().DeleteUser(gomock.Any(), int64(1)).Return(errors.New("oops"))
			},
		},
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

			// Act
			request := httptest.NewRequest(http.MethodDelete, test.route, test.body)
			test.setupExpectations(request, mockTokener, mockHasher, mockQuerier)
			router.ServeHTTP(responseRecorder, request)

			// result := responseRecorder.Result()
			// defer result.Body.Close()

			// Assert that the response recorder http.Response Status code matches the test expectations
			// assert.Equal(t, test.responseCode, result.StatusCode)
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got DeleteUserResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if got != test.want {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}

func TestUserHandler_Get(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              ListUsersResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier)
	}{
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "list handler responds with Status Code 400 given limit param is invalid int",
			responseCode: http.StatusBadRequest,
			route:        "/users/?limit=invalid",
			want:         ListUsersResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "list handler responds with Status Code 400 given offset param is invalid int",
			responseCode: http.StatusBadRequest,
			route:        "/users/?offset=invalid",
			want:         ListUsersResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "list handler responds with Status Code 500 given there is a server error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			want:         ListUsersResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				params := db.ListUsersParams{Limit: 50, Offset: 0}
				querier.EXPECT().ListUsers(gomock.Any(), params).Return([]db.User{}, errors.New("oops."))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "list handler responds with Status Code 200 given a valid request",
			responseCode: http.StatusOK,
			route:        "/users/",
			want:         ListUsersResponse{Users: []db.User{{ID: 69, Username: "foo", CreatedAt: sql.NullTime{}}}},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				params := db.ListUsersParams{Limit: 50, Offset: 0}
				querier.EXPECT().ListUsers(gomock.Any(), params).Return([]db.User{
					{ID: 69, Username: "foo", CreatedAt: sql.NullTime{}},
				}, nil)
			},
		},
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

			// Act
			request := httptest.NewRequest(http.MethodGet, test.route, test.body)
			test.setupExpectations(request, mockTokener, mockHasher, mockQuerier)
			router.ServeHTTP(responseRecorder, request)

			// result := responseRecorder.Result()
			// defer result.Body.Close()

			// Assert that the response recorder http.Response Status code matches the test expectations
			// assert.Equal(t, test.responseCode, result.StatusCode)
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			got := ListUsersResponse{}
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				// t.Fatalf("Falied to decode JSON response: '%+v'", responseRecorder.Body.String())
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}

func TestUserHandler_Patch(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              PatchUserResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier)
	}{
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"new\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/users/",
			want:         PatchUserResponse{Message: "Resource successfully patched"},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				reqID, reqUsername, reqCurrentPassword, reqNewPassword := int64(1), "user", "current", "new"
				user := db.User{
					ID:       1,
					Username: reqUsername,
					Password: "currenthashed",
				}
				querier.EXPECT().GetUser(gomock.Any(), reqID).Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, reqCurrentPassword).Return(nil)
				hasher.EXPECT().HashPassword(reqNewPassword).Return("newhashed", nil)
				params := db.UpdateUserPasswordParams{ID: int64(1), Password: "newhashed"}
				querier.EXPECT().UpdateUserPassword(gomock.Any(), params).Return(nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":69,\"username\":\"user\",\"wrong\":\"newpass\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update handler responds with Status Code 400 when invalid data supplied",
			responseCode: http.StatusBadRequest,
			route:        "/users/",
			want:         PatchUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"newpass\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update handler responds with Status Code 404 when record not found",
			responseCode: http.StatusNotFound,
			route:        "/users/",
			want:         PatchUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				reqID := int64(1)
				querier.EXPECT().GetUser(gomock.Any(), reqID).Return(db.User{}, sql.ErrNoRows)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"newpass\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update handler responds with Status Code 500 when server error on get before update",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			want:         PatchUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				reqID := int64(1)
				querier.EXPECT().GetUser(gomock.Any(), reqID).Return(db.User{}, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"new\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update handler responds with Status Code 500 when server error on update",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			want:         PatchUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				reqID, reqUsername, reqCurrentPassword, reqNewPassword := int64(1), "user", "current", "new"
				user := db.User{
					ID:       1,
					Username: reqUsername,
					Password: "currenthashed",
				}
				querier.EXPECT().GetUser(gomock.Any(), reqID).Return(user, nil)
				hasher.EXPECT().ComparePassword(user.Password, reqCurrentPassword).Return(nil)
				hasher.EXPECT().HashPassword(reqNewPassword).Return("newhashed", nil)
				params := db.UpdateUserPasswordParams{ID: int64(1), Password: "newhashed"}
				querier.EXPECT().UpdateUserPassword(gomock.Any(), params).Return(errors.New("db server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"username\":\"user\",\"currentPassword\":\"current\",\"newPassword\":\"new\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "update handler responds with Status Code 500 when server error on get user",
			responseCode: http.StatusInternalServerError,
			route:        "/users/",
			want:         PatchUserResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				querier.EXPECT().GetUser(gomock.Any(), int64(1)).Return(db.User{}, errors.New("connection refused"))
			},
		},
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

			// Act
			request := httptest.NewRequest(http.MethodPatch, test.route, test.body)
			test.setupExpectations(request, mockTokener, mockHasher, mockQuerier)
			router.ServeHTTP(responseRecorder, request)

			// result := responseRecorder.Result()
			// defer result.Body.Close()

			// Assert that the response recorder http.Response Status code matches the test expectations
			// assert.Equal(t, test.responseCode, result.StatusCode)
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got PatchUserResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if got != test.want {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}
