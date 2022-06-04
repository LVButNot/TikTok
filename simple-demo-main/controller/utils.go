package controller

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	"github.com/go-redis/redis"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var GLOBAL_DB *gorm.DB
var RDB *redis.Client
var RDB1 *redis.Client
var RDB2 *redis.Client
var RDB3 *redis.Client

func ConnectionSQL() {
	dsn := "root:wulingwei@tcp(175.178.126.39:3306)/tiktok?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:   false, //是否跳过默认事务
		DisableNestedTransaction: true,  //在 AutoMigrate 或 CreateTable 时，GORM 会自动创建外键约束，若要禁用该特性，可将其设置为 true
	})
	if err != nil {
		fmt.Println(db)
	}
	GLOBAL_DB = db
}

func ConnectionRedis() (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "175.178.126.39",
		Password: "wulingwei",
		DB:       0,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	RDB = rdb
	return err
}

func ConnectionRedis1() (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "175.178.126.39",
		Password: "wulingwei",
		DB:       1,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	RDB1 = rdb
	return err
}

func ConnectionRedis2() (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "175.178.126.39",
		Password: "wulingwei",
		DB:       2,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	RDB2 = rdb
	return err
}

func ConnectionRedis3() (err error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "175.178.126.39",
		Password: "wulingwei",
		DB:       3,
	})
	_, err = rdb.Ping().Result()
	if err != nil {
		return err
	}
	RDB3 = rdb
	return err
}

// GetSnapshot 生成视频缩略图并保存（作为封面）
func GetSnapshot(videoPath, snapshotPath string, frameNum int) (snapshotName string) {
	buf := bytes.NewBuffer(nil)
	err := ffmpeg.Input(videoPath).
		Filter("select", ffmpeg.Args{fmt.Sprintf("gte(n,%d)", frameNum)}).
		Output("pipe:", ffmpeg.KwArgs{"vframes": 1, "format": "image2", "vcodec": "mjpeg"}).
		WithOutput(buf, os.Stdout).
		Run()
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
	}

	img, err := imaging.Decode(buf)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
	}
	if len(snapshotPath) == 0 {
		paths, fileName := filepath.Split(videoPath)
		name := strings.Split(fileName, ".")[0]
		snapshotPath = filepath.Join(paths, name+".jpeg")
	}

	err = imaging.Save(img, snapshotPath)
	if err != nil {
		log.Fatal("生成缩略图失败：", err)
	}

	// 成功则返回生成的缩略图名
	// fmt.Println("snapshotName:", snapshotPath)
	return snapshotPath
}
