package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"

	"github.com/chentiangang/xlog"
	"github.com/gin-gonic/gin"
	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

const SuggestURL = "http://suggest3.sinajs.cn/suggest/name=info&key="

func isValidCode(code string) bool {
	return strings.Contains(code, "sh") || strings.Contains(code, "sz")
}

func Suggest(c *gin.Context) {
	input := c.Query("match")
	fmt.Println(input)

	ret := make([]map[string]string, 0, 1024)
	if len(input) == 0 {
		c.JSON(http.StatusOK, ret)
	}

	resp, err := grequests.Get(SuggestURL+input, nil)
	if err != nil {
		c.JSON(http.StatusOK, ret)
	}
	if resp.StatusCode != 200 {
		c.JSON(http.StatusOK, ret)
	}

	rawResult, err := ConvertGB2UTF8(resp.String())
	if err != nil {
		c.JSON(http.StatusOK, ret)
	}
	rawResult = strings.Split(rawResult, "=")[1]
	rawResult = strings.ReplaceAll(rawResult, "\"", "")
	parts := strings.Split(rawResult, ";")
	if len(parts) == 0 {
		c.JSON(http.StatusOK, ret)
	}
	for _, part := range parts {
		if len(part) == 0 {
			continue
		}
		seps := strings.Split(part, ",")
		if !isValidCode(seps[3]) {
			xlog.LogDebug("ignore, not valid sh or sz code [%s/%s]", seps[4], seps[3])
			continue
		}
		ret = append(ret, map[string]string{
			"name": seps[4],
			"code": ConvertCodeBack(seps[3]),
		})
	}
	c.JSON(http.StatusOK, ret)
}

func ConvertGB2UTF8(raw string) (string, error) {
	reader := transform.NewReader(bytes.NewReader([]byte(raw)), simplifiedchinese.GB18030.NewDecoder())
	rawBytes, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", errors.Wrapf(err, "failed to convert GB18030 to UTF8")
	}
	return string(rawBytes), nil
}

// ConvertCodeBack convert "sh000001" to sina accept code "sh.000001"
func ConvertCodeBack(code string) string {
	if len(code) < 2 {
		return code
	}
	return fmt.Sprintf("%s.%s", code[:2], code[2:])
}
