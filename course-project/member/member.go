package member

import (
	"course-project/Initdb"
	"course-project/types"
	"github.com/gin-gonic/gin"
)

//g.POST("/member/create")
//g.GET("/member")
//g.GET("/member/list")
//g.POST("/member/update")
//g.POST("/member/delete")

func GetMember(c *gin.Context) {
	tmember := types.TMember{}
	tmember.UserID = c.Query("UserID")
	db := Initdb.InitDB()
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

func DeleteMember(c *gin.Context) {
	//request := types.DeleteMemberRequest{}

	// 判断UserID是否存在
	tmember := types.TMember{}
	tmember.UserID = c.Query("UserID")
	db := Initdb.InitDB()
	if result := db.First(&tmember); result.Error != nil {
		response := types.DeleteMemberResponse{
			Code: types.UserNotExisted,
		}
		c.JSON(200, response)
		return
	} else {
		if tmember.UserStatus { //用户存在,删除，返回OK
			db.Delete(tmember)
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
