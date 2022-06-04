package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	var users []User
	GLOBAL_DB.Where("token=?", token).Find(&users)
	//数据库中查找token的记录行行数为0，表示用户不存在
	if len(users) == 0 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
	//token对应的记录行行数大于1，说明用户信息出错了
	if len(users) > 1 {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User error"})
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}
	title := c.PostForm("title")

	filename := filepath.Base(data.Filename)
	if len(filename) > 10 {
		filename = filename[0:9]
	}
	user := users[0]
	nowtime := time.Now().Unix()
	finalName := fmt.Sprintf("%d_%d_%s", user.Id, nowtime, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	snapshotPath := GetSnapshot(saveFile, "", 1)

	//制作video对象，存入数据库
	video := Video{
		UserId:        user.Id,
		PlayUrl:       saveFile,
		CoverUrl:      snapshotPath,
		FavoriteCount: 0,
		CommentCount:  0,
		Title:         title,
		LatestTime:    nowtime,
	}
	if err := GLOBAL_DB.Create(&video).Error; err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList
//总流程：得到两个请求参数用户id和token（相当于加密后的用户密码），通过userId在UserInfoTab表里查询token并与请求参数中的token比对，
//若不相同，返回状态码400，异常；若相同，验证通过，从video信息表里取出视频信息，放入响应参数中发送出去。

func PublishList(c *gin.Context) {
	//获取请求参数
	userId := c.Query("user_id")
	token := c.Query("token")

	//连接数据库
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个UserInfoTab和一个VideoListResponse
	var user User
	var video Video
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
		vlr.Response.StatusMsg = "用户token不匹配"
		c.JSON(http.StatusBadRequest, &vlr)
	} else {
		//验证成功，从video信息表里取出视频信息
		_ = GLOBAL_DB.AutoMigrate(&Video{})

		find := GLOBAL_DB.Where("user_id =", userId).Find(&video)
		if find.Error != nil {
			vlr.Response.StatusCode = 1
			vlr.Response.StatusMsg = "服务器内部错误"
			c.JSON(http.StatusInternalServerError, &vlr)
		}
		vlr.Response.StatusCode = 0
		vlr.Response.StatusMsg = "成功"
		//vlr.VideoList = video
		c.JSON(http.StatusOK, &vlr)
	}
	/*
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{
				StatusCode: 0,
			},
			VideoList: DemoVideos,
		})
	*/
}
