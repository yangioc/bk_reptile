package efish

import (
	"bk_reptile/tmptool"
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/YWJSonic/ycore/module/myhtml"
	"github.com/YWJSonic/ycore/types"
	"golang.org/x/net/html"
)

type Handle struct {
	TitleRow []string
	DataRow  [][]string
}

func GetHistory(fishId string, date_start, date_end time.Time) (*Handle, error) {
	yearl := strconv.Itoa(tmptool.YearConver(date_start.Year(), "tw"))
	monthl := strconv.Itoa(int(date_start.Month()))
	dayl := strconv.Itoa(date_start.Day())

	yeart := strconv.Itoa(tmptool.YearConver(date_end.Year(), "tw"))
	montht := strconv.Itoa(int(date_end.Month()))
	dayt := strconv.Itoa(date_end.Day())

	dateStrl := fmt.Sprintf("%s.%s.%s", yearl, monthl, dayl)
	dateStrt := fmt.Sprintf("%s.%s.%s", yeart, montht, dayt)

	data := url.Values{}
	data.Set("dateStrt", dateStrt)
	data.Set("dateStrl", dateStrl)
	data.Set("keyword", "")
	data.Set("marketStr", "")
	data.Set("marketType", "2") // 市場類型
	data.Set("calendarType", "tw")
	data.Set("yearl", yearl)
	data.Set("monthl", monthl)
	data.Set("dayl", dayl)
	data.Set("yeart", yeart)
	data.Set("montht", montht)
	data.Set("dayt", dayt)
	data.Set("pid", fishId) // 魚類編號
	data.Set("mt", "2")     // 2:全部 1:消費地 6:消費地(不含埔里) 0:生產地

	var err error
	var req *http.Request
	WebPage := "https://efish.fa.gov.tw/efish/statistics/daymultidaysinglefishmultimarket.htm"
	if req, err = http.NewRequest(http.MethodPost, WebPage, strings.NewReader(data.Encode())); err != nil {
		return nil, err
	}

	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	var hemlRes *http.Response
	if hemlRes, err = http.DefaultClient.Do(req); err != nil {
		return nil, err
	}
	defer hemlRes.Body.Close()

	loginHtml := html.NewTokenizer(hemlRes.Body)

	// 撈取網頁資料
	filters := map[types.TokenTypeName][]*myhtml.FilterObjnewSub{} // 網頁標籤資料
	filters["table"] = []*myhtml.FilterObjnewSub{
		{
			FiltAttrs: []html.Attribute{
				{
					Key: "border",
					Val: "1",
				},
			},
		},
	}

	myhtml.HtmlLoopFilterLevelSub(loginHtml, filters, nil)

	fmt.Println(filters)

	if len(filters["table"][0].Res) == 0 {
		return nil, errors.New("no data")
	}

	titles := []string{}
	datas := [][]string{}
	titleRow := filters["table"][0].Res[0].SubRes[0]
	dataRow := filters["table"][0].Res[0].SubRes[1]

	for i, count := 0, len(titleRow.SubRes[0].SubRes); i < count; i++ {
		titles = append(titles, myhtml.GetContext(titleRow.SubRes[0].SubRes[i]))
	}

	for i, count := 0, len(dataRow.SubRes); i < count; i += 2 { // tr
		SubRes := dataRow.SubRes[i]

		var dateStr string
		var rowStr []string
		for idx, SSubRes := range SubRes.SubRes { // TD
			if idx == 0 {
				dateStr = myhtml.GetContext(SSubRes)
				rowStr = append(rowStr, dateStr)
				continue
			}

			rowStr = append(rowStr, myhtml.GetContext(SSubRes))
		}
		datas = append(datas, rowStr)

		rowStr = []string{dateStr}
		SubRes = dataRow.SubRes[i+1]
		for _, SSubRes := range SubRes.SubRes { // TD
			rowStr = append(rowStr, myhtml.GetContext(SSubRes))
		}
		datas = append(datas, rowStr)
	}
	res := &Handle{
		TitleRow: titles,
		DataRow:  datas,
	}
	return res, nil
}

