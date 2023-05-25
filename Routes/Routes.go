package route

import(
	"first-api/Controllers"
	"first-api/Models"
	"github.com/gin-gonic/gin"
)

// setup Router .. Config routes

// AVOID redundancy in the names and paths

func SetupRouter(userStore *model.UserStore) *gin.Engine{
	r :=gin.Default()
	grp1:= r.Group("/user-api")
	{
		grp1.GET("user",controller.GetUsers(userStore))
		grp1.POST("user",controller.NewUserController(userStore))
		grp1.GET("user/:id",controller.GetUserByIDController(userStore))
		grp1.PUT("user/:id",controller.UpdateUser(userStore))
		grp1.DELETE("user/:id",controller.DeleteUser(userStore))
	}
	return r
}

// versioning concept 
// authentication
// naming conventions
// config 
// viper(reading configs)
