package stocktw

import (
	"bk_reptile/app/web/stocktw/stockdao"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/YWJSonic/ycore/module/myhtml"
	"github.com/YWJSonic/ycore/util"
	"golang.org/x/net/html"
)

// https://mops.twse.com.tw/server-java/t164sb01?step=1&CO_ID=1102&SYEAR=2023&SSEASON=1&REPORT_ID=C
// 財務報表_季度統計
// 單位：新臺幣仟元
// 每股盈餘單位：新台幣元
// 股數：股

type Financial_Statements struct {
	*Balance_sheet                  // 資產負債表
	*Comprehensive_income_statement // 綜合損益表
	*Cash_flow_statement            // 現金流量表
}

// 資產負債表
type Balance_sheet struct {
	*Property               // 資產
	*Liabilities_and_equity // 負債及權益
}

// 資產
type Property struct {
	*Current_assets     // 流動資產
	*Non_current_assets // 非流動資產

	Code1XXX string `json:"code_1XXX"` // 資產總計
	Code3997 string `json:"code_3997"` // 待註銷股本股數
	Code3998 string `json:"code_3998"` // 預收股款（權益項下）之約當發行股數
	Code3999 string `json:"code_3999"` // 母公司暨子公司所持有之母公司庫藏股股數（單位：股）
}

// 流動資產
type Current_assets struct {
	Code1100 string `json:"code_1100"` // 現金及約當現金
	Code1110 string `json:"code_1110"` // 透過損益按公允價值衡量之金融資產－流動
	Code1120 string `json:"code_1120"` // 透過其他綜合損益按公允價值衡量之金融資產－流動
	Code1136 string `json:"code_1136"` // 按攤銷後成本衡量之金融資產－流動
	Code1140 string `json:"code_1140"` // 合約資產－流動
	Code1150 string `json:"code_1150"` // 應收票據淨額
	Code1170 string `json:"code_1170"` // 應收帳款淨額
	Code1180 string `json:"code_1180"` // 應收帳款－關係人淨額
	Code1200 string `json:"code_1200"` // 其他應收款
	Code1220 string `json:"code_1220"` // 本期所得稅資產
	Code130X string `json:"code_130X"` // 存貨
	Code1410 string `json:"code_1410"` // 預付款項
	Code1470 string `json:"code_1470"` // 其他流動資產
	Code11XX string `json:"code_11XX"` // 流動資產合計

}

// 非流動資產
type Non_current_assets struct {
	Code1510 string `json:"code_1510"` // 透過損益按公允價值衡量之金融資產－非流動
	Code1517 string `json:"code_1517"` // 透過其他綜合損益按公允價值衡量之金融資產－非流動
	Code1535 string `json:"code_1535"` // 按攤銷後成本衡量之金融資產－非流動
	Code1550 string `json:"code_1550"` // 採用權益法之投資
	Code1600 string `json:"code_1600"` // 不動產、廠房及設備
	Code1755 string `json:"code_1755"` // 使用權資產
	Code1760 string `json:"code_1760"` // 投資性不動產淨額
	Code1780 string `json:"code_1780"` // 無形資產
	Code1840 string `json:"code_1840"` // 遞延所得稅資產
	Code1900 string `json:"code_1900"` // 其他非流動資產總和

	*Other_non_current_assets // 其他非流動資產細項

	Code15XX string `json:"code_15XX"` // 非流動資產合計

}

// 其他非流動資產
type Other_non_current_assets struct {
	Code1920 string `json:"code_1920"` // 存出保證金
	Code194D string `json:"code_194D"` // 長期應收融資租賃款淨額
	Code1975 string `json:"code_1975"` // 淨確定福利資產-非流動
	Code1990 string `json:"code_1990"` // 其他非流動資產－其他
	Code1995 string `json:"code_1995"` // 其他非流動資產－其他
}

// 負債及權益
type Liabilities_and_equity struct {
	*Liabilities // 負債
	*Equity      // 權益

	Code3X2X string `json:"code_3X2X"` // 負債及權益總計
}

