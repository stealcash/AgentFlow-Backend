package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/app/repo"
	"github.com/stealcash/AgentFlow/app/utils"
	"github.com/stealcash/AgentFlow/db"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptions.BadRequest(err.Error())
	}

	user, err := repo.GetUserByEmail(db.DB, input.Email)
	if err != nil {
		exceptions.Unauthorized("Invalid email or password")
	}

	if !utils.CheckPasswordHash(input.Password, user.PasswordHash) {
		exceptions.Unauthorized("Invalid email or password")
	}

	token, _ := utils.GenerateJWT(user.ID, user.UserType)

	// Do not return password hash!
	user.PasswordHash = ""

	responses.Success(c, "Login successful", gin.H{
		"token": token,
		"user":  user,
	})
}
