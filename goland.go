package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	"archive/zip"
)

const download_url = "http://idea.medeming.com/jets/images/jihuoma.zip"

func Goland(c *gin.Context) {
	res, err := http.Get(download_url)

	bs, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		panic(err)
	}
	num := int64(len(bs))
	zReader, _ := zip.NewReader(bytes.NewReader(bs), num)
	for _, i := range zReader.File {

		r, err := i.Open()
		if err != nil {
			panic(err)
		}
		_, err = io.Copy(c.Writer, r)
		if err != nil {
			panic(err)
		}

		err = r.Close()
		if err != nil {
			panic(err)
		}
		break
	}
}
