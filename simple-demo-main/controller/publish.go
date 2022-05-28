package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
)

type VideoListResponse struct {
	Response
	VideoList []Video `json:"video_list"`
}

// Publish check token then save upload file to public directory
func Publish(c *gin.Context) {
	token := c.PostForm("token")

	if _, exist := usersLoginInfo[token]; !exist {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
		return
	}

	data, err := c.FormFile("data")
	if err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	filename := filepath.Base(data.Filename)
	user := usersLoginInfo[token]
	finalName := fmt.Sprintf("%d_%s", user.Id, filename)
	saveFile := filepath.Join("./public/", finalName)
	if err := c.SaveUploadedFile(data, saveFile); err != nil {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  finalName + " uploaded successfully",
	})
}

// PublishList
//总流程：得到两个请求参数用户id和token（相当于加密后的用户密码），通过userId在UserInfoTab表里查询token并与请求参数中的token比对，
//若不相同，返回状态码400，异常；若相同，验证通过，从video信息表里取出视频信息，放入响应参数中发送出去。
//注：原demo这里有一句注释 “all users have same publish video list”， 所以这里不需要返回指定用户的视频，直接把表里所有的视频数据全部返回即可。

func PublishList(c *gin.Context) {
	//获取请求参数
	userId := c.Query("user_id")
	token := c.Query("token")
	//连接数据库
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&UserInfoTab{})
	//初始化一个UserInfoTab和一个VideoListResponse
	var userInfo UserInfoTab
	var video []Video
	vlr := VideoListResponse{Response{1, ""}, []Video{}}
	//通过userId在UserInfoTab表里查询token并与请求参数中的token比对
	verification := GLOBAL_DB.Select("token").Where("user_id = ?", userId).Find(&userInfo)
	if verification.Error != nil {
		vlr.Response.StatusCode = 1
		vlr.Response.StatusMsg = "服务器内部错误"
		c.JSON(http.StatusInternalServerError, &vlr)
	}
	if userInfo.Token != token {
		vlr.Response.StatusCode = 1
		vlr.Response.StatusMsg = "用户token不匹配"
		c.JSON(http.StatusBadRequest, &vlr)
	} else {
		//验证成功，从video信息表里取出视频信息
		_ = GLOBAL_DB.AutoMigrate(&Video{})

		find := GLOBAL_DB.Find(&video)
		if find.Error != nil {
			vlr.Response.StatusCode = 1
			vlr.Response.StatusMsg = "服务器内部错误"
			c.JSON(http.StatusInternalServerError, &vlr)
		}
		vlr.Response.StatusCode = 0
		vlr.Response.StatusMsg = "成功"
		vlr.VideoList = video
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
