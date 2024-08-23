package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
)

/*
@Author: OvO
@Date: 2024/8/21 14:35
*/

type MyTextFormatter struct {
}

func (f *MyTextFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := entry.Time.Format("2006-01-02 15:04:05")
	levelText := ""
	switch entry.Level {
	case logrus.InfoLevel:
		levelText = fmt.Sprintf("\u001B[32m[%s]\u001B[0m [%s] %s\n", strings.ToUpper(entry.Level.String()), timestamp, entry.Message)
	case logrus.ErrorLevel:
		levelText = fmt.Sprintf("\u001B[31m[%s]\u001B[0m [%s] %s\n", strings.ToUpper(entry.Level.String()), timestamp, entry.Message)
	case logrus.FatalLevel:
		levelText = fmt.Sprintf("\u001B[31m[%s]\u001B[0m [%s] %s\n", strings.ToUpper(entry.Level.String()), timestamp, entry.Message)
	case logrus.WarnLevel:
		levelText = fmt.Sprintf("\u001B[33m[%s]\u001B[0m [%s] %s\n", strings.ToUpper(entry.Level.String()), timestamp, entry.Message)
	default:
		levelText = fmt.Sprintf("\u001B[32m[%s]\u001B[0m [%s] %s\n", strings.ToUpper(entry.Level.String()), timestamp, entry.Message)
	}

	return []byte(levelText), nil
}

func init() {
	logrus.SetFormatter(&MyTextFormatter{})
}
