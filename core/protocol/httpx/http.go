package httpx

import (
	"net/http"
	"strings"

	"github.com/elazarl/goproxy"

	"HFish/core/report"
	"HFish/core/rpc/client"
	"HFish/utils/is"
	"HFish/utils/log"
)

func Start(address string) {
	proxy := goproxy.NewProxyHttpServer()

	var info string

	proxy.OnRequest().DoFunc(
		func(r *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
			info = "URL:" + r.URL.String() + "&&Method:" + r.Method + "&&RemoteAddr:" + r.RemoteAddr

			arr := strings.Split(r.RemoteAddr, ":")

			// 判断是否为 RPC 客户端
			if is.Rpc() {
				go client.ReportResult("HTTP", "HTTP代理蜜罐", arr[0], info, "0")
			} else {
				go report.ReportHttp("HTTP代理蜜罐", "本机", arr[0], info)
			}

			return r, nil
		})

	// proxy.OnResponse().DoFunc(
	//	func(r *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	//		input, _ := ioutil.ReadAll(r.Body)
	//		info += "Response Info&&||kon||&&Status:" + r.Status + "&&Body:" + string(input)
	//		return r
	//	})

	err := http.ListenAndServe(address, proxy)
	if err != nil {
		log.Warn("hop http start error: %v", err)
	}
}
