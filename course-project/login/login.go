package login

import (
	"course-project/Initdb"
	"course-project/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	request := types.LoginRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tmember := types.TMember{Username: request.Username, Password: request.Password}
	db := Initdb.InitDB()
	defer db.Close() //延时关闭
	//if result := db.First(&tmember); result.Error != nil {
	//查找用户并判断是否已被删除
	if result := db.Where("username = ? AND password = ?", tmember.Username, tmember.Password).Find(&tmember); result.Error != nil || tmember.UserStatus == false {
		//登录失败
		response := types.LoginResponse{
			Code: types.WrongPassword,
			Data: struct{ UserID string }{UserID: ""},
		}
		c.JSON(http.StatusOK, response)
	} else {
		//登录成功，设置cookie
		c.SetCookie("camp-session", tmember.UserID, 7*24*60*60, "/", "", false, true)
		response := types.LoginResponse{
			Code: types.OK,
			Data: struct{ UserID string }{UserID: tmember.UserID},
		}
		c.JSON(http.StatusOK, response)
	}
}

func Logout(c *gin.Context) {
	//设置maxAge<0即可
	c.SetCookie("camp-session", "", -1, "/", "", false, true)
	response := types.LogoutResponse{Code: types.OK}
	c.JSON(http.StatusOK, response)
}

func WhoAmI(c *gin.Context) {
	tmember := types.TMember{}
	campSession, err := c.Cookie("camp-session")
	if err != nil {
		//没有cookie
		response := types.WhoAmIResponse{
			Code: types.LoginRequired,
			Data: tmember,
		}
		c.JSON(http.StatusOK, response)
	} else {
		db := Initdb.InitDB()
		defer db.Close() //延时关闭
		db.Where("user_id = ?", campSession).Find(&tmember)
		response := types.WhoAmIResponse{
			Code: types.OK,
			Data: tmember,
		}
		c.JSON(http.StatusOK, response)
	}
}
