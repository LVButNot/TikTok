package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strconv"
	"time"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")
	ConnectionSQL()
	if token == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户信息不匹配",
		})
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
	var user User
	GLOBAL_DB.Where("token=?", token).Find(&user)
	nowTime := time.Now().Unix()
	finalName := fmt.Sprintf("%d_%d_%s", user.Id, nowTime, filename)
	saveFile := filepath.Join(".", "public", finalName+".mp4")
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	snapshotPath := GetSnapshot(saveFile, "", 1)

	//制作video对象，存入数据库
	video := Video{
		Author:   user,
		UserId:   user.Id,
		PlayUrl:  saveFile,
		CoverUrl: snapshotPath,
		Title:    title,
	}
	if err := GLOBAL_DB.Create(&video).Error; err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  "投稿成功",
	})
}

// PublishList
//总流程：得到两个请求参数用户id和token（相当于加密后的用户密码），通过userId在UserInfoTab表里查询token并与请求参数中的token比对，
//若不相同，返回状态码400，异常；若相同，验证通过，从video信息表里取出视频信息，放入响应参数中发送出去。

func PublishList(c *gin.Context) {
	//获取请求参数
	sUserId := c.Query("user_id")
	userId, _ := strconv.ParseInt(sUserId, 10, 64)

	//连接数据库
	ConnectionSQL()

	var videoList []Video
	GLOBAL_DB.Where("user_id = ?", userId).Find(&videoList)
	c.JSON(http.StatusOK, VideoListResponse{
		Response: Response{
			StatusCode: 0,
		},
		VideoList: videoList,
	})
}
