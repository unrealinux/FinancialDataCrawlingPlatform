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
	"fmt"
	// "math"
	// "time"
	//"log"
	//"log"
)

func init() {
	Wukuangguojixt.Register()
}

var Wukuangguojixt = &Spider{
	Name:        "五矿国际信托",
	Description: "五矿国际信托净值数据 [Auto Page] [http://www.mintrust.com/wkxtweb/product/page_networth]",
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
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  "http://www.mintrust.com/wkxtweb/product/page_networth",
							Rule: aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					//ss := query.Find("pages")
					//fmt.Println(ss)
					ss := query.Find("#_pageSum")
					var countStr string
					countStr = ss.Text()
					var count int
					count, error := strconv.Atoi(countStr)
					if error != nil {
						fmt.Println("string convert failed")
					}

					page := 0

					for i := 1; i < count+1; i++ {
						page++
						ctx.AddQueue(&request.Request{
							Url:  "http://www.mintrust.com/wkxtweb/product/page_networth?netWorthPage.pageSize=10&netWorthPage.pageNum=" + strconv.Itoa(i),
							Rule: "获取结果",
							Temp: map[string]interface{}{
								"level1pages": page,
							},
						})
					}
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

					ss := query.Find(".about .productListTab tbody").Find("tr")

					var page int
					page = ctx.GetTemp("level1pages", &page).(int)

					count := 0

					ss.Each(func(i int, goq *goquery.Selection) {

						titleLine := goq.Children().Eq(0).Text()
						if titleLine != "产品名称" {
							mingchen := goq.Children().Eq(0).Find("a").Text()
							jingzhi := goq.Children().Eq(1).Text()
							leijijingzhi := goq.Children().Eq(1).Text()
							guzhiriqi := goq.Children().Eq(3).Text()

							count++
							fundID := "XTWUKUANGGUOJI" + "P1" + strconv.Itoa(page) + "L" + strconv.Itoa(count)

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
