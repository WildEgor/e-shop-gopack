package core_dtos

import (
	fiberold "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v3"
	"reflect"
	"time"
)

const (
	OkStatus   = "ok"
	FailStatus = "fail"
)

// ErrorItem contains error's key and message
type ErrorItem struct {
	Key     int    `json:"key"`
	Message string `json:"message"`
}

type ResponseOptions struct {
	ctx    fiber.Ctx
	oldCtx *fiberold.Ctx
}

// ResponseDto generic response DTO
type ResponseDto struct {
	ResponseOptions

	Status      string        `json:"status"`
	Message     string        `json:"message"`
	Errors      []ErrorItem   `json:"errors"`
	Data        []interface{} `json:"data"`
	TmRequest   string        `json:"tm_req"`
	TmRequestSt time.Time     `json:"-"`
}

type Option func(o *ResponseOptions)

func WithOldContext(ctx *fiberold.Ctx) Option {
	return func(o *ResponseOptions) {
		o.oldCtx = ctx
	}
}

func WithContext(ctx fiber.Ctx) Option {
	return func(o *ResponseOptions) {
		o.ctx = ctx
	}
}

// NewResponse wrap context (compatible only with fiber v3)
func NewResponse(ctx fiber.Ctx) *ResponseDto {
	return &ResponseDto{
		Errors:      make([]ErrorItem, 0),
		Data:        make([]interface{}, 0),
		TmRequestSt: time.Now(),

		ResponseOptions: ResponseOptions{
			ctx: ctx,
		},
	}
}

// NewResp wrap context (compatible both for v2 and v3)
func NewResp(opts ...Option) *ResponseDto {
	options := &ResponseOptions{}

	for _, o := range opts {
		o(options)
	}

	return &ResponseDto{
		Errors:      make([]ErrorItem, 0),
		Data:        make([]interface{}, 0),
		TmRequestSt: time.Now(),

		ResponseOptions: *options,
	}
}

// SetError set key/msg error
func (r *ResponseDto) SetError(key int, msg string) *ResponseDto {
	r.Errors = append(r.Errors, ErrorItem{
		Key:     key,
		Message: msg,
	})
	return r
}

// SetMessage response message
func (r *ResponseDto) SetMessage(text string) {
	r.Message = text
}

// SetHeaders response headers
func (r *ResponseDto) SetHeaders(headers map[string]string) {
	if r.ResponseOptions.ctx != nil {
		for _, header := range headers {
			r.ctx.Set(header, headers[header])
		}
	}

	if r.ResponseOptions.oldCtx != nil {
		for _, header := range headers {
			r.oldCtx.Set(header, headers[header])
		}
	}
}

// SetStatus response status
func (r *ResponseDto) SetStatus(status int) {
	if r.ctx != nil {
		r.ctx.Status(status)
	}

	if r.oldCtx != nil {
		r.oldCtx.Status(status)
	}
}

// SetData response data
func (r *ResponseDto) SetData(data interface{}) {
	r.Data = r.reflectData(data)
}

// JSON finalize response and convert to Json
func (r *ResponseDto) JSON() error {
	r.TmRequest = time.Since(r.TmRequestSt).String()

	if len(r.Errors) > 0 {
		r.Status = FailStatus
	} else {
		r.Status = OkStatus
	}

	if r.ctx != nil {
		err := r.ctx.JSON(r)
		if err != nil {
			return err
		}
	}

	if r.oldCtx != nil {
		err := r.oldCtx.JSON(r)
		if err != nil {
			return err
		}
	}

	return nil
}

// reflectData convert data to slice
func (r *ResponseDto) reflectData(in interface{}) []interface{} {
	sType := reflect.ValueOf(in)

	if sType.Kind() != reflect.Slice {
		return []interface{}{}
	}

	ret := make([]interface{}, sType.Len())

	for i := 0; i < sType.Len(); i++ {
		ret[i] = sType.Index(i).Interface()
	}

	return ret
}
