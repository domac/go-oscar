package app

import (
	"github.com/codegangsta/cli"
)

//参数初始化
func FlagsInit() {

	AddFlagString(cli.StringFlag{
		Name:  "config",
		Usage: "the path of servers ip config file",
	})

	AddFlagString(cli.StringFlag{
		Name:  "error_record",
		Value: "/tmp/error_record.txt",
		Usage: "the path of error servers record",
	})

	AddFlagString(cli.StringFlag{
		Name:  "port",
		Value: "8000",
		Usage: "target service port",
	})

	AddFlagString(cli.StringFlag{
		Name:  "pkg",
		Usage: "install packages where store in /data/res/registry;\n \t1) eg: --pkg=kafka-python-0.9.0.tar.gz,requests-2.10.0.tar.gz \n \t2) if -pkg value is null, just request to app update else request to update packages ",
	})

	AddFlagString(cli.StringFlag{
		Name:  "server",
		Usage: "only update for single server; if -config param is being aready used, this param is no effect. eg -server 127.0.0.1",
	})

	AddFlagInt(cli.IntFlag{
		Name:  "w",
		Value: 10000,
		Usage: "num of worker num",
	})

	AddFlagInt(cli.IntFlag{
		Name:  "i",
		Value: 5,
		Usage: "interval of worker execute",
	})

	AddFlagInt(cli.IntFlag{
		Name:  "timeout",
		Value: 15000,
		Usage: "timeout of getting data",
	})

}
