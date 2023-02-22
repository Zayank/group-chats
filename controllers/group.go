package controllers

import (
	"io"
	"net/http"
	"time"

	"github.com/Massad/gin-boilerplate/helpers"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type GroupController struct{}

var groupModel = new(models.GroupModel)
var conversationModel = new(models.ConversationModel)

var roomManager = NewRoomManager()
var conversationManager = NewConversationManager()

type CreateGroupRequestStruct struct {
	Name string `form:"name" binding:"required,min=2,max=255"`
}

/*

Desc: API function to create a group

*/
func (ctrl GroupController) CreateGroup(c *gin.Context) {

	params := CreateGroupRequestStruct{}

	if err := c.ShouldBind(&params); err != nil {

		errors, _ := err.(validator.ValidationErrors)

		c.JSON(http.StatusBadRequest, gin.H{"errors": helpers.Descriptive(errors)})
		return

	}

	groupName := params.Name

	groupuuid := helpers.GenerateUuid()

	groupId := helpers.GetGroupId(groupuuid)

	groupData := models.GroupDetails{
		Id:   groupuuid,
		Name: groupName,
		Members: []models.GroupMember{
			{
				Id:          helpers.GetUserID(c),
				IsAdmin:     true,
				JoiningDate: time.Now(),
			},
		},
	}

	if exists := groupModel.CheckgroupExists(groupId); !exists {

		inserted := groupModel.AddGroupDetails(groupId, groupData)

		if inserted {

			conversationModel.InsertIntoConersations(helpers.GetUserID(c), groupuuid)

			c.JSON(http.StatusOK, gin.H{
				"status":  "success",
				"message": groupData,
			})

		} else {

			c.JSON(http.StatusInternalServerError, helpers.ErrDB)

		}

	} else {
		c.JSON(http.StatusBadRequest, helpers.DynamicError(http.StatusBadRequest, "group already exists"))

	}

}

/*

Desc: Load UI to create a group

*/
func (ctrl GroupController) CreateGroupForm(c *gin.Context) {

	helpers.DisplayHtmlUserOutputTemplate(c, "create_group", gin.H{})
}

/*

Desc: Load UI to add user to a group

*/
func (ctrl GroupController) AddUserToGroupForm(c *gin.Context) {

	helpers.DisplayHtmlUserOutputTemplate(c, "add_group_user", gin.H{"groupuuid": c.GetString("GroupUuid")})

}

/*

Desc: API function to add user to a group

*/
func (ctrl GroupController) AddUserToGroup(c *gin.Context) {

	groupuuid := c.GetString("GroupUuid")

	groupId := helpers.GetGroupId(groupuuid)

	userId := c.Param("userid")

	AlreadyExists := false

	groupDetails, _ := groupModel.GetgroupDetails(groupId)

	for i := range groupDetails.Members {

		if groupDetails.Members[i].Id == userId {

			AlreadyExists = true
			break

		}
	}

	if AlreadyExists {

		c.JSON(http.StatusBadRequest, helpers.DynamicError(http.StatusBadRequest, "user already exists"))
		return

	}

	member := models.GroupMember{
		Id:          userId,
		IsAdmin:     true,
		JoiningDate: time.Now(),
	}

	groupDetails.Members = append(groupDetails.Members, member)

	if groupModel.AddGroupDetails(groupId, groupDetails) {

		conversationModel.InsertIntoConersations(userId, groupuuid)

		c.JSON(http.StatusCreated, helpers.DynamicSuccessMessage(http.StatusCreated, "user added to "+groupDetails.Name))

	} else {

		c.JSON(http.StatusInternalServerError, helpers.ErrDB)

	}

}

/*

Desc: Load UI for showing chat history & retrieve new messages of a group

*/
func (ctrl GroupController) GetGroupMessages(c *gin.Context) {

	groupuuid := c.GetString("GroupUuid")

	groupId := helpers.GetGroupId(groupuuid)

	groupDetails, _ := groupModel.GetgroupDetails(groupId)

	userId := helpers.GetUserID(c)

	history, _ := conversationModel.GetConversation(groupuuid)

	helpers.DisplayHtmlUserOutputTemplate(c, "chat_room", gin.H{
		"history":   history,
		"groupuuid": groupuuid,
		"roomName":  groupDetails.Name,
		"userid":    userId,
	})

}

/*

Desc: Load UI for showing conversations tab

*/
func (ctrl GroupController) GetGroups(c *gin.Context) {

	helpers.DisplayHtmlUserOutputTemplate(c, "", gin.H{})

}

/*

Desc: API function to post a new message

*/
func (ctrl GroupController) PostMessageInGroup(c *gin.Context) {

	groupuuid := c.GetString("GroupUuid")
	userid := helpers.GetUserID(c)
	message := c.PostForm("message")
	time := c.PostForm("time")

	roomManager.Submit(userid, helpers.GetUserName(c), groupuuid, time, message)

	sendMessageEventToEveryOne(userid, groupuuid, time, message)

	conversationModel.InsertIntoMessages(userid, groupuuid, message)

	c.JSON(http.StatusCreated, helpers.DynamicSuccessMessage(http.StatusCreated, "message posted"))
}

/*

Desc: function to sent message events to everyone in a group except the initiater

*/
func sendMessageEventToEveryOne(sentByUserid string, groupuuid string, time string, message string) {

	groupDetails, _ := groupModel.GetgroupDetails(helpers.GetGroupId(groupuuid))

	for i := range groupDetails.Members {
		currentUserId := groupDetails.Members[i].Id
		if currentUserId == sentByUserid {
			continue
		}

		user, _ := userModel.One(currentUserId)
		conversationManager.PushConversation(currentUserId, user.Name, groupuuid, groupDetails.Name, time, message)

	}

}

/*

Desc: Authentication function

*/
func (ctrl GroupController) DeleteGroup(c *gin.Context) {

	groupuuid := c.GetString("GroupUuid")

	roomManager.DeleteBroadcast(groupuuid)

	deleted := groupModel.DeleteGroup(groupuuid)

	if deleted {

		c.JSON(http.StatusOK, helpers.DynamicSuccessMessage(http.StatusOK, "group deleted"))

	} else {

	}

}

/*

Desc: Push new message events to group listener

*/
func (ctrl GroupController) Stream(c *gin.Context) {

	listenerId := helpers.GetGroupId(c.GetString("GroupUuid"))

	sendMessageEvent(c, listenerId)
}

/*

Desc: Push conversation events to conversations listener

*/
func (ctrl GroupController) StreamConversations(c *gin.Context) {

	listenerId := helpers.GetConversationsId(helpers.GetUserID(c))

	sendMessageEvent(c, listenerId)
}

/*

Desc: Push message event to listener

*/
func sendMessageEvent(c *gin.Context, listenerId string) {

	listener := conversationManager.OpenListener(listenerId)
	defer conversationManager.CloseListener(listenerId, listener)

	clientGone := c.Request.Context().Done()
	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case message := <-listener:
			c.SSEvent("message", message)
			return true
		}
	})
}
