package secure

import (
	"UnixLog/model"
	"github.com/cheggaaa/pb/v3"
	"sync"
)

/*
@Author: OvO
@Date: 2024/8/22 20:16
*/

func Runner(wg *sync.WaitGroup, taskChan chan string, bar *pb.ProgressBar, sm *model.SyncMaps) {
	bar.Start()
	sm.ResultMaps = make(map[string]model.SecureLogData)
	for v := range taskChan {
		analysisSecure(v, sm, bar)
		bar.Increment()
	}
	wg.Done()
}
