package web591

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/YWJSonic/ycore/driver/protocol/rss/rssfeed"
	"github.com/YWJSonic/ycore/module/mydb"
	"github.com/YWJSonic/ycore/module/myhtml"
	"github.com/YWJSonic/ycore/module/mylog"
	"github.com/YWJSonic/ycore/util"

	"github.com/gorilla/feeds"
	"golang.org/x/net/html"
)

// var aMineReqCount = 60
var reqWebCount = 0

func GetData() {

	dbManager, err := mydb.NewArangoDB("http://10.140.0.2:8529", "", "", "WebData")
	if err != nil {
		return
	}

	startTime := time.Now()
	mylog.Infof("[start] Time: %v", startTime)
	myhtml.SetDefaultClient()

	authData := login("pc")

	urlStr := webpage
	payload := getList(authData.CsrfToken, urlStr)

	// lastReqTime := time.Now()
	homeViewMap := map[int64]struct{}{}
	homeResList := []string{}
	rssData := []*feeds.Item{}
	count := 0
	idx := 0
	customMax := 60 // 指定查尋數量
	var roomIds []int64
	for {

		// 列出篩選資料 簡易資料
		roomIds = []int64{}
		for _, roomInfo := range payload.Data.Data {

			if _, ok := homeViewMap[roomInfo.PostID]; ok {
				continue
			}

			homeViewMap[roomInfo.PostID] = struct{}{}
			idx = strings.Index(roomInfo.RoomStr, "開放式")
			if idx > 0 {
				continue
			}

			idx = strings.Index(roomInfo.RoomStr, "0廳")
			if idx > 0 {
				continue
			}

			roomIds = append(roomIds, roomInfo.PostID)
		}

		// 查詢詳細資料
		for _, roomId := range roomIds {

			// 防止過度攻擊
			// lastTime := time.Since(lastReqTime) / time.Second
			// lastReqTime = time.Now()
			// if lastTime < 2 {
			// 	time.Sleep(time.Second * time.Duration(rand.Int31n(3)+1))
			// }

			urlStr = util.Sprintf(objPage, roomId)
			detailInfo := getDetail(authData, urlStr)
			idx = strings.Index(detailInfo.Data.FavData.Layout, "衛")
			if idx < 0 {
				continue
			}
			bscount, _ := strconv.Atoi(detailInfo.Data.FavData.Layout[idx-1 : idx])
			if bscount < 2 {
				continue
			}

			// js, _ := util.Marshal(detailInfo.Data)
			// fmt.Println(string(js))

			cost := []string{}
			{
				if detailInfo.Data.GetDeposit() != "" {
					cost = append(cost, fmt.Sprintf("押金：[%s]", detailInfo.Data.GetDeposit()))
				}
				if detailInfo.Data.GetManageprice() != "" {
					cost = append(cost, fmt.Sprintf("管理費：[%s]", detailInfo.Data.GetManageprice()))
				}
			}
			layout := []string{}
			{
				if detailInfo.Data.GetLayout() != "" {
					layout = append(layout, detailInfo.Data.GetLayout())
				}
				if detailInfo.Data.GetArea() != "" {
					layout = append(layout, detailInfo.Data.GetArea())
				}
				if detailInfo.Data.GetFloor() != "" {
					layout = append(layout, detailInfo.Data.GetFloor())
				}
				if detailInfo.Data.GetShape() != "" {
					layout = append(layout, detailInfo.Data.GetShape())
				}
				if detailInfo.Data.Getkind() != "" {
					layout = append(layout, detailInfo.Data.Getkind())
				}
			}
			region := []string{}
			{
				if detailInfo.Data.GetLayout() != "" {
					region = append(region, detailInfo.Data.GetRegion())
				}
				if detailInfo.Data.GetArea() != "" {
					region = append(region, detailInfo.Data.GetSection())
				}
			}

			rssContextReplas := strings.NewReplacer(
				"$租金", detailInfo.Data.GetPrice(),
				"$額外成本", strings.Join(cost, ` `),
				"$格局", strings.Join(layout, ` `),
				"$區域", strings.Join(region, ` `),
			)
			description := rssContextReplas.Replace(rssTemplate)

			pageurl := util.Sprintf(targetPage, roomId)
			homeResList = append(homeResList, pageurl)
			rssData = append(rssData, &feeds.Item{
				Title: detailInfo.Data.Title,
				Link: &feeds.Link{
					Href: pageurl,
				},
				Description: description,
			})
			mylog.Infof("[newDetail] Time Spand:%v ReqCount: %v", time.Since(startTime), reqWebCount)
		}

		// 防止過度攻擊
		// lastTime := time.Since(lastReqTime) / time.Second
		// lastReqTime = time.Now()
		// if lastTime < 2 {
		// 	time.Sleep(time.Second * time.Duration(rand.Int31n(3)+1))
		// }

		// 下一頁
		count += 30

		// 沒有下一頁 or 到了指定查詢數量
		max, _ := strconv.Atoi(strings.ReplaceAll(payload.Records, ",", ""))
		mylog.Infof("[nextPage]Now Count: %v Max Count: %v", count, max)
		if count >= max || count >= customMax {
			break
		}

		urlStr = util.Sprintf(webpagelast, count, payload.Records)
		payload = getList(authData.CsrfToken, urlStr)
		mylog.Infof("[getList] Time Spand: %v ReqCount: %v", time.Since(startTime), reqWebCount)
	}

	rssfeed.NewFeed(rssData)

	time.Now().Format(time.RFC1123)
	_ = dbManager.Create(
		context.TODO(),
		"FilterHomeData",
		DBStruct{
			Time:     util.ServerTimeNow().Format("15:04:05 -07:00"),
			Date:     util.ServerTimeNow().Format("2006-01-02"),
			RoomList: homeResList,
		},
	)
}

