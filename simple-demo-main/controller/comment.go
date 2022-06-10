package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	token := c.Query("token")
	sVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(sVideoId, 10, 64)
	actionType := c.Query("action_type")
	context := c.Query("comment_text")
	sCommentId := c.Query("comment_id")
	commentId, _ := strconv.ParseInt(sCommentId, 10, 64)

	ConnectionSQL()
	if token == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户信息不匹配",
		})
	}
	if actionType == "1" {
		var user User
		GLOBAL_DB.Where("token=?", token).First(&user)
		timeFormat := time.Now().Format("01-02")
		comment := Comment{
			User:       user,
			VideoId:    videoId,
			UserId:     user.Id,
			Content:    context,
			CreateDate: timeFormat,
		}
		//开启数据库事务，在comments中添加记录，在videos中更改评论数目
		GLOBAL_DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&comment).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			tx.Model(&Video{}).Where("id=?", videoId).UpdateColumn("comment_count", gorm.Expr("comment_count + ?", 1))
			// 返回 nil 提交事务
			return nil
		})

		//文档中标明不需要拉取评论列表，数据库中的自增id无法获取
		//目前默认每次处理一条comment，所以数组只存入一条评论数据
		commentList := []Comment{comment}

		c.JSON(http.StatusOK, CommentListResponse{
			Response: Response{StatusCode: 0,
				StatusMsg: "发表成功"},
			CommentList: commentList,
		})
	} else if actionType == "2" {
		GLOBAL_DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("id=?", commentId).Delete(&Comment{}).Error; err != nil {
				return err
			}
			tx.Model(&Video{}).Where("id=?", videoId).UpdateColumn("comment_count", gorm.Expr("comment_count - ?", 1))
			return nil
		})
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "删除成功",
		})
	}
}

func CommentList(c *gin.Context) {

	token := c.Query("token")
	sVideoId := c.Query("video_id")
	videoId, _ := strconv.ParseInt(sVideoId, 10, 64)
	ConnectionSQL()
	if token == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "用户信息不存在",
		})
	}
	var commentList []Comment
	GLOBAL_DB.Preload("User").Where("video_id=?", videoId).Find(&commentList)

	c.JSON(http.StatusOK, CommentListResponse{
		Response:    Response{StatusCode: 0, StatusMsg: "success"},
		CommentList: commentList,
	})

}
