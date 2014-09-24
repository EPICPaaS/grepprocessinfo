package main

import (
	"log"
	"os/exec"
	"strings"
)

//执行top命令
func execTOP(result chan string) {
	top := exec.Command("top", "-bn2", "-d0.5")
	out, err := top.CombinedOutput()
	if err != nil {
		log.Println("execTOP", err)
		result <- ""
		return
	}
	s := string(out)
	result <- s
}

//分析给定进程ID的进程的CPU使用情况
func fetchCPU(pid string, topout string) (cpu string) {
	s := topout
	if len(strings.TrimSpace(s)) == 0 {
		return "-2.0"
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
		return "0.0"
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
	//for i, field := range newFields {
	//	log.Println(i, field)
	//}
	if len(newFields) > 8 {
		return newFields[8]
	}
	return "-1.0"

}