// 負債
type Liabilities struct {
	*Current_liabilities     // 流動負債
	*Non_current_liabilities // 非流動負債

	Code2XXX string `json:"code_2XXX"` // 負債總計
}

// 流動負債
type Current_liabilities struct {
	Code2100 string `json:"code_2100"` // 短期借款
	Code2110 string `json:"code_2110"` // 應付短期票券
	Code2120 string `json:"code_2120"` // 透過損益按公允價值衡量之金融負債－流動
	Code2130 string `json:"code_2130"` // 合約負債－流動
	Code2170 string `json:"code_2170"` // 應付帳款
	Code2180 string `json:"code_2180"` // 應付帳款－關係人
	Code2200 string `json:"code_2200"` // 其他應付款
	Code2216 string `json:"code_2216"` // 應付股利
	Code2219 string `json:"code_2219"` // 其他應付款－其他
	Code2230 string `json:"code_2230"` // 本期所得稅負債
	Code2250 string `json:"code_2250"` // 負債準備－流動
	Code2280 string `json:"code_2280"` // 租賃負債－流動
	Code2300 string `json:"code_2300"` // 其他流動負債
	Code2310 string `json:"code_2310"` // 預收款項
	Code2313 string `json:"code_2313"` // 遞延收入
	Code2320 string `json:"code_2320"` // 一年或一營業週期內到期長期負債
	Code2399 string `json:"code_2399"` // 其他流動負債－其他
	Code21XX string `json:"code_21XX"` // 流動負債合計

}

// 非流動負債
type Non_current_liabilities struct {
	Code2540 string `json:"code_2540"` // 長期借款
	Code2570 string `json:"code_2570"` // 遞延所得稅負債
	Code2580 string `json:"code_2580"` // 租賃負債－非流動
	Code2600 string `json:"code_2600"` // 其他非流動負債
	Code2670 string `json:"code_2670"` // 其他非流動負債－其他
	Code25XX string `json:"code_25XX"` // 非流動負債合計
}

// 權益
type Equity struct {
	*Parent_company_owner_equity // 歸屬於母公司業主之權益

	Code3XXX string `json:"code_3XXX"` // 權益總額
}

// 歸屬於母公司業主之權益
type Parent_company_owner_equity struct {
	*Capital_stock     // 股本
	*APIC              // 資本公積
	*Retained_earnings // 保留盈餘
	*Other_equity      // 其他權益

	Code31XX string `json:"code_31XX"` // 歸屬於母公司業主之權益合計
}

// 股本
type Capital_stock struct {
	Code3110 string `json:"code_3110"` // 普通股股本
	Code3100 string `json:"code_3100"` // 股本合計

}

// 資本公積 (Additional Paid-In Capital)
type APIC struct {
	Code3200 string `json:"code_3200"` // 資本公積合計
}

// 保留盈餘
type Retained_earnings struct {
	Code3300 string `json:"code_3300"` // 保留盈餘合計

}

// 其他權益
type Other_equity struct {
	Code3400 string `json:"code_3400"` // 其他權益合計
}

// 綜合損益表
type Comprehensive_income_statement struct {
	*Operating_income // 營業收入
	*Operating_cost   // 營業成本

	Code5900 string `json:"code_5900"` //  營業毛利（毛損）
	Code5950 string `json:"code_5950"` //  營業毛利（毛損）淨額

	*Operating_expenses // 營業費用

	Code6900 string `json:"code_6900"` //  營業利益（損失）

	*Non_operating_income_and_expenses // 營業外收入及支出

	Code7900 string `json:"code_7900"` //  繼續營業單位稅前淨利（淨損）

	*Income_tax_expense // 所得稅費用（利益）

	Code8000 string `json:"code_8000"` //  繼續營業單位本期淨利（淨損）
	Code8200 string `json:"code_8200"` //  本期淨利（淨損）

	*Other_comprehensive_income

	Code8500 string `json:"code_8500"` //  本期綜合損益總額

	*Net_profit_attributable_to                 // 淨利（損）歸屬於
	*Total_comprehensive_profit_attributable_to // 綜合損益總額歸屬於
	*Earnings_per_share                         // 每股盈餘(EPS)
}

