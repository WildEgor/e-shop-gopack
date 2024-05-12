package core_dtos_test

import (
	coreDtos "github.com/WildEgor/e-shop-gopack/pkg/core/dtos"
	fiberold "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"os"
	"testing"
)

var (
	fiberCtx    fiber.Ctx
	fiberoldCtx *fiberold.Ctx
	app         *fiber.App
	oldApp      *fiberold.App
)

func TestMain(m *testing.M) {
	app = fiber.New()
	defer app.Shutdown()
	fiberCtx = app.AcquireCtx(&fasthttp.RequestCtx{})

	oldApp = fiberold.New()
	defer oldApp.Shutdown()

	fiberoldCtx = oldApp.AcquireCtx(&fasthttp.RequestCtx{})

	code := m.Run()
	os.Exit(code)
}

func TestWrapResponse(t *testing.T) {
	resp := coreDtos.NewResponse(fiberCtx)

	err := resp.JSON()

	assert.Nil(t, err)
}

func TestOldWrapResponse(t *testing.T) {
	resp := coreDtos.NewResp(coreDtos.WithOldContext(fiberoldCtx))

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
