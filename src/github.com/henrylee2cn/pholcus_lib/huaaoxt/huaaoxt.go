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
	//"fmt"
	// "math"
	// "time"
	//"log"
	//"log"
	"fmt"
	"strings"
)

func init() {
	Huaaoxt.Register()
}

var Huaaoxt = &Spider{
	Name:        "华奥信托",
	Description: "华奥信托净值数据 [Auto Page] [http://www.huaao-trust.com/list/705/1.shtml]",
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

			webpage := 42

			var configs []string
			configs = strings.Split(Keys, ",") //各种配置按照key1=value1,key2=value2,...的形式解析

			for a := 0; a < len(configs); a++ {

				if strings.Contains(configs[a], "page=") {
					webpage, _ = strconv.Atoi(strings.TrimLeft(Keys, "page="))
					fmt.Println(webpage)
				}

			}

			ctx.Aid(map[string]interface{}{"loop": [2]int{1, webpage}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					page := 0

					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						page++

						ctx.AddQueue(&request.Request{
							Url:  "http://www.huaao-trust.com/list/705/" + strconv.Itoa(loop[0]) + ".shtml",
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

					ss := query.Find(".news_list ul").Children()

					var page1 int
					ctx.GetTemp("level1pages", &page1)

					page2 := 0
					ss.Each(func(i int, goq *goquery.Selection) {

						url, exist := goq.Find("a").Attr("href")
						if exist {
							page2++

							newUrl := strings.Replace(url, "-", "/", 2)
							newUrl1 := strings.Replace(newUrl, "_", "/", 1)
							ctx.AddQueue(&request.Request{
								Url:  "http://www.huaao-trust.com" + newUrl1,
								Rule: "获取结果",

								Temp: map[string]interface{}{
									"level1pages": page1,
									"level2pages": page2,
								},
							})
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

					ss := query.Find(".news_details tbody").Find("tr")

					count := 0

					var page int
					ctx.GetTemp("level1pages", &page)

					var page2 int
					ctx.GetTemp("level2pages", &page2)

					ss.Each(func(i int, goq *goquery.Selection) {

						titleLine := goq.Children().Eq(0).Text()
						if titleLine != "产品名称" {
							mingchen := goq.Children().Eq(0).Text()
							jingzhi := goq.Children().Eq(3).Text()
							leijijingzhi := goq.Children().Eq(3).Text()
							guzhiriqi := goq.Children().Eq(2).Text()

							count++
							fundID := "XTHUAAO" + "P1" + strconv.Itoa(page) + "P2" + strconv.Itoa(page2) + "L" + strconv.Itoa(count)

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