// 營業收入
type Operating_income struct {
	Code4000 string `json:"code_4000"` //  營業收入合計
}

// 營業成本
type Operating_cost struct {
	Code5000 string `json:"code_5000"` //  營業成本合計
}

// 營業費用
type Operating_expenses struct {
	Code6100 string `json:"code_6100"` //  推銷費用
	Code6200 string `json:"code_6200"` //  管理費用
	Code6300 string `json:"code_6300"` //  研究發展費用
	Code6450 string `json:"code_6450"` //  預期信用減損損失(利益)
	Code6000 string `json:"code_6000"` //  營業費用合計
}

// 營業外收入及支出
type Non_operating_income_and_expenses struct {
	*Interest_income                          // 利息收入
	*Other_income                             // 其他收入
	*Other_profits_and_losses                 // 其他利益及損失
	*Financial_cost                           // 財務成本
	*Associates_and_joint_ventures_net_amount // 採用權益法認列之關聯企業及合資損益之份額

	Code7000 string `json:"code_7000"` //  營業外收入及支出合計
}

// 利息收入
type Interest_income struct {
	Code7101 string `json:"code_7101"` //  銀行存款利息
	Code7100 string `json:"code_7100"` //  利息收入合計

}

// 其他收入
type Other_income struct {
	Code7010 string `json:"code_7010"` //  其他收入合計
}

// 其他利益及損失
type Other_profits_and_losses struct {
	Code7020 string `json:"code_7020"` //  其他利益及損失淨額
}

// 財務成本
type Financial_cost struct {
	Code7050 string `json:"code_7050"` //  財務成本淨額
}

// 採用權益法認列之關聯企業及合資損益之份額
type Associates_and_joint_ventures_net_amount struct {
	Code7060 string `json:"code_7060"` //  採用權益法認列之關聯企業及合資損益之份額淨額
}

// 所得稅費用（利益）
type Income_tax_expense struct {
	Code7950 string `json:"code_7950"` //  所得稅費用（利益）合計
}

// 其他綜合損益(淨額)
type Other_comprehensive_income struct {
	*None_Reclassification_adjustments // 不重分類至損益之項目
	*Reclassification_adjustments      // 後續可能重分類至損益之項目

	Code8300 string `json:"code_8300"` //  其他綜合損益（淨額）
}

// 不重分類至損益之項目
type None_Reclassification_adjustments struct {
	Code8316 string `json:"code_8316"` //  透過其他綜合損益按公允價值衡量之權益工具投資未實現評價損益
	Code8310 string `json:"code_8310"` //  不重分類至損益之項目總額
}

// 後續可能重分類至損益之項目
type Reclassification_adjustments struct {
	Code8361 string `json:"code_8361"` //  國外營運機構財務報表換算之兌換差額
	Code8399 string `json:"code_8399"` //  與可能重分類之項目相關之所得稅
	Code8360 string `json:"code_8360"` //  後續可能重分類至損益之項目總額
}

// 淨利（損）歸屬於
type Net_profit_attributable_to struct {
	Code8610 string `json:"code_8610"` //  母公司業主（淨利／損）
}

// 綜合損益總額歸屬於
type Total_comprehensive_profit_attributable_to struct {
	Code8710 string `json:"code_8710"` //  母公司業主（綜合損益）
}

// 每股盈餘(EPS)
type Earnings_per_share struct {
	Code9750 string `json:"code_9750"` //  基本每股盈餘合計
}

// 現金流量表
type Cash_flow_statement struct {
	*Cash_flow_from_operating_activities // 營業活動之現金流量－間接法

	CodeAAAA string `json:"code_AAAA"` // 營業活動之淨現金流入（流出）

	*Cash_flows_from_investing_activities // 投資活動之現金流量
	*Cash_flow_from_financing_activities  // 籌資活動之現金流量

	CodeDDDD   string `json:"code_DDDD"`   // 匯率變動對現金及約當現金之影響
	CodeEEEE   string `json:"code_EEEE"`   // 本期現金及約當現金增加（減少）數
	CodeE00100 string `json:"code_E00100"` // 期初現金及約當現金餘額
	CodeE00200 string `json:"code_E00200"` // 期末現金及約當現金餘額
	CodeE00210 string `json:"code_E00210"` // 資產負債表帳列之現金及約當現
}

