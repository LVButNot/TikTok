package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []Video `json:"video_list,omitempty"`
	NextTime  int64   `json:"next_time,omitempty"`
}

type VideoList struct {
	FeedResponse
}

// Feed same demo video list for every request
//不限制登录状态，返回投稿时间倒序的视频列表，视频数由服务端控制，单次最多30个
//总流程：无需验证登录状态，从video信息表里取出视频信息，并以创建时间（投稿时间）为序，生成视频列表，并返回。
func Feed(c *gin.Context) {

	//连接数据库，utils.go工具包
	ConnectionSQL()

	token := c.Query("token")
	var user User
	GLOBAL_DB.Where("token = ?", token).First(&user)

	var videoList []Video
	//从video数据库表里取出视频信息
	GLOBAL_DB.Preload("Author").Order("id desc").Limit(30).Find(&videoList)

	for i, v := range videoList {
		// 查找是否存在一条当前用户给该视频点赞的记录
		rows := GLOBAL_DB.Where("video_id=?", v.Id).Where("user_id=?", user.Id).Find(&Favorite{}).RowsAffected
		videoList[i].IsFavorite = rows > 0 // 查找到了点赞记录
	}

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: videoList,
		NextTime:  time.Now().Unix()})
}
