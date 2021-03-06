package member

import (
	"course-project/Initdb"
	"course-project/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

//g.POST("/member/create")
//g.GET("/member")
//g.GET("/member/list")
//g.POST("/member/update")
//g.POST("/member/delete")

// 获取成员信息，如果用户已删除请返回已删除状态码，不存在请返回不存在状态码

func GetMember(c *gin.Context) {
	tmember := types.TMember{}
	tmember.UserID = c.Query("UserID")
	db := Initdb.InitDB()
	defer db.Close() //延时关闭
	if result := db.First(&tmember); result.Error != nil { //用户不存在
		response := types.GetMemberResponse{
			Code: types.UserNotExisted,
			Data: tmember,
		}
		c.JSON(200, response)
	} else {
		if tmember.UserStatus { //用户存在
			response := types.GetMemberResponse{
				Code: types.OK,
				Data: tmember,
			}
			c.JSON(200, response)
		} else { //用户已删除
			response := types.GetMemberResponse{
				Code: types.UserHasDeleted,
				Data: tmember,
			}
			c.JSON(200, response)
		}
	}
}

// 批量获取成员信息

func GetMemberList(c *gin.Context) {
	request := types.GetMemberListRequest{}
	request.Offset, _ = strconv.Atoi(c.Query("Offset"))
	request.Limit, _ = strconv.Atoi(c.Query("Limit"))
	tmemberlist := []types.TMember{}
	db := Initdb.InitDB()
	defer db.Close() //延时关闭
	//if result := db.Find(&tmemberlist, "UserStatus=?", true); result.Error != nil {
	if result := db.Offset(request.Offset).Limit(request.Limit).Find(&tmemberlist, "user_status=?", true); result.Error != nil {
		response := types.GetMemberListResponse{
			Code: types.UserNotExisted,
			Data: struct{ MemberList []types.TMember }{
				MemberList: tmemberlist,
			},
		}
		//response.Data.MemberList = tmemberlist
		c.JSON(200, response)
	} else {
		response := types.GetMemberListResponse{
			Code: types.OK,
			Data: struct{ MemberList []types.TMember }{
				MemberList: tmemberlist,
			},
		}
		//response.Data.MemberList = tmemberlist
		c.JSON(200, response)
	}
}

//删除成员信息

func DeleteMember(c *gin.Context) {
	//request := types.DeleteMemberRequest{}

	// 判断UserID是否存在
	tmember := types.TMember{}
	//tmember.UserID = c.Query("UserID")
	request := types.DeleteMemberRequest{}
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	tmember.UserID = request.UserID
	db := Initdb.InitDB()
	defer db.Close() //延时关闭
	if result := db.First(&tmember); result.Error != nil {
		response := types.DeleteMemberResponse{
			Code: types.UserNotExisted,
		}
		c.JSON(200, response)
		return
	} else {
		if tmember.UserStatus { //用户存在,删除，返回OK
			//db.Delete(tmember)
			db.Model(&tmember).Update("user_status", false)
			response := types.DeleteMemberResponse{
				Code: types.OK,
			}
			c.JSON(200, response)
		} else { //用户已删除
			response := types.DeleteMemberResponse{
				Code: types.UserHasDeleted,
			}
			c.JSON(200, response)
		}
	}
}
