package main

import (
	"course-project/Initdb"
	_ "course-project/Initdb"
	"course-project/course_arrangement"
	_ "course-project/course_arrangement"
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
	defer db.Close() //延时关闭

	r := gin.Default()

	r.Handle("POST", "/course/create", course_arrangement.CreateCourse)
	r.Handle("POST", "/course/get", course_arrangement.GetCourse)
	r.Handle("GET", "/ping", sayHello)

	//r.Run()
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
