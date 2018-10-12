package main

import (
	"Go_Docker/core"
	"Go_Docker/util"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hhkbp2/go-logging"
)

var router *gin.Engine
var logger logging.Logger

func main() {
	// Set Gin to production mode
	gin.SetMode(gin.ReleaseMode)

	if err := logging.ApplyConfigFile("./go_logconfig.yml"); err != nil {
		log.Printf("[FATAL] main: error initializing logger %v\n", err.Error())
		return
	}

	logger = logging.GetLogger("main")

	if !util.SettingsFromConfFile() {
		logger.Error("main(): Failed to initialize settings. Exiting ...")
		return
	}

	fmt.Println(util.Config.RedisAddress)
	// Set the router as the default one provided by Gin

	accessLogFile, err := os.OpenFile(util.Config.AccessLogPath+"access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	if err != nil {
		logger.Error("main(): Failed to open access log file")
		return
	}

	router = gin.New()

	// Initialize the routes
	core.InitializeRoutes(router, accessLogFile)

	util.InitRedis()

	if !util.InitMongoClient() {
		logger.Error("Failed to initialize Mongo")
		return
	}

	core.LoadCSV()

	// Start serving the application
	err = router.Run(util.Config.ListenPort)

	if err != nil {
		logger.Errorf("Server could not start : %v", err.Error())
		return
	}
}
