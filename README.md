# grepprocessinfo


A linux process info fetch tool written in Golang. 

It can fetch the cpu usage, memory usage, net sent, net received, disk write, disk read, number of open net connection, number of open file PER every process that configed in [config.cfg] file




## How to use It
It's based on linux command: {top, iotop, lsof, nethogs}, so if you wanna run it, please install them  firstly.

	go get github.com/EPICPaaS/grepprocessinfo
	
	mkdir /opt/grepprocessinfo
	cp $GOPATH/bin/grepprocessinfo /opt/grepprocessinfo/
	cp $GOPATH/src/github.com/EPICPaaS/grepprocessinfo/config.cfg /opt/grepprocessinfo/
	cp $GOPATH/src/github.com/EPICPaaS/grepprocessinfo/grepprocessinfo.sh /opt/grepprocessinfo/
	cp $GOPATH/src/github.com/EPICPaaS/grepprocessinfo/grep_cron /etc/cron.d/

And then, vi /opt/grepprocessinfo/config.cfg, replace with your right info

	....
	[mysql]
	url=root:123456@tcp(127.0.0.1:3306)/platform
	...
	
And last, restart crond service

	service crond restart

The system will execute grepprocessinfo program every 10 seconds.




