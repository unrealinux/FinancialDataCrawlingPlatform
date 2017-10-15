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
	Kunlunxt.Register()
}

var Kunlunxt = &Spider{
	Name:        "昆仑信托",
	Description: "昆仑信托净值数据 [Auto Page] [http://www.kunluntrust.com/xinxipilu/chanpinxinxi/chanpinjingzhi/]",
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
							Url:  "http://www.kunluntrust.com/xinxipilu/chanpinxinxi/chanpinjingzhi/",
							Rule: aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					ss := query.Find(".cmspage ul").Find("li")

					page := 0

					ss.Each(func(i int, goq *goquery.Selection) {
						title := goq.Text()
						if title != "<" && title != "1" && title != ">" {

							url, exist := goq.Find("a").Attr("href")
							if exist {

								page++

								ctx.AddQueue(&request.Request{
									Url:  "http://www.kunluntrust.com/xinxipilu/chanpinxinxi/chanpinjingzhi/" + url,
									Rule: "获取结果",
									Temp: map[string]interface{}{
										"level1pages": page,
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
					query := ctx.GetDom()

					ss := query.Find(".aj_list").Find("dl")

					var page1 int
					ctx.GetTemp("level1pages", &page1)

					recordCount := 0

					ss.Each(func(i int, goq *goquery.Selection) {

						classattr, exist := goq.Attr("class")
						fmt.Println(classattr)
						if exist == false {

							ssaa := goq.Children()
							count := 0

							var mingchen string
							var jingzhi string
							var leijijingzhi string
							var guzhiriqi string

							ssaa.Each(func(i int, goqaa *goquery.Selection) {

								titleLine := goqaa.Children().Eq(0).Text()
								if titleLine != "产品名称" {

									count++

									if count%3 == 1 {
										mingchen = goqaa.Text()
									}

									if count%3 == 2 {
										jingzhi = goqaa.Text()
										leijijingzhi = goqaa.Text()
									}

									if count%3 == 0 {
										guzhiriqi = goqaa.Text()

										recordCount++
										fundID := "XTKUNLUN" + "P1" + strconv.Itoa(page1) + "L" + strconv.Itoa(recordCount)

										ctx.Output(map[int]interface{}{
											0: fundID,
											1: mingchen,
											2: jingzhi,
											3: leijijingzhi,
											4: guzhiriqi,
										})
									}

								}
							})
						}
					})
				},
			},
		},
	},
}