// func GetCostStr(data CostData) string {
// 	str := ""
// 	for _, info := range data.Data {
// 		str += fmt.Sprintf("%s:%s ", info.Name, info.Value)
// 	}
// 	if str == "" {
// 		str = "無"
// 	}
// 	return str
// }

// func GetInfoStr(datas []Info) string {
// 	str := ""
// 	for _, info := range datas {
// 		str += fmt.Sprintf("%s ", info.Value)
// 	}
// 	if str == "" {
// 		str = "無"
// 	}
// 	return str
// }

// func GetBreadcrumb(datas []Breadcrumb) string {
// 	str := ""
// 	for _, info := range datas {
// 		str += fmt.Sprintf("%s ", info.Name)
// 	}
// 	if str == "" {
// 		str = "無"
// 	}
// 	return str
// }

func getList(csrfToken string, url string) *HomeList {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-CSRF-TOKEN", csrfToken)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	reqWebCount++

	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	payload := &HomeList{}
	err = json.Unmarshal(sitemap, payload)
	if err != nil {
		mylog.Errorf("[getList] Unmarshal err: %v", err)
	}

	return payload
}

func getDetail(authData LoginData, url string) *HomeDetail {
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("X-CSRF-TOKEN", authData.CsrfToken)
	req.Header.Set("token", authData.Session)
	req.Header.Set("device", authData.Device)
	req.Header.Set("deviceid", authData.Deviceid)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	reqWebCount++

	sitemap, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}

	payload := &HomeDetail{}
	err = json.Unmarshal(sitemap, payload)
	if err != nil {
		mylog.Errorf("[getDetail] Unmarshal err: %v", err)
	}
	return payload
}

func login(device string) LoginData {
	data := LoginData{
		Device: device,
	}
	req, _ := http.NewRequest("GET", loginPage, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	for _, cookie := range res.Cookies() {
		if cookie.Name == "PHPSESSID" {
			data.Session = cookie.Value
		}
		if cookie.Name == "T591_TOKEN" {
			data.Deviceid = cookie.Value
		}
	}

	loginHtml := html.NewTokenizer(res.Body)

	filters := map[string][]*myhtml.FilterObj{
		"meta": {
			{
				FiltAttrs: []html.Attribute{
					{
						Key: "name",
						Val: "csrf-token",
					},
				},
			},
		},
	}
	myhtml.HtmlLoopFilterOne(loginHtml, filters)

	for _, htmlToken := range filters["meta"][0].Res {
		for _, attr := range htmlToken.Attr {
			if attr.Key == "content" {
				data.CsrfToken = attr.Val
			}
		}
	}

	return data
}
