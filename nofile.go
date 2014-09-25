package main

import (
	"log"
	"os/exec"
	"strconv"
	"strings"
)

var t_ipv4 string = " IPv4 "
var t_ipv6 string = " IPv6 "
var t_unix string = " unix "

//执行lsof命令
func execLSOF(result chan string) {
	top := exec.Command("lsof", "-bw")
	out, err := top.CombinedOutput()
	if err != nil {
		log.Println("execLSOF", err)
		result <- ""
		return
	}
	s := string(out)
	result <- s
}

//分析给定进程ID在lsof的output的当前进程，打开网络连接，打开文件数，系统总打开文件数情况
func fetchLSOF(pid string, lsofout string) (nonetconn, nofile, sysnofile string) {
	s := lsofout

	if len(strings.TrimSpace(s)) == 0 {
		return "-2", "-2", "-2"
	}

	sysopenfile := -1 //lsof 第一行排除
	openfile := 0
	opennetconn := 0

	lines := strings.Split(s, "\n")
	var nline string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) > 0 {
			sysopenfile += 1
		}
		if strings.Contains(line, " "+pid+" ") {
			nline = line
		} else {
			continue
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

		if len(newFields) > 5 && newFields[1] == pid {
			openfile += 1
			if strings.Count(nline, t_ipv4) > 0 || strings.Count(nline, t_ipv6) > 0 || strings.Count(nline, t_unix) > 0 {
				opennetconn += 1
			}
		}
	}

	return strconv.Itoa(opennetconn), strconv.Itoa(openfile), strconv.Itoa(sysopenfile)
}
