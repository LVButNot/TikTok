package controller

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

// FavoriteAction no practical effect, just check if token is valid
func FavoriteAction(c *gin.Context) {
	sUserId := c.Query("user_id")
	userId, _ := strconv.ParseInt(sUserId, 10, 64)
	token := c.Query("token")
	sVideoId := c.Query("video_id")
	actionType := c.Query("action_type")
	videoId, _ := strconv.ParseInt(sVideoId, 10, 64)
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

		//在favorites添加记录
		favorite := Favorite{
			UserId:  user.Id,
			VideoId: videoId,
		}

		//开启数据库事务，在favorites中添加记录，在videos中更改点赞数目
		GLOBAL_DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Create(&favorite).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			tx.Model(&Video{}).Where("id=?", videoId).UpdateColumn("favorite_count", gorm.Expr("favorite_count + ?", 1))
			// 返回 nil 提交事务
			return nil
		})
		c.JSON(http.StatusOK, Response{StatusCode: 0, StatusMsg: "点赞成功"})

	} else if actionType == "2" {
		GLOBAL_DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Where("user_id=?", userId).Where("video_id=?", videoId).Delete(&Favorite{}).Error; err != nil {
				// 返回任何错误都会回滚事务
				return err
			}
			tx.Model(&Video{}).Where("id=?", videoId).UpdateColumn("favorite_count", gorm.Expr("favorite_count - ?", 1))
			// 返回 nil 提交事务
			return nil
		})

		c.JSON(http.StatusOK, Response{
			StatusCode: 0,
			StatusMsg:  "取消成功",
		})
	}
}

// FavoriteList all users have same favorite video list
func FavoriteList(c *gin.Context) {
	sUserId := c.Query("user_id")
	userId, _ := strconv.ParseInt(sUserId, 10, 64)
	token := c.Query("token")

	ConnectionSQL()
	if token == "" {
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  "You haven't logged in yet",
		})
	}
	var favoriteList []Favorite
	var videoList []Video
	GLOBAL_DB.Where("user_id=?", userId).Find(&favoriteList)

	GLOBAL_DB.Table("favorites").Select("favorites.video_id,videos.*").
		Where("favorites.user_id=?", userId).
		Joins("LEFT JOIN videos ON favorites.video_id = videos.id").
		Find(&videoList)

	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0, StatusMsg: "拉取成功"},
		VideoList: videoList,
	})

}
