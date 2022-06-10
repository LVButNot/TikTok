package controller

import (
	"bytes"
	"fmt"
	"github.com/disintegration/imaging"
	ffmpeg "github.com/u2takey/ffmpeg-go"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var GLOBAL_DB *gorm.DB

func ConnectionSQL() {
	dsn := "lvjiayu:Lvjiayu20126493@tcp(rm-7xv2012nfi836h8908o.mysql.rds.aliyuncs.com:3306)/tiktok?charset=utf8mb4&parseTime=true&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:   false, //是否跳过默认事务
		DisableNestedTransaction: true,  //在 AutoMigrate 或 CreateTable 时，GORM 会自动创建外键约束，若要禁用该特性，可将其设置为 true
	})
	if err != nil {
		fmt.Println(db)
	}
	GLOBAL_DB = db
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
