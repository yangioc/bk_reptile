package stocktw

import (
	"bk_reptile/app/web/stocktw/stockdao"
	"log"
	"sort"
	"testing"
)

func Test_Analysis_financial_statements(t *testing.T) {

	datas := []stockdao.Company_financial_statements{}
	// if err := sqldb.Mysql().Where("year = ? AND company_id = ?", 2022, "1102").Find(&datas).Error; err != nil {
	// 	log.Fatalln(err)
	// }

	sort.Slice(datas, func(i, j int) bool {
		return datas[i].Session < datas[j].Session
	})

	if err := Analysis_Session_financial_statements(datas); err != nil {
		log.Fatalln(err)
	}

}
