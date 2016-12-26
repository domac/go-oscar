package app

import (
	"fmt"
	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
	"strings"
	"sync"
)

type FailFunction func(map[string]string)

//结果报告
type Report struct {
	mutex         sync.Mutex //互斥锁
	Data          map[string]*Result
	ErrServer     map[string]string
	errorCallback FailFunction
	BarNum        int
	Bar           *uiprogress.Bar
	sendChan      chan string
}

func NewReport() *Report {
	return &Report{
		Data:      make(map[string]*Result),
		ErrServer: make(map[string]string)}
}

func NewReportWithProgress(barnum int) *Report {
	r := NewReport()
	r.BarNum = barnum
	r.sendChan = make(chan string, barnum)
	p := uiprogress.New()
	p.Start()

	step := barnum - 1
	if step == 0 {
		step = 1
	}
	bar := p.AddBar(step).AppendCompleted()
	bar.Width = 80
	bar.PrependFunc(func(b *uiprogress.Bar) string {
		res := <-r.sendChan
		return strutil.Resize(res, 40)
	})
	r.Bar = bar
	return r
}

func (r *Report) AddResult(s *Result) {
	r.mutex.Lock()
	r.Data[s.Server] = s
	r.mutex.Unlock()
}

func (r *Report) AddErrorResult(s string, msg string) {
	r.mutex.Lock()
	r.ErrServer[s] = msg
	r.mutex.Unlock()
}

func (r *Report) IncrProgress(s string) {
	r.Bar.Incr()
	r.sendChan <- s
}

//结果输出
func (r *Report) ReportError() {
	fmt.Println("\n[FAIL RESULT]")
	for k, _ := range r.ErrServer {
		k = strings.Replace(k, "http://", "", 1)
		k = strings.Replace(k, ":8000/command", "", 1)
		fmt.Printf("%s", k)
		fmt.Println()
	}
}

//失效回调
func (r *Report) CallbackFail(f FailFunction) {
	f(r.ErrServer)
}
