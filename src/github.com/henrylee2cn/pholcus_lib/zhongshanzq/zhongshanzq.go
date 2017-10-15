package pholcus_lib

import (
	// 基础包
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
	// "github.com/henrylee2cn/pholcus/logs"           //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	// "encoding/json"

	// 字符串处理包
	// "regexp"
	//"strconv"
	"strings"
	// 其他包
	"fmt"
	// "math"
	// "time"
	//"log"
	//"log"
	//"reflect"
)

func init() {
	Zhongshanzq.Register()
}

var Zhongshanzq = &Spider{
	Name:        "中山证券私募基金",
	Description: "中山证券私募基金数据 [Auto Page] [http://ov.zszq.com/product/jh1/jh1-1.asp]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	NotDefaultField: true,

	Namespace: func(*Spider) string {
		return "zhengquan"
	},
	// 子命名空间相对于表名，可依赖具体数据内容，可选
	SubNamespace: func(self *Spider, dataCell map[string]interface{}) string {
		return "fund_src_nav"
	},

	EnableCookie: false,
	RuleTree: &RuleTree{

		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 2}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {

					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  "http://ov.zszq.com/product/jh1/jh1-1.asp",
							Rule: aid["Rule"].(string),
						})
					}

					return nil
				},

				ParseFunc: func(ctx *Context) {

					query := ctx.GetDom()

					ss := query.Find(".bg tbody").Find("tr")

					ss.Each(func(i int, goq *goquery.Selection) {

						url, exist := goq.Children().Eq(0).Find("a").Attr("href")
						if exist == true {

							strings.Replace(url, "1.asp", "4.asp", -1)

							ctx.AddQueue(&request.Request{
								Url:  "http://ov.zszq.com" + url,
								Rule: "获取结果",
							})
						}
					})
				},
			},

			"获取结果": {
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					ss := query.Find(".border tbody").Find("tr")
					fmt.Println(ss)
				},
			},
		},
	},
}
