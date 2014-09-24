package main

//后续版本默认编译不支持oracle 入库 2014-09-12 11:00:00
/**
import (
	"database/sql"
	_ "github.com/tgulacsi/goracle/godrv"
	//"github.com/tgulacsi/goracle/oracle"
	//"log"
	"strconv"
)

func Save2Oracle(infos []Info, url string) bool {
	//dsn := oracle.MakeDSN("10.180.120.95", 1521, "", "hatest")
	//log.Println(dsn)

	db, err := sql.Open("goracle", url)
	checkErr(err)
	defer db.Close()

	tx, err := db.Begin()
	checkErr(err)

	// Prepare statement for inserting data
	stmtIns, err := tx.Prepare("INSERT INTO t_logs(id,ip,port,pid,cpu,memory,module,time,diskread,diskwrite,netsent,netreceived,dubboport) VALUES( logs_seq.nextval, ?, ?,? ,?, ?,? ,?,?,?,?,?,?)")
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
		_, err = stmtIns.Exec(info.ip, port, pid, cpu, memory_mb, info.moduleName, info.createdTime.Format("2006-01-02 15:04:05"), diskread, diskwrite, netsent, netreceived, dubboport)
		checkErr(err)
	}
	err = tx.Commit()
	checkErr(err)
	return true
}
*/
