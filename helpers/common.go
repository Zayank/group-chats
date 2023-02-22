package helpers

import (
	"net/http"
	"os"

	"github.com/Massad/gin-boilerplate/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

/*

Desc: retrieves userid from session

*/
func GetUserID(c *gin.Context) (userID string) {
	session := sessions.Default(c)
	return session.Get("userId").(string)
}

/*

Desc: retrieves user name from session

*/
func GetUserName(c *gin.Context) (userName string) {
	session := sessions.Default(c)
	return session.Get("userName").(string)
}

/*

Desc: appends header and footer onto the output body

*/
func DisplayHtmlUserOutputTemplate(c *gin.Context, bodyTemplateName string, bodyTemplateData interface{}) {

	var conversationModel = new(models.ConversationModel)
	var groupModel = new(models.GroupModel)

	userId := GetUserID(c)
	groupuuid := c.Param("groupuuid")

	conversations := conversationModel.GetConversations(userId)

	c.HTML(http.StatusOK, "header.tmpl", gin.H{
		"conversations": groupModel.GroupsList(conversations),
		"userid":        userId,
		"groupuuid":     groupuuid,
	})

	c.HTML(http.StatusOK, bodyTemplateName+".tmpl", bodyTemplateData)
	c.HTML(http.StatusOK, "footer.tmpl", gin.H{})
}

/*

Desc: format groupuuid to get group id

*/
func GetGroupId(groupuuid string) string {
	return os.Getenv("COLUMN_GROUP_PREFIX") + groupuuid
}

/*

Desc: format userid to get conversation id for broadcast

*/
func GetConversationsId(userid string) string {
	return os.Getenv("BROADCAST_CONVERSATIONS_PREFIX") + userid
}

/*

Desc: generate a new random string

*/
func GenerateUuid() string {
	return uuid.New().String()
}
