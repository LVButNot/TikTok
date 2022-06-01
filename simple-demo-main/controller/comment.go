package controller

import (
	"github.com/gin-gonic/gin"
	"net/http"
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
	//comments表（新建）
	//comment表里有commentId、videoId、userId、commentContent、createTime
	//如果actionType=1，数据库添加评论；如果actionType=0，数据库删除评论
	//额外：实现线程锁解决并发问题

	//下面这个是原demo，我没有修改，做的时候删掉
	token := c.Query("token")
	actionType := c.Query("action_type")

	if user, exist := usersLoginInfo[token]; exist {
		if actionType == "1" {
			text := c.Query("comment_text")
			c.JSON(http.StatusOK, CommentActionResponse{Response: Response{StatusCode: 0},
				Comment: Comment{
					Id:         1,
					User:       user,
					Content:    text,
					CreateDate: "05-01",
				}})
			return
		}
		c.JSON(http.StatusOK, Response{StatusCode: 0})
	} else {
		c.JSON(http.StatusOK, Response{StatusCode: 1, StatusMsg: "User doesn't exist"})
	}
}

// CommentList all videos have same demo comment list
func CommentList(c *gin.Context) {
	//第一步，验证用户的那一套（注意这里只给了token，所以直接在user_info_tabs里面用token查，然后判断*gorm.DB.RowsAffected是否为0）（RowsAffected返回查到了几条数据）
	//第二步，通过videoId在comment表中查询响应中的数据，然后发送响应即可
}
