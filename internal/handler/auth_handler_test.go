package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/meads/firstly-api/internal/domain"
	"go.uber.org/mock/gomock"
)

func TestAuthLogin_Post(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              LoginResponse
		setupExpectations func(r *http.Request, authService *MockAuthServicer)
	}{
		{
			body:         bytes.NewBufferString("{\"username\":\"\",\"password\":\"\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "login handler responds with Status Code 400 given username and password are blank",
			responseCode: http.StatusBadRequest,
			route:        "/login/",
			want:         LoginResponse{},
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "login handler responds with Status Code 500 given login service returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/login/",
			want:         LoginResponse{},
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
				requestUsername, requestPassword := "newuser", "message"
				authService.EXPECT().Login(gomock.Any(), requestUsername, requestPassword).
					Return(&domain.LoginResult{}, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"username\":\"newuser\",\"password\":\"message\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "login handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/login/",
			want: LoginResponse{
				RefreshToken: "mockrefreshtoken",
				AccessToken:  "mockaccesstoken",
				SessionID:    "uuid",
				UserID:       1,
			},
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
				requestUsername, requestPassword := "newuser", "message"
				result := &domain.LoginResult{
					RefreshToken: "mockrefreshtoken",
					AccessToken:  "mockaccesstoken",
					SessionID:    "uuid",
					UserID:       1,
				}
				authService.EXPECT().Login(gomock.Any(), requestUsername, requestPassword).
					Return(result, nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockAuthService := NewMockAuthServicer(ctrl)
			authHandler := NewAuthHandler(mockAuthService)

			router := SetupRouter(authHandler, nil, nil, nil)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodPost, test.route, test.body)
			test.setupExpectations(request, mockAuthService)
			router.ServeHTTP(responseRecorder, request)

			// Assert

			// validate the response code
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			// validate the content type
			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got LoginResponse
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

func TestAuthLogout_Post(t *testing.T) {
	tests := []struct {
		contentType       string
		name              string
		responseCode      int
		route             string
		setupExpectations func(r *http.Request, authService *MockAuthServicer)
	}{
		{
			contentType:  "application/json; charset=utf-8",
			name:         "logout handler responds with Status Code 500 given auth service returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/logout/invalidsessionid",
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
				authService.EXPECT().Logout(gomock.Any(), "invalidsessionid").
					Return(errors.New("server error"))
			},
		},
		{
			contentType:  "",
			name:         "logout handler responds with Status Code 204 when valid data supplied",
			responseCode: http.StatusNoContent,
			route:        "/logout/sessionid",
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
				authService.EXPECT().Logout(gomock.Any(), "sessionid").
					Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockAuthService := NewMockAuthServicer(ctrl)
			authHandler := NewAuthHandler(mockAuthService)

			router := SetupRouter(authHandler, nil, nil, nil)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodPost, test.route, nil)
			test.setupExpectations(request, mockAuthService)
			router.ServeHTTP(responseRecorder, request)

			// Assert

			// validate the response code
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			// validate the content type
			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}
		})
	}
}

func TestRenewAccessToken_Post(t *testing.T) {
	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              RenewAccessTokenResponse
		setupExpectations func(r *http.Request, authService *MockAuthServicer)
	}{
		{
			body:         bytes.NewBufferString("{\"refreshToken\":\"\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "renew access token handler responds with Status Code 400 given refresh token is blank",
			responseCode: http.StatusBadRequest,
			route:        "/refresh/",
			want:         RenewAccessTokenResponse{},
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
			},
		},
		{
			body:         bytes.NewBufferString("{\"refreshToken\":\"token\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "renew access token handler responds with Status Code 200 given refresh token is valid",
			responseCode: http.StatusOK,
			route:        "/refresh/",
			want: RenewAccessTokenResponse{
				AccessToken: "newtoken",
			},
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
				authService.EXPECT().RenewAccessToken(gomock.Any(), "token").
					Return(&domain.RenewAccessTokenResult{
						AccessToken: "newtoken",
					}, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"refreshToken\":\"token\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "renew access token handler responds with Status Code 500 given auth service returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/refresh/",
			want:         RenewAccessTokenResponse{},
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
				authService.EXPECT().RenewAccessToken(gomock.Any(), "token").
					Return(nil, errors.New("server error"))
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockAuthService := NewMockAuthServicer(ctrl)
			authHandler := NewAuthHandler(mockAuthService)

			router := SetupRouter(authHandler, nil, nil, nil)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodPost, test.route, test.body)
			test.setupExpectations(request, mockAuthService)
			router.ServeHTTP(responseRecorder, request)

			// Assert

			// validate the response code
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			// validate the content type
			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got RenewAccessTokenResponse
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

func TestRevokeSession_Post(t *testing.T) {
	tests := []struct {
		contentType       string
		name              string
		responseCode      int
		route             string
		setupExpectations func(r *http.Request, authService *MockAuthServicer)
	}{
		{
			contentType:  "application/json; charset=utf-8",
			name:         "revoke session handler responds with Status Code 500 given auth service returns an error",
			responseCode: http.StatusInternalServerError,
			route:        "/revoke/invalidsessionid",
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
				authService.EXPECT().RevokeSession(gomock.Any(), "invalidsessionid").
					Return(errors.New("server error"))
			},
		},
		{
			contentType:  "",
			name:         "revoke session handler responds with Status Code 204 when valid data supplied",
			responseCode: http.StatusNoContent,
			route:        "/revoke/sessionid",
			setupExpectations: func(r *http.Request, authService *MockAuthServicer) {
				authService.EXPECT().RevokeSession(gomock.Any(), "sessionid").
					Return(nil)
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Arrange
			gin.SetMode(gin.TestMode)
			ctrl := gomock.NewController(t)

			mockAuthService := NewMockAuthServicer(ctrl)
			authHandler := NewAuthHandler(mockAuthService)

			router := SetupRouter(authHandler, nil, nil, nil)

			responseRecorder := httptest.NewRecorder()

			// Act
			request := httptest.NewRequest(http.MethodPost, test.route, nil)
			test.setupExpectations(request, mockAuthService)
			router.ServeHTTP(responseRecorder, request)

			// Assert

			// validate the response code
			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			// validate the content type
			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}
		})
	}
}
