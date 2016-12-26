go-oscar
===========

## 如何使用？

> 进入release文件夹进行下载 go-oscar 执行文件下载

目前提供 `mac` 和 `linux` 的版本

```sh
$ ./go-oscar -h

//参数输出：

NAME:
   go-oscar - a tool which help us to create tasks for executing commands

USAGE:
   go-oscar [global options] command [command options] [arguments...]

VERSION:
   1.0.0

COMMANDS:
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --error_record "/tmp/error_record.txt"       the path of error servers record
   -i "5"                                       interval of worker execute
   --pkg                                        install packages where store in /data/res/registry;
                                                1) eg: --pkg=kafka-python-0.9.0.tar.gz,requests-2.10.0.tar.gz
                                                2) if -pkg value is null, just request to app update else request to update packages
   --server                                     only update for single server; if -config param is being aready used, this param is no effect. eg -server 127.0.0.1
   -w "10000"                                   num of worker num
   --timeout "15000"                            timeout of getting data
   --config                                     the path of servers ip config file
   --port "8000"                               target service port
   --help, -h                                   show help
   --version, -v                                print the version
```

### 执行输出：

执行中状态

```
task executed start !
http://10.205.83.23:8000/command        [====================================================>-------------------------] 68%
```

执行完成状态

```
task executed start !
http://10.205.83.23:8000/command        [==============================================================================] 100%

Error records detail, please see : /tmp/error_record.txt


[EXECUTOR RESULT]
+-----------+---------------+---------------+--------------+----------------------+-----------------------+
| JOB COUNT | THREADS COUNT | SUCCESS COUNT | ERRORS COUNT | TASKS COST (SECONDS) | REPORT COST (SECONDS) |
+-----------+---------------+---------------+--------------+----------------------+-----------------------+
|       178 |           178 |             6 |          172 |            16.393963 |              0.001074 |
+-----------+---------------+---------------+--------------+----------------------+-----------------------+
```

如果执行过程中有错误的情况,会把服务器信息记录在 命令`error_record`指定的文件上 (默认为: /tmp/error_record.txt)