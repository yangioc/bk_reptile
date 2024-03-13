package app

import (
	"bk_reptile/app/web/coolpc"
	"bk_reptile/app/web/efish"
	"bk_reptile/app/web/stocktw"
	"bk_reptile/config"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/yangioc/bk_pack/log"
	"github.com/yangioc/bk_pack/proto/dtomsg"
	"github.com/yangioc/bk_pack/util"
)

func (self *Handle) handleReqMessage(msg *dtomsg.Dto_Msg) ([]byte, error) {
	return nil, errors.New("req message request not found.")
}

func (self *Handle) handleNoticMessage(msg *dtomsg.Dto_Msg) error {
	return errors.New("notic message request not found.")
}

func (self *Handle) GetCoolpc() {
	uuid := util.GenStrUUID(config.EnvInfo.NodeNum)
	datas, err := coolpc.GetWeb()
	if err != nil {
		panic(err)
	}

	for _, data := range datas {
		payload, err := util.Marshal(data)
		if err != nil {
			panic(err)
		}

		if err = self.dba.CreateCoolpcData(uuid, payload); err != nil {
			panic(err)
		}
	}
}

func (self *Handle) GetEfish(date time.Time) {
	uuid := util.GenStrUUID(config.EnvInfo.NodeNum)

	locations := []string{"F109", "F200", "F241", "F261", "F270", "F300", "F330", "F360", "F400", "F500", "F513", "F545", "F600", "F630", "F708", "F709", "F722", "F730", "F800", "F820", "F826", "F880", "F916", "F936", "F950"}

	locationGroup := map[string]interface{}{}
	for _, location := range locations {
		data, err := efish.GetDayFishByMarket(location, date)
		if err != nil {
			panic(err)
		}

		locationGroup[location] = data
	}

	dataGroup := map[string]interface{}{
		"_key":          date.Format("20060102"),
		"LocationGroup": locationGroup,
	}

	payload, _ := util.Marshal(dataGroup)
	if err := self.dba.CreateEfish(uuid, payload); err != nil {
		panic(err)
	}
}

func (self *Handle) GetOldEfish(date_start, date_end time.Time) {
	if date_start.After(date_end) {
		log.Error("GetOldEfish time error")
		return
	}

	currentTime := date_start
	for ; currentTime.Before(date_end); currentTime = currentTime.Add(time.Hour * 24) {
		self.GetEfish(currentTime)
	}
}

