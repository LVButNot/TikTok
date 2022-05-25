package controller

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync/atomic"
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

var userIdSequence = int64(1)

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
	username := c.Query("username")
	password := c.Query("password")

	token := username + password

	if _, exist := usersLoginInfo[token]; exist {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "User already exist"},
		})
	} else {
		atomic.AddInt64(&userIdSequence, 1)
		newUser := User{
			Id:   userIdSequence,
			Name: username,
		}
		usersLoginInfo[token] = newUser
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 0},
			UserId:   userIdSequence,
			Token:    username + password,
		})
	}
}

// UserInfoTab
//		对应数据库里的user_info_tabs表，用来存放用户的登录信息。
type UserInfoTab struct {
	UserId   int64
	Name     string
	Password string
	Token    string
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
	np := md5.Sum([]byte(username + password))
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
		us.Token = token
		c.JSON(http.StatusOK, &us)
	}
}

func UserInfo(c *gin.Context) {
	//获取请求参数
	userId := c.Query("user_id")
	//连接数据库
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})

	ur := UserResponse{Response{1, ""}, User{}}

	//sql语句,等同于select * from users where id = ?
	var user User
	find := GLOBAL_DB.Where("Id = ?", userId).Find(&user)
	if find.Error != nil {
		ur.Response.StatusCode = 1
		ur.Response.StatusMsg = "数据库查询错误"
		return
	}
	ur.Response.StatusCode = 0
	ur.User = user
	c.JSON(http.StatusOK, &ur)
	fmt.Println(ur.User)
}