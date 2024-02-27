package stocktw

import (
	"bk_reptile/app/web/stocktw/stockdao"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
)

type Financial_statements_data struct {
	Revenue                   float64 // 總營收
	Total_equity              float64 // 權益總額
	Non_controlling_interests float64 // 非控制權益
	Total_assets              float64 // 資產總計
	Total_share_capital       float64 // 股本合計

	Profit                                       float64 // 合併稅後淨收入
	Operating_income                             float64 // 營業利潤
	Gross_profit_from_operations                 float64 // 營業毛利
	Earnings_per_share                           float64 // 每股盈餘
	Profit_from_continuing_operations_before_tax float64 // 合併稅前淨利

	Net_cash_flows_from_operating_activities float64 // 營業活動之淨現金流入
	Net_cash_flows_from_investing_activities float64 // 投資活動之淨現金流入
	Interest_expense                         float64 // 利息費用
}

// Session
// 1~3 "綜合損益表"為季度資料, "現金流量表"為累計季度資料
// 4 "綜合損益表", "現金流量表"為年度資料

func Analysis_Session_financial_statements(datas []stockdao.Company_financial_statements) error {

	if len(datas) == 0 {
		return errors.New("no data")
	}

	DataMap := map[int]map[string]*stockdao.Company_analysis{}

	var currentSessionData *Financial_statements_data
	var preSessionData *Financial_statements_data
	var err error
	for i, count := 0, len(datas); i < count; i++ {
		dataMap := map[string][][]string{}
		if err := json.Unmarshal([]byte(datas[i].Data), &dataMap); err != nil {
			return err
		}

		currentSessionData, err = Unmarshal_Financial_statements(dataMap)
		if err != nil {
			return err
		}

		session := &stockdao.Company_analysis{
			Year:       datas[i].Year,
			Session:    datas[i].Session,
			Company_id: datas[i].Company_id,
			Data_type:  stockdao.Company_analysis_Data_type_session,
		}

		session_Sum := &stockdao.Company_analysis{
			Year:       datas[i].Year,
			Session:    datas[i].Session,
			Company_id: datas[i].Company_id,
			Data_type:  stockdao.Company_analysis_Data_type_session_sum,
		}

		switch i {
		case 0, 1, 2:

			// 總營收
			{
				var preData float64
				if _, ok := DataMap[i-1]; ok {
					preData = DataMap[i-1][stockdao.Company_analysis_Data_type_session_sum].Revenue
				}

				session.Revenue = currentSessionData.Revenue
				session_Sum.Revenue = currentSessionData.Revenue + preData
			}

			// 營業毛利
			{
				var preData float64
				if _, ok := DataMap[i-1]; ok {
					preData = DataMap[i-1][stockdao.Company_analysis_Data_type_session_sum].Revenue
				}

				session.Gross_profit_from_operations = currentSessionData.Gross_profit_from_operations
				session_Sum.Gross_profit_from_operations = currentSessionData.Gross_profit_from_operations + preData
			}

			// 合併稅後淨收入
			{

				var preData float64
				if _, ok := DataMap[i-1]; ok {
					preData = DataMap[i-1][stockdao.Company_analysis_Data_type_session_sum].Net_income
				}

				session.Net_income = currentSessionData.Profit
				session_Sum.Net_income = currentSessionData.Profit + preData
			}

			// 營運現金流
			{
				var preData float64
				if preSessionData != nil {
					preData = preSessionData.Net_cash_flows_from_operating_activities
				}
				session.Operation_cash_flow = currentSessionData.Net_cash_flows_from_operating_activities - preData
				session_Sum.Operation_cash_flow = currentSessionData.Net_cash_flows_from_operating_activities
			}

			// 資本支出
			{
				var preData float64
				if preSessionData != nil {
					preData = preSessionData.Net_cash_flows_from_investing_activities
				}
				session.Cap_spending = currentSessionData.Net_cash_flows_from_investing_activities - preData
				session_Sum.Cap_spending = currentSessionData.Net_cash_flows_from_investing_activities
			}

			// 毛利率 Gross Margin // 代表產品的成本與總收入的關係。(一項產品扣掉直接成本，可從售價獲利幾%)
			// 毛利率 = Get_營業毛利_Gross_profit_from_operations(dataMap) / Get_營業收入合計_Revenue(dataMap)
			{
				var preDataA, preDataB float64
				if _, ok := DataMap[i-1]; ok {
					preDataA = DataMap[i-1][stockdao.Company_analysis_Data_type_session_sum].Gross_profit_from_operations
					preDataB = DataMap[i-1][stockdao.Company_analysis_Data_type_session_sum].Revenue
				}

				session.Gross_margin = currentSessionData.Gross_profit_from_operations / currentSessionData.Revenue
				session_Sum.Gross_margin = (currentSessionData.Gross_profit_from_operations + preDataA) / (currentSessionData.Revenue + preDataB)
			}

			// 自由現金流
			{
				var preDataA, preDataB float64
				if preSessionData != nil {
					preDataA = preSessionData.Net_cash_flows_from_operating_activities
					preDataB = preSessionData.Net_cash_flows_from_investing_activities
				}
				session.Free_cash_flow = (currentSessionData.Net_cash_flows_from_operating_activities - preDataA) + (currentSessionData.Net_cash_flows_from_investing_activities - preDataB)
				session_Sum.Free_cash_flow = currentSessionData.Net_cash_flows_from_operating_activities + currentSessionData.Net_cash_flows_from_investing_activities
			}

			DataMap[i] = map[string]*stockdao.Company_analysis{
				stockdao.Company_analysis_Data_type_session:     session,
				stockdao.Company_analysis_Data_type_session_sum: session_Sum,
			}

		case 3:
			year := &stockdao.Company_analysis{
				Year:       datas[i].Year,
				Session:    datas[i].Session,
				Company_id: datas[i].Company_id,
				Data_type:  stockdao.Company_analysis_Data_type_session_sum,
			}

			// 總營收
			{
				var preData float64
				if _, ok := DataMap[i-1]; ok {
					preData = DataMap[i-1][stockdao.Company_analysis_Data_type_session_sum].Revenue
				}

				year.Revenue = currentSessionData.Revenue
				session_Sum.Revenue = currentSessionData.Revenue
				session.Revenue = currentSessionData.Revenue - preData
			}

			DataMap[i] = map[string]*stockdao.Company_analysis{
				stockdao.Company_analysis_Data_type_session:     session,
				stockdao.Company_analysis_Data_type_session_sum: session_Sum,
				stockdao.Company_analysis_Data_type_year:        year,
			}
		}

		preSessionData = currentSessionData
	}

	// if err := sqldb.Mysql().Create(obj).Error; err != nil {
	// 	return err
	// }

	return nil
}

