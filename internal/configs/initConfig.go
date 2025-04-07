package configs

import (
	"fmt"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
)

// 設定命令行標誌 (flags)，允許指定配置文件路徑
var (
	configFile = pflag.StringP("config", "c", "", "Path to configuration file")
)

// 應用配置結構\
// @Summary App configuration
// @Description App configuration
// @Tags App
// @ID app-conf
type AppConf struct {
	Version string `mapstructure:"version"`
	Mode    string `mapstructure:"mode"`
	Port    string `mapstructure:"port"`
	LogPath string `mapstructure:"logpath"`
	LogFile string `mapstructure:"logfile"`
	Ukey    string `mapstructure:"ukey"`
	Vtoken  string `mapstructure:"vtoken"`
	Salt    string `mapstructure:"salt"`
}

// metric user and password
// @Summary Metric configuration
// @Description Metric configuration
// @Tags Metric
// @ID metric-conf
type Metric struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Config struct {
	App       AppConf   `mapstructure:"app"`
	Metric    Metric    `mapstructure:"metric"`
	JwtSecret JwtSecret `mapstructure:"jwt"`
}

// 存儲全局配置
var Conf Config

// Init 初始化配置
func Init() {
	// 解析命令行參數
	pflag.Parse()

	// 創建新的 Viper 實例
	v := viper.New()
	setDefaults(v) // 設置默認配置
	v.AutomaticEnv()

	// 設定配置文件
	if *configFile != "" {
		v.SetConfigFile(*configFile)
	} else {
		v.SetConfigName("config")   // 預設文件名稱 config.yaml
		v.AddConfigPath("./config") // 尋找 config/ 目錄
		v.AddConfigPath(".")        // 當前目錄
	}
	v.SetConfigType("yaml")

	// 讀取配置
	if err := v.ReadInConfig(); err != nil {
		log.Fatalf("讀取配置文件失敗: %v", err)
	}

	// 解析配置到結構體
	if err := v.Unmarshal(&Conf); err != nil {
		log.Fatalf("解析配置文件失敗: %v", err)
	}

	// 確保加載成功
	if v.ConfigFileUsed() != "" {
		fmt.Printf("成功加載配置: %s\n", v.ConfigFileUsed())
	} else {
		log.Println("未找到任何有效的配置文件，將使用默認值")
	}
}

// 設置默認值
func setDefaults(v *viper.Viper) {
	v.SetDefault("app.version", "1.0.0")
	v.SetDefault("app.mode", "development")
	v.SetDefault("app.port", "8080")
	v.SetDefault("app.logpath", "./logs")
	v.SetDefault("app.logfile", "app.log")
	v.SetDefault("app.ukey", "")
	v.SetDefault("app.vtoken", "")
	v.SetDefault("app.salt", "")

	v.SetDefault("metric.user", "admin")
	v.SetDefault("metric.password", "password")
}
