package pholcus_lib

import (
	// 基础包
	//"github.com/henrylee2cn/pholcus/common/goquery"                        //DOM解析
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	// "github.com/henrylee2cn/pholcus/logs"           //信息输出
	. "github.com/henrylee2cn/pholcus/app/spider" //必需
	// . "github.com/henrylee2cn/pholcus/app/spider/common" //选用

	// net包
	// "net/http" //设置http.Header
	// "net/url"

	// 编码包
	// "encoding/xml"
	 "encoding/json"

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
	Zgdwjjmyxt.Register()
}

var Zgdwjjmyxt = &Spider{
	Name:        "中国对外经济贸易信托",
	Description: "中国对外经济贸易信托净值数据 [Auto Page] [http://www.fotic.com.cn/tabid/141/Default.aspx]",
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
			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 5}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {
				
				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
                                   
                    ctx.AddQueue(&request.Request{
                        Url:  "http://www.fotic.com.cn/tabid/141/Default.aspx",
                        Rule: aid["Rule"].(string),
                    })
					return nil
				},
				ParseFunc: func(ctx *Context) {
                    
					page := 0
                    for i:= 1; i <= 11400; i++ {

						page++
                       ctx.AddQueue(&request.Request{
                            Url: "http://www.fotic.com.cn/DesktopModules/ProductJZ/GetJsonResult.ashx?programName=&sDate=&eDate=&pageNo=" + strconv.Itoa(i) + "&pageSize=10",
                            Rule: "获取结果",
							Temp: map[string]interface{}{
								"pages": page,
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
                    
					jsonData := ctx.GetText()
					infos := map[string]interface{}{}
					err := json.Unmarshal([]byte(jsonData), &infos)
					if err != nil {
						return
					}

					var page int
					ctx.GetTemp("pages", &page)
					count := 0
					
					var jingzhiriqi string
					var danweijingzhi string
					var leijijingzhi string
					var mingchen string
					for _, v := range infos {
						switch vv := v.(type) {
						case []interface{}:
							for _, u := range vv {
								switch uu := u.(type) {
								case map[string]interface{}:
                                        {
                                            danweijingzhi = uu["netvalue"].(string)
                                            leijijingzhi = uu["netvalue"].(string)
                                            jingzhiriqi = uu["date"].(string)
                                            mingchen = uu["projectnameshort"].(string)
                                        }

										count++
										fundID := "XTDUIWAIJINGJIMAOYI" + "P" + strconv.Itoa(page) + "L" + strconv.Itoa(count)
                                        
                                        ctx.Output(map[int]interface{}{
											0: fundID,
                                            1: mingchen,
                                            2: danweijingzhi,
                                            3: leijijingzhi,
                                            4: jingzhiriqi,
                                        })
									default:
										fmt.Println("unknown type")
									}
								}
							default:
								fmt.Println("unknown type")
							}
						}
				},
			},
			
		},
	},
}
