package config

import (
	"flag"
	"gopkg.in/ini.v1"
	"log"
)

var (
	configFile string
	Conf       Config
)

type Config struct {
	Listen      string
	BucketSize  int
	BucketName  string
	RedisHost   string
	RedisPasswd string
	LogPath     string
	Timeout     int
	MaxTries    int
}

const (
	defaultListen      = "0.0.0.0:9527"
	defaultBucketSize  = 3
	defaultBucketName  = "dqc_bucket_%d%"
	defaultRedisHost   = "127.0.0.1:6379"
	defaultRedisPasswd = ""
	defaultLogPath     = "./log/delay.log"
	defaultTimeout     = 30
	defaultMaxTries    = 3
)

func init() {
	flag.StringVar(&configFile, "c", "", "./delay-queue -c /path/to/delay_queue_callback.conf");
	flag.Parse()

	if configFile == "" {
		initDefaultConfig()
		return
	}
	fp, err := ini.Load(configFile)
	if err != nil {
		log.Fatal("解析配置文件出错:", err)
	}
	section := fp.Section("")
	Conf.Listen = section.Key("listen").MustString(defaultListen)
	Conf.BucketSize = section.Key("bucket.size").MustInt(defaultBucketSize)
	Conf.BucketName = section.Key("bucket.name").MustString(defaultBucketName)
	Conf.RedisHost = section.Key("redis.host").MustString(defaultRedisHost)
	Conf.RedisPasswd = section.Key("redis.passwd").MustString(defaultRedisPasswd)
	Conf.LogPath = section.Key("logpath").MustString(defaultLogPath)
	Conf.Timeout = section.Key("timeout").MustInt(defaultTimeout)
	Conf.MaxTries = section.Key("max_tries").MustInt(defaultMaxTries)
}

func initDefaultConfig() {
	Conf.Listen = defaultListen
	Conf.BucketSize = defaultBucketSize
	Conf.BucketName = defaultBucketName
	Conf.RedisHost = defaultRedisHost
	Conf.RedisPasswd = defaultRedisPasswd
	Conf.LogPath = defaultLogPath
	Conf.Timeout = defaultTimeout
	Conf.MaxTries = defaultMaxTries
}