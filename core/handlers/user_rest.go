package handlers

import (
	"structure-golang/core/models"
	"structure-golang/core/services"
	"structure-golang/utils"

	"github.com/gofiber/fiber/v2"
)

type userHand struct {
	userSrv services.UserService
}

func NewUserHandler(userSrv services.UserService) userHand {
	return userHand{userSrv}
}

func (h userHand) Signin(c *fiber.Ctx) error {
	// parse body to model
	body := models.HandUserBodyModel{}
	if err := c.BodyParser(&body); err != nil {
		return utils.BodyParserFail(c)
	}

	// Signin service
	res, err := h.userSrv.Signin(body.Username, body.Password)
	if err != nil {
		appErr, ok := err.(utils.Err_Handler)
		if ok {
			return utils.ErrorFormat(c, appErr.Code, appErr.Message)
		}
	}

	resMsg := "signin success"
	return utils.SuccessFormat(c, fiber.StatusCreated, resMsg, res)
}
