package controllers

import (
	"log"
	"main/models"
	"main/utils/token"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CurrentUser(c *gin.Context) {
	user_id, err := token.ExtractTokenID(c)

	// log.Println("user_id is ", user_id)

	log.Println(c.Request.Header)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	u, err := models.GetUserByID(user_id)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	c.JSON(http.StatusOK, gin.H{"message": "success", "data": u})
}

// ----- ----- -----

type LoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}

	u := models.User{Name: input.Username, Password: input.Password}

	token, err := models.LoginCheck(u.Name, u.Password)

	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": "username or password is incorrect",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

// ----- ----- -----

type RegisterInput struct {
	// binding表示字段是必须的
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
	// 需要进行验证的输入
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	u := models.User{Name: input.Username, Password: input.Password}

	// 注册完保存用户
	_, err := u.SaveUser()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error(),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "validated!",
	})
}
