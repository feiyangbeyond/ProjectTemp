package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server ServerConfig `json:"server" yaml:"server"`
	Mysql  MysqlConfig  `json:"mysql" yaml:"mysql"`
	Redis  RedisConfig  `json:"redis" yaml:"redis"`
	Kafka  KafkaConfig  `json:"kafka" yaml:"kafka"`
}

type ServerConfig struct {
	Mode            string `json:"mode" yaml:"mode"`
	Port            int    `json:"port" yaml:"port"`
	ReadBufferSize  int    `json:"readBufferSized" yaml:"readBufferSize"`
	WriteBufferSize int    `json:"writeBufferSize" yaml:"writeBufferSize"`
	MaxMessageSize  int64  `json:"maxMessageSize" yaml:"maxMessageSize"`
}

type MysqlConfig struct {
	User        string `yaml:"user"`
	Password    string `yaml:"password"`
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Database    string `yaml:"database"`
	Timeout     string `yaml:"timeout"`
	MaxConn     int    `yaml:"maxConn"`
	MaxIdleConn int    `yaml:"maxIdleConn"`
}

type RedisConfig struct {
	Host     string
	Port     int
	Password string
	Db       int
}

type KafkaConfig struct {
	Host string
	Port int
}

func Scan(path string) (*Config, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var c Config
	err = yaml.Unmarshal(b, &c)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
