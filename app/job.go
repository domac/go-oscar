package app

import (
	"encoding/json"
	"github.com/domac/go-oscar/core"
	client "github.com/domac/go-oscar/httpclient"
	"io/ioutil"
)

var Wgroutinue int = 0

type Result struct {
	Code    int64       `json:code`
	Message string      `json:message`
	Object  interface{} `json:object`
	Success bool        `json:success`
	Server  string
}

//任务作业结构
type Job struct {
	Id          string
	ChannelCode string
	Server      string
	params      map[string]string
	Result      *Result
	fail        bool
	failMessage string
	Timeout     int
	Retries     int
}

func NewJob(id string, server string, params map[string]string, channel_code string, timeout int, retries int) *Job {
	return &Job{
		Id:          id,
		Server:      server,
		Timeout:     timeout,
		Retries:     retries,
		params:      params,
		ChannelCode: channel_code}
}

func (j *Job) SetFailure(message string) {
	j.fail = true
	j.failMessage = message
}

func (j *Job) Do() {
	Wgroutinue++
	//请求调用
	jobClient := client.NewHttpClient()
	jobClient.Defaults(client.Map{
		"opt_timeout_ms":        j.Timeout,
		"opt_connecttimeout_ms": j.Timeout,
	})
	resp, err := jobClient.Get(j.Server, j.params)

	if err != nil || resp == nil || resp.Body == nil {
		j.SetFailure(err.Error())
		return
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		j.SetFailure(err.Error())
		return
	}

	//作业执行结果
	jobRes := new(Result)
	if err = json.Unmarshal(res, jobRes); err != nil {
		j.SetFailure("the return value is parse error")
		return
	}
	resp.Body.Close()
	jobRes.Server = j.Server
	j.Result = jobRes
}

//上报方法回调
func GenerateMessageReportMethod(r *Report) core.MF {
	return func(task core.Task) {
		tj := (*task.TargetObj).(*Job)
		r.IncrProgress(tj.Server)
		if !tj.fail {
			r.AddResult(tj.Result)
		} else {
			r.AddErrorResult(tj.Server, tj.failMessage)
		}
	}
}
