package parser

import (
	"bytes"
	"compress/gzip"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"procontext/models"

	"github.com/brianvoe/gofakeit"
	"github.com/valyala/fasthttp"
	"golang.org/x/text/encoding/charmap"
)

const (
	acceptHeader         = "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7"
	acceptEncodingHeader = "gzip, deflate, br"
	acceptLanguageHeader = "en-US,en;q=0.9,ru;q=0.8"
	cacheControlHeader   = "max-age=0"
)

// Функция для получения курсов валют за дату, которая передается в url (использую fasthttp)
// Работает быстрее, чем GetXML
func GetFastXML(url string) (*models.ValCurs, error) {
	client := fasthttp.Client{}
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod("GET")
	req.Header.Add("accept", acceptHeader)
	req.Header.Add("accept-encoding", acceptEncodingHeader)
	req.Header.Add("accept-language", acceptLanguageHeader)
	req.Header.Add("cache-control", cacheControlHeader)
	req.Header.Add("user-agent", gofakeit.UserAgent())
	req.SetRequestURI(url)
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	err := client.Do(req, resp)
	if err != nil {
		return nil, err
	}
	contentEncoding := resp.Header.Peek("Content-Encoding")
	var body []byte
	if bytes.EqualFold(contentEncoding, []byte("gzip")) {
		body, _ = resp.BodyGunzip()
	} else {
		body = resp.Body()
	}
	vals := &models.ValCurs{}
	bfs := bytes.NewBuffer(body)
	xml := xml.NewDecoder(bfs)
	xml.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}
	err = xml.Decode(vals)
	if err != nil {
		return nil, err
	}
	return vals, nil
}

// Функция для получения курсов валют за дату, которая передается в url (использую net/http)
func GetXML(url string) (*models.ValCurs, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	//GetFastXML(url)
	req.Header.Add("accept", acceptHeader)
	req.Header.Add("accept-encoding", acceptEncodingHeader)
	req.Header.Add("accept-language", acceptLanguageHeader)
	req.Header.Add("cache-control", cacheControlHeader)
	req.Header.Add("user-agent", gofakeit.UserAgent())

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close() //подумать над этим

	var reader io.ReadCloser
	switch response.Header.Get("Content-Encoding") {
	case "gzip":
		reader, err = gzip.NewReader(response.Body)
		if err != nil {
			fmt.Printf("error reading xml: %v", err)
			return nil, err
		}
		defer reader.Close()
	default:
		reader = response.Body
	}

	xml := xml.NewDecoder(reader)
	xml.CharsetReader = func(charset string, input io.Reader) (io.Reader, error) {
		switch charset {
		case "windows-1251":
			return charmap.Windows1251.NewDecoder().Reader(input), nil
		default:
			return nil, fmt.Errorf("unknown charset: %s", charset)
		}
	}

	codes := &models.ValCurs{}

	err = xml.Decode(codes)
	if err != nil {
		return nil, err
	}

	return codes, nil

}
