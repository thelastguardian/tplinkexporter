package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type TPLINKSwitchClient interface {
	GetPortStats() ([]portStats, error)
	GetHost() string
	// Login()
}

type TPLINKSwitch struct {
	host     string
	username string
	password string
}

type portStats struct {
	State      int
	LinkStatus int
	PktCount   map[string]int
}

func (client *TPLINKSwitch) GetHost() string {
	return client.host
}

func (client *TPLINKSwitch) GetPortStats() ([]portStats, error) {
	type allInfo struct {
		State      []int
		LinkStatus []int
		Pkts       []int
	}
	// "http://IP/logon.cgi"
	resp, err := http.PostForm(fmt.Sprintf("http://%s/logon.cgi", client.host), url.Values{"username": {client.username}, "password": {client.password}, "logon": {"Login"}})
	if err != nil {
		// handle error
		return nil, err
	}
	defer resp.Body.Close()
	// fmt.Println(resp, err)
	// "http://IP/PortStatisticsRpm.htm")
	resp2, err := http.Get(fmt.Sprintf("http://%s/PortStatisticsRpm.htm", client.host))
	if err != nil {
		// handle error
		return nil, err
	}
	defer resp2.Body.Close()
	body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		// handle error
		return nil, err
	}
	// fmt.Println(string(body))
	var jbody string = strings.ReplaceAll(
		strings.ReplaceAll(
			strings.ReplaceAll(
				string(body), "link_status", `"linkStatus"`),
			"state", `"State"`),
		"pkts", `"Pkts"`)
	// fmt.Println(string(jbody))
	res := regexp.MustCompile(`all_info = ({[^;]*});`).FindStringSubmatch(jbody)
	if res == nil {
		// fmt.Println(jbody)
		return nil, errors.New("unexpected response for port statistics http call: " + jbody)
	}
	// fmt.Println(res[1])
	var jparsed allInfo
	json.Unmarshal([]byte(res[1]), &jparsed)
	// fmt.Println(jparsed)
	var portsInfos []portStats
	portcount := len(jparsed.State)
	for i := 0; i < portcount; i++ {
		var portInfo portStats
		portInfo.State = jparsed.State[i]
		portInfo.LinkStatus = jparsed.LinkStatus[i]
		if portInfo.State == 1 {
			portInfo.PktCount = make(map[string]int)
			portInfo.PktCount["TxGoodPkt"] = jparsed.Pkts[4*i]
			portInfo.PktCount["TxBadPkt"] = jparsed.Pkts[4*i+1]
			portInfo.PktCount["RxGoodPkt"] = jparsed.Pkts[4*i+2]
			portInfo.PktCount["RxBadPkt"] = jparsed.Pkts[4*i+3]
		}
		portsInfos = append(portsInfos, portInfo)
	}
	fmt.Println(portsInfos)
	return portsInfos, nil
}

/*
sample output of PortStatisticsRpm.htm call:
<script>
var max_port_num = 8;
var port_middle_num  = 16;
var all_info = {
state:[1,1,1,1,1,1,1,1,0,0],
link_status:[6,6,0,6,0,0,0,5,0,0],
pkts:[1901830310,0,1338131260,33254,4291149014,0,2311488878,564,0,0,0,0,1814018004,0,33552310,0,0,0,0,0,0,0,0,0,0,0,0,0,1678459124,0,1866169392,0,0,0]
};
var tip = "";
</script>
*/

func NewTPLinkSwitch(host string, username string, password string) *TPLINKSwitch {
	return &TPLINKSwitch{
		host:     host,
		username: username,
		password: password,
	}
}
