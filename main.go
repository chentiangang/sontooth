package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	xlsx "github.com/360EntSecGroup-Skylar/excelize"

	"github.com/gin-gonic/gin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func main() {
	router := gin.Default()

	router.GET("/api/get-stock-list", GetStockList)
	router.Run(":8081")
}

func shanghaiCompanies() ([]Company, error) {
	client := http.Client{}
	request, err := http.NewRequest("GET", "http://query.sse.com.cn/security/stock/downloadStockListFile.do?csrcCode=&stockCode=&areaName=&stockType=1", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36")
	request.Header.Add("Host", "query.sse.com.cn")
	request.Header.Add("Connection", "keep-alive")
	request.Header.Add("Accept", "*/*")
	request.Header.Add("Origin", "http://www.sse.com.cn")
	request.Header.Add("Referer", "http://www.sse.com.cn/assortment/stock/list/share/") //关键头 如果没有 则返回 错误
	request.Header.Add("Accept-Encoding", "gzip, deflate")
	request.Header.Add("Accept-Language", "zh-CN,zh;q=0.9")
	resp, _ := client.Do(request)
	defer resp.Body.Close()

	// 将GBK编码转为UTF8
	body := bufio.NewReader(resp.Body)
	utf8Reader := transform.NewReader(body, simplifiedchinese.GBK.NewDecoder())

	res, err := ioutil.ReadAll(utf8Reader)
	if err != nil {
		return nil, err
	}
	var companys []Company
	s := strings.Split(string(res), "\n")
	for _, i := range s {
		l := strings.Split(i, "\t")
		if len(l) > 2 {
			//fmt.Println(l[0], l[1])
			companys = append(companys, Company{Code: l[0], Name: strings.TrimSpace(l[1]), Exchange: "sh"})
		}
	}

	return companys, nil
}

func shenzhenCompanies() ([]Company, error) {
	resp, err := http.Get("http://www.szse.cn/api/report/ShowReport?SHOWTYPE=xlsx&CATALOGID=1110&TABKEY=tab1&random=0.5493995237987193")
	if err != nil {
		panic(err)
	}

	f, err := xlsx.OpenReader(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}

	rows := f.GetRows("A股列表")

	var companys []Company
	for _, row := range rows {
		companys = append(companys, Company{Code: row[4], Name: strings.TrimSpace(row[5]), Exchange: "sz", Industry: strings.TrimSpace(row[17])})
	}

	return companys, nil

}

type Company struct {
	Code     string `json:'code'`
	Name     string `json:"name"`
	Exchange string `json:"exchange"`
	Industry string `json:"industry"`
}

func GetStockList(c *gin.Context) {
	//fmt.Println(string(resp))

	companies, err := shanghaiCompanies()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	i, err := shenzhenCompanies()
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	companies = append(companies, i...)

	bs, err := json.Marshal(companies)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
	}

	c.String(http.StatusOK, string(bs))

}
