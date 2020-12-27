package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func (n *nessusImpl) Login(username, password string) error {
	if n.verbose {
		log.Printf("登录Nessus扫描器... %s\n", n.apiURL)
	}
	data := make(map[string]interface{})
	data["username"] = username
	data["password"] = password

	resp, err := n.Request("POST", "/session", data, []int{http.StatusOK})
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	reply := &loginResp{}
	if err = json.NewDecoder(resp.Body).Decode(&reply); err != nil {
		return err
	}
	n.authCookie = reply.Token
	return nil
}

func (n *nessusImpl) AuthCookie() string {
	return n.authCookie
}

func (n *nessusImpl) Request(method string, resource string, bodystr interface{}, wantStatus []int) (resp *http.Response, err error) {
	u, err := url.ParseRequestURI(n.apiURL)
	if err != nil {
		return nil, err
	}
	u.Path = resource
	urlStr := fmt.Sprintf("%v", u)

	jb, err := json.Marshal(bodystr)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, urlStr, bytes.NewBufferString(string(jb)))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")
	if n.authCookie != "" {
		req.Header.Add("X-Cookie", fmt.Sprintf("token=%s", n.authCookie))
		req.Header.Add("Content-Length", fmt.Sprintf("%d", len(string(jb))))
		req.Header.Add("X-API-Token", "7A2C4A96-4AEE-4706-9617-2EF643532628")
	}

	if n.verbose {
		db, err := httputil.DumpRequest(req, true)
		if err != nil {
			return nil, err
		}
		log.Println("发送数据...", string(db))
	}
	resp, err = n.client.Do(req)
	if err != nil {
		return nil, err
	}
	if n.verbose {
		if body, err := httputil.DumpResponse(resp, true); err == nil {
			log.Println(string(body))
		}
	}
	var statusFound bool
	for _, status := range wantStatus {
		if resp.StatusCode == status {
			statusFound = true
			break
		}
	}
	if !statusFound {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("[Error]异常响应错误代码, error code  %d wanted %v (%s)", resp.StatusCode, wantStatus, body)
	}
	return resp, nil
}
