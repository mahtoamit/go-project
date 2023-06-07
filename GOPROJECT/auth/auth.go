package auth

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
)

var jwtSecret = []byte("Secret")

func GenerateToken(userId string) (string, error) {
	utils.InitLogger()

	startTime := time.Now() // Declare and assign the start time
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["userId"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	t, err := token.SignedString(jwtSecret)

	if err != nil {
		utils.Log("ERROR", "auth", "Token", "GenerateToken", "", "Error in Generating Token: "+err.Error(), startTime, time.Now())
		return "", err
	}

	return t, nil
}

func Protected() func(*fiber.Ctx) error {
	utils.InitLogger()
	return func(c *fiber.Ctx) error {
		startTime := time.Now() // Declare and assign the start time

		authorization := c.Get("Authorization")

		if authorization == "" {
			utils.Log("ERROR", "auth", "TokenAuthentication", "Protected", "", "Unauthorized access, no Authorization header provided", startTime, time.Now())
			return c.Status(401).SendString("Unauthorized")
		}

		token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				utils.Log("ERROR", "auth", "TokenAuthentication", "Protected", "", "Invalid Token", startTime, time.Now())
				return nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid token")
			}
			return jwtSecret, nil
		})

		if err != nil {
			utils.Log("ERROR", "auth", "TokenAuthentication", "Protected", "", "Error in Parsing Token: "+err.Error(), startTime, time.Now())
			return c.Status(401).SendString(err.Error())
		}

		if !token.Valid {
			utils.Log("ERROR", "auth", "TokenAuthentication", "Protected", "", "Invalid Token", startTime, time.Now())
			return c.Status(401).SendString("Unauthorized")
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userId, _ := claims["userId"].(string)
			c.Locals("userId", userId) // Store the userId in the context so it can be accessed later
		}

		return c.Next()
	}
}

