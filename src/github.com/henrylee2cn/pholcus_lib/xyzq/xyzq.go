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
	//"reflect"
)

func init() {
	Xyzq.Register()
}

var Xyzq = &Spider{
	Name:        "兴业证券私募基金",
	Description: "兴业证券私募基金数据 [Auto Page] [http://www.xyzq.com.cn/xyzq/assetcustody/index.jsp?openflag=10016&lightflag=10877]",
	// Pausetime: 300,
	// Keyin:   KEYIN,
	// Limit:        LIMIT,
	NotDefaultField: true,

	Namespace: func(*Spider) string {
		return "zhengquan"
	},
	// 子命名空间相对于表名，可依赖具体数据内容，可选
	SubNamespace: func(self *Spider, dataCell map[string]interface{}) string {
		return "fund_src_nav"
	},

	EnableCookie: false,
	RuleTree: &RuleTree{

		Root: func(ctx *Context) {
			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 15}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {

					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						ctx.AddQueue(&request.Request{
							Url:  "http://www.xyzq.com.cn/xyzq/assetcustody/GetProduct.do?flag=01&page=" + strconv.Itoa(loop[0]) + "&name=",
							Rule: aid["Rule"].(string),
							Temp: map[string]interface{}{
								"level1pages": loop[0],
							},
						})
					}

					return nil
				},

				ParseFunc: func(ctx *Context) {

					jsonData := ctx.GetText()
					infos := map[string]interface{}{}
					err := json.Unmarshal([]byte(jsonData), &infos)
					if err != nil {
						return
					}

					var page int
					ctx.GetTemp("level1pages", &page)

					for _, v := range infos {
						switch vv := v.(type) {
						case []interface{}:
							for _, u := range vv {

								switch uu := u.(type) {
								case map[string]interface{}:
									mingchenValue := uu["name"].(string)
									fundNoValue := uu["fundNo"].(string)

									page := 0
									for i := 1; i < 11; i++ {

										page++
										ctx.AddQueue(&request.Request{
											Url:  "http://www.xyzq.com.cn/xyzq/assetcustody/GetProductValue.do?no=" + fundNoValue + "&page=" + strconv.Itoa(i),
											Rule: "获取结果",
											Temp: map[string]interface{}{
												"mingchen":    mingchenValue,
												"level1pages": page,
												"level2pages": i,
											},
										})
									}
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

			"获取结果": {
				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"基金ID",
					"名称",
					"净值",
					"累计净值",
					"净值日期",
				},
				ParseFunc: func(ctx *Context) {
					jsonData := ctx.GetText()
					infos := map[string]interface{}{}
					err := json.Unmarshal([]byte(jsonData), &infos)
					if err != nil {
						return
					}

					var page int
					ctx.GetTemp("level1pages", &page)
					var page2 int
					ctx.GetTemp("level2pages", &page2)
					count := 0

					var jingzhiriqi string
					var danweijingzhi float64
					var leijijingzhi float64
					var mingchen string
					for _, v := range infos {
						switch vv := v.(type) {
						case []interface{}:
							for _, u := range vv {
								switch uu := u.(type) {
								case map[string]interface{}:
									for _, noValue := range uu {
										switch noValueV := noValue.(type) {
										case map[string]interface{}:
											{
												danweijingzhi = noValueV["fundValue"].(float64)
												leijijingzhi = noValueV["fundValue1"].(float64)
												jingzhiriqi = noValueV["valueTime"].(string)
											}

											count++
											fundID := "XTDUIWAIJINGJIMAOYI" + "P1" + strconv.Itoa(page) + "P2" + strconv.Itoa(page2) + "L" + strconv.Itoa(count)

											ctx.GetTemp("mingchen", &mingchen)
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
						default:
							fmt.Println("unknown type")
						}
					}
				},
			},
		},
	},
}
