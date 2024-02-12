package server_test

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/mauricioabreu/keep/internal/db"
	"github.com/mauricioabreu/keep/internal/mocks"
	"github.com/mauricioabreu/keep/internal/server"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
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
	suite.noteHandler = server.NewNoteHandler(suite.noteStorer)
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
		ID:      1,
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
			"id":1,
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

func TestNoteHandler(t *testing.T) {
	suite.Run(t, new(NoteHandlerSuite))
}
