package course_arrangement

import (
	"course-project/InitRedis"
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
