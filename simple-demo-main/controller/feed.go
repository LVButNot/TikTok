package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
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

	//获取请求参数，在视频流中这不是必填信息

	// 	latest_time := c.Query("latest_time")
	// 	token := c.Query("token")

	//连接数据库，utils.go工具包
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&Video{})

	//初始化一个FeedResponse，用于存放本次的视频记录

	var video []Video
	var nextTime int64
	fr := FeedResponse{Response{1, ""}, []Video{}, nextTime}

	//从video数据库表里以时间戳大到小的顺序取出视频信息

	find := GLOBAL_DB.Order("lasted_time desc, id").Limit(30).Find(&video)

	// SELECT * FROM video ORDER BY create_time desc,id LIMIT 30;

	//检错报错

	if find.Error != nil {
		fr.Response.StatusCode = 1
		fr.Response.StatusMsg = "服务器内部错误"
		c.JSON(http.StatusInternalServerError, &fr)
	}

	fr.Response.StatusCode = 0
	fr.Response.StatusMsg = "成功"
	fr.VideoList = video
	fr.NextTime = nextTime
	c.JSON(http.StatusOK, &fr)

	//demo中的代码
	//c.JSON(http.StatusOK, FeedResponse{
	//	Response:  Response{StatusCode: 0},
	//	VideoList: DemoVideos,
	//	NextTime:  time.Now().Unix(),
	//})

}
