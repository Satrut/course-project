package course_arrangement

import (
	"course-project/Initdb"
	"course-project/types"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
)

func CreateCourse(c *gin.Context) {
	request := types.CreateCourseRequest{}
	c.BindJSON(&request)
	tcourse := types.TCourse{CourseID: uuid.NewRandom().String(), Name: request.Name, Cap: request.Cap}
	db := Initdb.InitDB()
	db.Create(&tcourse)
	response := types.CreateCourseResponse{
		Code: types.OK,
		Data: struct{ CourseID string }{CourseID: tcourse.CourseID},
	}
	c.JSON(200, response)
}

func GetCourse(c *gin.Context) {
	tcourse := types.TCourse{}
	tcourse.CourseID = c.Query("CourseID")
	db := Initdb.InitDB()
	if result := db.First(&tcourse); result.Error != nil {
		response := types.GetCourseResponse{
			Code: types.CourseNotExisted,
			Data: tcourse,
		}
		c.JSON(200, response)
	} else {
		response := types.GetCourseResponse{
			Code: types.OK,
			Data: tcourse,
		}
		c.JSON(200, response)
	}
}