func (self *Handle) GetStock() {
	thisDay := util.ServerTimeNow()

	switch thisDay.Weekday() {
	case time.Saturday, time.Sunday:
		return
	}

	dateStr := thisDay.Format("2006-01-02")
	stockAnalysis, err := stocktw.GetStockAnalysis(thisDay)
	if err != nil {
		panic(err)
	}

	stockAnalysisPayload := map[string]interface{}{
		"_key":  dateStr,
		"datas": stockAnalysis,
	}

	stockIndex, stockMarket, stockClosePrice, err := stocktw.GetStockIndex(thisDay)
	if err != nil {
		panic(err)
	}
	stockIndexPayload := map[string]interface{}{
		"_key":  dateStr,
		"datas": stockIndex,
	}
	stockMarketPayload := map[string]interface{}{
		"_key":  dateStr,
		"datas": stockMarket,
	}
	stockClosePricePayload := map[string]interface{}{
		"_key":  dateStr,
		"datas": stockClosePrice,
	}

	stockThreefoundationTotal, stockThreefoundationStockDay, err := stocktw.GetThreefoundation(thisDay)
	if err != nil {
		panic(err)
	}
	stockThreefoundationTotalPayload := map[string]interface{}{
		"_key":  dateStr,
		"datas": stockThreefoundationTotal,
	}
	stockThreefoundationStockDayPayload := map[string]interface{}{
		"_key":  dateStr,
		"datas": stockThreefoundationStockDay,
	}

	datas := map[string]interface{}{
		"StockAnalysis":           stockAnalysisPayload,
		"Index":                   stockIndexPayload,
		"Market":                  stockMarketPayload,
		"ClosePrice":              stockClosePricePayload,
		"ThreefoundationTotal":    stockThreefoundationTotalPayload,
		"ThreefoundationStockDay": stockThreefoundationStockDayPayload,
	}

	for key, data := range datas {
		uuid := util.GenStrUUID(config.EnvInfo.NodeNum)
		payload, err := json.Marshal(data)
		if err != nil {
			panic(err)
		}

		log.Infof("%s data size :%v", key, len(payload))

		switch key {
		case "StockAnalysis":
			if err := self.dba.CreateStockAnalysis(uuid, payload); err != nil {
				log.Errorf("CreateStockAnalysis error: %v", err)
				continue
			}
		case "Index":
			if err := self.dba.CreateStockIndex(uuid, payload); err != nil {
				log.Errorf("CreateStockIndex error: %v", err)
				continue
			}
		case "Market":
			if err := self.dba.CreateStockMarket(uuid, payload); err != nil {
				log.Errorf("CreateStockMarket error: %v", err)
				continue
			}
		case "ClosePrice":
			if err := self.dba.CreateStockClosePrice(uuid, payload); err != nil {
				log.Errorf("CreateStockClosePrice error: %v", err)
				continue
			}
		case "ThreefoundationTotal":
			if err := self.dba.CreateStockThreefoundationTotal(uuid, payload); err != nil {
				log.Errorf("CreateStockThreefoundationTotal error: %v", err)
				continue
			}
		case "ThreefoundationStockDay":
			if err := self.dba.CreateThreefoundationStockDay(uuid, payload); err != nil {
				log.Errorf("CreateThreefoundationStockDay error: %v", err)
				continue
			}
		}
	}
}

func (self *Handle) GetStockHistory(date_start, date_end time.Time) {
	dates := []time.Time{}
	count := util.CountTransDate(date_end, date_start)
	for i := 0; i < count; i++ {
		dates = append(dates, date_start.AddDate(0, 0, i))
	}

	for _, newDate := range dates {
		switch newDate.Weekday() {
		case time.Saturday, time.Sunday:
			continue
		}

		resStockAnalysis, err := stocktw.GetStockAnalysis(newDate)
		if err != nil {
			panic(err)
		}
		// resIndex, resMarket, resClosePrice, err := stocktw.GetStockIndex(newDate)
		// if err != nil {
		// 	panic(err)
		// }
		// resThreefoundationTotal, resThreefoundationStockDay, err := stocktw.GetThreefoundation(newDate)
		// if err != nil {
		// 	panic(err)
		// }

		datas := map[string]interface{}{
			"StockAnalysis": resStockAnalysis,
			// "Index":                   resIndex,
			// "Market":                  resMarket,
			// "ClosePrice":              resClosePrice,
			// "ThreefoundationTotal":    resThreefoundationTotal,
			// "ThreefoundationStockDay": resThreefoundationStockDay,
		}

		for key, data := range datas {
			uuid := util.GenStrUUID(config.EnvInfo.NodeNum)
			payload, err := json.Marshal(data)
			if err != nil {
				panic(err)
			}

			switch key {
			case "StockAnalysis":
				self.dba.CreateStockAnalysis(uuid, payload)
			case "Index":
				self.dba.CreateStockIndex(uuid, payload)
			case "Market":
				self.dba.CreateStockMarket(uuid, payload)
			case "ClosePrice":
				self.dba.CreateStockClosePrice(uuid, payload)
			case "ThreefoundationTotal":
				self.dba.CreateStockThreefoundationTotal(uuid, payload)
			case "ThreefoundationStockDay":
				self.dba.CreateThreefoundationStockDay(uuid, payload)
			}
		}

		fmt.Println("All Done")
		time.Sleep(time.Second * 5 * 1)
	}
}
