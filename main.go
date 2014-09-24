// grepprocessinfo project main.go

package main

import (
	"errors"
	"github.com/msbranco/goconfig"
	//"github.com/pelletier/go-toml"
	"log"
	"net"
	"os"
	"os/exec"
	"path"
	"runtime"
	"strings"
	"time"
)

var TIME_NOW time.Time

var LOCAL_IP string = "127.0.0.1"
var t_ipv4 string = " IPv4 "
var t_ipv6 string = " IPv6 "
var t_unix string = " unix "

var configFile *goconfig.ConfigFile

type Info struct {
	pid         string
	port        string
	cpu         string
	memory      string
	moduleName  string
	ip          string
	diskRead    string
	diskWrite   string
	netsent     string
	netreceived string
	nonetconn   string
	nofile      string
	sysnofile   string
	dubboport   string
	createdTime time.Time
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func localIP() (net.IP, error) {
	tt, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, t := range tt {
		aa, err := t.Addrs()
		if err != nil {
			return nil, err
		}
		for _, a := range aa {
			ipnet, ok := a.(*net.IPNet)
			if !ok {
				continue
			}
			v4 := ipnet.IP.To4()
			if v4 == nil || v4[0] == 127 { // loopback address
				continue
			}
			return v4, nil
		}
	}
	return nil, errors.New("cannot find local IP address")
}

func pipe_commands(commands ...*exec.Cmd) ([]byte, error) {
	for i, command := range commands[:len(commands)-1] {
		out, err := command.StdoutPipe()
		checkErr(err)
		command.Start()
		commands[i+1].Stdin = out
	}

	final, err := commands[len(commands)-1].Output()
	if err != nil {
		log.Println("fetch failed")
	}
	return final, nil
}

/**
func handleOtherJavaLine(line string) Info {
	fields := strings.Split(line, " ")
	var newFields []string
	//过滤掉空白field
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}
		newFields = append(newFields, field)
	}

	var info Info
	info.ip = LOCAL_IP
	info.createdTime = TIME_NOW //将同一次统计生成的数据时间统一掉

	//处理运行容器启动的JAVA进程
	for index, field := range newFields {
		//fmt.Println(index, " ", field)
		if index == 1 {
			info.pid = field
		}

		if index == 5 {
			info.memory = field
		}
	}

	info.moduleName = "unknown-java"
	return info
}
**/

func handleLine(line string, moduleName string, port string) Info {
	fields := strings.Split(line, " ")
	var newFields []string
	//过滤掉空白field
	for _, field := range fields {
		field = strings.TrimSpace(field)
		if len(field) == 0 {
			continue
		}
		newFields = append(newFields, field)
	}

	var info Info
	info.ip = LOCAL_IP
	info.createdTime = TIME_NOW //将同一次统计生成的数据时间统一掉

	info.moduleName = moduleName
	info.port = port

	isNoderunner := strings.Contains(line, "com.yuanxin.paas.mgmt.noderunner.Main")

	for index, field := range newFields {
		//fmt.Println(index, " ", field)
		if index == 1 {
			info.pid = field
		}

		if index == 5 {
			info.memory = field
		}

		if isNoderunner {
			if strings.Contains(field, "-Dmodule=") {
				info.moduleName = strings.Split(field, "=")[1]
			}
			if strings.Contains(field, "-Dport=") {
				info.port = strings.Split(field, "=")[1]
			}
			if strings.Contains(field, "-Ddubbo.protocol.port=") {
				info.dubboport = strings.Split(field, "=")[1]
			}
		}
	}

	//通过后期异步方式获取
	//info.cpu = execFetchCPU(info.pid)

	return info
}

/**

func execFetchCPU(pid string) (cpu string) {
	top := exec.Command("top", "-n", "1", "-bp", pid)
	out, err := top.CombinedOutput()
	if err != nil {
		log.Println("execFetchCPU", err)
		return "-1.0"
	}
	s := string(out)
	lines := strings.Split(s, "\n")
	var nline string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, pid+" ") {
			nline = line
		}
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
**/

func AssignMoreInfo(list []Info, topout, iotopout, lsofout string) []Info {
	var newList []Info
	for _, info := range list {
		info.cpu = fetchCPU(info.pid, topout)
		info.diskRead, info.diskWrite = fetchDISKIO(info.pid, iotopout)
		info.nonetconn, info.nofile, info.sysnofile = fetchLSOF(info.pid, lsofout)
		//log.Println(info)
		newList = append(newList, info)
	}
	return newList
}

func execFetchProcessList() []Info {
	ps := exec.Command("ps", "aux")
	grep2 := exec.Command("grep", "-v", "grep")

	out, err := pipe_commands(ps, grep2)
	checkErr(err)
	s := string(out)
	//fmt.Println(s)
	lines := strings.Split(s, "\n")

	grepProcessesString, err := configFile.GetString("main", "grepProcesses")
	var checkedList []string
	if err == nil {
		checkedList = strings.Split(grepProcessesString, ",")
	}

	var list []Info
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}

		hasHandled := false

		for _, item := range checkedList {
			item = strings.TrimSpace(item)
			if len(item) == 0 {
				continue
			}

			cmdkey, err := configFile.GetString(item, "cmdkey")
			if err != nil {
				continue
			}
			port, err := configFile.GetString(item, "port")
			if err != nil {
				continue
			}

			if strings.Contains(line, cmdkey) {
				r := handleLine(line, item, port)
				list = append(list, r)
				hasHandled = true
				break
			}
		}

		if !hasHandled && strings.Contains(line, "java") {
			r := handleLine(line, "unkown-java", "0")
			list = append(list, r)
			hasHandled = true
		}
	}
	return list
}

