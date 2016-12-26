package app

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/domac/go-oscar/core"
	"github.com/dustin/go-humanize"
	"github.com/olekukonko/tablewriter"
	"os"
	"sync"
	"time"
)

//应用版本
const APP_VERSION = "1.0.0"

//参数动作处理
func appAction(c *cli.Context) (err error) {

	var wg sync.WaitGroup
	var mpwg sync.WaitGroup

	timeout := c.Int("timeout")

	demoDataConfig := NewDemoDataConfig(c)
	ips, err := demoDataConfig.GetData()
	urls, requestParams, err := demoDataConfig.CreateCommads(ips)

	if nil != err {
		return err
	}

	//工作worker数
	work_num := c.Int("w")

	interval := c.Int("i")
	itvl := time.Duration(interval)

	d := core.NewDispatcherWithMQ(work_num, work_num, &wg, &mpwg)

	//r := NewReport()
	r := NewReportWithProgress(len(urls))
	d.SetMF(GenerateMessageReportMethod(r))

	//启动调度器
	d.RunWithLimiter(itvl * time.Millisecond)
	defer d.Stop()

	wg.Add(1)
	mpwg.Add(1)

	start := time.Now()

	go func() {
		fmt.Println("task executed start !")
		channel_code := GenChannelCode()
		for i, url := range urls {
			id := fmt.Sprintf("id-%d", i)
			job := NewJob(id, url, requestParams, channel_code, timeout, 3)
			t := core.CreateTask(job, "Do")
			//提交请求任务
			d.SubmitTask(t)
		}
		wg.Done()
		mpwg.Done()
	}()
	wg.Wait()
	mpwg.Wait()

	time.Sleep(500 * time.Millisecond)

	//调度耗时
	cost := time.Now().Sub(start).Seconds()

	//报告错误的服务
	report_start := time.Now()

	r.CallbackFail(demoDataConfig.HandleError)
	//错误上报耗时
	report_cost := time.Now().Sub(report_start).Seconds()

	//结果报表
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"job count", "threads count", "success count", "errors count", "tasks cost (seconds)", "report cost (seconds)"})

	data := []string{
		fmt.Sprintf("%d", len(urls)),
		fmt.Sprintf("%d", Wgroutinue),
		fmt.Sprintf("%d", len(r.Data)),
		fmt.Sprintf("%d", len(r.ErrServer)),
		humanize.Ftoa(cost),
		humanize.Ftoa(report_cost),
	}
	table.Append(data)
	fmt.Println("\n")
	fmt.Println("[EXECUTOR RESULT]")
	table.Render() // Send output
	return nil
}

func Main() {
	FlagsInit()
	app := cli.NewApp()
	app.Name = "go-oscar"
	app.Usage = "a tool which help us to create tasks for executing commands"
	app.Version = APP_VERSION
	app.Flags = GetAppFlags()
	app.Action = ActionWrapper(appAction)
	osArgs := os.Args
	if len(osArgs) == 1 {
		osArgs = append(osArgs, "-h")
	}
	app.Run(osArgs)
}
