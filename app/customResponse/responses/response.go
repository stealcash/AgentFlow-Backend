package responses

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app"
	"net/http"
)

// FailedApiResponse sends a standardized error response.
func FailedApiResponse(c *gin.Context, httpStatusCode int, err interface{}) {
	if httpStatusCode == 0 {
		httpStatusCode = http.StatusInternalServerError
	}

	var errors []string
	switch v := err.(type) {
	case error:
		errors = []string{v.Error()}
	case string:
		errors = []string{v}
	default:
		errors = []string{"An unexpected error occurred."}
	}

	c.AbortWithStatusJSON(httpStatusCode, gin.H{
		"status": 0,
		"errors": errors,
		"header": getHeader(c),
	})
}

// Success sends a simple success response.
func Success(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": msg,
		"data":    data,
		"header":  getHeader(c),
	})
}

// SuccessApiResponseWithParams sends a success response with dynamic params.
func SuccessApiResponseWithParams(c *gin.Context, baseMsg string, params map[string]string, data interface{}) {
	formattedMsg := baseMsg
	for k, v := range params {
		formattedMsg = fmt.Sprintf("%s %s=%s", formattedMsg, k, v)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  1,
		"message": formattedMsg,
		"data":    data,
		"header":  getHeader(c),
	})
}

// SuccessPublicResponse returns success with only status & data.
func SuccessPublicResponse(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"status": 1,
		"data":   data,
	})
}

// getHeader returns minimal CurrentUser context info.
func getHeader(c *gin.Context) map[string]interface{} {
	header := make(map[string]interface{})
	cu := app.GetCurrentUserFromContext(c)
	if cu != nil {
		header["user_id"] = cu.UserID
		header["user_type"] = cu.UserType
		header["parent_id"] = cu.ParentID
	}
	return header
}
