package logger

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

var Logger = logrus.New()

func InitLogger() {
	file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		Logger.Fatal("Error opening log file", err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	Logger.SetOutput(multiWriter)
	Logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	Logger.SetLevel(logrus.InfoLevel)
}
