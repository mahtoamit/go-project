package auth

import (
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/tutorialedge/go-fiber-tutorial/queries"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
)


func GenerateToken() (string) {
	
    requestID := uuid.New().String()
	
	return requestID
}

func Protected(c *fiber.Ctx)  error {
	utils.InitLogger()
	
		startTime := time.Now() // Declare and assign the start time

		authorization := c.Get("Authorization")

		if authorization == "" {
			utils.Log("ERROR", "auth", "TokenAuthentication", "Protected", "", "Unauthorized access, no Authorization header provided", startTime, time.Now())
			return c.Status(401).JSON(fiber.Map{"error":"Unauthorized"})
		}

        
		//validate the authentication token
		result := queries.GetAuthenticationCache(authorization)

		if result == "" {
			utils.Log("ERROR", "auth", "TokenAuthentication", "Protected", "", "Invalid Token", startTime, time.Now())
			return c.Status(401).JSON(fiber.Map{"error":"Unauthorized"})
		}

		c.Locals("userId", result) // Store the userId in the context so it can be accessed later


		return c.Next()
	}


