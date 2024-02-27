package stocktw

import (
	"bk_reptile/app/web/stocktw/stockdao"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/YWJSonic/ycore/util"
)

// 法人統計
type stock_threefoundation struct {
	Stat   string                       `json:"stat"`
	Date   string                       `json:"date"`
	Title  string                       `json:"title"`
	Fields []string                     `json:"fields"`
	Data   [][]interface{}              `json:"data"`
	Params stock_threefoundation_Params `json:"params"`
	Notes  []string                     `json:"notes"`
	Hints  string                       `json:"hints"`
}

type stock_threefoundation_Params struct {
	Type       string      `json:"type"`
	DayDate    string      `json:"dayDate"`
	Response   string      `json:"response"`
	Empty      string      `json:"_"`
	Controller string      `json:"controller"`
	Action     string      `json:"action"`
	Lang       string      `json:"lang"`
	MonthDate  interface{} `json:"monthDate"`
	WeekDate   interface{} `json:"weekDate"`
}

// 個股法人資料
type company_threefoundation struct {
	Stat       string          `json:"stat"`
	Date       string          `json:"date"`
	Title      string          `json:"title"`
	Hints      string          `json:"hints"`
	Fields     []string        `json:"fields"`
	Data       [][]interface{} `json:"data"`
	SelectType string          `json:"selectType"`
	Notes      []string        `json:"notes"`
	Total      int64           `json:"total"`
}

// 三大法人買賣資訊
func GetThreefoundation(date time.Time) ([]*stockdao.Stock_threefoundation, []*stockdao.Company_threefoundation, error) {

	res1, err := stock_three(date)
	if err != nil {
		if err != nil {
			return res1, []*stockdao.Company_threefoundation{}, err
		}
	}
	res2, err := company_three(date)
	if err != nil {
		return res1, res2, err
	}

	return res1, res2, nil
}

// 法人統計
func stock_three(date time.Time) ([]*stockdao.Stock_threefoundation, error) {
	res := []*stockdao.Stock_threefoundation{}
	url := fmt.Sprintf("https://www.twse.com.tw/rwd/zh/fund/BFI82U?type=day&dayDate=%s&response=json&_=%d", date.Format("20060102"), util.ServerTimeNow().UnixMilli())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(date)
		return res, err
	}

	req.Header.Add("Content-Type", "application/json")
	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return res, err
	}

	payload := &stock_threefoundation{}
	if err := json.Unmarshal(body, payload); err != nil {
		return res, err
	}

	for _, data := range payload.Data {
		stock := &stockdao.Stock_threefoundation{
			Date:             payload.Date,
			Stock_three_name: data[0].(string),
			Buy:              toInt(data[1]),
			Sell:             toInt(data[2]),
			Diff:             toInt(data[3]),
		}
		res = append(res, stock)
		// _ = sqldb.Stock_threefoundationUpdateOrCreate(stock)
	}

	fmt.Println("stock_three Done")
	return res, nil
}

// 個股法人買賣
func company_three(date time.Time) ([]*stockdao.Company_threefoundation, error) {
	res := []*stockdao.Company_threefoundation{}
	url := fmt.Sprintf("https://www.twse.com.tw/rwd/zh/fund/T86?date=%s&selectType=ALLBUT0999&response=json&_=%d", date.Format("20060102"), util.ServerTimeNow().UnixMilli())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(date)
		return res, err
	}

	req.Header.Add("Content-Type", "application/json")
	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return res, err
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return res, err
	}

	payload := &company_threefoundation{}
	if err := json.Unmarshal(body, payload); err != nil {
		fmt.Println(date, string(body))
		return res, err
	}

	for _, data := range payload.Data {
		if len(data) != 19 {
			fmt.Println("error", date, data)
			continue
		}
		stock := &stockdao.Company_threefoundation{
			Date:                          payload.Date,
			Company_id:                    data[0].(string),
			Company_name:                  data[1].(string),
			Global_china_stockbanker_buy:  toInt(data[2]),
			Global_china_stockbanker_sell: toInt(data[3]),
			Global_china_stockbanker_diff: toInt(data[4]),
			Global_stockbanker_buy:        toInt(data[5]),
			Global_stockbanker_sale:       toInt(data[6]),
			Global_stockbanker_diff:       toInt(data[7]),
			Stock_foundation_buy:          toInt(data[8]),
			Stock_foundation_sell:         toInt(data[9]),
			Stock_foundation_diff:         toInt(data[10]),
			Stockbanker_diff:              toInt(data[11]),
			Stockbanker_self_buy:          toInt(data[12]),
			Stockbanker_self_sell:         toInt(data[13]),
			Stockbanker_self_diff:         toInt(data[14]),
			Stockbanker_hedging_buy:       toInt(data[15]),
			Stockbanker_hedging_sell:      toInt(data[16]),
			Stockbanker_hedging_diff:      toInt(data[17]),
			Total_diff:                    toInt(data[18]),
		}
		res = append(res, stock)
		// _ = sqldb.Company_threefoundationUpdateOrCreate(stock)
	}

	// wg := sync.WaitGroup{}
	// group := intSplite(len(payload.Data), 30)
	// idx := 0
	// for _, count := range group {
	// 	wg.Add(1)
	// 	go func(dates [][]interface{}) {
	// 		for _, data := range dates {
	// 			if len(data) != 19 {
	// 				fmt.Println("error", date, data)
	// 				continue
	// 			}
	// 			stock := &stockdao.Company_threefoundation{
	// 				Date:                          payload.Date,
	// 				Company_id:                    data[0].(string),
	// 				Company_name:                  data[1].(string),
	// 				Global_china_stockbanker_buy:  toInt(data[2]),
	// 				Global_china_stockbanker_sell: toInt(data[3]),
	// 				Global_china_stockbanker_diff: toInt(data[4]),
	// 				Global_stockbanker_buy:        toInt(data[5]),
	// 				Global_stockbanker_sale:       toInt(data[6]),
	// 				Global_stockbanker_diff:       toInt(data[7]),
	// 				Stock_foundation_buy:          toInt(data[8]),
	// 				Stock_foundation_sell:         toInt(data[9]),
	// 				Stock_foundation_diff:         toInt(data[10]),
	// 				Stockbanker_diff:              toInt(data[11]),
	// 				Stockbanker_self_buy:          toInt(data[12]),
	// 				Stockbanker_self_sell:         toInt(data[13]),
	// 				Stockbanker_self_diff:         toInt(data[14]),
	// 				Stockbanker_hedging_buy:       toInt(data[15]),
	// 				Stockbanker_hedging_sell:      toInt(data[16]),
	// 				Stockbanker_hedging_diff:      toInt(data[17]),
	// 				Total_diff:                    toInt(data[18]),
	// 			}
	// 			_ = sqldb.Company_threefoundationUpdateOrCreate(stock)
	// 		}
	// 		wg.Done()
	// 	}(payload.Data[idx : idx+count])
	// 	idx += count
	// }
	// wg.Wait()

	fmt.Println("company_three Done")
	return res, nil
}
