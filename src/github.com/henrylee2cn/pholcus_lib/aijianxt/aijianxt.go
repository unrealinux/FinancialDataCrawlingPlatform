package pholcus_lib

import (
	// 基础包
	"github.com/henrylee2cn/pholcus/app/downloader/request" //必需
	. "github.com/henrylee2cn/pholcus/app/spider"           //必需
	// "github.com/henrylee2cn/pholcus/logs"           //信息输出
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
	"encoding/json"
	"fmt"
	"strings"
)

func init() {
	Aijianxt.Register()
}

var Aijianxt = &Spider{
	Name:        "爱建信托",
	Description: "爱建信托净值数据 [Auto Page] [http://www.ajxt.com.cn/Channel/3755?_tp_ptpro=1]",
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
							Url: "http://www.ajxt.com.cn/ajQuery/query.jsp?query=getNetvalue&channelId=77&query1=null",
							Rule: aid["Rule"].(string),
							Temp: map[string]interface{}{
								"level1pages": page,
							},
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {


					jsonData := strings.TrimSpace(ctx.GetText())
					fmt.Println(jsonData)

					var infos []map[string]interface{}
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
					page = ctx.GetTemp("level1pages", &page).(int)

					for _, uu := range infos {

								mingchen = uu["string1"].(string)
								danweijingzhi = uu["string2"].(string)
								leijijingzhi = uu["string2"].(string)
								jingzhiriqi = uu["string3"].(string)

								count++
								fundID := "XTAIJIAN" + "P1" + strconv.Itoa(page) + "L" + strconv.Itoa(count)

								ctx.Output(map[int]interface{}{
									0: fundID,
									1: mingchen,
									2: danweijingzhi,
									3: leijijingzhi,
									4: jingzhiriqi,
								})
						}
				},
			},
		},
	},
}
