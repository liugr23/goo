package main

import (
	"github.com/gin-gonic/gin"
	"./controllers"
	"github.com/gin-gonic/contrib/sessions"
	"log"
	"time"
	db "./database"
)

var uerController = new(controllers.UserController)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println(c.Request.URL.Path)
		path := c.Request.URL.Path
		if path != "/v1/user/hi" && path != "/v1/user/login" && path != "/v1/user/register" {
			session := sessions.Default(c)
			userID := session.Get("user_id")
			if userID == nil {
				log.Println(userID)
				c.JSON(403, gin.H{"message": "Please login first"})
				c.Abort()
				return
			}
		}
		c.Next()
	}
}

func main() {
	r := gin.Default()

	store := sessions.NewCookieStore([]byte("secret"))
	store.Options(sessions.Options{
		MaxAge: int(30 * time.Minute), //30min
		Path:   "/",
	})
	r.Use(sessions.Sessions("mysession", store))

	r.Use(AuthMiddleware())

	db.Init()

	v1 := r.Group("/v1")
	{
		user := new(controllers.UserController)
		v1.GET("/user/hi", user.Hi)
		v1.POST("/user/register", user.Signup)
	}

	r.Run(":9001")
}
