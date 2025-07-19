package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"github.com/stealcash/AgentFlow/app/entity"
	"github.com/stealcash/AgentFlow/app/repo"
	"github.com/stealcash/AgentFlow/app/utils"
	"github.com/stealcash/AgentFlow/db"
	"strings"
)

type SignupInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Company  string `json:"company_name"`
	UserType string `json:"user_type"` // admin or editor
	ParentID *int   `json:"parent_id"` // only for editor
}

func Signup(c *gin.Context) {
	var input SignupInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptions.BadRequest(err.Error())
	}

	// Normalize and validate user type
	userType := strings.ToLower(strings.TrimSpace(input.UserType))
	if userType == "" {
		userType = "editor"
	} else if userType != "admin" && userType != "editor" {
		exceptions.BadRequest("Invalid user_type. Must be 'admin' or 'editor'")
	}

	hashed, _ := utils.HashPassword(input.Password)

	user := entity.User{
		Email:        input.Email,
		PasswordHash: hashed,
		CompanyName:  input.Company,
		UserType:     userType,
		ParentID:     input.ParentID,
	}

	err := repo.CreateUser(db.DB, &user)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			exceptions.Conflict("Email already exists")
		}
		exceptions.Internal(err.Error())
	}

	token, _ := utils.GenerateJWT(user.ID, user.UserType)

	// Do not return password hash!
	user.PasswordHash = ""

	responses.Success(c, "Signup successful", gin.H{
		"token": token,
		"user":  user,
	})
}
