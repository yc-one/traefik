package middlewares

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"unsafe"
)

const (
	// AddPrefixKey is the key within the request context used to
	// store the added prefix
	SenAccessLogKey key = "SendAccessLog"
)

type SendAccessLog struct {
	Handler   http.Handler
	RemoteUrl string
}

// SetHandler sets handler
func (s *SendAccessLog) SetHandler(Handler http.Handler) {
	s.Handler = Handler
}

func (s *SendAccessLog) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// 获取参数: 请求时间，请求的服务， 请求体， 请求头信息
	headerJson, err := json.Marshal(req.Header)
	if err != nil {
		fmt.Println("解析请求头信息出错:/n", err.Error())
		return
	}
	body, err := ioutil.ReadAll(req.Body)
	req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	if err != nil {
		fmt.Println("解析请求体信息出错:/n", err.Error())
		return
	}
	headerStr := string(headerJson)
	fmt.Println("请求头信息:\n", headerStr)
	bodyStr := string(body)
	fmt.Println("请求体信息:\n", bodyStr)
	requestUri := req.URL.RequestURI()
	fmt.Println("请求的uri:\n", requestUri)

	requestParams := make(map[string]interface{})
	requestParams["headInfo"] = headerStr
	requestParams["bodyInfo"] = bodyStr
	requestParams["uriInfo"] = requestUri

	//url := "http://localhost/test"
	fmt.Println("准备发送post请求到sf-monitor, url:", "http://"+s.RemoteUrl)

	bytesData, err := json.Marshal(requestParams)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", "http://"+s.RemoteUrl, reader)

	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	client := http.Client{}
	resp, err := client.Do(request)

	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	str := (*string)(unsafe.Pointer(&respBytes))
	fmt.Println("响应信息：", *str)

	s.Handler.ServeHTTP(w, req)

}
