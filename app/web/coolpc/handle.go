package coolpc

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/YWJSonic/ycore/module/myhtml"
	"github.com/YWJSonic/ycore/types"
	"github.com/YWJSonic/ycore/util"

	"golang.org/x/net/html"
)

func GetWeb() ([]ItemInfo, error) {
	var err error
	var req *http.Request
	if req, err = http.NewRequest("GET", WebPage, nil); err != nil {
		return nil, err
	}

	var res *http.Response
	if res, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}
	defer res.Body.Close()

	loginHtml := html.NewTokenizer(res.Body)

	// 撈取網頁資料
	filters := map[types.TokenTypeName][]*myhtml.FilterObj{} // 網頁標籤資料
	for i, count := 1, 31; i < count; i++ {
		filter := &myhtml.FilterObj{}
		filter.FiltAttrs = append(filter.FiltAttrs,
			html.Attribute{
				Key: "name",
				Val: util.Sprintf("n%d", i),
			})
		filter.Operation = append(filter.Operation, myhtml.FilterOperation_GetSubToken, myhtml.FilterOperation_GetSubcContent)
		filters["select"] = append(filters["select"], filter)
	}
	filters["script"] = append(filters["script"], &myhtml.FilterObj{
		Operation: []int{
			myhtml.FilterOperation_GetContent,
		},
	})
	myhtml.HtmlLoopFilterOne(loginHtml, filters)

	var dataMap map[string]interface{}
	if len(filters["script"][0].Content) > 0 {
		data := filters["script"][0].Content[0] // 提取價格資料
		dataMap = spliteScriptData(data)        // 篩選價格資料
	}

	datetime := util.ParseJavaUnixSec(util.ServerTimeNow())
	dateStr := util.ServerTimeNow().Format("2006-01-02")

	item_slice := []ItemInfo{}
	for idx, filter := range filters["select"] {
		idx++
		key := util.Sprintf("c%d", idx)
		typeName := typeMap[idx]               // 商品類型
		price_slice := dataMap[key].([]string) // 價格陣列
		for subidx, token := range filter.SubRes {
			if token.Data == "option" {
				for _, attr := range token.Attr {
					v, _ := strconv.Atoi(attr.Val)
					price := price_slice[v]

					// 有價格的才是實際商品
					if price != "0" {
						spliteIdx := strings.LastIndex(filter.SubContent[subidx], ",")
						item := ItemInfo{
							Date:          dateStr,
							UpdateTime:    datetime,
							TypeName:      typeName,
							TypeId:        key,
							Price:         price,
							Content:       filter.SubContent[subidx][:spliteIdx],
							OriginContent: filter.SubContent[subidx],
						}

						item_slice = append(item_slice, item)
					}
				}
			}
		}
	}
	return item_slice, nil
}

func spliteScriptData(data string) map[string]interface{} {
	res := map[string]interface{}{}

	spliteStr := strings.Split(data, "\n")
	for _, dataStr := range spliteStr {
		// 判斷格式
		idx := strings.IndexByte(dataStr, '=')
		if idx < 0 {
			continue
		}

		// 確認資料結構
		key := dataStr[0:idx]

		// 排除 c類型以外資料
		if key[0] != 'c' {
			continue
		}

		valStr := dataStr[idx+1:]
		if len(valStr) < 2 || (valStr[0] != '[' && valStr[len(valStr)-1] != ']') {
			continue
		}

		// 資料解析
		values := []string{}
		valueSplite := strings.Split(dataStr[idx+2:len(dataStr)-1], ",")
		for _, valueStr := range valueSplite {
			values = append(values, valueStr)
		}

		res[key] = values
	}
	return res
}
