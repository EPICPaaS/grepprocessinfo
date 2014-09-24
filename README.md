# grepprocessinfo


A linux process info fetch tool written in Golang. 

It can fetch the cpu usage, memory usage, net sent, net received, disk write, disk read, number of open net connection, number of open file PER every process that configed in [config.cfg] file




## How to use It
It's based on linux command: {top, iotop, lsof, nethogs}, so if you wanna run it, please install them  firstly. And then following below steps:

### 1. Get me

	go get github.com/EPICPaaS/grepprocessinfo
	
	mkdir /opt/grepprocessinfo
	cp $GOPATH/bin/grepprocessinfo /opt/grepprocessinfo/
	cp $GOPATH/src/github.com/EPICPaaS/grepprocessinfo/config.cfg /opt/grepprocessinfo/
	cp $GOPATH/src/github.com/EPICPaaS/grepprocessinfo/grepprocessinfo.sh /opt/grepprocessinfo/
	cp $GOPATH/src/github.com/EPICPaaS/grepprocessinfo/grep_cron /etc/cron.d/

### 2. Init database
    
Create a database on your mysql server, and then execute following cmd

	mysql --host=$DB_HOST --port=$DB_PORT --user=$DB_USER --password=$DB_PASSWORD $DB_NAME < $GOPATH/src/github.com/EPICPaaS/grepprocessinfo/init.sql

### 3. Config

 vi /opt/grepprocessinfo/config.cfg, replace with your right info

	....
	[mysql]
	url=root:123456@tcp(127.0.0.1:3306)/platform
	...
	
### 4. Restart crond service

	service crond restart

The system will execute grepprocessinfo program every 10 seconds.




