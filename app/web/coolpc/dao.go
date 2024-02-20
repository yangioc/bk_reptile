package coolpc

type ItemInfo struct {
	Key           string `json:"_key,omitempty"`
	TypeId        string `json:"typeId,omitempty"`
	TypeName      string `json:"typeName,omitempty"`
	Price         string `json:"price,omitempty"`         // 價格
	Content       string `json:"content,omitempty"`       // 標示內容
	PriceTag      string `json:"priceTag,omitempty"`      // 價錢標籤(降價標示)
	OriginContent string `json:"originContent,omitempty"` // 原始資料
	UpdateTime    int64  `json:"updateTime,omitempty"`    // 更新時間
	Date          string `json:"date,omitempty"`          // 日期
}
