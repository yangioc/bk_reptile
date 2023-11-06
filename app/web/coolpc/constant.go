package coolpc

var WebPage = "https://www.coolpc.com.tw/evaluate.php"

const (
	CheckTimeTag = "Mdy"
)

const (
	N1_Type_Name  = "品牌小主機、AIO｜VR虛擬"
	N2_Type_Name  = "手機｜平板｜筆電｜穿戴"
	N3_Type_Name  = "酷！PC 套裝產線 "
	N4_Type_Name  = "處理器 CPU"
	N5_Type_Name  = "主機板 MB "
	N6_Type_Name  = "記憶體 RAM "
	N7_Type_Name  = "固態硬碟 M.2｜SSD"
	N8_Type_Name  = "傳統內接硬碟HDD"
	N9_Type_Name  = "外接硬碟｜隨身碟｜記憶卡"
	N10_Type_Name = "散熱器｜散熱墊｜散熱膏"
	N11_Type_Name = "封閉式｜開放式水冷"
	N12_Type_Name = "顯示卡VGA "
	N13_Type_Name = "螢幕｜投影機｜壁掛 "
	N14_Type_Name = "CASE 機殼(+電源) "
	N15_Type_Name = "電源供應器 "
	N16_Type_Name = "機殼風扇｜機殼配件"
	N17_Type_Name = " 鍵盤+鼠｜搖桿｜桌+椅"
	N18_Type_Name = "滑鼠｜鼠墊｜數位板"
	N19_Type_Name = "IP分享器｜網卡｜網通設備"
	N20_Type_Name = "網路NAS｜網路IPCAM"
	N21_Type_Name = "音效卡｜電視卡(盒)｜影音"
	N22_Type_Name = "喇叭｜耳機｜麥克風"
	N23_Type_Name = "燒錄器 CD/DVD/BD"
	N24_Type_Name = "USB週邊｜硬碟座｜讀卡機"
	N25_Type_Name = "行車紀錄器｜USB視訊鏡頭"
	N26_Type_Name = "UPS不斷電｜印表機｜掃描"
	N27_Type_Name = "介面擴充卡｜專業Raid卡"
	N28_Type_Name = "網路、傳輸線、轉頭｜KVM"
	N29_Type_Name = "OS+應用軟體｜禮物卡 "
	N30_Type_Name = "福利品出清"
)

var typeMap = map[int]string{
	1:  N1_Type_Name,
	2:  N2_Type_Name,
	3:  N3_Type_Name,
	4:  N4_Type_Name,
	5:  N5_Type_Name,
	6:  N6_Type_Name,
	7:  N7_Type_Name,
	8:  N8_Type_Name,
	9:  N9_Type_Name,
	10: N10_Type_Name,
	11: N11_Type_Name,
	12: N12_Type_Name,
	13: N13_Type_Name,
	14: N14_Type_Name,
	15: N15_Type_Name,
	16: N16_Type_Name,
	17: N17_Type_Name,
	18: N18_Type_Name,
	19: N19_Type_Name,
	20: N20_Type_Name,
	21: N21_Type_Name,
	22: N22_Type_Name,
	23: N23_Type_Name,
	24: N24_Type_Name,
	25: N25_Type_Name,
	26: N26_Type_Name,
	27: N27_Type_Name,
	28: N28_Type_Name,
	29: N29_Type_Name,
	30: N30_Type_Name,
}
