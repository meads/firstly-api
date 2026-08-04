package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gin-gonic/gin"
	db "github.com/meads/firstly-api/db/sqlc"
	"github.com/meads/firstly-api/security"
	"go.uber.org/mock/gomock"
)

func TestNotesHandler_Post(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              CreateNoteResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier)
	}{
		{
			body:         bytes.NewBufferString("{\"title\":\"note 1\",\"content\":\"this is a note\",\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/users/1/notes/",
			want:         CreateNoteResponse{ID: 1, Title: "note 1", Content: "this is a note", UserID: 1},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID, title, note := int64(1), "note 1", "this is a note"
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(true, nil)
				param := db.CreateNoteParams{UserID: userID, Title: title, Content: note}
				querier.EXPECT().CreateNote(gomock.Any(), param).Return(db.Note{
					ID: int64(1), Title: title, Content: note, UserID: userID,
				}, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"title\":\"note 1\",\"content\":\"this is a note\",\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with status code 500 given create note responds with an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID, title, note := int64(1), "note 1", "this is a note"
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(true, nil)
				param := db.CreateNoteParams{UserID: userID, Title: title, Content: note}
				querier.EXPECT().CreateNote(gomock.Any(), param).Return(db.Note{}, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"title\":\"note 1\",\"content\":\"this is a note\",\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with status code 400 given user not found",
			responseCode: http.StatusBadRequest,
			route:        "/users/1/notes/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				querier.EXPECT().UserExists(gomock.Any(), int64(1)).Return(false, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"title\":\"note 1\",\"content\":\"this is a note\",\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with status code 500 given userexists responds with an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				querier.EXPECT().UserExists(gomock.Any(), int64(1)).Return(false, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with status code 400 given request json is invalid",
			responseCode: http.StatusBadRequest,
			route:        "/users/1/notes/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
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
			var got CreateNoteResponse
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

func TestNotesHandler_Put(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              UpdateNoteResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier)
	}{
		{
			body:         bytes.NewBufferString("{\"id\":1,\"userId\":1,\"title\":\"note 1\",\"content\":\"this is a note\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/users/1/notes/",
			want:         UpdateNoteResponse{ID: int64(1), UserID: int64(1), Title: "note 1", Content: "this is a note"},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				id, userID, title, note := int64(1), int64(1), "note 1", "this is a note"
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(true, nil)
				param := db.UpdateNoteParams{ID: id, UserID: userID, Content: note, Title: title}
				querier.EXPECT().UpdateNote(gomock.Any(), param).Return(nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"userId\":1,\"title\":\"note 1\",\"content\":\"this is a note\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with status code 500 given update note responds with an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				id, userID, title, note := int64(1), int64(1), "note 1", "this is a note"
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(true, nil)
				param := db.UpdateNoteParams{ID: id, UserID: userID, Content: note, Title: title}
				querier.EXPECT().UpdateNote(gomock.Any(), param).Return(errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"userId\":1,\"title\":\"note 1\",\"content\":\"this is a note\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with status code 404 given user exists responds with false",
			responseCode: http.StatusNotFound,
			route:        "/users/1/notes/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				querier.EXPECT().UserExists(gomock.Any(), int64(1)).Return(false, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"id\":1,\"userId\":1,\"title\":\"note 1\",\"content\":\"this is a note\"}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with status code 500 given user exists responds with an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				querier.EXPECT().UserExists(gomock.Any(), int64(1)).Return(false, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with status code 400 given supplied JSON is invalid",
			responseCode: http.StatusBadRequest,
			route:        "/users/1/notes/",
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
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
			request := httptest.NewRequest(http.MethodPut, test.route, test.body)
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
			var got UpdateNoteResponse
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

func TestNotesHandler_Get(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              ListNotesResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier)
	}{
		{
			body:         bytes.NewBufferString("{\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/users/1/notes/",
			want: ListNotesResponse{
				Notes: []db.Note{
					{ID: int64(1), UserID: int64(1), Title: "note 1", Content: "this is a note"}},
			},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID := int64(1)
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(true, nil)
				querier.EXPECT().ListNotesByUserID(gomock.Any(), userID).Return([]db.Note{
					{ID: int64(1), UserID: int64(1), Title: "note 1", Content: "this is a note"},
				}, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 500 given list notes by user id responds with an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/",
			want:         ListNotesResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID := int64(1)
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(true, nil)
				querier.EXPECT().ListNotesByUserID(gomock.Any(), userID).Return([]db.Note{}, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 404 given user exists returns false",
			responseCode: http.StatusNotFound,
			route:        "/users/1/notes/",
			want:         ListNotesResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID := int64(1)
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(false, nil)
			},
		},
		{
			body:         bytes.NewBufferString("{\"userId\":1}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 500 given user exists responds with error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/",
			want:         ListNotesResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID := int64(1)
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(false, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString("{}"),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 400 given user id param invalid",
			responseCode: http.StatusBadRequest,
			route:        "/users/a/notes/",
			want:         ListNotesResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
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

			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got ListNotesResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}

func TestNotesHandler_Delete(t *testing.T) {

	tests := []struct {
		body              *bytes.Buffer
		contentType       string
		name              string
		responseCode      int
		route             string
		want              DeleteNoteResponse
		setupExpectations func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier)
	}{
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 200 when valid data supplied",
			responseCode: http.StatusOK,
			route:        "/users/1/notes/1",
			want:         DeleteNoteResponse{Message: "Resource successfully deleted"},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID := int64(1)
				noteID := int64(1)
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(true, nil)
				querier.EXPECT().DeleteNote(gomock.Any(), db.DeleteNoteParams{ID: noteID, UserID: userID}).Return(nil)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 500 when delete note responds with an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/1",
			want:         DeleteNoteResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID := int64(1)
				noteID := int64(1)
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(true, nil)
				querier.EXPECT().DeleteNote(gomock.Any(), db.DeleteNoteParams{ID: noteID, UserID: userID}).Return(
					errors.New("server error"),
				)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 404 when user exists returns false",
			responseCode: http.StatusNotFound,
			route:        "/users/1/notes/1",
			want:         DeleteNoteResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID := int64(1)
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(false, nil)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 500 when user exists responds with an error",
			responseCode: http.StatusInternalServerError,
			route:        "/users/1/notes/1",
			want:         DeleteNoteResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
				userID := int64(1)
				querier.EXPECT().UserExists(gomock.Any(), userID).Return(false, errors.New("server error"))
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 400 when invalid :noteid param",
			responseCode: http.StatusBadRequest,
			route:        "/users/1/notes/b",
			want:         DeleteNoteResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
			},
		},
		{
			body:         bytes.NewBufferString(""),
			contentType:  "application/json; charset=utf-8",
			name:         "notes handler responds with Status Code 400 when invalid :userid param",
			responseCode: http.StatusBadRequest,
			route:        "/users/a/notes/1",
			want:         DeleteNoteResponse{},
			setupExpectations: func(r *http.Request, tokener *security.MockTokener, hasher *security.MockHasher, querier *db.MockQuerier) {
				passClaimsMiddleware(r, tokener, hasher, querier)
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

			if responseRecorder.Code != test.responseCode {
				t.Fatalf("expected status code %d, got %d", test.responseCode, responseRecorder.Code)
			}

			if responseRecorder.Header().Get("Content-Type") != test.contentType {
				t.Fatalf("expected content type %s, got %s", test.contentType, responseRecorder.Header().Get("Content-Type"))
			}

			// Decode and verify the JSON Body
			var got DeleteNoteResponse
			err := json.NewDecoder(responseRecorder.Body).Decode(&got)
			if err != nil {
				t.Fatalf("Failed to decode JSON response: %v", err)
			}

			// Compare structural values
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Response mismatch!\n Want: %+v\n Got:  %+v", test.want, got)
			}
		})
	}
}
