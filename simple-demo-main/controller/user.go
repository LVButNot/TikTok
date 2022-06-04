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
/*type UserInfoTab struct {
	UserId   int64
	Name     string
	Password string
	Token    string
}*/

//从客户端接受username和password，然后先确认数据库里存不存在这个用户，若存在，返回"User already exist"；若不存在，将username和password、token保存到数据库中，userid自增1

func Register(c *gin.Context) {
	//获取请求参数
	username := c.Query("username")
	password := c.Query("password")
	token := username + password
	//连接数据库,封装在controller.utils包下
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//对password和token进行加密
	np := md5.Sum([]byte(token))
	tok := fmt.Sprintf("%X", np)
	fmt.Println(tok)
	pas := md5.Sum([]byte(password))
	pasmd5 := fmt.Sprintf("%X", pas)
	fmt.Println(pasmd5)
	//通过username从user表中查询数据
	var user User
	find := GLOBAL_DB.Where("name = ?", username).Find(&user)

	if find.RowsAffected != 0 {
		c.JSON(http.StatusOK, UserLoginResponse{
			Response: Response{StatusCode: 1, StatusMsg: "用户已存在"},
		})
	} else {
		newUser := User{
			Name:          username,
			FollowCount:   0,
			FollowerCount: 0,
			Password:      pasmd5,
			Token:         tok,
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
	//连接数据库,封装在controller.utils包下
	ConnectionSQL()
	_ = GLOBAL_DB.AutoMigrate(&User{})
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
	var user User
	find := GLOBAL_DB.Where("name = ?", username).Find(&user)
	fmt.Println(user)
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
	} else if user.Password != pasmd5 { //密码输入错误
		us.Response.StatusCode = 1
		us.Response.StatusMsg = "密码输入错误"
		fmt.Println(us)
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
	userId := c.Query("user_id")
	token := c.Query("token")
	//连接数据库
	ConnectionSQL()
	//刷新中的表数据库格，使其保持最新，即让数据库之前存储的记录的表格字段和程序中最新使用的表格字段保持一致（只增不减）
	_ = GLOBAL_DB.AutoMigrate(&User{})
	//初始化一个User和一个UserResponse
	var user User
	ur := UserResponse{Response{1, ""}, User{}}
	//通过userId在UserInfoTab表里查询token并与请求参数中的token比对
	verification := GLOBAL_DB.Select("token").Where("id = ?", userId).Find(&user)
	if verification.Error != nil {
		ur.Response.StatusCode = 1
		ur.Response.StatusMsg = "服务器内部错误"
		c.JSON(http.StatusInternalServerError, &ur)
	}
	if user.Token != token {
		ur.Response.StatusCode = 1
		ur.Response.StatusMsg = "用户token不匹配"
		c.JSON(http.StatusBadRequest, &ur)
	} else {
		//验证成功，通过userid从user信息表里取出用户信息
		//sql语句,等同于select * from users where id = ?
		var user User
		find := GLOBAL_DB.Select("id", "name", "follow_count", "follower_count").Where("Id = ?", userId).Find(&user)
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
		ur.Response.StatusCode = 0
		ur.Response.StatusMsg = "OK"
		ur.User = user
		c.JSON(http.StatusOK, &ur)
	}

}
