package utils

import (
	"fmt"
	"procontext/models"
	"sort"
	"time"
)

func CompleteURL(url string, t time.Time) string {
	return fmt.Sprintf("%sdate_req=%s", url, t.Format("02/01/2006"))
}

func PrintSortedMapByKeys(data map[string]models.ValueCalc) {
	keys := make([]string, 0)

	for k := range data {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	for _, k := range keys {
		fmt.Printf("%s | 1 %s | Max:%f RUB (%s) | Min:%f RUB (%s) | Avg:%f RUB\n", k, data[k].Name, data[k].Max, data[k].MaxDate, data[k].Min, data[k].MinDate, data[k].Avg/float64(data[k].Counter))
	}
}
