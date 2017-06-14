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
	Huarunxt.Register()
}

var Huarunxt = &Spider{
	Name:        "华润信托",
	Description: "华润信托净值数据 [Auto Page] [http://www.crctrust.com/servlet/json?funcNo=904005&page=1&numPerPage=10&type=30&name=&order=sxrq&sort=asc]",
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
			ctx.Aid(map[string]interface{}{"loop": [3]int{27, 30, 32}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {

					page := 0

					for _, item := range aid["loop"].([3]int) {
						page++

						ctx.AddQueue(&request.Request{
							Url:  "http://www.crctrust.com/servlet/json?funcNo=904005&page=1&numPerPage=10&type=" + strconv.Itoa(item) + "&name=&order=sxrq&sort=asc",
							Rule: aid["Rule"].(string),
							Temp: map[string]interface{}{
								"type":        strconv.Itoa(item),
								"level1pages": page,
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

					var fundType string
					ctx.GetTemp("type", &fundType)
					var totalPage int

					for _, v := range infos {
						switch pageresult := v.(type) {
						case []interface{}:
							for k, u := range pageresult {
								fmt.Println(k)
								switch uu := u.(type) {
								case map[string]interface{}:
									totalPage = int(uu["totalPages"].(float64))
								}

							}
						}
					}

					var page1 int
					ctx.GetTemp("level1pages", page1)

					page2 := 0

					for i := 0; i < totalPage; i++ {

						page2++

						ctx.AddQueue(&request.Request{
							Url:  "http://www.crctrust.com/servlet/json?funcNo=904005&page=" + strconv.Itoa(i+1) + "&numPerPage=10&type=" + fundType + "&name=&order=sxrq&sort=asc",
							Rule: "获取结果",
							Temp: map[string]interface{}{
								"level1pages": page1,
								"level2pages": page2,
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
					"累计增长率",
					"净值日期",
				},
				ParseFunc: func(ctx *Context) {
					jsonData := ctx.GetText()
					infos := map[string]interface{}{}
					err := json.Unmarshal([]byte(jsonData), &infos)
					if err != nil {
						return
					}

					var jingzhiriqi string
					var danweijingzhi string
					var leijijingzhi string
					var mingchen string

					count := 0

					var page int
					ctx.GetTemp("level1pages", &page)

					var page2 int
					ctx.GetTemp("level2pages", &page2)

					for _, v := range infos {
						switch fundDatas := v.(type) {
						case []interface{}:
							for _, u := range fundDatas {
								switch uu := u.(type) {
								case map[string]interface{}:
									{
										switch uuu := uu["data"].(type) {
										case []interface{}:
											for _, uuuu := range uuu {

												switch uuuuu := uuuu.(type) {
												case map[string]interface{}:
													mingchen = uuuuu["jjjc"].(string)
													danweijingzhi = uuuuu["nav"].(string)
													leijijingzhi = uuuuu["rate_real"].(string)
													jingzhiriqi = uuuuu["tradedate"].(string)

													count++
													fundID := "XTHUARUN" + "P1" + strconv.Itoa(page) + "P2" + strconv.Itoa(page2) + "L" + strconv.Itoa(count)

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
						}
					}
				},
			},
		},
	},
}
