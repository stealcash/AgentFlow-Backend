package app

import (
	"github.com/gin-gonic/gin"
	"github.com/stealcash/AgentFlow/app/utils/helpers"
)

type CurrentUser struct {
	UserID   int
	UserType string
	ParentID int
}

// CurrentUserId returns the user ID.
func (s *CurrentUser) CurrentUserId() int {
	return s.UserID
}

// SetUserDetail sets user data from a generic map.
func (s *CurrentUser) SetUserDetail(userData map[string]interface{}) {
	s.UserID = helpers.ConvertToInt(userData["user_id"])
	s.ParentID = helpers.ConvertToInt(userData["parent_id"])

	if userType, ok := userData["user_type"].(string); ok {
		s.UserType = userType
	} else {
		s.UserType = ""
	}
}

// GetCurrentUserFromContext extracts CurrentUser from gin.Context.
func GetCurrentUserFromContext(c *gin.Context) *CurrentUser {
	if v, ok := c.Get("current_user"); ok {
		if currentUser, ok := v.(*CurrentUser); ok {
			return currentUser
		}
	}
	return nil
}

// SetCurrentUserInReqContext adds CurrentUser to gin.Context.
func SetCurrentUserInReqContext(c *gin.Context, userData map[string]interface{}) *CurrentUser {
	cu := &CurrentUser{}
	cu.SetUserDetail(userData)
	c.Set("current_user", cu)
	return cu
}