func Unmarshal_Financial_statements(currentData map[string][][]string) (*Financial_statements_data, error) {
	res := &Financial_statements_data{}

	a, err := ToFloat64(Get_營業收入合計_Revenue(currentData))
	if err != nil {
		return nil, err
	}
	res.Revenue = a

	a, err = ToFloat64(Get_權益總額_Total_equity(currentData))
	if err != nil {
		return nil, err
	}
	res.Total_equity = a

	a, err = ToFloat64(Get_非控制權益_Non_controlling_interests(currentData))
	if err != nil {
		return nil, err
	}
	res.Non_controlling_interests = a

	a, err = ToFloat64(Get_資產總計_Total_assets(currentData))
	if err != nil {
		return nil, err
	}
	res.Total_assets = a

	a, err = ToFloat64(Get_股本合計_Total_share_capital(currentData))
	if err != nil {
		return nil, err
	}
	res.Total_share_capital = a

	///

	a, err = ToFloat64(Get_合併稅後淨利_Profit(currentData))
	if err != nil {
		return nil, err
	}
	res.Profit = a

	a, err = ToFloat64(Get_營業利潤_Operating_income(currentData))
	if err != nil {
		return nil, err
	}
	res.Operating_income = a

	a, err = ToFloat64(Get_營業毛利_Gross_profit_from_operations(currentData))
	if err != nil {
		return nil, err
	}
	res.Gross_profit_from_operations = a

	a, err = strconv.ParseFloat(strings.ReplaceAll(Get_每股盈餘_Earnings_per_share(currentData), ",", ""), 64)
	if err != nil {
		return nil, err
	}
	res.Earnings_per_share = a

	a, err = ToFloat64(Get_合併稅前淨利_Profit_from_continuing_operations_before_tax(currentData))
	if err != nil {
		return nil, err
	}
	res.Profit_from_continuing_operations_before_tax = a

	///

	a, err = ToFloat64(Get_營業活動之淨現金流入_Net_cash_flows_from_operating_activities(currentData))
	if err != nil {
		return nil, err
	}
	res.Net_cash_flows_from_operating_activities = a

	a, err = ToFloat64(Get_投資活動之淨現金流入_Net_cash_flows_from_investing_activities(currentData))
	if err != nil {
		return nil, err
	}
	res.Net_cash_flows_from_investing_activities = a

	a, err = ToFloat64(Get_利息費用_Interest_expense(currentData))
	if err != nil {
		return nil, err
	}
	res.Interest_expense = a

	return res, nil
}

