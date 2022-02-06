package course_arrangement

import (
	"course-project/Initdb"
	"course-project/type"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
)

func CreateCourse(c *gin.Context) {
	request := _type.CreateCourseRequest{}
	c.BindJSON(&request)
	tcourse := _type.TCourse{CourseID: uuid.NewRandom().String(), Name: request.Name, Cap: request.Cap}
	db := Initdb.InitDB()
	db.Create(&tcourse)
	response := _type.CreateCourseResponse{
		Code: _type.OK,
		Data: struct{ CourseID string }{CourseID: tcourse.CourseID},
	}
	c.JSON(200, response)
}

func GetCourse(c *gin.Context) {
	request := _type.GetCourseRequest{}
	c.BindJSON(&request)
	tcourse := _type.TCourse{}
	tcourse.CourseID = c.Query("CourseID")
	db := Initdb.InitDB()
	if result := db.First(&tcourse); result.Error != nil {
		response := _type.GetCourseResponse{
			Code: _type.CourseNotExisted,
			Data: tcourse,
		}
		c.JSON(200, response)
	} else {
		response := _type.GetCourseResponse{
			Code: _type.OK,
			Data: tcourse,
		}
		c.JSON(200, response)
	}
}
