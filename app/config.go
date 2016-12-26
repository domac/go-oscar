package app

import (
	"fmt"
	js "github.com/bitly/go-simplejson"
	"github.com/codegangsta/cli"
	client "github.com/domac/go-oscar/httpclient"
	"io/ioutil"
	"strings"
)

var ChannelCode = GenChannelCode()

//数据配置的接口
type DataConfig interface {
	GetData() ([]string, error)
	CreateCommads(servers []string) ([]string, map[string]string, error)
	HandleError(map[string]string)
}

//继承DataConfig接口,实现测试的功能
type DemoDataConfig struct {
	cli  *cli.Context
	pkgs string
	port string
}

func NewDemoDataConfig(c *cli.Context) *DemoDataConfig {
	dd := &DemoDataConfig{
		cli: c,
	}
	dd.pkgs = c.String("pkg")
	dd.port = ":" + c.String("port")
	return dd
}

func (dd *DemoDataConfig) GetData() (results []string, err error) {
	c := dd.cli
	configPath := c.String("config")
	server_ip := c.String("server")

	if configPath != "" {
		results, err = ReadLine(configPath)
		if err != nil {
			return nil, err
		}
	} else if server_ip != "" {
		results = strings.Split(server_ip, ",")
	}
	return results, nil
}

func (dd *DemoDataConfig) CreateCommads(ips []string) (urls []string, query map[string]string, err error) {
	action := ""
	query = make(map[string]string)
	query["channel_code"] = ChannelCode
	query["pwd"] = "qwertyuiopaz"

	if dd.pkgs == "" {
		query["command"] = "install_app"
		action = "command"
	} else {
		query["name"] = dd.pkgs
		action = "pkg/setup.do"
	}

	ips = RemoveDuplicatesAndEmpty(ips)

	for _, ip := range ips {
		url := fmt.Sprintf("http://%s%s/%s", ip, dd.port, action)
		urls = append(urls, url)
	}
	return urls, query, nil
}

//错误处理方法
func (dd *DemoDataConfig) HandleError(failservers map[string]string) {
	records := []string{}
	for k, _ := range failservers {
		k = strings.Replace(k, "http://", "", 1)
		k = strings.Replace(k, dd.port+"/command", "", 1)
		records = append(records, k)
	}

	//错误记录处理
	record_path := dd.cli.String("error_record")
	if len(records) > 0 {
		WriteIntoFile(record_path, records, WRITE_OVER)
		fmt.Printf("\nError records detail, please see : %s\n", record_path)
	} else {
		RemoveFile(record_path)
		fmt.Printf("\nDone ! \n")
	}

}
