package common

import (
	"bytes"
	"encoding/json"
	"ff14/config"
	"ff14/respdata"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func StrFirstToUpper(str string) string {
	if len(str) < 1 {
		return ""
	}
	strArry := []rune(str)
	if strArry[0] >= 97 && strArry[0] <= 122 {
		strArry[0] -= 32
	}
	return string(strArry)
}

func JsonEncode(data interface{}) string {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.Encode(data)
	return string(buffer.Bytes())
}

func TimeStamp(timeString string) int64 {
	timeLayout := "2006-01-02 15:04:05"  //转化所需模板
	loc, _ := time.LoadLocation("Local") //获取时区
	tmp, _ := time.ParseInLocation(timeLayout, timeString, loc)
	timestamp := tmp.Unix() //转化为时间戳 类型是int64
	return timestamp
}

func CurrentTime() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

func Get(apiUrl string, data url.Values, cookie ...http.Cookie) (response string) {
	parseURL, err := url.Parse(apiUrl)
	if err != nil {
		panic(err)
	}
	if data != nil {
		parseURL.RawQuery = data.Encode()
	}
	client := &http.Client{Timeout: 5 * time.Second}
	var req *http.Request
	req, _ = http.NewRequest("GET", parseURL.String(), nil)
	cookie1 := &http.Cookie{Name: "ff14risingstones", Value: config.Env.GetString("cookie.ff14risingstones")}
	req.AddCookie(cookie1)
	for _, v := range cookie {
		req.AddCookie(&v)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, error := client.Do(req)
	defer resp.Body.Close()
	if error != nil {
		panic(error)
	}
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	response = result.String()
	return
}

func Post(apiUrl string, data interface{}, contentType string, cookie ...http.Cookie) (response string) {
	var req *http.Request
	var err error
	value, ok := data.(url.Values)
	if ok {
		req, err = http.NewRequest("POST", apiUrl, strings.NewReader(value.Encode()))
	} else {
		body, _ := json.Marshal(data)
		req, err = http.NewRequest("POST", apiUrl, strings.NewReader(string(body)))
	}
	if contentType == "" {
		contentType = "application/x-www-form-urlencoded"
	}
	req.Header.Set("content-type", contentType)
	cookie1 := &http.Cookie{Name: "ff14risingstones", Value: config.Env.GetString("cookie.ff14risingstones")}
	req.AddCookie(cookie1)
	for _, v := range cookie {
		req.AddCookie(&v)
	}
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	client := &http.Client{Timeout: 5 * time.Second}
	resp, error := client.Do(req)
	if error != nil {
		panic(error)
	}
	defer resp.Body.Close()

	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
	}
	response = result.String()
	return
}

func FFXIVIsError(response string) error {
	var ffxivCode respdata.FFXIVCode
	json.Unmarshal([]byte(response), &ffxivCode)
	if ffxivCode.Code != 10000 {
		return fmt.Errorf("错误：%s", ffxivCode.Msg)
	}
	return nil
}
