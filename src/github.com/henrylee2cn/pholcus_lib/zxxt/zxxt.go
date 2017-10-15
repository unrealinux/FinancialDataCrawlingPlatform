package pholcus_lib

import (
	// 基础包
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析

	"fmt"
	"strconv"
	"strings"
)

func init() {
	Zxxt.Register()
}

var Zxxt = &Spider{
	Name:        "中信信托",
	Description: "中信信托阳光私募净值数据 [Auto Page] [http://trust.ecitic.com/XXPL_JZPL/index.jsp?type=1&pageNum=1&keword=]",
	// Pausetime: 300,
	Keyin: KEYIN,
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

			webpage := 7

			var configs []string
			configs = strings.Split(Keys, ",") //各种配置按照key1=value1,key2=value2,...的形式解析

			for a := 0; a < len(configs); a++ {

				if strings.Contains(configs[a], "page=") {
					webpage, _ = strconv.Atoi(strings.TrimLeft(Keys, "page="))
					fmt.Println(webpage)
				}

			}

			ctx.Aid(map[string]interface{}{"loop": [2]int{1, webpage}, "Rule": "生成请求", "count": 0}, "生成请求")
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
							Url:  "http://trust.ecitic.com/XXPL_JZPL/index.jsp?type=1&pageNum=" + strconv.Itoa(loop[0]) + "&keword=",
							Rule: aid["Rule"].(string),
							Temp: map[string]interface{}{
								"pages": page,
							},
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					ss := query.Find("#pageRight .info").Find("tbody").Find("tr")

					var page int
					ctx.GetTemp("pages", &page)

					count := 0
					ss.Each(func(i int, goq *goquery.Selection) {

						titleLine := goq.Children().Eq(1).Text()
						if titleLine != "产品名称" {
							mingchen := goq.Children().Eq(1).Text()
							jingzhi := strings.TrimSpace(goq.Children().Eq(3).Text())
							leijijingzhi := strings.TrimSpace(goq.Children().Eq(4).Text())
							guzhiriqi := goq.Children().Eq(7).Text()

							count++

							fundID := "XTZHONGXIN" + "P" + strconv.Itoa(page) + "L" + strconv.Itoa(count)

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
