package controller

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

func Register(c *gin.Context) {
	//获取请求参数
	username := c.Query("username")
	password := c.Query("password")
	token := username + password

	//连接数据库,封装在controller.utils包下
	ConnectionSQL()

	//对password和token进行加密
	np := md5.Sum([]byte(token))
	tok := fmt.Sprintf("%X", np)
	pas := md5.Sum([]byte(password))
	pasmd5 := fmt.Sprintf("%X", pas)

	//通过username从user表中查询数据
	var user User
	find := GLOBAL_DB.Where("name = ?", username).Find(&user)

	if find.RowsAffected != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已存在"},
		})
	} else {
		newUser := User{
			Name:     username,
			Password: pasmd5,
			Token:    tok,
		}
		userID := newUser.Id
		GLOBAL_DB.Create(&newUser)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userID,
			Token:    tok,
		})
	}
}

func Login(c *gin.Context) {
	//获取请求参数
	username := c.Query("username")
	password := c.Query("password")
	token := username + password

	//连接数据库,封装在controller.utils包下
	ConnectionSQL()

	us := UserLoginResponse{Response: Response{1, "error"}, UserId: 0, Token: "error"}

	//对token和password进行加密，使用标准包的md5算法
	np := md5.Sum([]byte(token))
	tok := fmt.Sprintf("%X", np)
	pas := md5.Sum([]byte(password))
	pasmd5 := fmt.Sprintf("%X", pas)

	//gorm的sql语句实现，等同于select * from user_info_tabs where username = ?
	var user User
	find := GLOBAL_DB.Where("name = ?", username).Find(&user)

	if find.Error != nil {
		us.Response.StatusCode = 1
		us.Response.StatusMsg = "服务器内部错误"
		c.JSON(500, &us)
	}

	if find.RowsAffected == 0 { //用户名不存在
		us.Response.StatusCode = 1
		us.Response.StatusMsg = "用户名未找到"
		c.JSON(http.StatusOK, &us)
	} else if user.Password != pasmd5 { //密码输入错误
		us.Response.StatusCode = 1
		us.Response.StatusMsg = "密码输入错误"
		c.JSON(http.StatusOK, &us)
	} else if user.Password == pasmd5 { //登录成功
		us.Response.StatusCode = 0
		us.UserId = user.Id
		us.Token = tok
		c.JSON(http.StatusOK, &us)
	}
}

//先通过userid验证token是否相同，若相同，返回用户信息

func UserInfo(c *gin.Context) {
	//获取请求参数
	sUserId := c.Query("user_id")
	userId, _ := strconv.ParseInt(sUserId, 10, 64)
	token := c.Query("token")

	//连接数据库
	ConnectionSQL()

	//初始化一个UserA,B和一个UserResponse
	ur := UserResponse{Response{1, ""}, User{}}
	var userA, userB User //userA是自己，userB是查看用户

	if token == "" {
		ur.Response.StatusCode = 1
		ur.Response.StatusMsg = "用户token不存在"
		c.JSON(http.StatusBadRequest, &ur)
	} else {
		//验证成功
		GLOBAL_DB.Where("token = ? ", token).Find(&userA)
		GLOBAL_DB.Where("id = ? ", userId).Find(&userB)
		find := GLOBAL_DB.Where("user_a_id = ? ", userA.Id).Where("user_b_id = ? ", userB.Id).Find(&Relation{}).RowsAffected
		userB.IsFollow = find > 0
		ur.Response.StatusCode = 0
		ur.User = userB
		c.JSON(http.StatusOK, &ur)
	}

}
