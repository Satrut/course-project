package course_arrangement

import (
	"course-project/Initdb"
	"course-project/types"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
)

func CreateCourse(c *gin.Context) {
	request := types.CreateCourseRequest{}                 //建立请求的格式（还没有数据）
	if c.BindJSON(&request) != nil || request.Name == "" { //将从前端获取的数据导入到request结构体中，以及如果参数不合法的处理情况
		response := types.CreateCourseResponse{
			Code: types.ParamInvalid,
			Data: struct{ CourseID string }{CourseID: ""},
		}
		c.JSON(200, response)
		return
	}
	tcourse := types.TCourse{Name: request.Name, Cap: request.Cap}
	//tcourse := types.TCourse{CourseID: uuid.NewRandom().String(), Name: request.Name, Cap: request.Cap} //对应数据库中的实体
	db := Initdb.InitDB() //与数据库的链接
	//查找是否已创建
	if result := db.Find(&tcourse, "name = ?", tcourse.Name); result.Error == nil {
		//课程已创建
		response := types.CreateCourseResponse{
			Code: types.ParamInvalid,
			Data: struct{ CourseID string }{CourseID: ""},
		}
		c.JSON(200, response)
	} else {
		tcourse.CourseID = uuid.NewRandom().String()
		db.Create(&tcourse)                     //在数据库中创建实体
		response := types.CreateCourseResponse{ //要返回给前端的数据
			Code: types.OK,
			Data: struct{ CourseID string }{CourseID: tcourse.CourseID},
		}
		c.JSON(200, response) //将数据返回
	}
}

func GetCourse(c *gin.Context) {
	tcourse := types.TCourse{}             //对应数据库中的实体类
	tcourse.CourseID = c.Query("CourseID") //从前端获取数据,get方法
	db := Initdb.InitDB()                  //与数据库的链接
	//if result := db.First(&tcourse); result.Error != nil { //从数据库中进行数据查找，这是如果没有成功找到的情况
	if result := db.Find(&tcourse, "course_id = ?", tcourse.CourseID); result.Error != nil {
		response := types.GetCourseResponse{
			Code: types.CourseNotExisted,
			Data: tcourse,
		}
		c.JSON(200, response)
	} else { //成功找到的情况
		response := types.GetCourseResponse{
			Code: types.OK,
			Data: tcourse,
		}
		c.JSON(200, response)
	}
}

func BindCourse(c *gin.Context) {
	request := types.BindCourseRequest{}
	if c.BindJSON(&request) != nil || request.TeacherID == "" || request.CourseID == "" { //如果参数不合法
		response := types.BindCourseResponse{
			Code: types.ParamInvalid,
		}
		c.JSON(200, response)
		return
	}
	tcourse := types.TCourse{CourseID: request.CourseID}
	db := Initdb.InitDB()
	if result := db.First(&tcourse); result.Error != nil { //如果课程不存在
		response := types.BindCourseResponse{
			Code: types.CourseNotExisted,
		}
		c.JSON(200, response)
		return
	} else if tcourse.TeacherID != "" { //如果课程已经被绑定过
		response := types.BindCourseResponse{
			Code: types.CourseHasBound,
		}
		c.JSON(200, response)
		return
	} else {
		db.Model(&tcourse).Update("TeacherID", request.TeacherID)
		response := types.BindCourseResponse{
			Code: types.OK,
		}
		c.JSON(200, response)
	}
}

func UnbindCourse(c *gin.Context) {
	request := types.UnbindCourseRequest{}
	if c.BindJSON(&request) != nil || request.TeacherID == "" || request.CourseID == "" { //如果参数不合法
		response := types.UnbindCourseResponse{
			Code: types.ParamInvalid,
		}
		c.JSON(200, response)
		return
	}
	tcourse := types.TCourse{CourseID: request.CourseID}
	db := Initdb.InitDB()
	if result := db.First(&tcourse); result.Error != nil { //如果课程不存在
		response := types.UnbindCourseResponse{
			Code: types.CourseNotExisted,
		}
		c.JSON(200, response)
		return
	} else if tcourse.TeacherID != "" && tcourse.TeacherID != request.TeacherID { //如果课程已经被绑定过
		response := types.UnbindCourseResponse{
			Code: types.UnknownError,
		}
		c.JSON(200, response)
		return
	} else if tcourse.TeacherID == "" {
		response := types.UnbindCourseResponse{
			Code: types.CourseNotBind,
		}
		c.JSON(200, response)
	} else {
		db.Model(&tcourse).Update("TeacherID", "")
		response := types.UnbindCourseResponse{
			Code: types.OK,
		}
		c.JSON(200, response)
	}
}

func GetTeacherCourse(c *gin.Context) {
	teacher_id := c.Query("TeacherID")
	var tcourses []*types.TCourse
	if teacher_id == "" { //如果参数不合法
		response := types.GetTeacherCourseResponse{
			Code: types.ParamInvalid,
			Data: struct{ CourseList []*types.TCourse }{CourseList: tcourses},
		}
		c.JSON(200, response)
		return
	}
	db := Initdb.InitDB()
	db.Where(`teacher_id = ?`, teacher_id).Find(&tcourses)
	response := types.GetTeacherCourseResponse{
		Code: types.OK,
		Data: struct{ CourseList []*types.TCourse }{CourseList: tcourses},
	}
	c.JSON(200, response)
}

//排课求解器
func ScheduleCourse(c *gin.Context) {
	request := types.ScheduleCourseRequest{}
	res := make(map[string]string)
	if c.BindJSON(&request) != nil { //如果参数不合法
		response := types.ScheduleCourseResponse{
			Code: types.ParamInvalid,
			Data: res,
		}
		c.JSON(200, response)
		return
	}
	relationShip := request.TeacherCourseRelationShip
	used := make(map[string]bool)
	pre := make(map[string]string)
	//遍历原始数据
	for k, v := range relationShip { //不一定是插入顺序
		res[k] = ""
		for _, courseID := range v {
			used[courseID] = false
			pre[courseID] = ""
		}
	}
	for teacherID, _ := range relationShip {
		for courseID, _ := range used {
			used[courseID] = false
		}
		dfs(relationShip, res, used, pre, teacherID)
	}
	response := types.ScheduleCourseResponse{
		Code: types.OK,
		Data: res,
	}
	c.JSON(200, response)
}

func dfs(relationShip map[string][]string, res map[string]string, used map[string]bool, pre map[string]string, teacherID string) bool {
	courseList := relationShip[teacherID]
	for _, courseID := range courseList {
		if !used[courseID] {
			used[courseID] = true
			if pre[courseID] == "" || dfs(relationShip, res, used, pre, pre[courseID]) {
				pre[courseID] = teacherID
				res[teacherID] = courseID
				return true
			}
		} else {
			continue
		}
	}
	return false
}
