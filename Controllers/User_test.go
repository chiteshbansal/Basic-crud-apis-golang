package Controllers

import (
	"first-api/Models"
	"testing"
	"github.com/stretchr/testify/assert"
)



func TestValidateUserData(t *testing.T){
	assert  := assert.New(t)

	testService := Models.User{Name:"",Email:"test@gmail.com",Phone:"9999999999",Address:"abcd efgh ijkl"}

	err := testService.Validate();
	// fmt.Println(err)

	assert.NotNil(t,err);

	testService.Name = "testing"
	testService.Email = "test"
	err = testService.Validate();
	// fmt.Println(err)

	assert.NotNil(t,err);


	testService.Phone = "999999999" 
	err = testService.Validate();
	// fmt.Println(err)

	assert.NotNil(t,err);

	testService.Address = ""
	err = testService.Validate()
	// fmt.Println(err)

	assert.NotNil(t,err);

}

