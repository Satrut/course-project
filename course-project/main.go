package main

import (
	"course-project/Initdb"
	"course-project/course_arrangement"
	"course-project/login"
	"course-project/member"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func sayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello golang!",
	})
}

func main() {
	db := Initdb.InitDB()
	//tmember := types.TMember{
	//	Nickname:   "admin",
	//	Username:   "JudgeAdmin",
	//	Password:   "JudgePassword2022",
	//	UserType:   types.Admin,
	//	UserStatus: true,
	//	UserID:     uuid.NewRandom().String(),
	//}
	//db.Create(&tmember)
	defer db.Close() //延时关闭

	r := gin.Default()
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create", member.CreateMember)
	g.GET("/member", member.GetMember)
	g.GET("/member/list", member.GetMemberList)
	g.POST("/member/update", member.UpdateMember)
	g.POST("/member/delete", member.DeleteMember)

	// 登录

	g.POST("/auth/login", login.Login)
	g.POST("/auth/logout", login.Logout)
	g.GET("/auth/whoami", login.WhoAmI)

	// 排课
	g.POST("/course/create", course_arrangement.CreateCourse)
	g.GET("/course/get", course_arrangement.GetCourse)

	g.POST("/teacher/bind_course")
	g.POST("/teacher/unbind_course")
	g.GET("/teacher/get_course")
	g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")

	//r.Run()
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
