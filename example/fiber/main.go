package main

import (
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"

	syncpool "github.com/anacpe9/sync-pool"
)

const emptyString = ""

type DTOMiddleware[T interface{ Reset() }] struct {
	dtoPool *syncpool.Pool[T]
}

type ErrorResponse struct {
	FailedField string
	Tag         string
	Value       interface{}
	Error       string
}

type ResponseHTTP struct {
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Errors  []string    `json:"errors,omitempty"`
	Code    int         `json:"code"`
	OK      bool        `json:"ok,omitempty"`
}

type OtpDTO struct {
	TargetId  string `json:"targetId"  validate:"required,min=2,max=50"`
	CommandId string `json:"commandId" validate:"required"`
	Timestamp string `json:"timestamp" validate:"required"`
	OTP       string `json:"otp"       validate:"required"`
}

func (r *ResponseHTTP) Reset() {
	r.Data = nil
	r.Message = emptyString
	r.Error = emptyString
	r.Errors = nil
	r.Code = 0
	r.OK = false
}

func (j *OtpDTO) Reset() {
	j.CommandId = ""
	j.TargetId = ""
	j.Timestamp = ""
	j.OTP = ""
}

var validate = validator.New()
var responseHTTPPool = syncpool.GetPool[ResponseHTTP]()

func NewDTOMiddleware[T interface{ Reset() }]() *DTOMiddleware[T] {
	pool := syncpool.GetPool[T]() // when call GetPoo[T], it's call new (T) every time
	return &DTOMiddleware[T]{
		dtoPool: pool,
	}
}

// ref: https://docs.gofiber.io/guide/validation/
func ValidateStruct[T interface{}](user T) []*ErrorResponse {
	var errors []*ErrorResponse
	err := validate.Struct(user)
	if err != nil {
		if valErrors, ok := err.(validator.ValidationErrors); ok {
			for _, err := range valErrors {
				var element ErrorResponse
				element.FailedField = err.Field() // err.StructNamespace()
				element.Tag = err.Tag()
				element.Value = err.Value() // err.Param()
				// element.Error = err.Error()
				element.Error = fmt.Sprintf("'%s' has a value of '%v' which does not satisfy '%s'.", err.Field(), err.Value(), err.Tag())
				errors = append(errors, &element)

				// fmt.Printf("'%s' has a value of '%v' which does not satisfy '%s'.\n", err.Field(), err.Value(), err.Tag())
			}
			// } else if ive, ok := err.(validator.InvalidValidationError); ok {
			// 	var element ErrorResponse
			// 	element.Error = ive.Error()
			// 	errors = append(errors, &element)
		} else {
			var element ErrorResponse
			element.Error = err.Error()
			errors = append(errors, &element)
		}
	}
	return errors
}

func (mid *DTOMiddleware[T]) DTOValidate() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		resMsg := responseHTTPPool.Get()
		dtoBody := mid.dtoPool.Get()
		defer func() {
			dto := *dtoBody
			dto.Reset()
			mid.dtoPool.Put(&dto)

			resMsg.Reset()
			responseHTTPPool.Put(resMsg)
		}()

		// for body
		if err := ctx.BodyParser(dtoBody); err != nil {
			resMsg.Code = fiber.StatusBadRequest
			resMsg.Error = err.Error()
			resMsg.Message = resMsg.Error

			return ctx.
				Status(fiber.StatusBadRequest).
				JSON(resMsg)
		}

		errors := ValidateStruct(*dtoBody)
		if errors != nil {
			resMsg.Code = fiber.StatusBadRequest

			var errMsg string
			errMsgs := []string{}
			for idx, err := range errors {
				if idx == 0 {
					errMsg = err.Error
				}

				errMsgs = append(errMsgs, err.Error)
			}

			resMsg.Error = errMsg
			resMsg.Errors = errMsgs

			return ctx.
				Status(fiber.StatusBadRequest).
				JSON(resMsg)
		}

		// share dto to next chain
		ctx.Locals("body", *dtoBody)
		return ctx.Next()
	}
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	otpMid := NewDTOMiddleware[*OtpDTO]()
	app.Post(
		"/otp",
		otpMid.DTOValidate(),
		otpCtl)

	app.Listen(":3000")
}

// Receive OTP
//
//	@Summary		Receive OTP
//	@Description	Receive OTP and forward to a client.
//	@Tags			agents
//	@Success		200
//	@Accept			json
//	@Produce		json
//
//	@Param			node_id		path	string	true	"Node ID"	pattern('^\d{2,4}$')
//	@Param			target_id	path	string	true	"Target ID"
//	@Param			data		body	OtpDTO	true	"OTP from Agent"
//
//	@Router			/{node_id}/{target_id}/otp [post]
func otpCtl(ctx *fiber.Ctx) error {
	lBody := ctx.Locals("body")
	lDto := lBody.(*OtpDTO)

	fmt.Printf("[DEBUG] OtpDTO >>> %+v\n", lDto.OTP)

	return ctx.Status(http.StatusCreated).JSON(&struct {
		Data    interface{} `json:"data,omitempty"`
		Message string      `json:"message,omitempty"`
		Error   string      `json:"error,omitempty"`
		Errors  []string    `json:"errors,omitempty"`
		Code    int         `json:"code"`
		OK      bool        `json:"ok,omitempty"`
	}{
		Message: "Success",
		Code:    http.StatusCreated,
		OK:      true,
	})
}
