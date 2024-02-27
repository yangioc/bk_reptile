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

type CompanyStockAnalysis struct {
	Stat       string          `json:"stat"`
	Date       string          `json:"date"`
	Title      string          `json:"title"`
	Fields     []string        `json:"fields"`
	Data       [][]interface{} `json:"data"`
	SelectType string          `json:"selectType"`
	Total      int64           `json:"total"`
}

func GetStockAnalysis(date time.Time) ([]*stockdao.Company_stock_analysis, error) {
	res := []*stockdao.Company_stock_analysis{}
	url := fmt.Sprintf("https://www.twse.com.tw/rwd/zh/afterTrading/BWIBBU_d?date=%s&selectType=ALL&response=json&_=%d", date.Format("20060102"), util.ServerTimeNow().UnixMilli())
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(date)
		return res, err
	}

	req.Header.Add("Content-Type", "application/json")

	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(date)
		return res, err
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return res, err
	}

	payload := &CompanyStockAnalysis{}
	if err := json.Unmarshal(body, payload); err != nil {
		return res, err
	}

	for _, data := range payload.Data {
		stock := &stockdao.Company_stock_analysis{
			Date:                payload.Date,
			Company_id:          data[0].(string),
			Company_name:        data[1].(string),
			Yield:               toFloat32(data[2]),
			Yield_year:          fmt.Sprint(data[3]),
			PE:                  toFloat32(data[4]),
			PBR:                 toFloat32(data[5]),
			Analysis_year_month: data[6].(string),
		}
		res = append(res, stock)
		// _ = sqldb.Company_stock_analysisUpdateOrCreate(stock)
	}

	fmt.Println("stock_analysis Done")
	return res, nil
}
