package stocktw

import (
	"bk_reptile/app/web/stocktw/stockdao"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/YWJSonic/ycore/util"
)

type StockIndex struct {
	Tables []Table `json:"tables"`
	Params Params  `json:"params"`
	Stat   string  `json:"stat"`
	Date   string  `json:"date"`
}

type Params struct {
	Date       string `json:"date"`
	Type       string `json:"type"`
	Response   string `json:"response"`
	Empty      string `json:"_"`
	Controller string `json:"controller"`
	Action     string `json:"action"`
	Lang       string `json:"lang"`
}

type Table struct {
	Title  string          `json:"title"`
	Fields []string        `json:"fields"`
	Data   [][]interface{} `json:"data"`
	Hints  *string         `json:"hints,omitempty"`
	Notes  []string        `json:"notes"`
	Groups []Group         `json:"groups"`
}

type Group struct {
	Start int64  `json:"start"`
	Span  int64  `json:"span"`
	Title string `json:"title"`
}

// 每日 指數, 大盤統計, 收盤行情
func GetStockIndex(date time.Time) ([]*stockdao.Stock_index, []*stockdao.Stock_price, []*stockdao.Company_stock, error) {
	res1 := []*stockdao.Stock_index{}
	res2 := []*stockdao.Stock_price{}
	res3 := []*stockdao.Company_stock{}

	url := fmt.Sprintf("https://www.twse.com.tw/rwd/zh/afterTrading/MI_INDEX?date=%s&type=ALL&response=json&_=%d", date.Format("20060102"), util.ServerTimeNow().UnixMilli())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(date)
		return res1, res2, res3, err
	}

	req.Header.Add("Content-Type", "application/json")

	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return res1, res2, res3, err
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return res1, res2, res3, err
	}

	payload := &StockIndex{}
	if err := json.Unmarshal(body, payload); err != nil {
		return res1, res2, res3, err
	}

	wg := sync.WaitGroup{}
	for _, table := range payload.Tables {
		switch {
		case strings.Contains(table.Title, "價格指數"), strings.Contains(table.Title, "報酬指數"):
			wg.Add(1)
			go func(table Table) {
				for _, data := range table.Data {
					stock := &stockdao.Stock_index{
						Date:         payload.Date,
						Index_name:   data[0].(string),
						Index_number: toFloat32(data[1]),
					}
					res1 = append(res1, stock)
					// _ = sqldb.Stock_indexUpdateOrCreate(stock)
				}
				wg.Done()
			}(table)

		case strings.Contains(table.Title, "大盤統計資訊"):
			wg.Add(1)
			go func(table Table) {
				for _, data := range table.Data {
					stock := &stockdao.Stock_price{
						Date:               payload.Date,
						Stock_name:         data[0].(string),
						Stock_total_amount: toInt(data[1]),
						Stock_total_number: toInt(data[2]),
						Stock_total_count:  toInt(data[3]),
					}
					res2 = append(res2, stock)
					// _ = sqldb.Stock_priceUpdateOrCreate(stock)
				}
				wg.Done()
			}(table)

		case strings.Contains(table.Title, "每日收盤行情"):
			group := intSplite(len(table.Data), 30)
			idx := 0
			for _, count := range group {
				wg.Add(1)
				go func(dates [][]interface{}) {
					for _, data := range dates {
						stock := &stockdao.Company_stock{
							Date:               payload.Date,
							Company_id:         data[0].(string),
							Company_name:       data[1].(string),
							Transaction_number: toInt(data[2]),
							Transaction_count:  toInt(data[3]),
							Transaction_amount: toInt(data[4]),
							Price_open:         toFloat32(data[5]),
							Price_close:        toFloat32(data[8]),
							Price_max:          toFloat32(data[6]),
							Price_min:          toFloat32(data[7]),
						}
						res3 = append(res3, stock)
						// _ = sqldb.Company_stockUpdateOrCreate(stock)
					}
					wg.Done()
				}(table.Data[idx : idx+count])
				idx += count
			}
		}
	}

	wg.Wait()
	fmt.Println("GetStockIndex Done")
	return res1, res2, res3, nil
}

func toInt(idata interface{}) int {
	res := 0
	switch data := idata.(type) {
	case string:
		data = strings.ReplaceAll(data, ",", "")
		i, err := strconv.Atoi(data)
		if err != nil {
			panic(err)
		}
		res = i
	case float64:
		res = int(data)
	case int:
		res = data
	}
	return res
}
func toFloat32(idata interface{}) float32 {
	var res float32
	switch data := idata.(type) {
	case string:
		if data == "--" || data == "-" {
			return -1
		}
		data = strings.ReplaceAll(data, ",", "")
		f, err := strconv.ParseFloat(data, 32)
		if err != nil {
			panic(err)
		}
		res = float32(f)
	case float64:
		res = float32(data)
	case float32:
		res = data
	}

	return res
}
func intSplite(number int, spliteCount int) []int {
	res := []int{}
	if spliteCount > number {
		return append(res, number)
	} else if spliteCount == number {
		for i := 0; i < spliteCount; i++ {
			res = append(res, 1)
		}
		return res
	}

	groupContent := number / spliteCount

	for !(number-groupContent < 0) {
		res = append(res, groupContent)
		number -= groupContent
	}
	res = append(res, number)
	return res
}
