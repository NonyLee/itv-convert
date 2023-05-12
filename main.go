package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"

	"github.com/xuri/excelize/v2"
)

type Channel struct {
	name         string
	bitRate      float64
	audioQuality string
	category     string
	channelId    string
	url          string
}

func filterBitRate(cs []Channel) []Channel {
	ls := make([]Channel, 0)
	sort.SliceStable(cs, func(i, j int) bool {
		return cs[i].bitRate > cs[j].bitRate
	})
	for _, c := range cs {
		lsLen := len(ls)
		if lsLen > 0 && ls[lsLen-1].bitRate != c.bitRate {
			break
		}
		ls = append(ls, c)
	}

	return ls
}

func selectChannel(cs []Channel) Channel {
	priority := map[string]int{
		"易视腾": 10,
		"百视通": 9,
		"华数":  8,
	}
	sort.SliceStable(cs, func(i, j int) bool {
		return priority[cs[i].category] > priority[cs[j].category]
	})

	return cs[0]
}

func main() {
	f, err := excelize.OpenFile("itv.xlsx")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer func() {
		// Close the spreadsheet.
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		fmt.Println(err)
		return
	}
	m := make(map[string][]Channel)
	for i, row := range rows {
		if i == 0 || row[3] == "湖南移动" || row[3] == "宁夏移动" || row[3] == "咪咕视频" {
			continue
		}
		compileRegex := regexp.MustCompile("(.*?)M")
		matchArr := compileRegex.FindStringSubmatch(row[1])
		bitRate, _ := strconv.ParseFloat(matchArr[1], 32)

		c := Channel{
			name:         row[0],
			bitRate:      bitRate,
			audioQuality: row[2],
			category:     row[3],
			channelId:    row[4],
			url:          row[5],
		}

		if m[c.name] == nil {
			m[c.name] = make([]Channel, 0)
		}

		m[c.name] = append(m[c.name], c)
	}

	channels := make([]Channel, 0)
	for _, v := range m {
		hcs := filterBitRate(v)
		c := selectChannel(hcs)
		channels = append(channels, c)
	}
}
