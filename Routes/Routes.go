package Routes

import(
	"first-api/Controllers"
	"first-api/Models"
	"github.com/gin-gonic/gin"
)

// setup Router .. Config routes
func SetupRouter(userStore *Models.UserStore) *gin.Engine{
	r :=gin.Default()
	grp1:= r.Group("/user-api")
	{
		grp1.GET("user",Controllers.GetUsers)
		grp1.POST("user",Controllers.NewUserController(userStore))
		grp1.GET("user/:id",Controllers.GetUserByID)
		grp1.PUT("user/:id",Controllers.UpdateUser)
		grp1.DELETE("user/:id",Controllers.DeleteUser)
	}
	return r
}

