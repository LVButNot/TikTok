package controller

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

//这个usersLoginInfo是demo里测试用的，应该删掉，但是用户注册接口还没有完成，为了不爆红先留着，注册功能完成后删掉。
var usersLoginInfo = map[string]User{
	"zhangleidouyin": {
		Id:            1,
		Name:          "zhanglei",
		FollowCount:   10,
		FollowerCount: 5,
		IsFollow:      true,
	},
}

type UserLoginResponse struct {
	Response
	UserId int64  `json:"user_id,omitempty"`
	Token  string `json:"token"`
}

type UserResponse struct {
	Response
	User User `json:"user"`
}

// UserInfoTab
//		对应数据库里的user_info_tabs表，用来存放用户的登录信息。
type UserInfoTab struct {
	UserId   int64 `json:"userid,omitempty" gorm:"primary_key;type:bigint(20);not null;auto_increment"`
	Name     string
	Password string
	Token    string
}

func Register(c *gin.Context) {
	//获取请求参数
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	//连接数据库,封装在controller.utils包下
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&UserInfoTab{})
	//对password和token进行加密
	np := md5.Sum([]byte(token))
	tok := fmt.Sprintf("%X", np)
	//fmt.Println(tok)
	pas := md5.Sum([]byte(password))
	pasmd5 := fmt.Sprintf("%X", pas)
	//fmt.Println(pasmd5)
	//通过username从user_info_tabs表中查询数据
	var uit UserInfoTab
	find := GLOBAL_DB.Where("name = ?", username).Find(&uit)

	if find.RowsAffected != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已存在"},
		})
	} else {

		newUser := User{
			//Id:            userIdSequence,
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
			IsFollow:      false,
		}

		GLOBAL_DB.Create(&newUser)
		newUserInfo := UserInfoTab{
			//UserId:   userID,
			Name:     username,
			Password: pasmd5,
			Token:    tok,
		}
		userID := newUser.Id

		GLOBAL_DB.Create(&newUserInfo)
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userID,
			Token:    tok,
		})
	}
}

func Login(c *gin.Context) {
	//连接数据库,封装在controller.utils包下
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&UserInfoTab{})
	//获取请求参数
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	us := UserLoginResponse{Response: Response{1, "error"}, UserId: 0, Token: "error"}

	//用户注册，这段你放开，就能使用客户端抖音往数据库添加数据*/
	/*np := md5.Sum([]byte(username + password))
	tok := fmt.Sprintf("%X", np)
	pas := md5.Sum([]byte(password))
	pasmd5 := fmt.Sprintf("%X", pas)
	u := UserInfoTab{
		Name:     username,
		Password: pasmd5,
		Token:    tok,
	}
	GLOBAL_DB.Create(&u)*/

	//对token和password进行加密，使用标准包的md5算法
	np := md5.Sum([]byte(token))
	tok := fmt.Sprintf("%X", np)
	fmt.Println(tok)
	pas := md5.Sum([]byte(password))
	pasmd5 := fmt.Sprintf("%X", pas)
	fmt.Println(pasmd5)

	//gorm的sql语句实现，等同于select * from user_info_tabs where username = ?
	var uit UserInfoTab
	find := GLOBAL_DB.Where("name = ?", username).Find(&uit)
	fmt.Println(uit)
	fmt.Println(find)

	if find.Error != nil {
		us.Response.StatusCode = 1
		us.Response.StatusMsg = "服务器内部错误"
		c.JSON(500, &us)
	}

	if find.RowsAffected == 0 { //用户名不存在
		fmt.Println(find.Error)
		us.Response.StatusCode = 1
		us.Response.StatusMsg = "用户名未找到"
		fmt.Println(us)
		c.JSON(http.StatusOK, &us)
	} else if uit.Password != pasmd5 { //密码输入错误
		us.Response.StatusCode = 1
		us.Response.StatusMsg = "密码输入错误"
		fmt.Println(us)
		c.JSON(http.StatusOK, &us)
	} else if uit.Password == pasmd5 { //登录成功
		us.Response.StatusCode = 0
		us.UserId = uit.UserId
		us.Token = tok
		c.JSON(http.StatusOK, &us)
	}
}

//先通过userid验证token是否相同，若相同，返回用户信息

func UserInfo(c *gin.Context) {
	//获取请求参数
	userId := c.Query("user_id")
	token := c.Query("token")
	//连接数据库
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个UserInfoTab和一个UserResponse
	var userInfo UserInfoTab
	ur := UserResponse{Response{1, ""}, User{}}
	//通过userId在UserInfoTab表里查询token并与请求参数中的token比对
	verification := GLOBAL_DB.Select("token").Where("user_id = ?", userId).Find(&userInfo)
	if verification.Error != nil {
		ur.Response.StatusCode = 1
		ur.Response.StatusMsg = "服务器内部错误"
		c.JSON(http.StatusInternalServerError, &ur)
	}
	if userInfo.Token != token {
		ur.Response.StatusCode = 1
		ur.Response.StatusMsg = "用户token不匹配"
		c.JSON(http.StatusBadRequest, &ur)
	} else {
		//验证成功，通过userid从user信息表里取出用户信息
		//sql语句,等同于select * from users where id = ?
		var user User
		find := GLOBAL_DB.Select("id", "name", "follow_count", "follower_count", "is_follow").Where("Id = ?", userId).Find(&user)
		if find.Error != nil {
			ur.Response.StatusCode = 1
			ur.Response.StatusMsg = "服务器内部错误"
			c.JSON(http.StatusInternalServerError, &ur)
		}
		if find.RowsAffected == 0 {
			ur.Response.StatusCode = 1
			ur.Response.StatusMsg = "未找到用户"
			c.JSON(http.StatusBadRequest, &ur)
			fmt.Println(ur)
		}
		ur.Response.StatusCode = 0
		ur.Response.StatusMsg = "成功"
		ur.User = user
		c.JSON(http.StatusOK, &ur)
	}

	/*
		var userInfo UserInfoTab
		verification := GLOBAL_DB.Select("token").Where("user_id = ?", userId).Find(&userInfo)

		ur := UserResponse{Response{1, ""}, User{}}
		if verification.Error != nil {
			ur.Response.StatusCode = 1
			ur.Response.StatusMsg = "服务器内部错误"
			c.JSON(http.StatusInternalServerError, &ur)
		}
		if verification.RowsAffected == 0 {
			ur.Response.StatusCode = 1
			ur.Response.StatusMsg = "用户token未匹配"
			c.JSON(http.StatusOK, &ur)
			fmt.Println(ur)
		}

		//sql语句,等同于select * from users where id = ?
		var user User
		find := GLOBAL_DB.Select("id", "name", "follow_count", "follower_count", "is_follow").Where("Id = ?", userId).Find(&user)
		if find.Error != nil {
			ur.Response.StatusCode = 1
			ur.Response.StatusMsg = "服务器内部错误"
			c.JSON(http.StatusInternalServerError, &ur)
		}
		if find.RowsAffected == 0 {
			ur.Response.StatusCode = 1
			ur.Response.StatusMsg = "未找到用户"
			c.JSON(http.StatusOK, &ur)
			fmt.Println(ur)
		}
		fmt.Println(find.RowsAffected)
		ur.Response.StatusCode = 0
		ur.Response.StatusMsg = "成功"
		ur.User = user
		c.JSON(http.StatusOK, &ur)
		fmt.Println(ur.User)
	*/
}
