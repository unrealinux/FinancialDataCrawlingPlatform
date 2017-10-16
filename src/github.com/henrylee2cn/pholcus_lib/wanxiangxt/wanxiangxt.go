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
	Wanxiangxt.Register()
}

var Wanxiangxt = &Spider{
	Name:        "万向信托",
	Description: "万向信托净值数据 [Auto Page] [http://www.siti.com.cn/product.php?fid=23&fup=3]",
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

			webpage := 5

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

					ctx.AddQueue(&request.Request{
						Url:  "http://www.wxtrust.com/c55ef089-7d41-4e6c-8439-e5694694365e/index.html",
						Rule: aid["Rule"].(string),
					})

					page := 0
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						page++
						ctx.AddQueue(&request.Request{
							Url:  "http://www.wxtrust.com/c55ef089-7d41-4e6c-8439-e5694694365e/index_" + strconv.Itoa(loop[0]) + ".html",
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
					fmt.Println(query.Text())

					ss := query.Find(".right_info_1 dl")
					//fmt.Println(ss)

					var page2 int
					ctx.GetTemp("level1pages", &page2)

					page := 0
					ss.Each(func(i int, goq *goquery.Selection) {

						page++
						//fmt.Println(goq.Children())

						goq.Children().Each(func(i int, goqchild *goquery.Selection) {

							//fmt.Println(goqchild.Text())

							goqchild.Children().Each(func(i int, goqchildchild *goquery.Selection) {
								//fmt.Println(goqchildchild)
								//fmt.Println(goqchildchild.Text())

								scripts := strings.Split(goqchildchild.Text(), ";")

								dateString := strings.TrimSpace(scripts[0])

								beginIndex := strings.Index(dateString, "\"")
								endIndex := strings.Index(dateString, "-")
								dateString_1 := dateString[beginIndex+1 : endIndex]
								//fmt.Println(dateString_1)

								hrefString := strings.TrimSpace(scripts[1])
								hrefs := strings.Split(hrefString, "+")

								beginIndex1 := strings.Index(hrefs[0], "/")
								endIndex1 := strings.LastIndex(hrefs[0], "'")
								tempString_ := hrefs[0][beginIndex1:endIndex1]
								//fmt.Println(tempString_)

								beginIndex2 := strings.Index(hrefs[2], "/")
								endIndex2 := strings.Index(hrefs[2], "\"")
								tempString__ := hrefs[2][beginIndex2:endIndex2]
								//fmt.Println(tempString__)

								beginIndex3 := strings.Index(hrefs[2], ">")
								endIndex3 := strings.Index(hrefs[2], "<")
								tempString___ := hrefs[2][beginIndex3+1 : endIndex3]
								//fmt.Println(tempString___)

								fmt.Println("http://www.wxtrust.com" + tempString_ + dateString_1 + tempString__)

								ctx.AddQueue(&request.Request{
									Url:  "http://www.wxtrust.com" + tempString_ + dateString_1 + tempString__,
									Rule: "获取结果",
									Temp: map[string]interface{}{
										"mingcheng": tempString___,
										"level1pages": page2,
										"level2pages": page,
									},
								})
							})
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
					fmt.Print(query.Text())

					titleMingcheng := strings.TrimSpace(query.Find(".right_title").Text())

					ss_ := query.Find("#customers tbody")
					if ss_ == nil {
						return
					}

					ss := ss_.Find("tr")

					//var titleMingcheng string
					//ctx.GetTemp("mingcheng", &titleMingcheng)
					//fmt.Println(titleMingcheng)

					var page int
					page = ctx.GetTemp("level1pages", &page).(int)

					var page2 int
					page2 = ctx.GetTemp("level2pages", &page2).(int)

					mingchen := titleMingcheng

					count := 0
					ss.Each(func(i int, goq *goquery.Selection) {

						title := strings.TrimSpace(goq.Children().Eq(0).Text())
						fmt.Println(title)

						if title != "估值日期" {
							jingzhi := strings.TrimSpace(goq.Children().Eq(1).Text())
							leijijingzhi := strings.TrimSpace(goq.Children().Eq(2).Text())
							guzhiriqi := strings.TrimSpace(goq.Children().Eq(0).Text())

							count++
							fundID := "XTWANXIANG" + "P1" + strconv.Itoa(page) + "P2" + strconv.Itoa(page2) + "L" + strconv.Itoa(count)

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
