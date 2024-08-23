package pkg

import (
	"UnixLog/model"
	"UnixLog/pkg/utils"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"github.com/samber/lo"
	"sync"
)

/*
@Author: OvO
@Date: 2024/8/21 16:33
*/

func RunTask(filename string, sm *model.SyncMaps, f model.Abc) (interface{}, []string) {
	var (
		taskChan   = make(chan string)
		taskLength = 0
		wg         = sync.WaitGroup{}
		bar        = pb.New(-1)
	)
	go utils.ReadFileByLine(filename, &taskLength, taskChan)
	go func() {
		for {
			bar.SetTotal(int64(taskLength))
		}
	}()

	for i := 0; i < 40; i++ {
		wg.Add(1)
		go f(&wg, taskChan, bar, sm)
	}
	wg.Wait()
	s := lo.MapToSlice(sm.Headers, func(k string, value bool) string {
		return fmt.Sprintf("%s", k)
	})
	bar.Finish()
	return sm.ResultMaps, s
}
