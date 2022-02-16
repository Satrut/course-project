package course_arrangement

import (
	"course-project/Initdb"
	"course-project/types"
)

const maxMessageNum = 20000

var SpikeCourseChannel = make(chan types.BookCourseRequest, maxMessageNum) //有缓存的channel

var isConsumerRun = false

func RunSpikeCourseConsumer() {
	// Only Run one consumer.
	if !isConsumerRun {
		go SpikeCourseConsumer() //开启一个消费者goroutune，作用是接收redis的改动信息，更新数据库
		isConsumerRun = true
	}
}

func SpikeCourseConsumer() {
	for {
		message := <-SpikeCourseChannel
		bookCourse := types.BookCourse{StudentID: message.StudentID, CourseID: message.CourseID}
		db := Initdb.InitDB() //与数据库的链接
		defer db.Close()      //延时关闭
		db.Create(&bookCourse)
	}
}
