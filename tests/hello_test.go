package tests

import (
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"ledger-app/controllers"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLedger(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/api", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := controllers.HelloWorld(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "Merhaba, dünya!", rec.Body.String())
	}
}