// func process() {
// 	// 總營收
// 	{
// 		a, err := ToFloat64(Get_營業收入合計_Revenue(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Revenue = a
// 	}

// 	// 合併稅後淨收入
// 	{
// 		a, err := ToFloat64(Get_合併稅後淨利_Profit(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Net_income = a
// 	}

// 	// 營運現金流
// 	{
// 		a, err := ToFloat64(Get_營業活動之淨現金流入_Net_cash_flows_from_operating_activities(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Operation_cash_flow = a
// 	}

// 	// 資本支出
// 	{
// 		a, err := ToFloat64(Get_投資活動之淨現金流入_Net_cash_flows_from_investing_activities(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Cap_spending = a
// 	}

// 	// 自由現金流
// 	{
// 		a, err := ToFloat64(Get_營業活動之淨現金流入_Net_cash_flows_from_operating_activities(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		b, err := ToFloat64(Get_投資活動之淨現金流入_Net_cash_flows_from_investing_activities(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Free_cash_flow = a + b
// 	}

// 	// 毛利率 Gross Margin // 代表產品的成本與總收入的關係。(一項產品扣掉直接成本，可從售價獲利幾%)
// 	// 毛利率 = Get_營業毛利_Gross_profit_from_operations(dataMap) / Get_營業收入合計_Revenue(dataMap)
// 	{
// 		a, err := ToFloat64(Get_營業毛利_Gross_profit_from_operations(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		b, err := ToFloat64(Get_營業收入合計_Revenue(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Gross_margin = a / b
// 	}

// 	// 營業利潤 Operating Income 公司透過商業活動所獲得的收入，通常是提供產品及服務而來。
// 	// 營業利潤 = Get_營業利潤_Operating_income(dataMap)
// 	{
// 		a, err := ToFloat64(Get_營業利潤_Operating_income(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Operating_income = a
// 	}

// 	// 營業利率 Operating Margin // 營業利率則是營業利潤與總營收的關係。(一項產品扣掉所有成本，可從售價獲利幾%)
// 	// 營業利率 = Get_營業利益_Operating_income(dataMap) / Get_營業收入合計_Revenue(dataMap)
// 	{
// 		a, err := ToFloat64(Get_營業利潤_Operating_income(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		b, err := ToFloat64(Get_營業收入合計_Revenue(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Operation_margin = a / b
// 	}

// 	// EPS 是公司的獲利指標，而通常EPS與股價也會有非常大的關連性。 *發行股數與淨收入存在隱藏資料無法從財報內計算
// 	// EPS = Get_每股盈餘_Earnings_per_share(dataMap)
// 	{
// 		a, err := strconv.ParseFloat(strings.ReplaceAll(Get_每股盈餘_Earnings_per_share(currentData), ",", ""), 64)
// 		if err != nil {
// 			return err
// 		}
// 		analy.Earnings_per_share = a
// 	}

