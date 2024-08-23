package cmd

import (
	"UnixLog/model"
	"UnixLog/pkg"
	"UnixLog/pkg/auditd"
	"UnixLog/pkg/secure"
	"UnixLog/pkg/utils"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"sync"
)

/*
@Author: OvO
@Date: 2024/8/20 15:56
*/

var (
	rootCMD = &cobra.Command{
		Short: "Analysis log",
	}
	files []string
	sm    = &model.SyncMaps{}

	// audit审计日志
	auditCMD = &cobra.Command{
		Use:   "audit",
		Short: "Analysis audit log",
		Run: func(cmd *cobra.Command, ags []string) {
			lo.Filter(files, func(f string, index int) bool {
				if !lo.IsEmpty(f) {
					sm.Headers = make(map[string]bool)
					sm.Lock = &sync.RWMutex{}
					sm.ResultMaps = make([]map[string]string, 0)
					auditMaps, s := pkg.RunTask(f, sm, auditd.Runner)
					utils.WriteInCsv(auditMaps.([]map[string]string), f+".csv", s)
					return true
				}
				return false
			})
		},
	}

	// secure日志
	secureCMD = &cobra.Command{
		Use:   "secure",
		Short: "Analysis secure log",
		Run: func(cmd *cobra.Command, args []string) {
			lo.Filter(files, func(f string, index int) bool {
				if !lo.IsEmpty(f) {
					sm.Lock = &sync.RWMutex{}
					secureMap, _ := pkg.RunTask(f, sm, secure.Runner)
					utils.WriteInXlsx(secureMap.(map[string]model.SecureLogData), f, []string{"Hostname", "ProcessName", "Pid", "OpenTime", "CloseTime", "IP", "Description"})
					return true
				}
				return false
			})
		},
	}
)

func init() {
	rootCMD.AddCommand(auditCMD)
	rootCMD.AddCommand(secureCMD)
	rootCMD.PersistentFlags().StringSliceVarP(&files, "file", "f", []string{}, "multiple file")
}

func Execute() {
	rootCMD.Execute()
}
