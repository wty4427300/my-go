package main

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func wse_demo1() {

	appid := "5cbe80cb"
	apikey := ""
	curtime := strconv.FormatInt(time.Now().Unix(),10)

	param := make(map[string]string)
	param["aue"] = "raw"
	param["language"] = "en_us"
	param["category"] = "read_sentence"

	tmp, _ := json.Marshal(param)
	base64_param := base64.StdEncoding.EncodeToString(tmp)

	w := md5.New()
	io.WriteString(w, apikey+curtime+base64_param)
	checksum := fmt.Sprintf("%x", w.Sum(nil))

	f, _ := ioutil.ReadFile("./test.pcm")
	audio := base64.StdEncoding.EncodeToString(f)
	text := "Good  morning , Ladies and gentlemen,We are honored to welcome you aboard Air China, a proud Star Alliance member. Kindly store all your carry-on luggage securely in the overhead bins or under the seat in front of you. Please take your assigned seats as quickly as possible ."

	data := url.Values{}
	data.Add("audio", audio)
	data.Add("text", text)
	body := data.Encode()

	client := &http.Client{}
	req, _ := http.NewRequest("POST", "http://api.xfyun.cn/v1/service/v1/ise", strings.NewReader(body))
	req.Header.Set("Content-Type","application/x-www-form-urlencoded")
	req.Header.Set("X-Appid", appid)
	req.Header.Set("X-CurTime", curtime)
	req.Header.Set("X-Param", base64_param)
	req.Header.Set("X-CheckSum", checksum)

	res, _ := client.Do(req)
	defer res.Body.Close()

	resp_body, _ := ioutil.ReadAll(res.Body)
	a:=string(resp_body)
	value := gjson.Get(a, "data.read_sentence.rec_paper.read_sentence.total_score")
	fmt.Println(value.String())
	fmt.Print(string(resp_body))
}

func main(){
	wse_demo1()
}