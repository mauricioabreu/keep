package server_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/mauricioabreu/keep/internal/db"
	"github.com/mauricioabreu/keep/internal/mocks"
	"github.com/mauricioabreu/keep/internal/server"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"
)

type NoteHandlerSuite struct {
	suite.Suite
	ctrl        *gomock.Controller
	Echo        *echo.Echo
	noteStorer  *mocks.MockNoteStorer
	noteHandler *server.NoteHandler
}

func (suite *NoteHandlerSuite) SetupTest() {
	suite.ctrl = gomock.NewController(suite.T())
	suite.noteStorer = mocks.NewMockNoteStorer(suite.ctrl)
	suite.noteHandler = server.NewNoteHandler(suite.noteStorer, zap.NewNop().Sugar())
	suite.Echo = echo.New()
}

func (suite *NoteHandlerSuite) TearDownTest() {
	suite.ctrl.Finish()
}

func (suite *NoteHandlerSuite) TestCreateNoteSuccess() {
	reqBody := `
	{
		"title": "Test Title",
		"content": "Test Content"
	}`
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)

	suite.noteStorer.EXPECT().CreateNote(gomock.Any(), db.CreateNoteParams{
		Title:   "Test Title",
		Content: "Test Content",
	}).Return(db.Note{
		ID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Title:   "Test Title",
		Content: "Test Content",
	}, nil)

	err := suite.noteHandler.CreateNote(c)
	suite.NoError(err)
	suite.Equal(http.StatusCreated, rec.Code)

	expectedBody := `
	{
		"message": "Note created",
		"data": {
			"id":"123e4567-e89b-12d3-a456-426614174000",
			"title":"Test Title",
			"content":"Test Content"
		}
	}`
	suite.JSONEq(expectedBody, rec.Body.String())
}

func (suite *NoteHandlerSuite) TestCreateNoteFailureStore() {
	reqBody := `
	{
		"title": "Test Title",
		"content": "Test Content"
	}`
	req := httptest.NewRequest(http.MethodPost, "/notes", bytes.NewBufferString(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)

	suite.noteStorer.EXPECT().CreateNote(gomock.Any(), gomock.Any()).
		Return(db.Note{}, errors.New("failed to create note"))

	err := suite.noteHandler.CreateNote(c)
	suite.NoError(err)
	suite.Equal(http.StatusInternalServerError, rec.Code)
}

func (suite *NoteHandlerSuite) TestGetNoteSuccess() {
	req := httptest.NewRequest(http.MethodGet, "/notes/123e4567-e89b-12d3-a456-426614174000", nil)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000")

	suite.noteStorer.EXPECT().GetNote(gomock.Any(), uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")).Return(db.Note{
		ID:      uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		Title:   "Test Note",
		Content: "Test Content",
	}, nil)

	err := suite.noteHandler.GetNote(c)
	suite.NoError(err)
	suite.Equal(http.StatusOK, rec.Code)

	expectedBody := `
    {
        "message": "Note found",
        "data": {
            "id":"123e4567-e89b-12d3-a456-426614174000",
            "title":"Test Note",
            "content":"Test Content"
        }
    }`
	suite.JSONEq(expectedBody, rec.Body.String())
}

func (suite *NoteHandlerSuite) TestGetNoteNotFound() {
	req := httptest.NewRequest(http.MethodGet, "/notes/123e4567-e89b-12d3-a456-426614174000", nil)
	rec := httptest.NewRecorder()
	c := suite.Echo.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("123e4567-e89b-12d3-a456-426614174000")

	suite.noteStorer.EXPECT().GetNote(gomock.Any(), uuid.MustParse("123e4567-e89b-12d3-a456-426614174000")).Return(db.Note{}, pgx.ErrNoRows)

	err := suite.noteHandler.GetNote(c)
	suite.NoError(err)
	suite.Equal(http.StatusNotFound, rec.Code)

	expectedBody := `
    {
        "error": {
            "code": "NOT_FOUND",
            "message": "Note not found"
        }
    }`
	suite.JSONEq(expectedBody, rec.Body.String())
}

func TestNoteHandler(t *testing.T) {
	suite.Run(t, new(NoteHandlerSuite))
}
