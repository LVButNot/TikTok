package controller

import (
	"time"
)

type Response struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type Video struct {
	Id            int64  `json:"id,omitempty"`
	UserId        int64  `json:"author" `
	PlayUrl       string `json:"play_url" json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	//IsFavorite    bool   `json:"is_favorite,omitempty" gorm:"-"`
	Title      string `json:"title,omitempty"`
	LatestTime int64  `json:"latest_time"`
}

type Comment struct {
	Id         int64  `json:"id,omitempty"`
	UserId     int64  `json:"user"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type User struct {
	Id            int64     `json:"id,omitempty" gorm:"primary_key;type:bigint(20);not null;auto_increment" `
	CreatedAt     time.Time `json:"create_at,omitempty"`
	UpdatedAt     time.Time `json:"updated_at,omitempty"`
	Name          string    `json:"name,omitempty"`
	FollowCount   int64     `json:"follow_count,omitempty"`
	FollowerCount int64     `json:"follower_count,omitempty"`
	//IsFollow      bool      `json:"is_follow,omitempty" gorm:"-"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}