// 營業活動之現金流量－間接法
type Cash_flow_from_operating_activities struct {
	CodeA00010 string `json:"code_A00010"` //  繼續營業單位稅前淨利（淨損）
	CodeA10000 string `json:"code_A10000"` //  本期稅前淨利（淨損）

	*Adjustment_item // 調整項目

	CodeA33000 string `json:"code_A33000"` // 營運產生之現金流入（流出）
	CodeA33100 string `json:"code_A33100"` // 收取之利息
	CodeA33300 string `json:"code_A33300"` // 支付之利息
	CodeA33500 string `json:"code_A33500"` // 退還（支付）之所得稅
}

// 調整項目
type Adjustment_item struct {
	*Income_expense_item                     // 收益費損項目
	*Changes_related_to_operating_activities // 與營業活動相關之資產／負債變動數

	CodeA20000 string `json:"code_A20000"` // 調整項目合計
}

// 收益費損項目
type Income_expense_item struct {
	CodeA20100 string `json:"code_A20100"` // 折舊費用
	CodeA20200 string `json:"code_A20200"` // 攤銷費用
	CodeA20300 string `json:"code_A20300"` // 預期信用減損損失（利益）數／呆帳費用提列（轉列收入）數
	CodeA20900 string `json:"code_A20900"` // 利息費用
	CodeA21200 string `json:"code_A21200"` // 利息收入
	CodeA22300 string `json:"code_A22300"` // 採用權益法認列之關聯企業及合資損失（利益）之份額
	CodeA20010 string `json:"code_A20010"` // 收益費損項目合
}

// 與營業活動相關之資產／負債變動數
type Changes_related_to_operating_activities struct {
	*Net_change_in_assets_related_to_operating_activities      // 與營業活動相關之資產之淨變動
	*Net_change_in_liabilities_related_to_operating_activities // 與營業活動相關之負債之淨變動

	CodeA30000 string `json:"code_A30000"` // 與營業活動相關之資產及負債之淨變動合計
}

// 與營業活動相關之資產之淨變動
type Net_change_in_assets_related_to_operating_activities struct {
	CodeA31150 string `json:"code_A31150"` // 應收帳款（增加）減少
	CodeA31180 string `json:"code_A31180"` // 其他應收款（增加）減少
	CodeA31230 string `json:"code_A31230"` // 預付款項（增加）減少
	CodeA31240 string `json:"code_A31240"` // 其他流動資產（增加）減少
	CodeA31990 string `json:"code_A31990"` // 其他營業資產（增加）減少
	CodeA31000 string `json:"code_A31000"` // 與營業活動相關之資產之淨變動合
}

// 與營業活動相關之負債之淨變動
type Net_change_in_liabilities_related_to_operating_activities struct {
	CodeA32125 string `json:"code_A32125"` // 合約負債增加（減少）
	CodeA32180 string `json:"code_A32180"` // 其他應付款增加（減少）
	CodeA32230 string `json:"code_A32230"` // 其他流動負債增加（減少）
	CodeA32000 string `json:"code_A32000"` // 與營業活動相關之負債之淨變動合計
}

// 投資活動之現金流量
type Cash_flows_from_investing_activities struct {
	CodeB00040 string `json:"code_B00040"` // 取得按攤銷後成本衡量之金融資產
	CodeB01800 string `json:"code_B01800"` // 取得採用權益法之投資
	CodeB04500 string `json:"code_B04500"` // 取得無形資產
	CodeBBBB   string `json:"code_BBBB"`   // 投資活動之淨現金流入（流出）
}

