package course_arrangement

import (
	"course-project/InitRedis"
	"course-project/Initdb"
	"course-project/types"
	"github.com/gin-gonic/gin"
)

func BookCourse(c *gin.Context) {
	request := types.BookCourseRequest{}
	if c.BindJSON(&request) != nil || request.StudentID == "" || request.CourseID == "" { //如果参数不合法
		response := types.BookCourseResponse{
			Code: types.ParamInvalid,
		}
		c.JSON(200, response)
		return
	}

	res, err := InitRedis.EvalSHA(InitRedis.SpikeCourseSHA, []string{request.StudentID, request.CourseID})
	if err != nil {
		response := types.BookCourseResponse{
			Code: types.UnknownError,
		}
		c.JSON(200, response)
		return
	}

	// 该lua脚本应当返回int值
	BookCourseRes, ok := res.(int64)
	if !ok {
		response := types.BookCourseResponse{
			Code: types.UnknownError,
		}
		c.JSON(200, response)
		return
	}
	switch {
	case BookCourseRes == -1:
		response := types.BookCourseResponse{
			Code: types.CourseNotExisted,
		}
		c.JSON(200, response)
		return
	case BookCourseRes == -2:
		response := types.BookCourseResponse{
			Code: types.CourseNotAvailable,
		}
		c.JSON(200, response)
		return
	case BookCourseRes == -3:
		response := types.BookCourseResponse{
			Code: types.StudentHasCourse,
		}
		c.JSON(200, response)
		return
	case BookCourseRes == 1: //抢课成功
		SpikeCourseChannel <- request //将信息放入channel中，加给消费者goroutine完成数据库更新。
		response := types.BookCourseResponse{
			Code: types.OK,
		}
		c.JSON(200, response)
		return
	default:
		{
			response := types.BookCourseResponse{
				Code: types.UnknownError,
			}
			c.JSON(200, response)
			return
		}
	}
}

func GetStudentCourse(c *gin.Context) {
	StudentID := c.Query("StudentID")
	db := Initdb.InitDB()
	tmember := types.TMember{}
	var tcourses []types.TCourse
	if result := db.Find(&tmember, "user_id = ?", StudentID); result.Error != nil || tmember.UserType != 2 {
		response := types.GetStudentCourseResponse{
			Code: types.StudentNotExisted,
			Data: struct{ CourseList []types.TCourse }{CourseList: tcourses},
		}
		c.JSON(200, response)
		return
	}
	var courseID []types.BookCourse
	db.Where(`student_id = ?`, StudentID).Find(&courseID)
	tcoursess := make([]types.TCourse, len(courseID))
	x := 0
	for _, v := range courseID {
		db.Where("course_id = ?", v.CourseID).Find(&tcoursess[x])
		x += 1
	}
	response := types.GetStudentCourseResponse{
		Code: types.OK,
		Data: struct{ CourseList []types.TCourse }{CourseList: tcoursess},
	}
	c.JSON(200, response)
	return
}
