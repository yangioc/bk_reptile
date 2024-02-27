package stocktw

import (
	"bk_reptile/app/web/stocktw/stockdao"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
	"time"
)

func Test_Get_twse_financial_statements(t *testing.T) {

	sessions := []string{"1", "2", "3", "4"}
	for i := 2019; i <= 2023; i++ {
		for _, session := range sessions {
			if err := Get_twse_financial_statements("3711", i, session); err != nil {
				log.Fatalln(err)
			}
			time.Sleep(time.Second * 3)
		}
	}

}

func Test_Unmarshal_twse_financial_statements(t *testing.T) {

	data, err := os.ReadFile(".\\twse_financial_statements_error_mock.html")
	if err != nil {
		log.Fatalln(err)
	}

	res, err := Unmarshal_twse_financial_statements(data)
	if err != nil {
		log.Fatalln(err)
	}
	for k, v := range res {
		fmt.Println(k)
		for _, vv := range v {
			fmt.Println(vv)
		}
	}

	js, _ := json.Marshal(res)
	fmt.Println(string(js))

	payload := &stockdao.Company_financial_statements{
		Year:         2023,
		Session:      "1",
		Company_id:   "1102",
		Data:         string(js),
		Data_version: stockdao.Company_financial_statements_DataVersion,
	}
	fmt.Println(payload)
}
