package stocktw

import (
	"testing"
)

func Test_Get_twse_base(t *testing.T) {
	// sqldb.InitMysql("root@tcp(127.0.0.1:3306)/stock?charset=utf8mb4&parseTime=True&loc=UTC")
	// err := Get_twse_base("5283")
	// fmt.Println(err)
	Unmarshal_twse_base([]byte(twse_base_mock))
}

const twse_base_mock = `

<html>
<head>
	<title>公開資訊觀測站</title>
	
	
	<link href="css/css2.css" rel="stylesheet" type="text/css" Media="Screen"/> 
	
	<script type="text/javascript" src="js/mops2.js"></script>
</head>

<body>
<!_co_id_hhc=1101__>
<table class='noBorder' width='80%'>
<tr><td align='center' class='compName'>
<b>
本資料由　<span style='color:blue;'>(上市公司)
台泥</span>　公司提供</b>
</td></tr></table>
<table class='hasBorder' width=95%>
<tr>
<th class='dColor nowrap' style='text-align:left !important;'>股票代號</th><td class='lColor'>1101</td>
<th class='dColor nowrap' style='text-align:left !important;'>產業類別</th>
<td nowrap class='lColor'>水泥工業                &nbsp;</td><th class='dColor nowrap' style='text-align:left !important;'>外國企業註冊地國</th>
<td class='lColor nowrap' colspan='2' style='text-align:center !important; font-family:arial;'>－</td>
<tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>公司名稱</th>
<td class='lColor'>臺灣水泥股份有限公司</td><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>總機</th>
<td class='lColor' colspan='2'>(02)2531-7099</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>地址</th>
<td colspan=5 class='lColor'>台北市中山北路2段113號</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>董事長</th>
<td class='lColor'>張安平</td><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>總經理</th>
<td class='lColor' colspan='2'>程耀輝</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>發言人</th>
<td class='lColor'>葉毓君</td><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>發言人職稱</th>
<td class='lColor' colspan='2'>永續長</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>發言人電話</th>
<td class='lColor'>2531-7099轉20265</td><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>代理發言人</th>
<td class='lColor' colspan='2'>賴家柔</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>主要經營業務</th>
<td colspan=5 class='lColor'>有關水泥及水泥製品之生產及運銷<BR>有關水泥原料及水泥製品原料之開採製造運銷及附屬礦石之開採經銷<BR>經營有關水泥工業及其附屬事業</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>公司成立日期</th>
<td class='lColor'>39/12/29</td>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>營利事業統一編號</th>
<td class='lColor' colspan='2'>11913502</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>實收資本額</th>
<td class='lColor'>     73,561,817,420元</td>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>上市日期</th>
<td class='lColor' colspan='2'>51/02/09</td>
<tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>上櫃日期</th>
<td class='lColor'>&nbsp</td>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>興櫃日期</th>
<td class='lColor' colspan='2'>&nbsp</td>
<tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>公開發行日期</th>
<td class='lColor'>&nbsp</td>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>普通股每股面額</th>
<td class='lColor' colspan='2'>新台幣                 10.0000元</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>已發行普通股數或TDR原股發行股數</th>
<td class='lColor'>      7,156,181,742股
(含私募                  0股)
</td>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>特別股</th>
<td class='lColor' colspan='2'>        200,000,000股</td>
<tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>普通股盈餘分派或虧損撥補頻率</th>
<td class='lColor'>每年</td>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>普通股年度 (含第4季或後半年度)現金股息及紅利決議層級</th>
<td class='lColor' colspan='2'>股東會</td>
<tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>股票過戶機構</th>
<td class='lColor'>中國信託商業銀行代理部</td><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>電話</th>
<td class='lColor' colspan='2'>66365566</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>過戶地址</th>
<td colspan=5 class='lColor'>台北市重慶南路一段83號5樓</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>簽證會計師事務所</th>
<td colspan=5 class='lColor'>勤業眾信聯合會計師事務所</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>簽證會計師1</th>
<td colspan=5 class='lColor'>翁雅玲</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>簽證會計師2</th>
<td colspan=5 class='lColor'>黃惠敏</td><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>備註</th>
<td colspan=5 class='lColor'>&nbsp;</td><tr>
<th class='dColor nowrap' style='text-align:left !important;'>本公司</th>
<td class='lColor'>有</td>
<th class='dColor nowrap' style='text-align:left !important;'>特別股發行</th>
<th class='dColor nowrap' style='text-align:left !important;'>本公司</th>
<td class='lColor'>有</td>
<th class='dColor nowrap' style='text-align:left !important;' colspan='2'>公司債發行</th><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>英文簡稱</th>
<td colspan=5 class='lColor'>&nbsp;TCC</td><tr><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>英文全名</th>
<td colspan=5 class='lColor'>&nbsp;Taiwan Cement Corp.</td><tr><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>英文通訊地址(街巷弄號)</th>
<td colspan=1 class='lColor' colspan='2'>&nbsp;No.113, Sec.2, Zhongshan N. Rd.,</td><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>英文通訊地址(縣市國別)</th>
<td colspan=2 class='lColor'>&nbsp;Taipei City 104,Taiwan (R.O.C.)</td></tr><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>傳真機號碼</th>
<td class='lColor'>&nbsp;(02)2531-6529</td><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>電子郵件信箱</th>
<td class='lColor' colspan='2'>&nbsp;finance@taiwancement.com</td><tr><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>公司網址</th>
<td colspan=5 class='lColor'>&nbsp;<a href=http://www.taiwancement.com target=NEW_WINDOW>http://www.taiwancement.com</a></td><tr><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>投資人關係聯絡人</th>
<td colspan=1 class='lColor' colspan='2'>&nbsp;永續辦公室</td><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>投資人關係聯絡人職稱</th>
<td colspan=2 class='lColor'>&nbsp;永續辦公室</td></tr><tr>
<th colspan=2 class='dColor nowrap' style='text-align:left !important;'>投資人關係聯絡電話</th>
<td class='lColor'>&nbsp;02-25317099</td><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>投資人關係電子郵件</th>
<td class='lColor' colspan='2'>&nbsp;ir@taiwancement.com</td><tr><th colspan=2 class='dColor nowrap' style='text-align:left !important;'>公司網站內利害<br>關係人專區網址</th>
<td colspan=5 class='lColor'>&nbsp;<a href=http://www.taiwancement.com/tw/csr/csr5-1.html target=NEW_WINDOW>http://www.taiwancement.com/tw/csr/csr5-1.html</a></td><tr><th class='dColor nowrap' style='text-align:left !important;'>變更前名稱</th>
<td class='lColor'>&nbsp;</td><th class='dColor nowrap' style='text-align:left !important;'>變更前簡稱</th>
<td class='lColor'>&nbsp;</td><th class='dColor nowrap' style='text-align:left !important;'>公司名稱變更核准日期</th>
<td class='lColor' colspan='2'>0</td>
</tr>
</table>
<br><br>
<table class='hasBorder' width=95%>
<th class='dColor nowrap' style='text-align:left !important;'>本公司採</th>
<td class='lColor'>&nbsp;</td><th colspan=3 class='dColor nowrap' style='text-align:left !important;'>月制會計年度（空白表曆年制）</th>
<tr>
<th class='dColor nowrap' style='text-align:left !important;'>本公司於</th>
<td class='lColor'>&nbsp;</td><th class='dColor nowrap' style='text-align:left !important;'>之前採</th>
<td>&nbsp;</td><th class='dColor nowrap' style='text-align:left !important;'>月制會計年度</th>
</tr>
<tr>
<th colspan='2' class='dColor nowrap' style='text-align:left !important;'>編製財務報告類型</th>
<td colspan='3' class='lColor' style='text-align:left !important;'>
●合併○個別
</td>
</tr>
</table><br><br>

</body>
</html>`
