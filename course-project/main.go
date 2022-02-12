package main

import (
	"course-project/Initdb"
	"course-project/course_arrangement"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"net/http"
)

func sayHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello golang!",
	})
}

func main() {
	db := Initdb.InitDB()
	defer db.Close() //延时关闭

	r := gin.Default()
	g := r.Group("/api/v1")

	// 成员管理
	g.POST("/member/create")
	g.GET("/member")
	g.GET("/member/list")
	g.POST("/member/update")
	g.POST("/member/delete")

	// 登录

	g.POST("/auth/login")
	g.POST("/auth/logout")
	g.GET("/auth/whoami")

	// 排课
	g.POST("/course/create", course_arrangement.CreateCourse)
	g.GET("/course/get", course_arrangement.GetCourse)

	g.POST("/teacher/bind_course", course_arrangement.BindCourse)
	g.POST("/teacher/unbind_course", course_arrangement.UnbindCourse)
	g.GET("/teacher/get_course", course_arrangement.GetTeacherCourse)
	g.POST("/course/schedule")

	// 抢课
	g.POST("/student/book_course")
	g.GET("/student/course")

	g.GET("/hello", func(c *gin.Context) {
		name := c.Query("name")
		c.String(http.StatusOK, fmt.Sprintf("hello %s", name))
	})
	r.Handle("GET", "/say_hello", func(c *gin.Context) {
		name := c.Query("name")
		c.String(http.StatusOK, fmt.Sprintf("hello %s", name))
	})
	//r.Run()
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
