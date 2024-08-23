package auditd

import (
	"UnixLog/model"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
	"time"
)

/*
@Author: OvO
@Date: 2024/8/21 14:32
*/

func parseAudit(line string, sm *model.SyncMaps) map[string]string {
	var AuditMap = make(map[string]string)

	partsFirst := line[:strings.Index(line, "):")+1] // 时间戳前

	partsSecond := line[strings.Index(line, "):")+2:] // 时间戳后
	if strings.Contains(partsSecond, "msg='") {
		partsSecond = strings.TrimRight(strings.Replace(partsSecond, "msg='", "", -1), "'")
	}

	partBySpace := strings.Split(partsFirst+partsSecond, " ")
	for _, val := range partBySpace {
		keyValue := strings.SplitN(val, "=", 2)
		if len(keyValue) != 2 {
			logrus.Errorf("Parsing %s error.", keyValue)
			continue
		}
		// 时间戳解析
		if strings.HasPrefix(keyValue[1], "audit(") {
			keyValue[1] = strings.TrimSuffix(strings.TrimPrefix(keyValue[1], "audit("), ")")
			seconds, err := strconv.ParseInt(strings.Split(keyValue[1], ".")[0], 10, 64)
			if err != nil {
				logrus.Error(err)
			}
			keyValue[1] = time.Unix(seconds, 0).Format("2006-01-02 15:04:05")
		}
		AuditMap[keyValue[0]] = keyValue[1]
		sm.Lock.Lock()
		sm.Headers[keyValue[0]] = true
		sm.Lock.Unlock()

	}
	return AuditMap
}
