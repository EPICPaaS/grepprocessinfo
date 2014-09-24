package main

import (
	"log"
	"os/exec"
	"strings"
)

//执行iotop命令
func execIOTOP(result chan string) {
	//top := exec.Command("sudo", "iotop", "-kboPn2", "-d0.5")
	top := exec.Command("iotop", "-kboPn2", "-d0.5")
	out, err := top.CombinedOutput()
	if err != nil {
		log.Println("execIOTOP", err)
		result <- ""
		return
	}
	s := string(out)
	result <- s
}

//分析给定进程ID在iotop 的output的磁盘io读写情况
func fetchDISKIO(pid string, iotopout string) (diskRead string, diskWrite string) {
	s := iotopout

	if len(strings.TrimSpace(s)) == 0 {
		return "-2.0", "-2.0"
	}

	lines := strings.Split(s, "\n")
	var nline string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, pid+" ") {
			nline = line
		}
	}

	if len(nline) == 0 {
		return "0.0", "0.0"
	}

	fields := strings.Split(nline, " ")
	var newFields []string
	//过滤掉空白field
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}
		newFields = append(newFields, field)
	}

	if len(newFields) > 5 {
		return newFields[3], newFields[5]
	}
	return "-1.0", "-1.0"

}
