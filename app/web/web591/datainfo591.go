package web591

type DBStruct struct {
	Time     string   `json:"time"`
	Date     string   `json:"date"`
	RoomList []string `json:"roomList"`
}

type LoginData struct {
	CsrfToken string
	Session   string
	Deviceid  string
	Device    string
}

type HomeList struct {
	BluekaiData      BluekaiData   `json:"bluekai_data"`
	Data             Data          `json:"data"`
	DealRecom        []interface{} `json:"deal_recom"`
	IsRecom          int64         `json:"is_recom"`
	OnlineSocialUser int64         `json:"online_social_user"`
	Recommend        []interface{} `json:"recommend"`
	Records          string        `json:"records"` // 總筆數
	SEO              SEO           `json:"seo"`
	Status           int64         `json:"status"`
}

type BluekaiData struct {
	Kind             string `json:"kind"`
	MrtCity          string `json:"mrt_city"`
	MrtLine          string `json:"mrt_line"`
	Page             string `json:"page"`
	RegionID         string `json:"region_id"`
	RentalPrice      int64  `json:"rental_price"`
	Room             string `json:"room"`
	SalePrice        int64  `json:"sale_price"`
	SectionID        string `json:"section_id"`
	Shape            string `json:"shape"`
	Tag              int64  `json:"tag"`
	Type             string `json:"type"`
	UnitPricePerPing string `json:"unit_price_per_ping"`
}

type Data struct {
	Biddings []interface{} `json:"biddings"`
	Data     []Datum       `json:"data"`
	Page     string        `json:"page"`
	TopData  []TopDatum    `json:"topData"`
}

type Datum struct {
	Area             string      `json:"area"`
	CasesID          interface{} `json:"cases_id"`
	Community        string      `json:"community"`
	Contact          string      `json:"contact"`
	DiscountPriceStr string      `json:"discount_price_str"`
	FloorStr         string      `json:"floor_str"`
	Hurry            int64       `json:"hurry"`
	IsCombine        int64       `json:"is_combine"`
	IsSocail         int64       `json:"is_socail"`
	IsVideo          int64       `json:"is_video"`
	IsVip            int64       `json:"is_vip"`
	KindName         KindName    `json:"kind_name"`
	Location         string      `json:"location"`
	PhotoList        []string    `json:"photo_list"`
	PostID           int64       `json:"post_id"`
	Preferred        int64       `json:"preferred"`
	Price            string      `json:"price"`
	PriceUnit        PriceUnit   `json:"price_unit"`
	RefreshTime      RefreshTime `json:"refresh_time"`
	RentTag          []RentTag   `json:"rent_tag"`
	RoleName         RoleName    `json:"role_name"`
	RoomStr          string      `json:"room_str"`
	SectionName      string      `json:"section_name"`
	StreetName       string      `json:"street_name"`
	Surrounding      Surrounding `json:"surrounding"`
	Title            string      `json:"title"`
	Type             string      `json:"type"`
	YesterdayHit     int64       `json:"yesterday_hit"`
}

type RentTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Surrounding struct {
	Desc     string `json:"desc"`
	Distance string `json:"distance"`
	Type     Type   `json:"type"`
}

type TopDatum struct {
	Area        string      `json:"area"`
	Community   string      `json:"community"`
	IsVideo     int64       `json:"is_video"`
	PhotoList   []string    `json:"photo_list"`
	PostID      int64       `json:"post_id"`
	Preferred   int64       `json:"preferred"`
	Price       string      `json:"price"`
	PriceUnit   PriceUnit   `json:"price_unit"`
	RentTag     []RentTag   `json:"rent_tag"`
	RoomStr     string      `json:"room_str"`
	SectionName string      `json:"section_name"`
	StreetName  string      `json:"street_name"`
	Surrounding Surrounding `json:"surrounding"`
	Title       string      `json:"title"`
	Type        int64       `json:"type"`
}

type SEO struct {
	Description string `json:"description"`
	Keywords    string `json:"keywords"`
	Title       string `json:"title"`
}

type KindName string

// const (
// 	分租套房 KindName = "分租套房"
// 	整層住家 KindName = "整層住家"
// 	獨立套房 KindName = "獨立套房"
// 	車位   KindName = "車位"
// 	雅房   KindName = "雅房"
// )

type PriceUnit string

// const (
// 	元月 PriceUnit = "元/月"
// )

type RefreshTime string

// const (
// 	The14分鐘內 RefreshTime = "14分鐘內"
// 	The4分鐘內  RefreshTime = "4分鐘內"
// 	The9分鐘內  RefreshTime = "9分鐘內"
// )

type RoleName string

// const (
// 	代理人 RoleName = "代理人"
// 	仲介  RoleName = "仲介"
// 	屋主  RoleName = "屋主"
// )

type Type string

// const (
// 	BusStation    Type = "bus_station"
// 	SubwayStation Type = "subway_station"
// )

// type CasesID struct {
// 	Integer *int64
// 	String  *string
// }
