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
)

func init() {
	Hrxt.Register()
}

var Hrxt = &Spider{
	Name:        "华融信托",
	Description: "华融信托净值数据 [Auto Page] [http://www.huarongtrust.com.cn/am/zjxtczjzpl.aspx]",
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
			ctx.Aid(map[string]interface{}{"loop": [2]int{1, 5}, "Rule": "生成请求"}, "生成请求")
		},

		Trunk: map[string]*Rule{

			"生成请求": {

				//注意：有无字段语义和是否输出数据必须保持一致
				ItemFields: []string{
					"名称",
					"净值",
					"累计净值",
					"估值日期",
				},

				AidFunc: func(ctx *Context, aid map[string]interface{}) interface{} {
					for loop := aid["loop"].([2]int); loop[0] < loop[1]; loop[0]++ {
						/*
							ctx.AddQueue(&request.Request{
								Url:  "http://www.huarongtrust.com.cn/am/zjxtczjzpl.aspx?"+"__VIEWSTATE=/wEPDwULLTIwNjA2NTU3MDUPZBYCAgMPZBYEAgsPFgIeC18hSXRlbUNvdW50AhAWIGYPZBYCZg8VDQExBzIwMTIyMDEM5q2j5byY5peX6IOcDOato+W8mOaXl+iDnAkyMDE2LjQuMjIBLQEtBjEuMDAwMAYxLjAwMDAFMC4wMCUGLTEuNTUlCTIwMTYuNS4xMQcyMDEyMjAxZAIBD2QWAmYPFQ0BMgcyMDEyMjAyC+iejeaxhzUw5Y+3C+iejeaxhzUw5Y+3CTIwMTYuNC4yMAEtAS0GMC45OTc3BjAuOTk3NwYtMC4yMyUGLTIuMDAlCTIwMTYuNS4xMQcyMDEyMjAyZAICD2QWAmYPFQ0BMwcyMDEyMTk0Fuebm+S4luaZr+aWsOetlueVpTHlj7cV55ub5LiW5pmv5paw562W55WlLi4uCTIwMTUuOS4xOAEtAS0GMC45NzUxBjAuOTc1MQYtMi40OSUGLTUuOTYlCTIwMTYuNS4xMQcyMDEyMTk0ZAIDD2QWAmYPFQ0BNAcyMDEyMTkxDeWbveS/oUZPRjLlj7cN5Zu95L+hRk9GMi4uLggyMDE1LjguNwEtAS0GMS4wMjY3BjEuMDI2NwUyLjY3JQctMjIuMTklCTIwMTYuNS4xMQcyMDEyMTkxZAIED2QWAmYPFQ0BNQcyMDEyMTkzDOato+W8mOmUkOaEjwzmraPlvJjplJDmhI8JMjAxNS44LjI1AS0BLQYwLjk4NTMGMC45ODUzBi0xLjQ3JQYtMS43NCUJMjAxNi41LjExBzIwMTIxOTNkAgUPZBYCZg8VDQE2BzIwMTIxODcQ6L+q6b6Z5oqV6LWEMeWPtxDov6rpvpnmipXotYQx5Y+3CDIwMTUuNy4zAS0BLQYwLjk2OTMGMC45NjkzBi0zLjA3JQctMjAuOTglCTIwMTYuNS4xMQcyMDEyMTg3ZAIGD2QWAmYPFQ0BNwcyMDEyMTg4EOmTtumZhei1hOacrDLlj7cQ6ZO26ZmF6LWE5pysMuWPtwkyMDE1LjcuMjIBLQEtBjAuOTU5NAYwLjk1OTQGLTQuMDYlBy0yNy42NCUJMjAxNi41LjExBzIwMTIxODhkAgcPZBYCZg8VDQE4BzIwMTIxODQQ5rW36KW/5pmf5Lm+N+WPtxDmtbfopb/mmZ/kub435Y+3CTIwMTUuNi4xMgEtAS0GMC42MDI1BjAuNjAyNQctMzkuNzUlBy00My42MSUJMjAxNi41LjExBzIwMTIxODRkAggPZBYCZg8VDQE5BzIwMTIxODMK5q2j5byYMuWPtwrmraPlvJgy5Y+3CTIwMTUuNi4xMAEtAS0GMC44OTUzBjAuODk1MwctMTAuNDclBy00Mi45NSUJMjAxNi41LjExBzIwMTIxODNkAgkPZBYCZg8VDQIxMAcyMDEyMTgwD+W3peihjOWbveS/oUZPRhHlt6XooYzlm73kv6FGTy4uLgkyMDE1LjUuMjYBLQEtBjEuMDUwNwYxLjA1MDcFNS4wNyUHLTQwLjY4JQkyMDE2LjUuMTEHMjAxMjE4MGQCCg9kFgJmDxUNAjExBzIwMTIxNzcK5rGH6KOVM+acnwrmsYfoo5Uz5pyfCTIwMTUuNS4yMgEtAS0GMC44ODMxBjAuODgzMQctMTEuNjklBy0zNy40NSUJMjAxNi41LjExBzIwMTIxNzdkAgsPZBYCZg8VDQIxMgcyMDEyMTc4EOS4ieaZuuWkqem4vzHlj7cQ5LiJ5pm65aSp6bi/MeWPtwkyMDE1LjUuMjIBLQEtBjEuMDk2NQYxLjA5NjUFOS42NSUHLTM3LjQ1JQkyMDE2LjUuMTEHMjAxMjE3OGQCDA9kFgJmDxUNAjEzBzIwMTIxNjUL6J6N5rGHNDDlj7cL6J6N5rGHNDDlj7cIMjAxNS40LjcBLQEtBjAuOTM5MgYwLjkzOTIGLTYuMDglBy0yNi40NiUJMjAxNi41LjExBzIwMTIxNjVkAg0PZBYCZg8VDQIxNAcyMDEyMTc0EuawuOi1oui1hOS6p+S4gOacnxLmsLjotaLotYTkuqfkuIDmnJ8JMjAxNS40LjI4AS0BLQYwLjkwOTUGMC45MDk1Bi05LjA1JQctMzQuOTIlCTIwMTYuNS4xMQcyMDEyMTc0ZAIOD2QWAmYPFQ0CMTUHMjAxMjE3MQrph5HmmZ815Y+3CumHkeaZnzXlj7cJMjAxNS40LjIyAS0BLQYxLjI2MDUGMS4yNjA1BjI2LjA1JQctMzMuNzclCTIwMTYuNS4xMQcyMDEyMTcxZAIPD2QWAmYPFQ0CMTYHMjAxMjE2NhLkuJzmupDlmInnm4jkuIPmnJ8S5Lic5rqQ5ZiJ55uI5LiD5pyfCTIwMTUuNC4xNwEtAS0GMC43Nzc5BjAuNzc3OQctMjIuMjElBy0zMi4wNSUJMjAxNi41LjExBzIwMTIxNjZkAg0PDxYCHgtSZWNvcmRjb3VudAInZGRky/e0YsCMNstu7DOsubeC6i0jePs=&__VIEWSTATEGENERATOR=4A2FE269&__EVENTTARGET=AspNetPager1&__EVENTARGUMENT=" + strconv.Itoa(loop[0]) + "&__EVENTVALIDATION=/wEWBgKzra2bAwLs0bLrBgLs0fbZDALs0Yq1BQLs0e58AoznisYGxMAJZLTGiQCuIQrGDG1PtaH8ubc=&TextBox1=&TextBox2=&TextBox3=&TextBox4=",
								Rule: aid["Rule"].(string),
								Method: "POST",
							})
						*/
						ctx.AddQueue(&request.Request{
							Url:      "http://www.huarongtrust.com.cn/am/zjxtczjzpl.aspx",
							Rule:     aid["Rule"].(string),
							Method:   "POST",
							PostData: "__VIEWSTATE=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&__VIEWSTATEGENERATOR=4A2FE269&__EVENTTARGET=AspNetPager1&__EVENTARGUMENT=" + strconv.Itoa(loop[0]) + "&__EVENTVALIDATION=/wEWBgKzra2bAwLs0bLrBgLs0fbZDALs0Yq1BQLs0e58AoznisYGxMAJZLTGiQCuIQrGDG1PtaH8ubc=&TextBox1=&TextBox2=&TextBox3=&TextBox4=",
						})
					}
					return nil
				},
				ParseFunc: func(ctx *Context) {
					query := ctx.GetDom()

					ss := query.Find("#myTab3_Content0 dl").Find("dd")

					count := 0
					var mingchen string
					var jingzhi string
					var leijijingzhi string
					var guzhiriqi string
					ss.Each(func(i int, goq *goquery.Selection) {

						count = count + 1

						if count%9 == 2 {
							mingchen = goq.Text()
						}

						if count%9 == 4 {
							jingzhi = goq.Text()
						}

						if count%9 == 5 {
							leijijingzhi = goq.Text()
						}

						if count%9 == 8 {
							guzhiriqi = goq.Text()
						}

						if count >= 9 && count%9 == 0 {
							ctx.Output(map[int]interface{}{
								0: mingchen,
								1: jingzhi,
								2: leijijingzhi,
								3: guzhiriqi,
							})
						}
					})
				},
			},
		},
	},
}
