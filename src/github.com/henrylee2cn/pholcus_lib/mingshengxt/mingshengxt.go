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
	"fmt"
	"strings"
)

func init() {
	Mingshengxt.Register()
}

var Mingshengxt = &Spider{
	Name:        "中国民生信托",
	Description: "中国民生信托净值数据 [Auto Page] [http://www.msxt.com/?jzgb.html]",
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

			webpage := 6

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

				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"基金ID",
					"名称",
					"净值",
					"累计净值",
					"估值日期",
				},

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {

					page := 0

					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						page++

						ctx.AddQueue(&request.Request{
							Url:  "http://www.msxt.com/?jzgb/page/" + strconv.Itoa(loop[0]) + ".html",
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

					ss := query.Find(".jzgb table tbody").Find("tr")
					fmt.Println(ss)

					var page1 int
					page1 = ctx.GetTemp("level1pages", &page1).(int)

					count := 0

					ss.Each(func(i int, goq *goquery.Selection) {

						titleLine := goq.Children().Eq(0).Text()
						if titleLine != "序号" {
							mingchen := goq.Children().Eq(1).Find("a").Text()
							jingzhi := goq.Children().Eq(4).Text()
							leijijingzhi := goq.Children().Eq(4).Text()
							guzhiriqi := goq.Children().Eq(3).Text()

							count++
							fundID := "XTMINGSHENG" + "P1" + strconv.Itoa(page1) + "L" + strconv.Itoa(count)

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
