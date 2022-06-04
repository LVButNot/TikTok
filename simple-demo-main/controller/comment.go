package controller

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

type CommentListResponse struct {
	Response
	CommentList []Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	Response
	Comment Comment `json:"comment,omitempty"`
}

// CommentAction no practical effect, just check if token is valid
func CommentAction(c *gin.Context) {
	sUserId := c.PostForm("user_id")
	token := c.PostForm("token")
	videoId := c.PostForm("video_id")
	actionType := c.PostForm("action_type")
	commentText := c.PostForm("comment_text")
	sCommentId := c.PostForm("comment_id")
	commentId, _ := strconv.ParseInt(sCommentId, 10, 64)
	userId, _ := strconv.ParseInt(sUserId, 10, 64)

	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个User和一个VideoListResponse
	var user User
	car := CommentActionResponse{Response{1, ""}, Comment{}}
	//通过userId在UserInfoTab表里查询token并与请求参数中的token比对
	verification := GLOBAL_DB.Select("token").Where("id = ?", userId).Find(&user)
	if verification.Error != nil {
		car.Response.StatusCode = 1
		car.Response.StatusMsg = "服务器内部错误"
		c.JSON(http.StatusInternalServerError, &car)
	}
	if user.Token != token {
		car.Response.StatusCode = 1
		car.Response.StatusMsg = "用户不存在"
		c.JSON(http.StatusOK, &car)
	} else {
		//comments表（新建）
		//comment表里有commentId、videoId、userId、commentText、createTime
		//如果actionType=1，数据库添加评论；如果actionType=0，数据库删除评论
		//额外：实现线程锁解决并发问题
		ConnectionRedis1()
		if actionType == "1" {
			var user User
			user.Id = userId
			createDate := strconv.FormatInt(time.Now().Unix(), 10)
			var comment = Comment{Id: commentId, UserId: userId, Content: commentText, CreateDate: createDate}
			data, _ := json.Marshal(comment)
			RDB1.HSet(videoId, sCommentId, data)

			//在video表中使commentCount + 1
			video := Video{}
			GLOBAL_DB.Where("id = ", videoId).Take(&video)
			video.CommentCount += 1
			GLOBAL_DB.Save(video)

			car.StatusCode = 0
			car.StatusMsg = "评论成功"
			car.Comment.Content = commentText
			c.JSON(http.StatusOK, &car)
		} else if actionType == "2" {
			RDB1.HDel(videoId, sCommentId)

			//在video表中使commentCount - 1
			video := Video{}
			GLOBAL_DB.Where("id = ", videoId).Take(&video)
			video.CommentCount -= 1
			GLOBAL_DB.Save(video)

			car.StatusCode = 0
			car.StatusMsg = "删除评论成功"
			c.JSON(http.StatusOK, &car)
		}
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {

	token := c.Query("token")
	videoId := c.Query("video_id")
	//第一步，验证用户的那一套（注意这里只给了token，所以直接在user里面用token查，然后判断*gorm.DB.RowsAffected是否为0）（RowsAffected返回查到了几条数据）
	//第二步，通过videoId在comment表中查询响应中的数据，然后发送响应即可
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个User
	var user User
	clr := CommentListResponse{Response{StatusCode: 1, StatusMsg: ""}, []Comment{}}
	//通过token在User表里查询，判断返回数据个数
	verification := GLOBAL_DB.Where("token = ?", token).Find(&user)
	if verification.Error != nil {
		clr.StatusMsg = "数据库查询错误"
		c.JSON(http.StatusInternalServerError, &clr)
	}
	if verification.RowsAffected == 0 {
		clr.StatusMsg = "用户不存在"
		c.JSON(http.StatusOK, &clr)
	} else {
		err := ConnectionRedis1()
		if err != nil {
			clr.Response.StatusCode = 1
			clr.Response.StatusMsg = "Redis查询失败"
			c.JSON(http.StatusInternalServerError, &clr)
		}
		//如何取出json数据并转化？？
		var commentList []Comment
		finds, _ := RDB1.HVals(videoId).Result()
		for _, find := range finds {
			comment := Comment{}
			err := json.Unmarshal([]byte(find), &comment)
			if err != nil {
				clr.StatusMsg = "反序列化失败"
				c.JSON(http.StatusOK, &clr)
			}
			commentList = append(commentList, comment)
		}
		clr.Response.StatusCode = 0
		clr.Response.StatusMsg = "拉取成功"
		clr.CommentList = commentList
		c.JSON(http.StatusOK, &clr)
	}
}
