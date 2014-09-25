package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

//执行nethogs命令, 因nethogs有bug，第一次采集不到数据，所以至少需要2秒以上才能采集到数据
func execNethogs(result chan []string) {
	filepath := os.TempDir() + "/" + strconv.FormatInt(time.Now().UnixNano(), 10)
	//defer os.Remove(filepath)

	niName, err := configFile.GetString("main", "networkInterfaceName")
	if err != nil {
		niName = "eth0"
	}
	top := exec.Command("nethogs", niName, "-d1")
	tee := exec.Command("tee", filepath)

	out, err := top.StdoutPipe()
	if err != nil {
		log.Println("execNethogs", err)
		result <- []string{}
		return
	}
	err = top.Start()
	if err != nil {
		log.Println("execNethogs|nethogs", err)
		result <- []string{}
		return
	}
	tee.Stdin = out
	time.Sleep(time.Second * 2)

	top.Process.Kill()
	//log.Println(out)
	err = tee.Start()
	if err != nil {
		log.Println("execNethogs|tee", err)
		result <- []string{}
		return
	}
	time.Sleep(time.Millisecond * 100)
	tee.Process.Kill()

	lines := filterNethogsOutput(filepath)

	result <- lines
}

//分析给定进程ID在nethogs 的output的网络IO发送和接收情况
func fetchNethogs(pid string, nethogsOut []string) (netsent string, netreceived string) {

	if len(nethogsOut) == 0 {
		return "-2.0", "-2.0"
	}

	var nline string
	for _, line := range nethogsOut {
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

	nl := len(newFields)

	if nl > 4 {
		return newFields[nl-3], newFields[nl-2]
	}
	return "-1.0", "-1.0"

}

//预处理nethogs输出
func filterNethogsOutput(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		return []string{}
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var a string
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		line = strings.Replace(line, "\x1B", "", -1)    //去掉ESC
		line = strings.Replace(line, "(B[m", "\n", -1)  //换行
		line = strings.Replace(line, "(B[0;7m", "", -1) //去掉

		a = a + line + "\n"
	}

	lines := strings.Split(a, "\n")
	var r []string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "[?") || strings.HasPrefix(line, "PID") || strings.HasPrefix(line, "TOTAL") || strings.Contains(line, "NetHogs version") || strings.Contains(line, "Ethernet link detected") {
			continue
		}

		s := strings.Index(line, "[")
		e := strings.Index(line, "H")

		for s > 0 && e > 0 {
			//fmt.Println("##>" + line)
			line = line[:s] + " " + line[e+1:]
			//fmt.Println(s, e, "==>"+line)
			s = strings.Index(line, "[")
			e = strings.Index(line, "H")
		}

		s2 := strings.Index(line, "[")
		e2 := strings.Index(line, "G")
		for s2 > 0 && e2 > 0 {
			//fmt.Println("###>" + line)
			line = line[:s2] + " " + line[e2+1:]
			//fmt.Println(s2, e2, "===>"+line)
			s2 = strings.Index(line, "[")
			e2 = strings.Index(line, "G")
		}

		r = append(r, line)
	}

	return r
}
