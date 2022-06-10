package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {

	token := c.Query("token")
	sToUserId := c.Query("to_user_id")
	toUserId, _ := strconv.ParseInt(sToUserId, 10, 64)
	actionType := c.Query("action_type")

	ConnectionSQL()
	if token == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户信息不存在",
		})
	}
	if actionType == "1" {
		var user User
		GLOBAL_DB.Where("token=?", token).First(&user)
		relation := Relation{
			UserAId: user.Id,
			UserBId: toUserId,
		}
		GLOBAL_DB.Transaction(func(tx *gorm.DB) error {
			// 创建关注关系
			if err := tx.Create(&relation).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			// 增加关注数、被关注数
			tx.Model(&User{}).Where("id=?", user.Id).UpdateColumn("follow_count", gorm.Expr("follow_count + ?", 1))
			tx.Model(&User{}).Where("id=?", toUserId).UpdateColumn("follower_count", gorm.Expr("follower_count + ?", 1))
			// 返回 nil 提交事务
			return nil
		})
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "关注成功",
		})
	} else if actionType == "2" {
		var user User
		GLOBAL_DB.Where("token = ?", token).First(&user)
		GLOBAL_DB.Transaction(func(tx *gorm.DB) error {
			// 删除关注关系
			if err := tx.Where("user_a_id=?", user.Id).Where("user_b_id=?", toUserId).Delete(&Relation{}).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			// 减少关注数、被关注数
			tx.Model(&User{}).Where("id=?", user.Id).UpdateColumn("follow_count", gorm.Expr("follow_count - ?", 1))
			tx.Model(&User{}).Where("id=?", toUserId).UpdateColumn("follower_count", gorm.Expr("follower_count - ?", 1))
			// 返回 nil 提交事务
			return nil
		})
		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "取关成功",
		})
	}
}

func FollowList(c *gin.Context) {

	token := c.Query("token")
	ConnectionSQL()
	if token == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户信息不匹配",
		})
	}
	var user User
	GLOBAL_DB.Where("token=?", token).First(&user)
	var followList []Relation
	// 加载 UserB 即加载当前用户关注的用户
	GLOBAL_DB.Preload("UserB").Where("user_a_id=?", user.Id).Find(&followList)
	// 这里直接暴力复制了，不知道 Go 语言有无更好的方法可以提取结构体数组中的元素
	followUserList := make([]User, len(followList))
	for i, f := range followList {
		followUserList[i] = f.UserB
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followUserList,
	})
}

func FollowerList(c *gin.Context) {
	token := c.Query("token")
	ConnectionSQL()
	if token == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户信息不匹配",
		})
	}
	var user User
	GLOBAL_DB.Where("token=?", token).First(&user)
	var followerList []Relation
	// 加载 UserA 即加载当前用户的粉丝
	GLOBAL_DB.Preload("UserA").Where("user_b_id=?", user.Id).Find(&followerList)

	followerUserList := make([]User, len(followerList))
	for i, f := range followerList {
		followerUserList[i] = f.UserA
	}
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: followerUserList,
	})
}