func GetToday(marketId string) (*Handle, error) {
	// nowTime := time.Now()
	// year := strconv.Itoa(tmptool.YearConver(nowTime.Year(), "tw"))
	// month := strconv.Itoa(int(nowTime.Month()))
	// day := strconv.Itoa(nowTime.Day())

	// // dateStrl := fmt.Sprintf("%s.%s.%s", year, month, day)
	// dateStrl := fmt.Sprintf("%s.%s.%s", "113", "1", "31")

	// data := url.Values{}
	// data.Set("dateStr", dateStrl)
	// data.Set("calendarType", "tw")
	// data.Set("year", year)
	// data.Set("month", month)
	// data.Set("day", day)
	// data.Set("mid", marketId)
	// data.Set("numbers", "999") // 魚類編號 999=全種類
	// data.Set("orderby", "i")

	// var err error
	// var req *http.Request
	// WebPage := "https://efish.fa.gov.tw/efish/statistics/daysinglemarketmultifish.htm"
	// if req, err = http.NewRequest(http.MethodPost, WebPage, strings.NewReader(data.Encode())); err != nil {
	// 	return nil, err
	// }

	// req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36")
	// req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// fmt.Println(req.PostForm)

	// var res *http.Response
	// if res, err = http.DefaultClient.Do(req); err != nil {
	// 	return nil, err
	// }
	// defer res.Body.Close()

	// body, _ := io.ReadAll(res.Body)

	// os.WriteFile("./efishToday.log", body, 0777)

	// return nil, nil

	body, _ := os.ReadFile("./efishToday.log")
	htmlRes := bytes.NewBuffer(body)

	loginHtml := html.NewTokenizer(htmlRes)

	// 撈取網頁資料
	filters := map[types.TokenTypeName][]*myhtml.FilterObjnewSub{} // 網頁標籤資料
	filters["table"] = []*myhtml.FilterObjnewSub{
		{
			FiltAttrs: []html.Attribute{
				{
					Key: "id",
					Val: "ltable",
				},
			},
		},
	}

	HtmlLoopFilterLevelSub_CoustomToken(loginHtml, filters, nil, []string{"meta", "link", "input", "img", "br"})

	if len(filters["table"][0].Res) == 0 {
		return nil, errors.New("no data")
	}

	titles := []string{}
	datas := [][]string{}
	titleRow := filters["table"][0].Res[0].SubRes[0]
	dataRow := filters["table"][0].Res[0].SubRes[1]

	for i, count := 0, len(titleRow.SubRes[0].SubRes); i < count; i++ {
		titles = append(titles, myhtml.GetContext(titleRow.SubRes[0].SubRes[i]))
	}

	for i, count := 0, len(dataRow.SubRes); i < count; i++ { // tr
		SubRes := dataRow.SubRes[i]

		var rowStr []string
		for _, SSubRes := range SubRes.SubRes { // TD
			tmpval := strings.TrimSpace(myhtml.GetContext(SSubRes))
			rowStr = append(rowStr, tmpval)
		}
		datas = append(datas, rowStr)

	}
	res := &Handle{
		TitleRow: titles,
		DataRow:  datas,
	}
	return res, nil
}

func HtmlLoopFilterLevelSub_CoustomToken(tokenizer *html.Tokenizer, filterMap map[types.TokenTypeName][]*myhtml.FilterObjnewSub, pageFix func(nodeDepth int, next, current, previous **myhtml.TokenObjSub) (int, bool), selfClosingTagToken []string) {
	tool := myhtml.HtmlLoopTool{
		SelfClosingTagToken: selfClosingTagToken,
	}
	tool.HtmlLoopFilterLevelSub(tokenizer, filterMap, pageFix)
}
