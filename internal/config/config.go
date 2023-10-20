// config/config.go
package config

import (
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

var AppConfig struct {
	Port string
	// Add other configuration options here
}

func InitLogger() *logrus.Logger {
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{})
	log.SetOutput(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    50, // MB
		MaxBackups: 3,
		MaxAge:     7, // days
	})
	return log
}
