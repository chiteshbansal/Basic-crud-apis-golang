package Controllers
import(
	"first-api/Models"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)


// get all users



// type UserStore interface {
// 	CreateUser(user *Models.User) error
// 	Validate(user *models.User) error
// }

type UserStorer interface{
	CreateUser(user *Models.User) error
	Validate(user Models.User) error
	GetAllUsers(users *[]Models.User) error
	GetUserByID(user *Models.User ,id string ) error
	UpdateUser(user *Models.User ,id string) error
	DeleteUser(user *Models.User, id string) error 
}

func GetUserByIDController(store UserStorer) gin.HandlerFunc{
	return func (c *gin.Context){
		id := c.Params.ByName("id")
		var user Models.User

		err:= store.GetUserByID(&user,id);
		if(err !=nil){
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)
		}else{
			c.JSON(http.StatusOK,user)
		}
	}
}

func GetUsers(store UserStorer) gin.HandlerFunc{
	return func (c *gin.Context){
		var users []Models.User
		err := store.GetAllUsers(&users);
		if err!= nil{
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)
		}else{
			c.JSON(http.StatusOK,users)
		}
	
	}
}


func NewUserController(store UserStorer) gin.HandlerFunc {
	return func(c *gin.Context) {
		var user Models.User
		c.BindJSON(&user)

		valErr := store.Validate(user)
		if valErr != nil {
			fmt.Println(valErr) 
			c.JSON(http.StatusNotFound, valErr)
			return 
		}

		err := store.CreateUser(&user)
		if err != nil {
			fmt.Println(err.Error())
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			fmt.Println(user.Name)
			c.JSON(http.StatusOK, user)
		}
	}
}
// update user data
func UpdateUser(store UserStorer) gin.HandlerFunc{
	return func (c *gin.Context){
	id:= c.Params.ByName("id")
	var user Models.User
	err:=store.GetUserByID(&user,id)
	if err!=nil{
		fmt.Println(err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}


	c.BindJSON(&user)// putting the user data from the body to the user variable
	valErr := user.Validate();
	if valErr != nil {
		fmt.Println(valErr) 
		c.JSON(http.StatusNotFound,valErr)
		return 
	}
	err = store.UpdateUser(&user,id)

	if(err !=nil){
		c.AbortWithStatus(http.StatusNotFound)
	}else{
		c.JSON(http.StatusOK,user)
	}

}}
// delete user

func DeleteUser (store UserStorer) gin.HandlerFunc{
	return func(c *gin.Context){
	var user Models.User
	id := c.Params.ByName("id")

	err:= store.GetUserByID(&user,id);
	if(err !=nil){
		fmt.Println(err.Error())
		  c.JSON(http.StatusNotFound,"REcord not found")
		  return 
	}

	c.BindJSON(&user)
	err = store.DeleteUser(&user,id)

	if err!=nil{
		c.AbortWithStatus(http.StatusNotFound)
	}else{
		c.JSON(http.StatusOK,gin.H{"id" + id: "is deleted"})
	}
}

}



