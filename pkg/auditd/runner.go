package auditd

import (
	"UnixLog/model"
	"github.com/cheggaaa/pb/v3"
	"sync"
)

/*
@Author: OvO
@Date: 2024/8/21 16:34
*/

func Runner(wg *sync.WaitGroup, taskChan chan string, bar *pb.ProgressBar, sm *model.SyncMaps) {
	bar.Start()

	for v := range taskChan {
		auditMap := parseAudit(v, sm)
		sm.Lock.Lock()
		sm.ResultMaps = append(sm.ResultMaps.([]map[string]string), auditMap)
		sm.Lock.Unlock()
		bar.Increment()
	}
	wg.Done()

}
