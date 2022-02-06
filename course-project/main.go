package main

import (
	"course-project/Initdb"
	"course-project/types_1"
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
	types_1.RegisterRouter(r)

	//r.Run()
	panic(r.Run()) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
