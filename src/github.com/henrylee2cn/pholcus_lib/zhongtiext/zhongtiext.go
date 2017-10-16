package pholcus_lib

import (
	// 基础包
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	//"github.com/henrylee2cn/pholcus/common/goquery"         //DOM解析
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
	//"strconv"
	//"strings"
	// 其他包
	// "fmt"
	// "math"
	// "time"
	//"log"
	//"log"
	"encoding/json"
	"strconv"
	"fmt"
)

func init() {
	Zhongtiext.Register()
}

var Zhongtiext = &Spider{
	Name:        "中铁信托",
	Description: "中铁信托净值数据 [Auto Page] [http://www.crtrust.com/front/getProductsXTZQbyPage_195401472991.jhtml]",
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
			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 2}, "Rule": "生成请求"}, "生成请求")
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
							Url:  "http://www.crtrust.com/front/getProductsXTZQbyPage_195401472991.jhtml",
							Rule: aid["Rule"].(string),
							Temp: map[string]interface{}{
								"pages": page,
							},
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {

					jsonData := ctx.GetText()
					fmt.Println(jsonData)

					infos := map[string]interface{}{}
					err := json.Unmarshal([]byte(jsonData), &infos)
					if err != nil {
						return
					}
					fmt.Println(infos)

					var page int
					page = ctx.GetTemp("pages", &page).(int)

					var mingchen string
					var jingzhi string
					var leijijingzhi string
					var guzhiriqi string


					count := 0
					for _, v := range infos {
						switch fundDatas := v.(type) {
						case []interface{}:
							fmt.Println(v)
							for _, u := range fundDatas{
								fmt.Println(u)
								switch uu := u.(type) {
								case map[string]interface{}:

									fmt.Println(uu)



									if uu["Wjz"] != nil && uu["Wjzrq"] != nil{

										mingchen = uu["Wcpmc"].(string)
										jingzhi = uu["Wjz"].(string)
										leijijingzhi = uu["Wjz"].(string)
										guzhiriqi = uu["Wjzrq"].(string)

										count++
										fundID := "XTZHONGTIE" + "P1" + strconv.Itoa(page) + "L" + strconv.Itoa(count)

										ctx.Output(map[int]interface{}{
											0: fundID,
											1: mingchen,
											2: jingzhi,
											3: leijijingzhi,
											4: guzhiriqi,
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
		},
	},
}