// 	// 每股淨值 (BPS) 企業把所有資產變現，並償還債務後所剩餘的部份。
// 	// 每股淨值 = (Get_權益總額_Total_equity(dataMap) - Get_非控制權益_Non_controlling_interests(dataMap)) / (Get_股本合計_Total_share_capital(dataMap) / 10)
// 	{
// 		a, err := ToFloat64(Get_權益總額_Total_equity(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		b, err := ToFloat64(Get_非控制權益_Non_controlling_interests(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		c, err := ToFloat64(Get_股本合計_Total_share_capital(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Book_value_per_share = (a - b) / c / 10
// 	}

// 	// 股東權益報酬率 ROE
// 	// 股東權益報酬率 = Get_合併稅後淨利_Profit(dataMap) / Get_權益總額_Total_equity(dataMap)
// 	{
// 		a, err := ToFloat64(Get_合併稅後淨利_Profit(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		b, err := ToFloat64(Get_權益總額_Total_equity(currentData))
// 		if err != nil {
// 			return err
// 		}

// 		analy.Return_on_equity = a / b
// 	}

// 	// 資產報酬率 ROA
// 	// 資產報酬率 = Get_合併稅後淨利_Profit(dataMap) / Get_資產總計_Total_assets(dataMap)
// 	{
// 		a, err := ToFloat64(Get_合併稅後淨利_Profit(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		b, err := ToFloat64(Get_資產總計_Total_assets(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Return_on_assets = a / b
// 	}

// 	// 利息保障倍數(TIE) 公司從總營收中計算淨收入的比例
// 	// 利息保障倍數 = (Get_合併稅前淨利_Profit_from_continuing_operations_before_tax(dataMap) - Get_利息費用_Interest_expense(dataMap)) / Get_利息費用_Interest_expense(dataMap)
// 	{
// 		a, err := ToFloat64(Get_合併稅前淨利_Profit_from_continuing_operations_before_tax(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		b, err := ToFloat64(Get_利息費用_Interest_expense(currentData))
// 		if err != nil {
// 			return err
// 		}

// 		analy.Interest_coverage = (a + b) / b
// 	}

// 	// 淨利率% 公司從總營收中計算淨收入的比例
// 	// 淨利率 = Get_合併稅後淨利_Profit(dataMap) / Get_營業收入合計_Revenue(dataMap)
// 	{
// 		a, err := ToFloat64(Get_合併稅後淨利_Profit(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		b, err := ToFloat64(Get_營業收入合計_Revenue(currentData))
// 		if err != nil {
// 			return err
// 		}
// 		analy.Net_margin = a / b
// 	}
// }

/// 資產負債表 ----- 基礎資料為當日資料

func Get_資產總計_Total_assets(dataMap map[string][][]string) string {
	if _, ok := dataMap["資產負債表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["資產負債表"] {
		if datas[0] == "36XX" {
			return datas[3]
		}
	}

	return ""
}

func Get_負債總計_Total_liabilities(dataMap map[string][][]string) string {
	if _, ok := dataMap["資產負債表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["資產負債表"] {
		if datas[0] == "2XXX" {
			return datas[3]
		}
	}

	return ""
}

func Get_股本合計_Total_share_capital(dataMap map[string][][]string) string {
	if _, ok := dataMap["資產負債表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["資產負債表"] {
		if datas[0] == "3100" {
			return datas[3]
		}
	}

	return ""
}

func Get_權益總額_Total_equity(dataMap map[string][][]string) string {
	if _, ok := dataMap["資產負債表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["資產負債表"] {
		if datas[0] == "3XXX" {
			return datas[3]
		}
	}

	return ""
}

func Get_非控制權益_Non_controlling_interests(dataMap map[string][][]string) string {
	if _, ok := dataMap["資產負債表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["資產負債表"] {
		if datas[0] == "36XX" {
			return datas[3]
		}
	}

	return ""
}

