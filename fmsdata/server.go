package fmsdata

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type DataBase struct {
	url             string
	database        string
	user            string
	pass            string
	token           string
	quit            chan struct{}
	isAuthenticated bool
	tokenTime       time.Time
	sync.RWMutex
}

func NewDataBase(url, database, user, pass string) DataBase {
	return DataBase{url: url, user: user, pass: pass, database: database}
}

func getClient() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	return client
}

func (d *DataBase) Login() error {
	if d.isAuthenticated {
		return nil
	}
	url := d.url + "/fmi/data/v1/databases/" + d.database + "/sessions"

	b := new(bytes.Buffer)
	b.WriteString("{}")
	req, _ := http.NewRequest("POST", url, b)
	req.Header.Set("Content-Type", "application/json")
	req.SetBasicAuth(d.user, d.pass)
	resp, err := getClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	/*buf := bytes.Buffer{}
	buf.ReadFrom(resp.Body)
	fmt.Println(string(buf.Bytes()))
	return errors.New("Testing")*/

	loginResult := struct {
		Response struct {
			Token string `json:"token"`
		} `json:"response"`
		Messages []struct {
			Code    string `json:"code"`
			Message string `json:"message"`
		} `json:"messages"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&loginResult); err != nil {
		return err
	}
	//fmt.Println(loginResult)
	if loginResult.Messages[0].Code == "0" {
		d.tokenTime = time.Now()
		d.token = loginResult.Response.Token
		d.isAuthenticated = true
		// We have a token, keep it and renew when needed
		ticker := time.NewTicker(14 * time.Minute)
		d.quit = make(chan struct{})
		go func() {
			for {
				select {
				case <-ticker.C:
					d.Lock()
					d.Logout()
					d.Login()
					d.Unlock()

				case <-d.quit:
					ticker.Stop()
					return
				}
			}
		}()
	} else {
		return errors.New("Could not log in")
	}

	return nil
}

func (d *DataBase) Logout() error {
	url := d.url + "/fmi/data/v1/databases/" + d.database + "/sessions/" + d.token
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := getClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	logoutInfo := struct {
		Result int `json:"result"`
	}{}
	if err := json.NewDecoder(resp.Body).Decode(&logoutInfo); err != nil {
		return err
	}
	if logoutInfo.Result != 0 {
		return errors.New("Failed to logout")
	}
	d.isAuthenticated = false
	//close(d.quit)
	return nil
}

func (d *DataBase) makeCall(url, method string, inf interface{}) ([]byte, error) {
	//fmt.Println(url)
	var err error = nil
	req := &http.Request{}
	if inf != nil {
		b := new(bytes.Buffer)
		err := json.NewEncoder(b).Encode(inf)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, url, b)
	} else {
		req, err = http.NewRequest(method, url, nil)
	}
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+d.token)
	req.Header.Add("Content-Type", "application/json")
	resp, err := getClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}
