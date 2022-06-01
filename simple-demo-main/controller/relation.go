package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserListResponse struct {
	Response
	UserList []User `json:"user_list"`
}

// RelationAction no practical effect, just check if token is valid
func RelationAction(c *gin.Context) {
	//新建一个relations表
	//表中要有userId、toUserId
	//若actionType为1，在relation表中添加数据,同时在user表使userid对应用户的followCount + 1，toUserId对应用户的followerCount + 1
	//反之，在relation表中删除一条数据（取关），同时在user表使userid对应用户的followCount - 1，toUserId对应用户的followerCount - 1
	//问题：isFollow怎么解决，不太理解
	//用线程锁解决并发问题，用户在同时点击时不能出现数据错误

	//原demo
	token := c.Query("token")

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {
	//用userId作为relation表中的userId，查询toUserId数据，作为userList返回
	//问题：返回的列表里是不是应该存放username而不是userid，如何实现。

	//原demo
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	//用userId作为relation表中的toUserId，查询UserId数据，作为userList返回
	//问题：返回的列表里是不是应该存放username而不是userid，如何实现。

	//原demo
	c.JSON(http.StatusOK, UserListResponse{
		Response: Response{
			StatusCode: 0,
		},
		UserList: []User{DemoUser},
	})
}
