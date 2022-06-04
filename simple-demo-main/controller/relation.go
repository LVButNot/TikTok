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

	userId := c.PostForm("user_id")
	token := c.PostForm("token")
	toUserId := c.PostForm("to_user_id")
	actionType := c.PostForm("action_type")

	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个UserInfoTab
	var user User
	//通过userId在User表里查询token并与请求参数中的token比对
	verification := GLOBAL_DB.Select("token").Where("id = ?", userId).Find(&user)
	if verification.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "服务器内部错误"})
	}
	if user.Token != token {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户不存在"})
	} else {
		//新建一个relations表
		//表中要有userId、toUserId
		//若actionType为1，在relation表中添加数据,同时在user表使userid对应用户的followCount + 1，toUserId对应用户的followerCount + 1
		//反之，在relation表中删除一条数据（取关），同时在user表使userid对应用户的followCount - 1，toUserId对应用户的followerCount - 1
		//问题：isFollow怎么解决，不太理解
		//用线程锁解决并发问题，用户在同时点击时不能出现数据错误
		err := ConnectionRedis2()
		if err != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "连接Redis数据库出错"})
		}

		err1 := ConnectionRedis3()
		if err1 != nil {
			c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "连接Redis数据库出错"})
		}
		if actionType == "1" {
			RDB2.SAdd(userId, toUserId)
			RDB3.SAdd(toUserId, userId)
			//在user表中使userid的followCount + 1,to_user_id的followerCount + 1
			user := User{}
			GLOBAL_DB.Where("id = ", userId).Take(&user)
			user.FollowCount += 1
			GLOBAL_DB.Save(user)

			toUser := User{}
			GLOBAL_DB.Where("id = ", toUserId).Take(&toUser)
			toUser.FollowerCount += 1
			GLOBAL_DB.Save(toUser)

			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "关注成功"})
		} else if actionType == "2" {
			RDB2.SRem(userId, toUserId)
			RDB3.SRem(toUserId, userId)
			//在user表中使userid的followCount - 1,to_user_id的followerCount - 1
			user := User{}
			GLOBAL_DB.Where("id = ", userId).Take(&user)
			user.FollowCount -= 1
			GLOBAL_DB.Save(user)

			toUser := User{}
			GLOBAL_DB.Where("id = ", toUserId).Take(&toUser)
			toUser.FollowerCount -= 1
			GLOBAL_DB.Save(toUser)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "取关成功"})
		}

	}
}

// FollowList all users have same follow list
func FollowList(c *gin.Context) {

	userId := c.Query("user_id")
	token := c.Query("token")

	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个User和一个VideoListResponse
	var user User
	ulr := UserListResponse{Response{1, ""}, []User{}}
	//通过userId在UserInfoTab表里查询token并与请求参数中的token比对
	verification := GLOBAL_DB.Select("token").Where("id = ?", userId).Find(&user)
	if verification.Error != nil {
		ulr.Response.StatusCode = 1
		ulr.Response.StatusMsg = "服务器内部错误"
		c.JSON(http.StatusInternalServerError, &ulr)
	}
	if user.Token != token {
		ulr.Response.StatusCode = 1
		ulr.Response.StatusMsg = "用户不存在"
		c.JSON(http.StatusOK, &ulr)
	} else {
		//用userId作为relation表中的userId，查询toUserId数据，作为userList返回
		err := ConnectionRedis2()
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserListResponse{Response{StatusCode: 1, StatusMsg: "Redis连接失败"}, []User{}})
		}
		finds, _ := RDB2.SMembers(userId).Result()
		var userList []User
		for _, find := range finds {
			var user User
			verification = GLOBAL_DB.Select("name").Where("id = ", find).Take(&user)
			if verification != nil {
				ulr.Response.StatusCode = 1
				ulr.Response.StatusMsg = "查询出现错误"
				c.JSON(http.StatusInternalServerError, &ulr)
			}

			userList = append(userList, user)
		}
		c.JSON(http.StatusOK, UserListResponse{Response{StatusCode: 0, StatusMsg: "拉取成功"}, userList})
	}
}

// FollowerList all users have same follower list
func FollowerList(c *gin.Context) {
	//查询toUserId数据，作为userList返回
	//问题：返回的列表里是不是应该存放username而不是userid，如何实现。

	userId := c.Query("user_id")
	token := c.Query("token")

	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个User和一个VideoListResponse
	var user User
	ulr := UserListResponse{Response{1, ""}, []User{}}
	//通过userId在UserInfoTab表里查询token并与请求参数中的token比对
	verification := GLOBAL_DB.Select("token").Where("id = ?", userId).Find(&user)
	if verification.Error != nil {
		ulr.Response.StatusCode = 1
		ulr.Response.StatusMsg = "服务器内部错误"
		c.JSON(http.StatusInternalServerError, &ulr)
	}
	if user.Token != token {
		ulr.Response.StatusCode = 1
		ulr.Response.StatusMsg = "用户不存在"
		c.JSON(http.StatusOK, &ulr)
	} else {
		//查询toUserId数据，作为userList返回
		err := ConnectionRedis3()
		if err != nil {
			c.JSON(http.StatusInternalServerError, UserListResponse{Response{StatusCode: 1, StatusMsg: "Redis连接失败"}, []User{}})
		}
		finds, _ := RDB3.SMembers(userId).Result()
		var userList []User
		for _, find := range finds {
			var user User
			verification = GLOBAL_DB.Select("name").Where("id = ", find).Take(&user)
			if verification != nil {
				ulr.Response.StatusCode = 1
				ulr.Response.StatusMsg = "查询出现错误"
				c.JSON(http.StatusInternalServerError, &ulr)
			}

			userList = append(userList, user)
		}
		c.JSON(http.StatusOK, UserListResponse{Response{StatusCode: 0, StatusMsg: "拉取成功"}, userList})
	}
}
