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
	Gmxt.Register()
}

var Gmxt = &Spider{
	Name:        "国民信托",
	Description: "国民信托净值数据 [Auto Page] [http://www.natrust.cn/fe/equity/catList.gsp?rd=0.10877497288215976]",
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
			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 15}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					page := 0
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						page++
						ctx.AddQueue(&request.Request{
							Url:  "http://www.natrust.cn/fe/equity/catList.gsp?rd=0.10877497288215976",
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

					ss := query.Find(".job_list ul")

					var page1 int
					ctx.GetTemp("level1pages", &page1)

					page2 := 0

					ss.Each(func(i int, goq *goquery.Selection) {

						var title string
						var ok bool
						if title, ok = goq.Find("a").Attr("title"); ok {

						}

						if url, ok := goq.Find("a").Attr("href"); ok {

							for i := 0; i < 10; i++ {
								page2++

								ctx.AddQueue(&request.Request{
									Url:  "http://www.natrust.cn/fe/equity/" + url + "&page=" + strconv.Itoa(i+1),
									Rule: "获取结果",
									Temp: map[string]interface{}{
										"title":       title,
										"level1pages": page1,
										"level2pages": page2,
									},
								})
							}
						}

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

					queryResult := ctx.GetDom()
					ssResult := queryResult.Find(".jz_table tbody").Find("tr")

					count := 0

					var page int
					ctx.GetTemp("level1pages", &page)

					var page2 int
					ctx.GetTemp("level2pages", &page2)

					ssResult.Each(func(i int, goq *goquery.Selection) {

						titleLineResult := goq.Children().Eq(0).Text()

						var titleMingCheng string
						ctx.GetTemp("title", &titleMingCheng)

						if titleLineResult != "日期" && titleLineResult != "" {

							mingchen := titleMingCheng
							jingzhi := goq.Children().Eq(1).Text()
							leijijingzhi := goq.Children().Eq(2).Text()
							guzhiriqi := goq.Children().Eq(0).Text()

							count++
							fundID := "XTGUOMING" + "P1" + strconv.Itoa(page) + "P2" + strconv.Itoa(page2) + "L" + strconv.Itoa(count)

							ctx.Output(map[int]interface{}{
								0: fundID,
								1: mingchen,
								2: jingzhi,
								3: leijijingzhi,
								4: guzhiriqi,
							})
						}

					})
				},
			},
		},
	},
}
