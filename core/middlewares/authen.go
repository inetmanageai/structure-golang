package middlewares

import (
	"strings"
	"structure-golang/common/authorization"
	"structure-golang/config"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ApiKey(c *fiber.Ctx) error {
	if config.Env.Apikey != "" && config.Env.Apikey != c.Get("apikey") {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"status":  false,
			"message": "unauthorized",
			"data":    "",
		})
	}

	return c.Next()
}

func AccessToken(c *fiber.Ctx) error {
	authorizationHeader := c.Get("Authorization") // get from header Bearer
	cookie := c.Cookies("Accesstoken")            // get from cookie

	access_token := ""
	fields := strings.Fields(authorizationHeader)

	if len(fields) == 2 && fields[0] == "Bearer" {
		access_token = fields[1]
	} else {
		access_token = cookie
	}

	if access_token == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"status":  false,
			"message": "unauthorized",
			"data":    "",
		})
	}

	jwtHS256 := authorization.NewJWT_HS256(viper.GetString("app.signature.key"), viper.GetDuration("app.signature.expired"))

	sub := authorization.AppAuthorizationClaim{}
	err := jwtHS256.ValidateToken(access_token, &sub)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"code":    fiber.StatusUnauthorized,
			"status":  false,
			"message": err.Error(),
			"data":    "",
		})
	}

	return c.Next()
}
