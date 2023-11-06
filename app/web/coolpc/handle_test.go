package coolpc

import (
	"bytes"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/YWJSonic/ycore/module/myhtml"
	"github.com/YWJSonic/ycore/types"
	"github.com/YWJSonic/ycore/util"
	"golang.org/x/net/html"
)

func Test_Get(t *testing.T) {
	mockdata, _ := os.ReadFile(`C:\Users\sony7\Documents\Soho\bunker_space\bk_system\bk_reptile\app\web\coolpc\evaluate.php`)

	utf8Str, _ := util.Utf8ToBig5(mockdata)
	loginHtml := html.NewTokenizer(bytes.NewBuffer(utf8Str))

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
	date := util.ServerTimeNow().Format("2006-01-02")

	item_slice := []ItemInfo{}
	for idx, filter := range filters["select"] {
		idx++
		key := util.Sprintf("c%d", idx)
		typeName := typeMap[idx]            // 商品類型
		price_slice := dataMap[key].([]int) // 價格陣列
		for subidx, token := range filter.SubRes {
			if token.Data == "option" {
				for _, attr := range token.Attr {
					v, _ := strconv.Atoi(attr.Val)
					price := price_slice[v]

					// 有價格的才是實際商品
					if price > 0 {
						spliteIdx := strings.LastIndex(filter.SubContent[subidx], ",")
						item := ItemInfo{
							Date:          date,
							UpdateTime:    datetime,
							TypeName:      typeName,
							TypeId:        idx,
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
}
