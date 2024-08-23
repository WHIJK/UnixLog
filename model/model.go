package model

import (
	"github.com/cheggaaa/pb/v3"
	"sync"
)

/*
@Author: OvO
@Date: 2024/8/22 21:01
*/

type SecureLogData struct {
	Hostname    string
	ProcessName string
	Pid         string
	OpenTime    string
	CloseTime   string
	IP          string
	Description string
}

type SyncMaps struct {
	Lock       *sync.RWMutex
	Headers    map[string]bool //
	ResultMaps interface{}
}

type Abc = func(wg *sync.WaitGroup, taskChan chan string, bar *pb.ProgressBar, sm *SyncMaps)