func init() {
	os.Chdir(path.Dir(os.Args[0]))
	AppPath := path.Dir(os.Args[0])
	AppConfigPath := path.Join(AppPath, "config.cfg")
	var err error
	configFile, err = goconfig.ReadConfigFile(AppConfigPath)
	checkErr(err)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	ip, err := localIP()
	checkErr(err)
	LOCAL_IP = ip.String()
	TIME_NOW = time.Now()
	ioresult := make(chan string)
	topresult := make(chan string)
	lsofresult := make(chan string)

	go execIOTOP(ioresult)
	go execTOP(topresult)
	go execLSOF(lsofresult)

	dbType, err := configFile.GetString("main", "dbType")
	checkErr(err)

	list := execFetchProcessList()

	topout := <-topresult
	iotopout := <-ioresult
	lsofout := <-lsofresult
	//log.Println(lsofout)

	list = AssignMoreInfo(list, topout, iotopout, lsofout)

	//for _, info := range list {
	//	log.Println(info)
	//}

	if dbType == "mysql" {
		url, err := configFile.GetString("mysql", "url")
		checkErr(err)
		Save2MySQL(list, url)
		/*
			} else if dbType == "oracle" {
				url, err := configFile.GetString("oracle", "url")
				checkErr(err)
				Save2Oracle(list, url)
		*/
	} else {
		/**
		DB_ENTRY_DBNAME=platform
		DB_ENTRY_USER_PASSWORD=123456
		DB_ENTRY=10.180.120.64:3306
		DB_ENTRY_USER=root
		**/
		dbEntry := os.Getenv("DB_ENTRY")
		dbName := os.Getenv("DB_ENTRY_DBNAME")
		dbUser := os.Getenv("DB_ENTRY_USER")
		dbPassword := os.Getenv("DB_ENTRY_USER_PASSWORD")

		if dbEntry == "" || dbName == "" || dbUser == "" || dbPassword == "" {
			log.Println("can not get db env info")
		} else {
			url := dbUser + ":" + dbPassword + "@tcp(" + dbEntry + ")/" + dbName
			Save2MySQL(list, url)
		}

	}

	log.Println("Finished!")
}
