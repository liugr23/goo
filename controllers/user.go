package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/contrib/sessions"
	"fmt"
	"../models"
	"../forms"
	"log"
)

type UserController struct{}

var userModel = new(models.UserModel)

func (ctrl UserController) Hi(c *gin.Context) {
	fmt.Println("hi")
	c.JSON(200, gin.H{"message": "hi,jason."})
}

func getUserId(c *gin.Context) int64 {
	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId != nil {
		return models.ConvertToInt64(userId)
	}
	return 0
}

func getSessionUserInfo(c *gin.Context) (userSessionInfo models.UserSessionInfo) {
	session := sessions.Default(c)
	userId := session.Get("user_id")
	if userId != nil {
		userSessionInfo.ID = models.ConvertToInt64(userId)
		userSessionInfo.Name = session.Get("user_name").(string)
		userSessionInfo.Email = session.Get("user_email").(string)
	}
	return userSessionInfo
}

func (ctrl UserController) Signin(c *gin.Context) {
	var signinForm forms.SigninForm
	if c.BindJSON(&signinForm) != nil {
		c.JSON(406, gin.H{"message": "Invalid form", "form": signinForm})
		c.Abort()
		return
	}

	user, err := userModel.Signin(signinForm)
	if err == nil {
		session := sessions.Default(c)
		session.Set("user_id", user.Id)
		session.Set("user_email", user.Email)
		session.Set("user_name", user.Name)
		session.Save()

		c.JSON(200, gin.H{"message": "User signed in", "user": user})
	} else {
		c.JSON(406, gin.H{"message": "Invalid signin details", "error": err.Error()})
	}
}

func (ctrl UserController) Signup(c *gin.Context) {
	var signupForm forms.SignupForm
	log.Println(c.Request)
	if c.BindJSON(&signupForm) != nil {
		c.JSON(406, gin.H{"message": "Invalid form", "form": signupForm})
		c.Abort()
		return
	}

	user, err := userModel.Signup(signupForm)

	if err != nil {
		c.JSON(406, gin.H{"message": err.Error()})
		c.Abort()
		return
	}

	if user.Id > 0 {
		session := sessions.Default(c)
		session.Set("user_id", user.Id)
		session.Set("user_email", user.Email)
		session.Set("user_name", user.Name)
		session.Save()
		c.JSON(200, gin.H{"message": "Success signup", "user": user})
	} else {
		c.JSON(406, gin.H{"message": "Could not signup this user", "error": err.Error()})
	}

}
