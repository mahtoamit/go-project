package configs

import (
	"github.com/joho/godotenv"
	"github.com/tutorialedge/go-fiber-tutorial/utils"
	"log"
	"os"
	"time"
)

func EnvDBURI(key string) (string)  {
     utils.InitLogger()
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error while loading .env file")
		startTime := time.Now()
		endTime := time.Now()
		utils.Log("ERROR", "configs", "env.go", "", "","Error while loading.envfile"+ err.Error(),startTime, endTime)
		os.Exit(1)
	}
	return os.Getenv(key)
	

}
