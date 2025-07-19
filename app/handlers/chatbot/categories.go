package chatbot

import (
	"fmt"
	"github.com/stealcash/AgentFlow/app/customResponse/exceptions"
	"github.com/stealcash/AgentFlow/app/customResponse/responses"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/db"
)

type CategoryInput struct {
	Name     string `json:"name" binding:"required"`
	ParentID *int   `json:"parent_id"`
}

func CreateCategory(c *gin.Context) {
	userID := c.GetInt("user_id")
	chatbotID := c.Param("id")

	var exists bool
	err := db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM chatbots WHERE id = $1 AND user_id = $2)`,
		chatbotID, userID).Scan(&exists)

	if err != nil || !exists {
		exceptions.Unauthorized("Invalid chatbot ID")
	}

	var input CategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		exceptions.BadRequest(err.Error())
	}

	var id int
	err = db.DB.QueryRow(
		`INSERT INTO categories (chatbot_id, parent_id, name)
		 VALUES ($1, $2, $3) RETURNING id`,
		chatbotID, input.ParentID, input.Name,
	).Scan(&id)

	if err != nil {
		exceptions.Internal("Failed to create category")
	}

	responses.Success(c, "Category created", gin.H{"id": id})
}

func GetCategories(c *gin.Context) {
	userID := c.GetInt("user_id")
	chatbotID := c.Param("id")

	var exists bool
	err := db.DB.QueryRow(`SELECT EXISTS(SELECT 1 FROM chatbots WHERE id = $1 AND user_id = $2)`,
		chatbotID, userID).Scan(&exists)

	if err != nil || !exists {
		exceptions.Unauthorized("Invalid chatbot ID")
	}

	rows, err := db.DB.Query(`
		SELECT c.id, c.name, c.parent_id, m.image_path
		FROM categories c
		LEFT JOIN category_media m ON c.id = m.category_id
		WHERE c.chatbot_id = $1
	`, chatbotID)

	if err != nil {
		exceptions.Internal("Could not fetch categories")
	}
	defer rows.Close()

	var categories []gin.H
	for rows.Next() {
		var id int
		var name string
		var parentID *int
		var imagePath *string

		rows.Scan(&id, &name, &parentID, &imagePath)

		categories = append(categories, gin.H{
			"id":         id,
			"name":       name,
			"parent_id":  parentID,
			"image_path": imagePath,
		})
	}

	responses.Success(c, "Categories fetched", gin.H{
		"categories": categories,
	})
}

func UploadCategoryImage(c *gin.Context) {
	userID := c.GetInt("user_id")

	categoryIDStr := c.PostForm("category_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		exceptions.BadRequest("Invalid category ID")
	}

	var chatbotID int
	err = db.DB.QueryRow(`
		SELECT c.chatbot_id
		FROM categories c
		JOIN chatbots b ON c.chatbot_id = b.id
		WHERE c.id = $1 AND b.user_id = $2
	`, categoryID, userID).Scan(&chatbotID)

	if err != nil {
		exceptions.NotFound("Category not found for this user")
	}

	file, err := c.FormFile("image")
	if err != nil {
		exceptions.BadRequest("Image file is required")
	}

	uploadDir := fmt.Sprintf("uploads/user_%d/chatbot_%d/categories", userID, chatbotID)
	os.MkdirAll(uploadDir, os.ModePerm)

	filename := fmt.Sprintf("%d_%s", categoryID, filepath.Base(file.Filename))
	filePath := filepath.Join(uploadDir, filename)

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		exceptions.Internal("Upload failed")
	}

	_, err = db.DB.Exec(`
		INSERT INTO category_media (category_id, image_path)
		VALUES ($1, $2)
		ON CONFLICT (category_id) DO UPDATE SET image_path = EXCLUDED.image_path
	`, categoryID, filePath)

	if err != nil {
		exceptions.Internal("Failed to save image path")
	}

	responses.Success(c, "Image uploaded", gin.H{
		"image_path": filePath,
	})
}

func DeleteCategory(c *gin.Context) {
	userID := c.GetInt("user_id")
	categoryIDStr := c.Param("cat_id")
	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		exceptions.BadRequest("Invalid category ID")
	}

	var chatbotID int
	err = db.DB.QueryRow(`
		SELECT c.chatbot_id
		FROM categories c
		JOIN chatbots b ON c.chatbot_id = b.id
		WHERE c.id = $1 AND b.user_id = $2
	`, categoryID, userID).Scan(&chatbotID)

	if err != nil {
		exceptions.Unauthorized("Not authorized to delete this category")
	}

	_, err = db.DB.Exec(`DELETE FROM categories WHERE id = $1`, categoryID)
	if err != nil {
		exceptions.Internal("Failed to delete category")
	}

	responses.Success(c, "Category deleted", nil)
}
