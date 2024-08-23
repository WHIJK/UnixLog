package utils

import (
	"UnixLog/model"
	"UnixLog/pkg/secure"
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"os"
)

/*
@Author: OvO
@Date: 2024/8/21 14:34
*/

/*
ReadFileByLine
@Description: 逐行读取
@param filename
@param lenFile
@param taskChan
*/
func ReadFileByLine(filename string, lenFile *int, taskChan chan string) {
	file, err := os.Open(filename)

	if err != nil {
		logrus.Fatal(err)
	}
	read := bufio.NewScanner(file)

	if read.Err() != nil {
		logrus.Fatal(read.Err())
	}

	for read.Scan() {
		taskChan <- read.Text()
		*lenFile++
	}
	close(taskChan)
}

/*
WriteInCsv
@Description: Audit写入csv
@param data
@param output
@param headers
*/
func WriteInCsv(data []map[string]string, output string, headers []string) {
	file, err := os.Create(output)
	if err != nil {
		logrus.Error(err)
		return
	}
	defer file.Close()

	w := csv.NewWriter(file)
	defer w.Flush()

	if errW := w.Write(headers); errW != nil {
		logrus.Error(err)
		return
	}

	for _, recordMap := range data {
		record := make([]string, len(headers))
		for i, colName := range headers {
			if value, ok := recordMap[colName]; ok {
				record[i] = value
			} else {
				record[i] = "" // 如果没有匹配的键，写入空值
			}
		}

		if errww := w.Write(record); errww != nil {
			logrus.Error(errww)
			return
		}
	}
	fmt.Printf("\nOutput: %s.csv \n", output)
}

/*
WriteInXlsx
@Description: 写入Secure分析结果至XLSX
@param mapper
@param filename
@param headers
*/
func WriteInXlsx(mapper map[string]model.SecureLogData, filename string, headers []string) {
	// 创建一个新的 XLSX 文件
	f := excelize.NewFile()

	// 创建一个工作表
	sheet := "Sheet1"
	index, _ := f.NewSheet(sheet)

	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", string('A'+i), 1) // 从 A1 开始
		f.SetCellValue(sheet, cell, header)
	}

	// 填充数据
	row := 2 // 从第二行开始
	for _, v := range mapper {
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), v.Hostname)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), v.ProcessName)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), v.Pid)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), parseTime(v.OpenTime))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), parseTime(v.CloseTime))
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), v.IP)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), v.Description)
		row++
	}
	// 失败统计
	FailedSheet := "Failure statistics"
	f.NewSheet(FailedSheet)
	headersFailed := []string{"IP", "Failure Count"}
	for i, header := range headersFailed {
		cell := fmt.Sprintf("%s%d", string('A'+i), 1) // 从 A1 开始
		f.SetCellValue(FailedSheet, cell, header)
	}
	FailedRow := 2 // 从第二行开始
	for ip, count := range secure.FailedCount {
		f.SetCellValue(FailedSheet, fmt.Sprintf("A%d", FailedRow), ip)
		f.SetCellValue(FailedSheet, fmt.Sprintf("B%d", FailedRow), count)
		FailedRow++
	}

	// 设置工作表为活动工作表
	f.SetActiveSheet(index)

	// 保存 XLSX 文件
	if err := f.SaveAs(filename + ".xlsx"); err != nil {
		fmt.Println("Error saving file:", err)
	}
	fmt.Printf("\nOutput: %s.xlsx \n", filename)
}
