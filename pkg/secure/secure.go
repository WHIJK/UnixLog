package secure

import (
	"UnixLog/model"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"regexp"
	"slices"
	"strings"
)

/*
@Author: OvO
@Date: 2024/8/22 20:14
*/

// 获取日志中如下进程
var processList = []string{"sshd", "useradd", "pkexec", "sudo", "userdel", "su", "groupadd"}

// 统计登录失败的次数

var FailedCount = make(map[string]int)

/*
parseSecure
@Description: 解析单行日志
*/
func parseSecure(line string) (string, string, string, string, string) {
	parts := strings.Fields(line)
	if len(parts) < 4 {
		fmt.Println("Invalid log line format")
		return "", "", "", "", ""
	}

	// 组成时间戳
	// 组合时间戳
	month := parts[0]
	day := parts[1]
	time := parts[2]
	timestamp := fmt.Sprintf("%s %s %s", month, day, time)
	rest := strings.Join(parts[3:], " ")

	// 找到进程信息和消息的分隔符位置
	colonIndex := strings.Index(rest, ":")
	if colonIndex == -1 {
		fmt.Println("Invalid log line format")
		return "", "", "", "", ""
	}

	// 提取主机名
	hostname := strings.SplitN(rest, " ", 2)[0]

	// 提取进程信息
	processInfo := rest[len(hostname)+1 : colonIndex]
	// 解析进程名和 PID
	var processName string
	var pid string

	// 查找方括号
	leftBracket := strings.Index(processInfo, "[")
	rightBracket := strings.Index(processInfo, "]")

	if leftBracket != -1 && rightBracket != -1 && rightBracket > leftBracket {
		// 有 PID
		processName = processInfo[:leftBracket]
		pid = processInfo[leftBracket+1 : rightBracket]
	} else {
		// 没有 PID
		processName = processInfo
		pid = ""
	}

	message := rest[colonIndex+2:] // 从 ": " 后面的部分开始是消息内容

	return timestamp, hostname, processName, pid, message
}

/*
analysisSecure
@Description:
@param v
@param sm
@param bar
*/
func analysisSecure(v string, sm *model.SyncMaps, bar *pb.ProgressBar) {
	if lo.IsEmpty(v) {
		return
	}

	timestamp, hostname, process, pid, message := parseSecure(v)

	if !slices.Contains(processList, process) || lo.IsEmpty(timestamp) {
		return
	}
	// pid空白, 则生成key
	var pidRand string
	if lo.IsEmpty(pid) {
		pidRand = uuid.New().String()
	} else {
		pidRand = pid
	}
	sm.Lock.Lock()
	logData, _ := sm.ResultMaps.(map[string]model.SecureLogData)[pidRand]
	logData.Hostname = hostname
	logData.ProcessName = process
	logData.Pid = pid
	switch strings.ToLower(process) {
	case "sshd": // ssh会话在线时间
		if strings.Contains(message, "Accepted") {
			var re2 = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`) // 匹配IP
			logData.IP = re2.FindAllString(message, -1)[0]
			logData.Description = message
			logData.OpenTime = timestamp
		} else if strings.Contains(strings.ToLower(message), "session opened") {
			if logData.OpenTime == "" {
				logData.OpenTime = timestamp
			}

			if logData.Description == "" {
				logData.Description = message
			}

		} else if strings.Contains(strings.ToLower(message), "session closed") {

			if logData.CloseTime == "" {
				logData.CloseTime = timestamp
			}
			if logData.Description == "" {
				logData.Description = message
			}
		} else if strings.Contains(strings.ToLower(message), "failed password") {
			var re2 = regexp.MustCompile(`\d+\.\d+\.\d+\.\d+`) // 匹配IP
			faiedIP := re2.FindAllString(message, -1)[0]
			FailedCount[faiedIP]++
			// 失败的不写入文件
			delete(sm.ResultMaps.(map[string]model.SecureLogData), pidRand)
			goto Skip
		} else {
			// 如果不符合则删除map, 跳过赋值
			delete(sm.ResultMaps.(map[string]model.SecureLogData), pidRand)
			goto Skip
		}

	default:
		logData.Description = message
		logData.OpenTime = timestamp

	}

	sm.ResultMaps.(map[string]model.SecureLogData)[pidRand] = logData
Skip:
	sm.Lock.Unlock()
}
