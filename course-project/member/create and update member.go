package member

import (
	"course-project/Initdb"
	"course-project/types"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"unicode"
)

// 只有管理员才有操作权限

// CreateMember 添加成员
func CreateMember(c *gin.Context) {
	// 判断参数是否合法
	request := types.CreateMemberRequest{}
	if c.BindJSON(&request) != nil || len(request.Nickname) < 4 || len(request.Nickname) > 20 || !JudgeUsername(request.Username) || !JudgePassWord(request.Password) {
		response := types.CreateMemberResponse{
			Code: types.ParamInvalid,
			Data: struct{ UserID string }{UserID: ""},
		}
		c.JSON(200, response)
		return
	}

	// 判断UserName是否存在
	tmember := types.TMember{}
	tmember.Username = c.Query("Username")
	db := Initdb.InitDB()
	if result := db.First(&tmember); result.Error == nil {
		response := types.CreateMemberResponse{
			Code: types.UserHasExisted,
			Data: struct{ UserID string }{UserID: ""},
		}
		c.JSON(200, response)
		return
	}

	// 添加成员
	tmember.UserID = uuid.NewRandom().String()
	tmember.Nickname = request.Nickname
	tmember.Username = request.Username
	tmember.UserType = request.UserType

	db.Create(&tmember)
	response := types.CreateMemberResponse{
		Code: types.OK,
		Data: struct{ UserID string }{UserID: tmember.UserID},
	}
	c.JSON(200, response)
}

func JudgeUsername(s string) bool {
	cnt := 0
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
		cnt++
	}
	if cnt >= 8 && cnt <= 20 {
		return true
	}
	return false
}

func JudgePassWord(s string) bool {
	upper, lower, number := false, false, false
	cnt := 0
	for _, r := range s {
		if unicode.IsLower(r) {
			lower = true
		} else if unicode.IsUpper(r) {
			upper = true
		} else if unicode.IsNumber(r) {
			number = true
		} else {
			return false
		}
		cnt++
	}
	if cnt >= 8 && cnt <= 20 && upper && lower && number {
		return true
	}
	return false
}

// UpdateMember  更新成员
func UpdateMember(c *gin.Context) {
	request := types.UpdateMemberRequest{}
	if c.BindJSON(&request) != nil {
		response := types.UpdateMemberResponse{
			Code: types.ParamInvalid,
		}
		c.JSON(200, response)
	}

	// 判断UserID是否存在
	tmember := types.TMember{}
	tmember.Username = c.Query("UserID")
	db := Initdb.InitDB()
	if result := db.First(&tmember); result.Error != nil {
		response := types.UpdateMemberResponse{
			Code: types.UserNotExisted,
		}
		c.JSON(200, response)
		return
	}
	db.Model(&tmember).Update("Nickname", request.Nickname)
	response := types.UpdateMemberResponse{
		Code: types.OK,
	}
	c.JSON(200, response)
}
