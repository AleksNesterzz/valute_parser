package parser

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"procontext/models"

	"procontext/pkg/utils"
)

const (
	Days       = 90
	DataURL    = "https://cbr.ru/scripts/XML_daily.asp?"
	DynamicURL = "https://cbr.ru/scripts/XML_dynamic.asp?"
)

func Parse() map[string]models.ValueCalc {
	t := time.Now()
	mapping := make(map[string]models.ValueCalc, Days)
	var mx sync.Mutex
	var wg sync.WaitGroup

	wg.Add(Days)
	for i := 0; i < Days; i++ {
		go func(mp map[string]models.ValueCalc, v int) {
			t1 := t.Add(-time.Hour * time.Duration(v) * 24)
			tStr := t1.Format("02/01/2006")
			respUrl := utils.CompleteURL(DataURL, t.Add(-time.Hour*time.Duration(v)*24))
			resp, err := GetFastXML(respUrl)
			if err != nil {
				fmt.Printf("error getting daily values:%v, programm shutdown", err)
				os.Exit(3)
			}
			values := resp.Valutes
			for _, v := range values {
				v.Value = strings.ReplaceAll(v.Value, ",", ".")
				price, err := strconv.ParseFloat(v.Value, 64)
				if err != nil {
					fmt.Printf("error parsing float:%v", err)
					return
				}
				mx.Lock()
				if calc, ok := mapping[v.CharCode]; ok {
					if calc.Max <= price/float64(v.Nominal) {
						calc.Max = price / float64(v.Nominal)
						calc.MaxDate = tStr
					}
					if calc.Min >= price/float64(v.Nominal) {
						calc.Min = price / float64(v.Nominal)
						calc.MinDate = tStr
					}
					calc.Avg += price / float64(v.Nominal)
					calc.Counter++
					mapping[v.CharCode] = calc
				} else {
					calc := models.ValueCalc{}
					calc.Max = price / float64(v.Nominal)
					calc.MaxDate = t.Format("02/01/2006")
					calc.Min = price / float64(v.Nominal)
					calc.MinDate = t.Format("02/01/2006")
					calc.Name = v.Name
					calc.Avg = price / float64(v.Nominal)
					calc.Counter++
					mapping[v.CharCode] = calc
				}
				mx.Unlock()
			}
			wg.Done()
		}(mapping, i)
	}
	wg.Wait()
	return mapping
}
