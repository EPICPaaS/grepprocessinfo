package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
)

func Save2MySQL(infos []Info, url string) bool {
	db, err := sql.Open("mysql", url)
	checkErr(err)
	defer db.Close()

	tx, err := db.Begin()
	checkErr(err)

	// Prepare statement for inserting data
	stmtIns, err := tx.Prepare("INSERT INTO t_logs(ip,port,pid,cpu,memory,module,time,diskread,diskwrite,netsent,netreceived,dubboport,nonetconn,nofile,sysnofile) VALUES( ?, ?,? ,?, ?,? ,?,?,?,?,?,?,?,?,?)")
	checkErr(err)
	defer stmtIns.Close()

	for _, info := range infos {
		port, _ := strconv.Atoi(info.port)
		pid, _ := strconv.Atoi(info.pid)
		cpu, _ := strconv.ParseFloat(info.cpu, 32)
		diskread, _ := strconv.ParseFloat(info.diskRead, 32)
		diskwrite, _ := strconv.ParseFloat(info.diskWrite, 32)
		netsent, _ := strconv.ParseFloat(info.netsent, 32)
		netreceived, _ := strconv.ParseFloat(info.netreceived, 32)
		dubboport, _ := strconv.Atoi(info.dubboport)

		memory, _ := strconv.Atoi(info.memory)
		memory_mb := memory / 1024

		nonetconn, _ := strconv.Atoi(info.nonetconn)
		nofile, _ := strconv.Atoi(info.nofile)
		sysnofile, _ := strconv.Atoi(info.sysnofile)

		_, err = stmtIns.Exec(info.ip, port, pid, cpu, memory_mb, info.moduleName, info.createdTime.Format("2006-01-02 15:04:05"), diskread, diskwrite, netsent, netreceived, dubboport, nonetconn, nofile, sysnofile)
		checkErr(err)
	}

	err = tx.Commit()
	checkErr(err)
	return true
}
