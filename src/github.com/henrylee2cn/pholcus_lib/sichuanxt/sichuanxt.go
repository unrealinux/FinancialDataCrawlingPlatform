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
	"strconv"
	// "strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
	//"log"
	//"log"
)

func init() {
	Sichuanxt.Register()
}

var Sichuanxt = &Spider{
	Name:        "四川信托",
	Description: "四川信托净值数据 [Auto Page] [http://www.schtrust.com/index.php?m=content&c=index&a=lists&catid=76]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	NotDefaultField: true,

	Namespace: func(*Spider) string {
		return "xintuo"
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
					page := 0
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						page++
						ctx.AddQueue(&request.Request{
							Url:  "http://www.schtrust.com/index.php?m=content&c=index&a=lists&catid=76",
							Rule: aid["Rule"].(string),
							Temp: map[string]interface{}{
								"level1pages": page,
							},
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					//ss := query.Find("pages")
					//fmt.Println(ss)
					ss := query.Find(".search-dl").Children().Eq(0)

					var page1 int
					ctx.GetTemp("level1pages", &page1)

					page2 := 0

					ss.Each(func(i int, goq *goquery.Selection) {
						ssaa := goq.Find("dd").Children()
						ssaa.Each(func(i int, goq *goquery.Selection) {
							url, exist := goq.Attr("href")
							if exist {
								page2++

								mingchengTitle := goq.Text()
								ctx.AddQueue(&request.Request{
									Url:  "http://www.schtrust.com" + url,
									Rule: "获取结果",
									Temp: map[string]interface{}{
										"mingcheng":   mingchengTitle,
										"level1pages": page1,
										"level2pages": page2,
									},
								})
							}

						})
					})
				},
			},

			"获取结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"基金ID",
					"名称",
					"净值",
					"累计净值",
					"估值日期",
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					ss := query.Find(".product-net-con tbody").Find("tr")

					count := 0

					var page int
					page = ctx.GetTemp("level1pages", &page).(int)

					var page2 int
					page2 = ctx.GetTemp("level2pages", &page2).(int)

					var titleMingcheng string
					mingchen := ctx.GetTemp("mingcheng", &titleMingcheng).(string)
					//mingchen := titleMingcheng
					ss.Each(func(i int, goq *goquery.Selection) {

						jingzhi := goq.Children().Eq(1).Text()
						leijijingzhi := goq.Children().Eq(2).Text()
						guzhiriqi := goq.Children().Eq(0).Text()

						count++
						fundID := "XTSICHUANGUOJI" + "P1" + strconv.Itoa(page) + "P2" + strconv.Itoa(page2) + "L" + strconv.Itoa(count)

						ctx.Output(map[int]interface{}{
							0: fundID,
							1: mingchen,
							2: jingzhi,
							3: leijijingzhi,
							4: guzhiriqi,
						})
					})
				},
			},
		},
	},
}
