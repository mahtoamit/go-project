package handler

import (
	"fmt"
	"time"
	"github.com/gofiber/fiber/v2"
	"github.com/tutorialedge/go-fiber-tutorial/auth"
	"github.com/tutorialedge/go-fiber-tutorial/database"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
	"github.com/tutorialedge/go-fiber-tutorial/configs"
)



func Login(c *fiber.Ctx) error {
	utils.InitLogger()
	request := new(models.LoginRequest)
	startTime := time.Now()

	if err := c.BodyParser(request); err != nil {
		utils.Log("ERROR", "handler", configs.Url_login,"","Login", "Error parsing JSON in Login:"+err.Error(),startTime, time.Now())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}
	// In real world application, you should hash the password
	user := models.SignupRequest{}
	database.Database.Where("email = ? AND password = ?", request.Email, request.Password).First(&user)
	if user.Email == "" {
		utils.Log("ERROR", "handler", configs.Url_login,"", "Login", "Invalid login credentials in Login",startTime, time.Now())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid login credentials",
		})
	}
	token, err := auth.GenerateToken(user.Email)
	if err != nil {
		utils.Log("ERROR", "handler", configs.Url_login,"", "Login", "Error in token generation in Login:"+err.Error(),startTime, time.Now())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not login",
		})
	}
	utils.Log("INFO","handler",configs.Url_login,"","Login","Login Succesfully.",startTime, time.Now())
	return c.JSON(fiber.Map{
		"token": token,
	})
}

func Signup(c *fiber.Ctx) error {
	utils.InitLogger()
	request := new(models.SignupRequest)
	startTime := time.Now()

	if err := c.BodyParser(request); err != nil {
		utils.Log("ERROR", "handler",configs.Url_signup,"","Signup", "Error parsing JSON in Signup:"+err.Error(), startTime, time.Now())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// In real world application, you should hash the password
	user := models.SignupRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Location: request.Location,
	}
	fmt.Println("data", user)
	database.Database.Create(&user)
	utils.Log("INFO", "handler", configs.Url_signup,"", "Signup", "User created successfully in Signup", startTime, time.Now())
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	utils.InitLogger()
	startTime := time.Now()
	utils.Log("INFO", "handler", configs.Url_logout,"", "Logout", "Processing logout Request.", startTime, time.Now())
	return c.JSON(fiber.Map{
		"message": "Logout successful.",
	})

}
