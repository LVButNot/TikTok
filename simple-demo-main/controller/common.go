package controller

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty" gorm:"primaryKey"`
	Author        User   `json:"author" gorm:"foreignKey:UserId"`
	UserId        int64  `gorm:"not null"` // 视频对应的用户 Id
	PlayUrl       string `json:"play_url" json:"play_url,omitempty" gorm:"not null"`
	CoverUrl      string `json:"cover_url,omitempty" gorm:"not null"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
	Title         string `json:"title,omitempty"` // demo 里没有 title
}

type Comment struct {
	Id         int64  `json:"id,omitempty" gorm:"primaryKey"`
	User       User   `json:"user" gorm:"foreignKey:UserId"`
	UserId     int64  `gorm:"not null"` // 评论对应的用户 Id
	Video      Video  `gorm:"foreignKey:VideoId"`
	VideoId    int64  `gorm:"not null"` // 评论对应的视频 Id
	Content    string `json:"content,omitempty" gorm:"not null"`
	CreateDate string `json:"create_date,omitempty" gorm:"not null"`
}

type User struct {
	Id            int64  `json:"id,omitempty" gorm:"primaryKey"`
	Name          string `json:"name,omitempty" gorm:"unique; not null"`
	Password      string `gorm:"not null"`
	Token         string `json:"token,omitempty" gorm:"not null"`
	FollowCount   int64  `json:"follow_count,omitempty"`
	FollowerCount int64  `json:"follower_count,omitempty"`
	IsFollow      bool   `json:"is_follow,omitempty"`
}

// Favorite 记录用户点赞的视频
type Favorite struct {
	Id      int64 `gorm:"primaryKey"`
	User    User  `gorm:"foreignKey:UserId"`
	UserId  int64 `gorm:"not null"`
	Video   Video `gorm:"foreignKey:VideoId"`
	VideoId int64 `gorm:"not null"`
}

// Relation 用于维护用户关注关系，待之后完善
// 一行数据代表 "UserA 关注了 UserB"
type Relation struct {
	Id      int64 `gorm:"primaryKey"`
	UserA   User  `gorm:"foreignKey:UserAId"`
	UserAId int64 `gorm:"notnull"`
	UserB   User  `gorm:"foreignKey:UserBId"`
	UserBId int64 `gorm:"notnull"`
}
