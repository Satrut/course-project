package course_arrangement

import (
	"course-project/Initdb"
	"course-project/types_2"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
)

func CreateCourse(c *gin.Context) {
	request := types_2.CreateCourseRequest{}
	c.BindJSON(&request)
	tcourse := types_2.TCourse{CourseID: uuid.NewRandom().String(), Name: request.Name, Cap: request.Cap}
	db := Initdb.InitDB()
	db.Create(&tcourse)
	response := types_2.CreateCourseResponse{
		Code: types_2.OK,
		Data: struct{ CourseID string }{CourseID: tcourse.CourseID},
	}
	c.JSON(200, response)
}

func GetCourse(c *gin.Context) {
	request := types_2.GetCourseRequest{}
	c.BindJSON(&request)
	tcourse := types_2.TCourse{CourseID: request.CourseID}
	db := Initdb.InitDB()
	if result := db.First(&tcourse); result.Error != nil {
		response := types_2.GetCourseResponse{
			Code: types_2.CourseNotExisted,
			Data: tcourse,
		}
		c.JSON(200, response)
	} else {
		response := types_2.GetCourseResponse{
			Code: types_2.OK,
			Data: tcourse,
		}
		c.JSON(200, response)
	}
}
