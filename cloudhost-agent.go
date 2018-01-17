package main

import (
	"encoding/json"
	"fmt"
	"github.com/sparrc/go-ping"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	//    "net/url"
	"strings"
)

type Config struct {
	Config   []Target
	Url      string
	Instance string
}

type Target struct {
	Target string
	Module string
}

func (t *Target) icmpTarget() {
	var status string
	pinger, err := ping.NewPinger(t.Target)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		return
	}
	pinger.OnFinish = func(stats *ping.Statistics) {
		if stats.PacketLoss != 0 {
			status = "0"
		} else {
			status = "1"
		}
		data := t.Module + "_" + t.Target + " " + status + "\n"
		pushData(strings.Replace(data, ".", "_", -1))
		return
	}

	pinger.Count = 1
	pinger.SetPrivileged(true)
	pinger.Run()

}

func (t *Target) tcpTarget() {
	var status string
	conn, err := net.Dial("tcp", t.Target)
	if err != nil {
		status = "0"
		return
	} else {
		status = "1"
		conn.Close()
	}
	data := t.Module + "_" + strings.Replace(t.Target, ":", "_", -1) + " " + status + "\n"
	pushData(strings.Replace(data, ".", "_", -1))
	return

}

func (t *Target) httpTarget() {
	var status string
	resp, err := http.Get("http://" + t.Target)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode == 200 {
		status = "1"
	} else {
		status = "0"
		return
	}
	data := t.Module + "_" + t.Target + " " + status + "\n"
	pushData(strings.Replace(data, ".", "_", -1))
	return
}

func pushData(s string) {
	//url := c.Url + "/metrics/job/cloudagent/instance/bzgxjh"
	fmt.Println(s)
	resp, err := http.Post(url, "text/plain", strings.NewReader(s))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp)
}

func loadConfig() (c *Config) {
	configFile := "./config.json"
	buf, err := ioutil.ReadFile(configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "File Error: %s\n", err)
	}
	err = json.Unmarshal(buf, &c)
	return c
}

var c = loadConfig()
var url = c.Url + "/metrics/job/cloudagent/instance/" + c.Instance

func main() {
	ts := c.Config
	for _, t := range ts {
		if t.Module == "icmp" {
			t.icmpTarget()
		} else if t.Module == "tcp_connect" {
			t.tcpTarget()
		} else if t.Module == "http_2xx" {
			t.httpTarget()
		}

	}
}
