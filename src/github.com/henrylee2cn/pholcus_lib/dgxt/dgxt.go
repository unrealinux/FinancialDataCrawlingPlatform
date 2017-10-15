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
	"strings"
)

func init() {
	Dgxt.Register()
}

var Dgxt = &Spider{
	Name:        "东莞信托",
	Description: "东莞信托净值数据 [Auto Page] [http://www.dgxt.com/xthxxlt/index.html]",
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

			Keys := ctx.GetKeyin()
			fmt.Println(Keys)

			webpage := 20

			var configs []string
			configs = strings.Split(Keys, ",") //各种配置按照key1=value1,key2=value2,...的形式解析

			for a := 0; a < len(configs); a++ {

				if strings.Contains(configs[a], "page=") {
					webpage, _ = strconv.Atoi(strings.TrimLeft(Keys, "page="))
					fmt.Println(webpage)
				}

			}

			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 2}, "Rule": "生成请求", "webpage" : webpage}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {

					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  "http://www.dgxt.com/xthxxlt/index.html",
							Rule: aid["Rule"].(string),
							Temp: map[string]interface{}{
								"webpage":aid["webpage"],
							},
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					//fmt.Println(ss)

					ss := query.Find(".xmcnfds a")
					fmt.Println(ss)

					var webpage int
					webpage = ctx.GetTemp("webpage", webpage).(int)

					page1 := 0

					ss.Each(func(i int, goq *goquery.Selection) {

						if url, ok := goq.Attr("href"); ok {
							page1++

							ctx.AddQueue(&request.Request{
								Url:  "http://www.dgxt.com/" + url,
								Rule: "净值详情",
								Temp: map[string]interface{}{
									"level1pages": page1,
									"level1url" : "http://www.dgxt.com/" + url,
									"webpage" : ctx.GetTemp("webpage", webpage),
								},
							})
						}

					})
				},
			},

			"净值详情": {

				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					ss := query.Find("ul").Find("li")

					var page1 int
					page1 = ctx.GetTemp("level1pages", &page1).(int)
					var level2url string
					level2url = ctx.GetTemp("level1url", &level2url).(string)

					var webpage int
					webpage = ctx.GetTemp("level1pages", &webpage).(int)

					page2 := 0

					for i:= 0; i < webpage; i++{

						page2++

						ctx.AddQueue(&request.Request{
							Url:  level2url + "&pageNo=" + strconv.Itoa(i),
							Rule: "获取结果",
							Temp: map[string]interface{}{
								"level1pages": page1,
								"level2pages": page2,
							},
						})


					}

					ss.Each(func(i int, goq *goquery.Selection) {

						mingcheng := goq.Children().Eq(0).Find("a").Text()
						urlink := goq.Children().Eq(4).Find("a")

						urlink.Each(func(i int, goq *goquery.Selection) {

							if url, ok := goq.Attr("href"); ok {

								page2++

								ctx.AddQueue(&request.Request{
									Url:  "http://www.dgxt.com/" + url,
									Rule: "获取结果",
									Temp: map[string]interface{}{
										"mingcheng":   mingcheng,
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

					ss := query.Find(".cot_fl ul").Find("li")

					titleMingCheng := query.Find(".uynb").Text();

					count := 0

					var page int
					ctx.GetTemp("level1pages", &page)

					var page2 int
					ctx.GetTemp("level2pages", &page2)

					ss.Each(func(i int, goq *goquery.Selection) {

						divDetail := goq.Children().Eq(0)

						mingchen := titleMingCheng
						jingzhi := divDetail.Children().Eq(1).Text()
						leijijingzhi := divDetail.Children().Eq(2).Text()
						guzhiriqi := divDetail.Children().Eq(0).Text()

						count++
						fundID := "XTDONGGUAN" + "P1" + strconv.Itoa(page) + "P2" + strconv.Itoa(page2) + "L" + strconv.Itoa(count)

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
