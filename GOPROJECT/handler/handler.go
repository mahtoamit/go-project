package handler

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tutorialedge/go-fiber-tutorial/auth"
	"github.com/tutorialedge/go-fiber-tutorial/constants"
	"github.com/tutorialedge/go-fiber-tutorial/models"
	"github.com/tutorialedge/go-fiber-tutorial/queries"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
)



func Login(c *fiber.Ctx) error {
	utils.InitLogger()
	request := new(models.LoginRequest)
	startTime := time.Now()

	if err := c.BodyParser(request); err != nil {
		utils.Log("ERROR", "handler", constants.Url_login,"","Login", "Error parsing JSON in Login:"+err.Error(),startTime, time.Now())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// validation for  email 
	valid_email := models.ValidateEmail(request.Email)

	if !valid_email {
		utils.Log("ERROR", "handler",constants.Url_signup,"","Login", "Email is not valid", startTime, time.Now())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is not valid "},)
	}
    // validation for password 
	Valid_password := models.ValidatePassword(request.Password)

	if !Valid_password {
		utils.Log("ERROR", "handler",constants.Url_signup,"","Login", "password is not valid", startTime, time.Now())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is not valid "},)
	}


	// In real world application, you should hash the password
	user := queries.Login(request.Email,request.Password)
	
	if user.Email == "" {
		utils.Log("ERROR", "handler", constants.Url_login,"", "Login", "Invalid login credentials in Login",startTime, time.Now())
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid login credentials",
		})
	}
	token:= auth.GenerateToken()

	if token == ""{
		utils.Log("ERROR", "handler", constants.Url_login,"", "Login", "Error in token generation in Login:",startTime, time.Now())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not login",
		})
	}

	//push the token to redis
	result,err := queries.SetAuthenticationCache(token,user.Email)
	
	if err != nil && !result {
		utils.Log("ERROR", "handler", constants.Url_login,"", "Login", "Error in token generation in Login:"+err.Error(),startTime, time.Now())
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Could not login",
		})
	}
	utils.Log("INFO","handler",constants.Url_login,"","Login","Login Succesfully.",startTime, time.Now())
	return c.JSON(fiber.Map{
		"token": token,
	})
}

func Signup(c *fiber.Ctx) error {
	utils.InitLogger()
	request := new(models.SignupRequest)
	startTime := time.Now()

	if err := c.BodyParser(request); err != nil {
		utils.Log("ERROR", "handler",constants.Url_signup,"","Signup", "Error parsing JSON in Signup:"+err.Error(), startTime, time.Now())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// validation for  email 
	valid_email := models.ValidateEmail(request.Email)

	if !valid_email {
		utils.Log("ERROR", "handler",constants.Url_signup,"","Signup", "Email is not valid", startTime, time.Now())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is not valid "},)
	}
    // validation for password 
	Valid_password := models.ValidatePassword(request.Password)

	if !Valid_password {
		utils.Log("ERROR", "handler",constants.Url_signup,"","Signup", "password is not valid", startTime, time.Now())
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "password is not valid "},)
	}

	// In real world application, you should hash the password
	user := models.SignupRequest{
		Name:     request.Name,
		Email:    request.Email,
		Password: request.Password,
		Location: request.Location,
	}
	fmt.Println("data", user)
	err := queries.CreateUser(user)
	if err!= nil{
		utils.Log("ERROR", "handler",constants.Url_signup,"","Signup", err.Error())
		return c.Status(253).JSON(fiber.Map{"msg":"User Not created"})
	}
	utils.Log("INFO", "handler", constants.Url_signup,"", "Signup", "User created successfully in Signup", startTime, time.Now())
	return c.Status(200).JSON(fiber.Map{"data":user})
}

func Logout(c *fiber.Ctx) error {
	utils.InitLogger()
	startTime := time.Now()
	authorization := c.Get("Authorization")
	result := queries.RedisDeleteAuthentication(authorization)
	if !result {
		utils.Log("ERROR", "handler", constants.Url_logout,"", "Logout", "Error while log out .", time.Now())
	return c.JSON(fiber.Map{
		"error": "Not able to log out.",
	})
	}
	utils.Log("INFO", "handler", constants.Url_logout,"", "Logout", "Processing logout Request.", startTime, time.Now())
	return c.JSON(fiber.Map{
		"message": "Logout successful.",
	})

}