type Cash_flow_from_financing_activities struct {
	CodeC01700 string `json:"code_C01700"` // 償還長期借款
	CodeC04020 string `json:"code_C04020"` // 租賃本金償還
	CodeCCCC   string `json:"code_CCCC"`   // 籌資活動之淨現金流入（流出）
}

// 財務報表_季度統計
//
// @params string 個股編號
//
// @params int 季度 1~4
func Get_twse_financial_statements(stockId string, year int, season string) error {
	url := fmt.Sprintf("https://mops.twse.com.tw/server-java/t164sb01?step=1&CO_ID=%s&SYEAR=%d&SSEASON=%s&REPORT_ID=C", stockId, year, season)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "text/html; charset=big5")

	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	nres, err := util.Big5ToUtf8(body)
	if err != nil {
		return err
	}

	data, err := Unmarshal_twse_financial_statements([]byte(nres))
	if err != nil {
		return err
	}

	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	payload := &stockdao.Company_financial_statements{
		Year:         year,
		Session:      season,
		Company_id:   stockId,
		Data:         string(js),
		Data_version: stockdao.Company_financial_statements_DataVersion,
	}

	fmt.Println(payload)
	// err = sqldb.Mysql().Create(payload).Error
	return err
}

func Unmarshal_twse_financial_statements(payload []byte) (map[string][][]string, error) {
	loginHtml := html.NewTokenizer(bytes.NewBuffer(payload))
	filtersMap := map[string][]*myhtml.FilterObjnewSub{
		"h4": {
			{
				FiltAttrs: []html.Attribute{},
			},
		},
		"table": {
			{
				FiltAttrs: []html.Attribute{},
			},
		},
	}

	myhtml.HtmlLoopFilterLevelSub(loginHtml, filtersMap, nil)

	statusNode := filtersMap["h4"]
	if len(statusNode[0].Res) > 0 {
		return nil, fmt.Errorf("檔案不存在")
	}

	codeMap := map[string][][]string{}
	for _, filters := range filtersMap {
		for _, filter := range filters {
			for _, res := range filter.Res { // table
				// table 資料名稱
				tableNameZHT, _ := nameNode(res.SubRes[0].SubRes[0])
				switch tableNameZHT {
				case "資產負債表", "綜合損益表", "現金流量表":
					dataSlice := [][]string{}

					// 資料標籤補完
					dataKey := []string{"Code", "CodeNameZHT", "CodeNameEN"}
					if len(res.SubRes[1].SubRes) > 2 {
						dateColumn := res.SubRes[1].SubRes[2]
						// for _, data := range dateColumn {
						_, dateEN := nameNode(dateColumn)
						dataKey = append(dataKey, dateEN)
						// }
					}
					dataSlice = append(dataSlice, dataKey)

					// talbe 資料內容
					dataRows := res.SubRes[2:]
					for _, row := range dataRows {
						data := make([]string, 0, len(dataKey))

						code := myhtml.GetContext(row.SubRes[0])
						dataNameZTH, dataNameEN := nameNode(row.SubRes[1])
						if code == "" { // 子集合標籤
							continue
						}

						valueOfDate := []string{}
						for _, value := range row.SubRes[2:] {
							valueOfDate = append(valueOfDate, myhtml.GetContext(value))
						}

						data = append(data, code, dataNameZTH, dataNameEN)
						data = append(data, valueOfDate...)
						dataSlice = append(dataSlice, data)
					}
					codeMap[tableNameZHT] = dataSlice
				case "當期權益變動表":
					dataSlice := [][]string{}
					dataKey := []string{"SubCode", "SubCodeNameZHT", "SubCodeNameEN"}

					if len(res.SubRes[1].SubRes)-1 > 0 {
						for i, count := 1, len(res.SubRes[1].SubRes); i < count; i++ {
							code := strings.TrimSpace(myhtml.GetContext(res.SubRes[1].SubRes[i]))
							dataKey = append(dataKey, code)
						}
					}
					dataSlice = append(dataSlice, dataKey)

					// talbe 資料內容
					dataRows := res.SubRes[3:]
					for _, row := range dataRows {
						data := make([]string, 0, len(dataKey))

						subCode := myhtml.GetContext(row.SubRes[0])
						subCodeNameZHT, subCodeEN := nameNode(row.SubRes[1])
						data = append(data, subCode, subCodeNameZHT, subCodeEN)
						datas := row.SubRes[2:]

						for keyidx := range datas {
							context := myhtml.GetContext(datas[keyidx])
							context = trimNewline(context)
							data = append(data, context)
						}
						dataSlice = append(dataSlice, data)
					}
					codeMap[tableNameZHT] = dataSlice
				case "被投資公司名稱、所在地區…等相關資訊":
					dataSlice := [][]string{}
					dataKey := []string{"投資公司名稱", "被投資公司名稱", "地區別代號", "所在地區", "主要營業項目", "原始投資金額_本期期末", "原始投資金額_去年年底", "股數", "比率", "帳面金額", "被投資公司本期損益", "本期認列之投資損益", "備註"}
					dataSlice = append(dataSlice, dataKey)

					dataRows := res.SubRes[3:]
					for _, row := range dataRows {
						data := make([]string, 0, len(dataKey))
						for _, subRes := range row.SubRes {
							for _, ssubRes := range subRes.SubRes {
								if len(ssubRes.SubRes) > 1 {
									for _, sssubRes := range ssubRes.SubRes {
										context := myhtml.GetContext(sssubRes)
										context = trimNewline(context)
										data = append(data, context)
									}
								} else {
									context := myhtml.GetContext(ssubRes)
									context = trimNewline(context)
									data = append(data, context)
								}
							}
						}
						dataSlice = append(dataSlice, data)
					}
					codeMap[tableNameZHT] = dataSlice

				case "轉投資大陸地區之事業相關資訊":
					dataSlice := [][]string{}
					dataKey := []string{"大陸被投資公司名稱", "主要營業項目", "實收資本額", "投資方式", "本期期初自台灣匯出累積投資金額", "本期匯出或收回投資金額_匯出", "本期匯出或收回投資金額_收回", "本期期末自台灣匯出累積投資金額", "被投資公司本期損益", "本公司直接或間接投資之持股比例", "本期認列投資損益", "期末投資帳面價值", "截至本期止已匯回台灣之投資收益", "備註"}
					dataSlice = append(dataSlice, dataKey)

					dataRows := res.SubRes[3:]
					for _, row := range dataRows {
						data := make([]string, 0, len(dataKey))
						for _, subRes := range row.SubRes {
							for _, ssubRes := range subRes.SubRes {
								if len(ssubRes.SubRes) > 1 {
									for _, sssubRes := range ssubRes.SubRes {
										context := myhtml.GetContext(sssubRes)
										context = trimNewline(context)
										data = append(data, context)
									}
								} else {
									context := myhtml.GetContext(ssubRes)
									context = trimNewline(context)

									if context == "備註" {
										goto skip
									}

									data = append(data, context)
								}
							}
						}
						dataSlice = append(dataSlice, data)
					skip:
					}
					codeMap[tableNameZHT] = dataSlice
				}
			}
		}
	}

	return codeMap, nil
}

func nameNode(node *myhtml.TokenObjSub) (string, string) {
	switch len(node.SubRes) {
	case 2:
		dataNameZTH := trimNewline(myhtml.GetContext(node.SubRes[0]))
		dataNameEN := trimNewline(myhtml.GetContext(node.SubRes[1]))
		return dataNameZTH, dataNameEN
	case 1:
		dataNameZTH := trimNewline(myhtml.GetContext(node.SubRes[0]))
		return dataNameZTH, ""
	}
	return "", ""
}

func trimNewline(source string) string {
	newStr := strings.Split(source, "\n")
	for idx := range newStr {
		newStr[idx] = strings.TrimSpace(newStr[idx])
	}

	return strings.Join(newStr, " ")
}

// 用於判斷是否反向資料
func IsInverse(data string) bool {
	if len(data) == 0 {
		return false
	}

	return data[0] == '('
}
