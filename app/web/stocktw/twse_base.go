package stocktw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/YWJSonic/ycore/module/myhtml"
	"github.com/YWJSonic/ycore/util"
	"golang.org/x/net/html"
)

// 公司基本資料

// https://mops.twse.com.tw/mops/web/ajax_t05st03

// params
// encodeURIComponent:1
// step:1
// firstin:1
// off:1
// queryName:co_id
// inpuType:co_id
// TYPEK:all
// co_id:6531

// 公開資訊觀測站_公司基本資料
type TwseBase struct {
	Id                        string `json:"id"`                           // 股票代號
	IndustryType              string `json:"industry_type"`                // 產業類型
	Name                      string `json:"name"`                         // 公司名稱
	TelPhone                  string `json:"telphone"`                     // 總機
	Address                   string `json:"address"`                      // 地址
	Chairman                  string `json:"chairman"`                     // 董事長
	GeneralManager            string `json:"general_manager"`              // 總經理
	Spokesman                 string `json:"spokesman"`                    // 發言人
	SpokespersonTitle         string `json:"spokesperson_title"`           // 發言人職稱
	SpokespersonPhone         string `json:"spokesperson_phone"`           // 發言人電話
	ActingSpokesperson        string `json:"acting_spokesperson"`          // 代理發言人
	MainBusiness              string `json:"main_business"`                // 主要經營業務
	CompanyEstablishmentDate  string `json:"company_establishment_date"`   // 公司成立日期
	UniformNumbers            string `json:"uniform_numbers"`              // 營利事業統一編號
	Paid_in_capital           int    `json:"paid_in_capital"`              // 實收資本額
	ListingDate               string `json:"listing_date"`                 // 上市日期
	OTCDate                   string `json:"otc_date"`                     // 上櫃日期
	OpeningDate               string `json:"opening_date"`                 // 興櫃日期
	PublicReleaseDate         string `json:"public_release_date"`          // 公開發行日期
	PreStockPrice             int    `json:"pre_stock_price"`              // 普通股每股面額
	ReleaseStockAmount        int    `json:"release_stock_amount"`         // 已發行普通股數或TDR原股發行股數
	ReleasePrivateStockAmount int    `json:"release_private_stock_amount"` // 已發行普通股數或TDR原股發行股數(私募股票)
	SpecialStockAmount        int    `json:"special_stock_amount"`         // 特別股
	StockTransferAgency       string `json:"stock_transfer_agency"`        // 股票過戶機構
	VisaAccountingFirm        string `json:"visa_accounting_firm"`         // 簽證會計師事務所
	CompanyWebsite            string `json:"company_website"`              // 公司網址
}

