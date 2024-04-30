package core_dtos_test

import (
	coreDtos "github.com/WildEgor/g-core/pkg/core/dtos"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"os"
	"testing"
)

var (
	fiberCtx *fiber.Ctx
	app      *fiber.App
)

func TestMain(m *testing.M) {
	app = fiber.New()
	defer app.Shutdown()

	fiberCtx = app.AcquireCtx(&fasthttp.RequestCtx{})

	code := m.Run()
	os.Exit(code)
}

func TestWrapResponse(t *testing.T) {
	resp := coreDtos.NewResponse(fiberCtx)

	err := resp.JSON()

	assert.Nil(t, err)
}

func TestResponseDto(t *testing.T) {
	resp := coreDtos.NewResponse(fiberCtx)

	data := []string{"test"}

	resp.SetStatus(200)
	resp.SetData(data)
	err := resp.JSON()

	assert.Nil(t, err)
	assert.Equal(t, 1, len(resp.Data))
	assert.Equal(t, data[0], resp.Data[0])
	assert.Equal(t, "ok", resp.Status)
	assert.Equal(t, 0, len(resp.Errors))
}
