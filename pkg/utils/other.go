package utils

import (
	"fmt"
	"time"
)

/*
@Author: OvO
@Date: 2024/8/22 20:19
*/

// 解析日期字符串并返回标准格式的日期时间
func parseTime(inputTime string) string {
	if inputTime == "" {
		return ""
	}

	// 定义输入时间格式
	const inputLayout = "Jan 2 15:04:05"
	const outputLayout = "2006-01-02 15:04:05"

	now := time.Now()
	currentYear := now.Year()

	// 解析日期字符串
	logTime, err := time.Parse(inputLayout, inputTime)
	if err != nil {
		return fmt.Sprintf("Error parsing input time: %v", err)
	}

	// 创建一个新的时间对象，使用当前年份
	logTime = time.Date(currentYear, logTime.Month(), logTime.Day(), logTime.Hour(), logTime.Minute(), logTime.Second(), +8, time.Local)

	// 如果 logTime 在当前时间之后，调整为前一年
	if logTime.After(now) {
		logTime = logTime.AddDate(-1, 0, 0)
	}

	// 返回标准格式的日期时间
	return logTime.Format(outputLayout)
}
