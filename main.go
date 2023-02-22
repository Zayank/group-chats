package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Massad/gin-boilerplate/controllers"
	"github.com/Massad/gin-boilerplate/db"
	"github.com/Massad/gin-boilerplate/forms"
	"github.com/Massad/gin-boilerplate/helpers"
	"github.com/Massad/gin-boilerplate/models"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error: failed to load the env file")
	}

	router := gin.Default()

	binding.Validator = new(forms.DefaultValidator)

	initializeSession(router)

	initializeDB()

	initializeResourses(router)

	user := new(controllers.UserController)

	router.GET("/user/login", user.LoginForm)

	router.POST("/user/login", user.Login)
	router.POST("/user/register", user.Register)
	router.GET("/user/logout", user.Logout)

	group := new(controllers.GroupController)

	router.GET("/group/create-group", SessionAuthMiddleware(), group.CreateGroupForm)
	router.GET("/add-user-group/:groupuuid", SessionAuthMiddleware(), group.AddUserToGroupForm)

	router.POST("/group", SessionAuthMiddleware(), group.CreateGroup)
	router.POST("/group/:groupuuid/user/:userid", SessionAuthMiddleware(), bindGroupUuid(), group.AddUserToGroup)
	router.GET("/group", SessionAuthMiddleware(), group.GetGroups)
	router.GET("/group/:groupuuid", SessionAuthMiddleware(), hasAccessToGroup(false), bindGroupUuid(), group.GetGroupMessages)
	router.POST("/group/:groupuuid", SessionAuthMiddleware(), hasAccessToGroup(false), bindGroupUuid(), group.PostMessageInGroup)
	router.DELETE("/group/:groupuuid", SessionAuthMiddleware(), hasAccessToGroup(true), bindGroupUuid(), group.DeleteGroup)

	router.GET("/events/conversations/:groupuuid", SessionAuthMiddleware(), hasAccessToGroup(false), bindGroupUuid(), group.Stream)
	router.GET("/events/conversations", SessionAuthMiddleware(), group.StreamConversations)

	router.Run(":8091")
}

func initializeSession(router *gin.Engine) {

	store := cookie.NewStore([]byte(os.Getenv("COOKIE_SECRET")))

	router.Use(sessions.Sessions(os.Getenv("SESSION_NAME"), store))

}

func initializeDB() {

	db.Init()

	db.InitRedis(1)

}

func initializeResourses(router *gin.Engine) {

	router.LoadHTMLGlob("templates/*.tmpl")

	router.Static("templates/assets", "./templates/assets")

}

func bindGroupUuid() gin.HandlerFunc {

	return func(c *gin.Context) {

		groupuuid := c.Param("groupuuid")

		c.Set("GroupUuid", groupuuid)

		c.Next()
	}

}

/*

Desc: Authentication function

*/
func SessionAuthMiddleware() gin.HandlerFunc {

	return func(c *gin.Context) {

		session := sessions.Default(c)

		userId := session.Get("userId")

		if userId != nil {

			c.Next()

		} else {

			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrUnauthorzied)
		}
	}

}

/*

Desc: Authorization function

*/
func hasAccessToGroup(shouldBeAdmin bool) gin.HandlerFunc {

	var groupModel = new(models.GroupModel)

	return func(c *gin.Context) {

		groupuuid := c.Param("groupuuid")

		groupId := helpers.GetGroupId(groupuuid)

		if groupModel.CheckgroupExists(groupId) {

			userId := helpers.GetUserID(c)

			groupDetails, _ := groupModel.GetgroupDetails(groupId)

			validUser := true

			for i := range groupDetails.Members {

				if groupDetails.Members[i].Id == userId {

					if shouldBeAdmin && groupDetails.Members[i].IsAdmin {

						validUser = true

						break
					}
				}
			}

			if validUser {

				c.Next()

			} else {

				c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrForbidden)

			}

		} else {

			c.AbortWithStatusJSON(http.StatusUnauthorized, helpers.ErrUnauthorzied)

		}
	}
}
