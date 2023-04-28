package utils

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

var DB *gorm.DB

func InitConfig() {
	viper.SetConfigName("app")
	viper.AddConfigPath("config")
	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("app inited ...")
}

func InitMySQL() {
	//自定义日志模板 打印SQL语句
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), //Writer
		logger.Config{ //Config
			SlowThreshold: time.Second, //慢SQL阈值
			LogLevel:      logger.Info, //日志级别
			Colorful:      true,        //彩色
		})
	//DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{})

	DB, _ = gorm.Open(mysql.Open(viper.GetString("mysql.dns")), &gorm.Config{Logger: newLogger})
	fmt.Println("mysql inited ...")
	//user := models.UserBasic{}
	//DB.Find(&user)
	//fmt.Println(user)
}
