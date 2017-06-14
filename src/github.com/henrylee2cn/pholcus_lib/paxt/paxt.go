package pholcus_lib

import (
	// 基础包
	"github.com/henrylee2cn/pholcus/common/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
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
)

func init() {
	Paxt.Register()
}

var Paxt = &Spider{
	Name:        "平安信托",
	Description: "平安信托净值数据 [Auto Page] [http://trust.pingan.com/xintuochanpinjingzhi/index.shtml]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	NotDefaultField: true,
	
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
							Url:  "http://trust.pingan.com/xintuochanpinjingzhi/index.shtml",
							Rule: aid["Rule"].(string),
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {

					page := 0

                    for i:= 1; i < 34; i++{
						page++

                        ctx.AddQueue(&request.Request{
                            Url:  "http://trust.pingan.com/xintuochanpinjingzhi/index.shtml?trustNo=&trustName=&currentPageNo="+ strconv.Itoa(i) +"&cooID=&type=",
                            Rule: "获取结果",
							Temp: map[string]interface{}{
								"level1pages" : page,
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
					
					ss := query.Find(".not_special_key_word tbody").Find("tr")

					var page1 int
					ctx.GetTemp("level1pages", &page1)

					count := 0
							
                    ss.Each(func(i int, goq *goquery.Selection) {
                        					
                            	
						titleLine := goq.Children().Eq(0).Text()
						if titleLine != "产品名称" {
							mingchen := goq.Children().Eq(0).Text()
							jingzhi := goq.Children().Eq(2).Text()
							leijijingzhi := goq.Children().Eq(3).Text()
							guzhiriqi := goq.Children().Eq(5).Text()

							count++
							fundID := "XTPINGAN" + "P1" + strconv.Itoa(page1) + "L" + strconv.Itoa(count)
						
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