func Get_twse_base(stockId string) error {

	req, err := http.NewRequest(http.MethodPost, "https://mops.twse.com.tw/mops/web/ajax_t05st03", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err := req.ParseForm(); err != nil {
		return err
	}
	req.PostForm.Add("encodeURIComponent", "1")
	req.PostForm.Add("step", "1")
	req.PostForm.Add("firstin", "1")
	req.PostForm.Add("off", "1")
	req.PostForm.Add("queryName", "co_id")
	req.PostForm.Add("inpuType", "co_id")
	req.PostForm.Add("TYPEK", "all")
	req.PostForm.Add("co_id", stockId)
	httpRes, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(httpRes.Body)
	if err != nil {
		return err
	}

	fmt.Println(body)
	Unmarshal_twse_base(body)
	return nil
}

func Unmarshal_twse_base(xml_payload []byte) (*TwseBase, error) {
	loginHtml := html.NewTokenizer(bytes.NewBuffer(xml_payload))
	filtersMap := map[string][]*myhtml.FilterObjnewSub{
		"tr": {
			{
				FiltAttrs: []html.Attribute{},
			},
		},
	}

	pageFix := func(nodeDepth int, next, current, previous **myhtml.TokenObjSub) (int, bool) {

		if (*next).Res.Data == "tr" && (*current).Res.Data == "tr" {
			if previous != nil {
				(*current), (*previous) = (*current).Previous, (*current).Previous.Previous
			}
			nodeDepth--
			return nodeDepth, true
		} else {
			return nodeDepth, false
		}
	}

	obj := &TwseBase{}
	myhtml.HtmlLoopTool{SelfClosingTagToken: []string{"meta", "link", "input", "img", "br"}}.HtmlLoopFilterLevelSub(loginHtml, filtersMap, pageFix)

	for _, filters := range filtersMap["tr"] {
		for _, res := range filters.Res {
			count := len(res.SubRes)
			switch count % 2 {
			case 0:
				for i := 0; i < count; i += 2 { // td
					key := myhtml.GetContext(res.SubRes[i])
					value := myhtml.GetContext(res.SubRes[i+1])
					switch key {
					case "股票代號":
						obj.Id = value
					case "產業類別":
						obj.IndustryType = value
					case "外國企業註冊地國":
					case "公司名稱":
						obj.Name = value
					case "總機":
						obj.TelPhone = value
					case "地址":
						obj.Address = value
					case "董事長":
						obj.Chairman = value
					case "總經理":
						obj.GeneralManager = value
					case "發言人":
						obj.Spokesman = value
					case "發言人職稱":
						obj.SpokespersonTitle = value
					case "發言人電話":
						obj.SpokespersonPhone = value
					case "代理發言人":
						obj.ActingSpokesperson = value
					case "主要經營業務":
						obj.MainBusiness = value
					case "公司成立日期":
						obj.CompanyEstablishmentDate = value
					case "營利事業統一編號":
						obj.UniformNumbers = value
					case "實收資本額":
						number, err := strconv.Atoi(strings.Join(util.GetNumberInString(value), ""))
						if err != nil {
							return nil, err
						}
						obj.Paid_in_capital = number

					case "上市日期":
						obj.ListingDate = value
					case "上櫃日期":
						obj.OTCDate = value
					case "興櫃日期":
						obj.OpeningDate = value
					case "公開發行日期":
						obj.PublicReleaseDate = value
					case "普通股每股面額":
						numberStr := util.GetNumberInString(value)
						if len(numberStr) <= 0 {
							return nil, fmt.Errorf("Error 普通股每股面額 value:%v", value)
						}
						number, err := strconv.Atoi(numberStr[0])
						if err != nil {
							return nil, err
						}
						obj.PreStockPrice = number

					case "已發行普通股數或TDR原股發行股數":
						amount := strings.Split(value, "\n")
						number, err := strconv.Atoi(strings.Join(util.GetNumberInString(amount[0]), ""))
						if err != nil {
							return nil, err
						}
						obj.ReleaseStockAmount = number

						number, err = strconv.Atoi(strings.Join(util.GetNumberInString(amount[1]), ""))
						if err != nil {
							return nil, err
						}
						obj.ReleasePrivateStockAmount = number

					case "特別股":
						number, err := strconv.Atoi(strings.Join(util.GetNumberInString(value), ""))
						if err != nil {
							return nil, err
						}
						obj.SpecialStockAmount = number

					case "普通股盈餘分派或虧損撥補頻率":
					case "普通股年度 (含第4季或後半年度)現金股息及紅利決議層級":
					case "股票過戶機構":
						obj.StockTransferAgency = value
					case "電話":
					case "過戶地址":
					case "簽證會計師事務所":
						obj.VisaAccountingFirm = value
					case "簽證會計師1":
					case "簽證會計師2":
					case "備註":
					case "本公司":
					case "特別股發行":
					case "有":
					case "英文簡稱":
					case "英文全名":
					case "英文通訊地址(街巷弄號)":
					case "英文通訊地址(縣市國別)":
					case "傳真機號碼":
					case "電子郵件信箱":
					case "公司網址":
						obj.CompanyWebsite = value
					case "投資人關係聯絡人":
					case "投資人關係聯絡人職稱":
					case "投資人關係聯絡電話":
					case "投資人關係電子郵件":
					case "公司網站內利害關係人專區網址":
					case "變更前名稱":
					case "變更前簡稱":
					case "公司名稱變更核准日期":
					case "編製財務報告類型":
					}
				}
			case 1:
				continue
			default:
				for i := 0; i < count; i += 2 { // td
					fmt.Println(myhtml.GetContext(res.SubRes[i]))
					fmt.Println(myhtml.GetContext(res.SubRes[i+1]))
				}
			}
		}
	}

	js, _ := json.Marshal(obj)
	fmt.Println(string(js))

	return obj, nil
}
