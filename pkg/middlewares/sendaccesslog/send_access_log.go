package sendaccesslog

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/containous/traefik/v2/pkg/config/dynamic"
	"github.com/containous/traefik/v2/pkg/log"
	"github.com/containous/traefik/v2/pkg/middlewares"
	"io/ioutil"
	"net/http"
	"unsafe"
)

const (
	typeName = "SendAccessLog"
)

type sendAccessLog struct {
	next      http.Handler
	name      string
	remoteUrl string
}

func New(ctx context.Context, next http.Handler, config dynamic.SendAccessLog, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName)).Debug("Creating send-access-log middleware")
	return &sendAccessLog{
		name:      name,
		next:      next,
		remoteUrl: config.RemoteUrl,
	}, nil
}

func (s *sendAccessLog) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	// 获取参数: 请求时间，请求的服务， 请求体， 请求头信息
	headerJson, err := json.Marshal(req.Header)
	if err != nil {
		fmt.Println("解析请求头信息出错:/n", err.Error())
		return
	}
	body, err := ioutil.ReadAll(req.Body)
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
	fmt.Println("准备发送post请求到sf-monitor, url:", s.remoteUrl)

	bytesData, err := json.Marshal(requestParams)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	reader := bytes.NewReader(bytesData)
	request, err := http.NewRequest("POST", s.remoteUrl, reader)

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

}
