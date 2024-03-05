package post_handler_test

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/AtIasShrugged/antisocial/internal/domain/models"
	post_handler "github.com/AtIasShrugged/antisocial/internal/http/handler/post"
	postRepo "github.com/AtIasShrugged/antisocial/internal/repository/post"
	repoMock "github.com/AtIasShrugged/antisocial/internal/repository/post/mocks"
	"github.com/AtIasShrugged/antisocial/internal/service/post"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	reqID := 1
	repo := repoMock.NewMockPostRepository(ctrl)
	exp := models.Post{
		ID:       1,
		AuthorID: 1,
		Body:     "test",
	}
	repo.EXPECT().GetByID(ctx, reqID).Return(exp, nil).Times(1)

	service := post.New(repo, log)
	handler := post_handler.New(service, log)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/posts/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(reqID))

	expected :=
		`{"id":1,"author_id":1,"body":"test"}` + "\n"

	if assert.NoError(t, handler.GetByID(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, rec.Body.String())
	}
}

func TestGetByIDBadParams(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	reqID := "err"
	repo := repoMock.NewMockPostRepository(ctrl)

	service := post.New(repo, log)
	handler := post_handler.New(service, log)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/posts/:id")
	c.SetParamNames("id")
	c.SetParamValues(reqID)

	if assert.NoError(t, handler.GetByID(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `"bad params: strconv.Atoi: parsing \"err\": invalid syntax"`+"\n", rec.Body.String())
	}
}

func TestGetByIDNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	reqID := 1
	repo := repoMock.NewMockPostRepository(ctrl)
	repo.EXPECT().GetByID(ctx, reqID).Return(models.Post{}, postRepo.ErrPostNotFound).Times(1)

	service := post.New(repo, log)
	handler := post_handler.New(service, log)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/posts/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(reqID))
	if assert.NoError(t, handler.GetByID(c)) {
		assert.Equal(t, http.StatusNotFound, rec.Code)
		assert.Equal(t, "post not found", rec.Body.String())
	}
}

func TestGetByIDRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	reqID := 1
	repo := repoMock.NewMockPostRepository(ctrl)
	repo.EXPECT().GetByID(ctx, reqID).Return(models.Post{}, fmt.Errorf("db is down")).Times(1)

	service := post.New(repo, log)
	handler := post_handler.New(service, log)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/posts/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(reqID))

	if assert.NoError(t, handler.GetByID(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `"db is down"`+"\n", rec.Body.String())
	}
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	repo := repoMock.NewMockPostRepository(ctrl)
	postBody := models.Post{
		AuthorID: 1,
		Body:     "test",
	}
	postId := 1
	repo.EXPECT().Create(ctx, postBody).Return(postId, nil).Times(1)

	service := post.New(repo, log)
	handler := post_handler.New(service, log)

	postJson := `{"author_id":1,"body":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/posts/create", strings.NewReader(postJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.Create(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, strconv.Itoa(postId)+"\n", rec.Body.String())
	}
}

func TestCreateBadJson(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	repo := repoMock.NewMockPostRepository(ctrl)

	service := post.New(repo, log)
	handler := post_handler.New(service, log)

	badPostJson := `{"author_id":1,"body":"test}`
	req := httptest.NewRequest(http.MethodPost, "/posts/create", strings.NewReader(badPostJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.Create(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `"bad json: code=400, message=unexpected EOF, internal=unexpected EOF"`+"\n", rec.Body.String())
	}
}

func TestCreateRepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	e := echo.New()
	ctx := context.Background()
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	postBody := models.Post{
		AuthorID: 1,
		Body:     "test",
	}
	repo := repoMock.NewMockPostRepository(ctrl)
	repo.EXPECT().Create(ctx, postBody).Return(0, fmt.Errorf("db is down")).Times(1)

	service := post.New(repo, log)
	handler := post_handler.New(service, log)

	postJson := `{"author_id":1,"body":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/posts/create", strings.NewReader(postJson))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)

	if assert.NoError(t, handler.Create(c)) {
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.Equal(t, `"db is down"`+"\n", rec.Body.String())
	}
}