/// 綜合損益表 ----- 基礎資料為當季資料

// 公司尚未扣除任何費用的總收入。
func Get_營業收入合計_Revenue(dataMap map[string][][]string) string {
	if _, ok := dataMap["綜合損益表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["綜合損益表"] {
		if datas[0] == "4000" {
			return datas[3]
		}
	}

	return ""
}

func Get_營業毛利_Gross_profit_from_operations(dataMap map[string][][]string) string {
	if _, ok := dataMap["綜合損益表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["綜合損益表"] {
		if datas[0] == "5900" {
			return datas[3]
		}
	}
	return ""
}

// 公司透過商業活動所獲得的收入，通常是提供產品及服務而來。
func Get_營業利潤_Operating_income(dataMap map[string][][]string) string {
	if _, ok := dataMap["綜合損益表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["綜合損益表"] {
		if datas[0] == "6900" {
			return datas[3]
		}
	}

	return ""
}

func Get_合併稅前淨利_Profit_from_continuing_operations_before_tax(dataMap map[string][][]string) string {
	if _, ok := dataMap["綜合損益表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["綜合損益表"] {
		if datas[0] == "7900" {
			return datas[3]
		}
	}

	return ""
}

// 是評估公司獲利能力的關鍵。 (淨收入)
func Get_合併稅後淨利_Profit(dataMap map[string][][]string) string {
	if _, ok := dataMap["綜合損益表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["綜合損益表"] {
		if datas[0] == "8200" {
			return datas[3]
		}
	}

	return ""
}

// 稅後淨利
func Get_母公司業主稅後淨利_Profit_attributable_to_owners_of_parent(dataMap map[string][][]string) string {
	if _, ok := dataMap["綜合損益表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["綜合損益表"] {
		if datas[0] == "8610" {
			return datas[3]
		}
	}

	return ""
}

func Get_每股盈餘_Earnings_per_share(dataMap map[string][][]string) string {
	if _, ok := dataMap["綜合損益表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["綜合損益表"] {
		if datas[0] == "9750" {
			return datas[3]
		}
	}

	return ""
}

/// 現金流量表 ----- 基礎資料為季累積資料

// 利息費用
func Get_利息費用_Interest_expense(dataMap map[string][][]string) string {
	if _, ok := dataMap["現金流量表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["現金流量表"] {
		if datas[0] == "A20900" {
			return datas[3]
		}
	}

	return ""
}

// 營業活動之淨現金流入
func Get_營業活動之淨現金流入_Net_cash_flows_from_operating_activities(dataMap map[string][][]string) string {
	if _, ok := dataMap["現金流量表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["現金流量表"] {
		if datas[0] == "AAAA" {
			return datas[3]
		}
	}

	return ""
}

// 投資活動之淨現金流入
func Get_投資活動之淨現金流入_Net_cash_flows_from_investing_activities(dataMap map[string][][]string) string {
	if _, ok := dataMap["現金流量表"]; !ok {
		return ""
	}

	for _, datas := range dataMap["現金流量表"] {
		if datas[0] == "BBBB" {
			return datas[3]
		}
	}

	return ""
}

/// 當期權益變動表 ----- 基礎資料為當日資料

/// 工具 -----

// 字串轉為福點數
//
// @params string 數字字串
//
// @params []bool idx 0: 數值是否為正號, 無設定為正號
//
// @return float64 結果
//
// @return error 錯誤資訊
func ToFloat64(s string, isNormalSign ...bool) (float64, error) {

	if s == "" {
		return 0, errors.New("no Data")
	}

	isNegative := false
	if s[0] == '(' {
		isNegative = true
	}

	s = strings.NewReplacer("(", "", ")", "", ",", "").Replace(s)
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	// 數值代表的正負
	if len(isNormalSign) > 0 && !isNormalSign[0] {
		v *= -1
	}

	// 是否有反向符號
	if isNegative {
		v *= -1
	}

	return v, nil
}
