package Controllers
import(
	"first-api/Models"
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)


// get all users

func GetUsers(c *gin.Context){
	var user []Models.User
	err := Models.GetAllUsers(&user);
	if err!= nil{
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}else{
		c.JSON(http.StatusOK,user)
	}

}

func CreateUser(c *gin.Context){
	var user Models.User
	c.BindJSON(&user)
	valErr := user.Validate();
	if valErr != nil {
		fmt.Println(valErr) 
		c.JSON(http.StatusNotFound,valErr)
		return 
	}

	err := Models.CreateUser(&user)
	fmt.Println("user:",user);
	if err!=nil {
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}else{
		fmt.Println(user.Name)
		c.JSON(http.StatusOK,user);
	}
}

// update user data
func UpdateUser(c *gin.Context){
	id:= c.Params.ByName("id")
	var user Models.User
	err:=Models.GetUserByID(&user,id)
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
	err = Models.UpdateUser(&user,id)

	if(err !=nil){
		c.AbortWithStatus(http.StatusNotFound)
	}else{
		c.JSON(http.StatusOK,user)
	}

}

func GetUserByID(c *gin.Context){
	id := c.Params.ByName("id")
	var user Models.User

	err:= Models.GetUserByID(&user,id);
	if(err !=nil){
		fmt.Println(err.Error())
		c.AbortWithStatus(http.StatusNotFound)
	}else{
		c.JSON(http.StatusOK,user)
	}
}


// delete user

func DeleteUser (c *gin.Context){
	var user Models.User
	id := c.Params.ByName("id")

	err:= Models.GetUserByID(&user,id);
	if(err !=nil){
		fmt.Println(err.Error())
		  c.JSON(http.StatusNotFound,"REcord not found")
		  return 
	}

	c.BindJSON(&user)
	err = Models.DeleteUser(&user,id)

	if err!=nil{
		c.AbortWithStatus(http.StatusNotFound)
	}else{
		c.JSON(http.StatusOK,gin.H{"id" + id: "is deleted"})
	}
}