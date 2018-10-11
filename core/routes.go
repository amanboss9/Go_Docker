package core

import (
	"Go_Docker/util"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gopkg.in/mgo.v2/bson"
)

var allowedOrigins = []string{"http://localhost:8080"}

func customLogger(out io.Writer) gin.HandlerFunc {

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		fmt.Fprintf(out, "%v | %3d | %13v | %15s | %-7s %s\n",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

func InitializeRoutes(router *gin.Engine, writer io.Writer) {
	router.Use(customLogger(writer))
	router.Use(gin.Recovery())

	// var store redis.Store
	// var err error
	// if util.IsDevSetup() {
	// 	store, _ = redis.NewStore(10, "tcp", strings.Split(util.Config.RedisAddress, ",")[0], "", []byte("secret"))
	// } else {
	// 	store, err = redis.NewStore(10, "tcp", strings.Split(util.Config.RedisAddress, ",")[0], "", []byte("secret"))

	// 	if err != nil {
	// 		fmt.Println("the error is ", err.Error())
	// 	}
	// }

	// store.Options(sessions.Options{MaxAge: util.Config.SessionTimeout, Path: "/", HttpOnly: true, Secure: false})
	// router.Use(sessions.Sessions("mysession", store))

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowCredentials = true
	config.AllowHeaders = []string{"Origin", "Content-Type"}
	config.AllowOriginFunc = func(origin string) bool { return true }
	router.Use(cors.New(config))

	router.Use(func(c *gin.Context) {
		allowOriginVal := c.Request.Header["Origin"]

		if len(allowOriginVal) > 0 {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowOriginVal[0])
		}
		c.Next()
	})

	// router.LoadHTMLGlob("./static/**/*")

	// Handle the index route
	router.GET("/:key", handleone)
	// router.GET("/two", loginHandler)
}

func handleone(c *gin.Context) {
	key := strings.TrimPrefix(c.Request.URL.Path, "/")

	session := util.MongoSession.Copy()
	defer session.Close()

	collection := session.DB(util.Config.DbName).C("csvload")

	result := Mongo{}
	err := collection.Find(bson.M{"key": key}).One(&result)

	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Not found!",
		})
	} else {
		c.JSON(200, result)
	}

}
