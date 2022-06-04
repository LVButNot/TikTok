package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	userId := c.PostForm("user_id")
	token := c.PostForm("token")
	videoId := c.PostForm("video_id")
	actionType := c.PostForm("action_type")

	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个User
	var user User
	//通过userId在UserInfoTab表里查询token并与请求参数中的token比对
	verification := GLOBAL_DB.Select("token").Where("id = ?", userId).Find(&user)
	if verification.Error != nil {
		c.JSON(http.StatusInternalServerError, Response{StatusCode: 1, StatusMsg: "服务器内部错误"})
	}
	if user.Token != token {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "用户不存在"})
	} else {
		/*两张表：videos表和favorite表
		 *favorite表存放用户id和视频id
		 *若用户未曾点赞，则在favorite表中增加一条数据，同时在video表中使favoriteCount + 1
		 *若点赞过，则删除这条数据（即取消点赞），同时在video表中使favoriteCount — 1
		 *
		 *需要解决的问题：
		 *用线程锁解决并发问题，用户在同时点击时不能出现增加或减少等数据错误
		 *在video表中还有一列is_favorite数据，如何处理
		 */
		err := ConnectionRedis()
		if err != nil {
			c.JSON(http.StatusInternalServerError, &Response{StatusCode: 1, StatusMsg: "Redis查询失败"})
		}
		if actionType == "1" {
			//在favorite表中增加一条数据
			RDB.SAdd(userId, videoId)

			//在video表中使favoriteCount + 1
			video := Video{}
			GLOBAL_DB.Where("id = ", videoId).Take(&video)
			video.FavoriteCount += 1
			GLOBAL_DB.Save(video)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "点赞成功"})
		} else if actionType == "2" {
			//在favorite表中删除一条数据
			RDB.SRem(userId, videoId)

			//在video表中使favoriteCount - 1
			video := Video{}
			GLOBAL_DB.Where("id = ", videoId).Take(&video)
			video.FavoriteCount -= 1
			GLOBAL_DB.Save(video)
			c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "取消成功"})
		}
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	userId := c.Query("user_id")
	token := c.Query("token")

	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个UserInfoTab和一个VideoListResponse
	var user User
	vlr := VideoListResponse{Response{1, ""}, []Video{}}
	//通过userId在UserInfoTab表里查询token并与请求参数中的token比对
	verification := GLOBAL_DB.Select("token").Where("id = ?", userId).Find(&user)
	if verification.Error != nil {
		vlr.Response.StatusCode = 1
		vlr.Response.StatusMsg = "服务器内部错误"
		c.JSON(http.StatusInternalServerError, &vlr)
	}
	if user.Token != token {
		vlr.Response.StatusCode = 1
		vlr.Response.StatusMsg = "用户不存在"
		c.JSON(http.StatusOK, &vlr)
	} else {
		//favorite表中通过userid查询到当前用户所有的点赞视频并作为响应返回
		//修改一下video表，该表目前还缺一个user数据，外键关联。
		err := ConnectionRedis()
		if err != nil {
			vlr.Response.StatusCode = 1
			vlr.Response.StatusMsg = "Redis查询失败"
			c.JSON(http.StatusInternalServerError, &vlr)
		}

		finds, err := RDB.SMembers(userId).Result()
		var videoList []Video
		for _, find := range finds {
			var video Video
			verification := GLOBAL_DB.Where("id = ?", find).Find(&video)
			if verification.Error != nil {
				vlr.Response.StatusCode = 1
				vlr.Response.StatusMsg = "服务器内部错误"
				c.JSON(http.StatusInternalServerError, &vlr)
			}

			videoList = append(videoList, video)
		}
		vlr.Response.StatusCode = 0
		vlr.Response.StatusMsg = "拉取成功"
		vlr.VideoList = videoList
		c.JSON(http.StatusOK, &vlr)
	}
}
