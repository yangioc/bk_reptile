package web591

var loginPage = "https://rent.591.com.tw"
var webpage = "https://rent.591.com.tw/home/search/rsList?is_format_data=1&is_new_list=1&type=1&kind=1&multiRoom=2,3,4&rentprice=15000,40000&order=posttime&orderType=desc&showMore=0&multiNotice=not_cover"
var webpagelast = "https://rent.591.com.tw/home/search/rsList?is_format_data=1&is_new_list=1&type=1&kind=1&multiRoom=2,3,4&rentprice=15000,40000&order=posttime&orderType=desc&showMore=0&multiNotice=not_cover&firstRow=%d&totalRows=%d"
var objPage = "https://bff.591.com.tw/v1/house/rent/detail?id=%d"
var targetPage = "https://rent.591.com.tw/home/%d"
var rssTemplate = `租金：[$租金] 額外成本：$額外成本 格局：[$格局] 區域：[$區域]`
