package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	info "manageiq-exchange/models/info"
	meta "manageiq-exchange/models/metadata"
	user "manageiq-exchange/models/user"
	"net"
	"net/http"
	"time"
	"bufio"
	"strconv"
	"errors"
)

var netTransport = &http.Transport{
	Dial: (&net.Dialer{
		Timeout: 5 * time.Second,
	}).Dial,
	TLSHandshakeTimeout: 5 * time.Second,
}

var netClient = &http.Client{
	Timeout:   time.Second * 10,
	Transport: netTransport,
}

type Api struct {
	Server string
	Port   int
	Client *http.Client
	Data   DataApi
}

type DataApi struct {
	Data interface{}   `json:"data"`
	Meta meta.Metadata `json:"meta"`
}

func (a *Api) Init(server string, port int) {
	a.Server = server
	a.Port = port
	a.Client = netClient
}

func (a *Api) CheckConnectionServer() bool {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", a.Server, strconv.Itoa(a.Port)))
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		return false
	}
	fmt.Fprintf(conn, "GET / HTTP/1.0\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Printf("Fatal error: %s\n", err.Error())
		return false
	}
	fmt.Printf("Status connection server: %s", status)
	return true
}

func (a *Api) URL() string {
	url := fmt.Sprintf("http://%s", a.Server)
	if a.Port > 0 {
		url += fmt.Sprintf(":%d", a.Port)
	}
	return url
}

func (a *Api) GetInfo() info.Info {
	a.Request("GET", "", nil)
	var info info.Info
	info.Init(a.Data.Data.(map[string]interface{}))
	return info
}

func (a *Api) GetUsers(expand bool) user.UserCollection {
	var path string
	if path = "/v1/users"; expand {
		path = "/v1/users?expand=resources"
	}
	err := a.Request("GET", path, nil)
	if err != nil {
		fmt.Printf("%+v", err)
	}
	var users user.UserCollection
	users.Init(a.Data.Data.([]interface{}))
	return users
}

func (a *Api) Request(method string, path string, data io.Reader) error {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", a.URL(), path), data)
	if err != nil {
		return err
	}
	resp, err := a.Client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(strconv.Itoa(resp.StatusCode))
	} else {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		jsonErr := json.Unmarshal(body, &a.Data)
		if jsonErr != nil {
			return jsonErr
		}
	}
	return nil
}
