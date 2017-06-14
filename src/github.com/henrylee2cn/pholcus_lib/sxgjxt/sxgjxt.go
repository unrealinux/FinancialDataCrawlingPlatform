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
	 "fmt"
	// "math"
	// "time"
	//"log"
	//"log"
)

func init() {
	Sxgjxt.Register()
}

var Sxgjxt = &Spider{
	Name:        "陕西国际信托",
	Description: "陕西国际信托净值数据 [Auto Page] [http://www.siti.com.cn/product.php?fid=23&fup=3]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	NotDefaultField: true,
	
	Namespace: func(*Spider) string{
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
							Url:  "http://www.siti.com.cn/product.php?fid=23&fup=3",
							Rule: aid["Rule"].(string),
							Temp: map[string]interface{}{
								"level1pages" : page,
							},
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()
					
                    //ss := query.Find("pages")
                    //fmt.Println(ss)
					ss := query.Find(".pages a")
                    fmt.Println(ss)

					var page1 int
					ctx.GetTemp("level1pages", &page1)

					page2 := 0
							
					ss.Each(func(i int, goq *goquery.Selection) {
                                
						if url, ok := goq.Attr("href"); ok {
							page2++

							ctx.AddQueue(&request.Request{
								Url:  "http://www.siti.com.cn/" + url,
								Rule: "获取结果",
								Temp: map[string]interface{}{
									"level1pages" : page1,
									"level2pages" : page2,
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
					
					ss := query.Find(".proDeatil tbody").Find("tr")
							
					count := 0

					var page int
					ctx.GetTemp("level1pages", &page)

					var page2 int
					ctx.GetTemp("level2pages", &page2)

                    ss.Each(func(i int, goq *goquery.Selection) {
                        						
							mingchen := goq.Children().Eq(1).Text()
							jingzhi := goq.Children().Eq(4).Text()
							leijijingzhi := goq.Children().Eq(5).Text()
							guzhiriqi := goq.Children().Eq(9).Text()

							count++
							fundID := "XTSHANXIGUOJI" + "P1" + strconv.Itoa(page) + "P2" + strconv.Itoa(page2) + "L" + strconv.Itoa(count)
						
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