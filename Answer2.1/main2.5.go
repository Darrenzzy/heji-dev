package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

// 1.将上交所的股票和深交所的标的 (WindCode) 分割成两个文件分别存储，且只保留列 WindCode,Time, Price, Volume, Turnover, BSFlag, FunctinCode.
//
// 2. 统计上交所和深交所分别出现了多少个不同的标的 (WindCode)
// 3. 加总上海和深圳交易所分别的总成交金额 (需要使用 FunctionCode, Turnover), 以及主买主卖 (BSFlag) 分别的成交金额
// 4. 统计每只股票的撤单次数 (FunctionCode)
// 5. 每只股票成交的次数和总成交量 (FunctinCode, Volume), 分主卖主买 (BSFlag) 分别统计的次数和成交量 (Volume)
// 6. 每只股票的最高成交价格和最低成交价格，以及最新成交价格 (Price).
// 最后比较 awk 计算需要的时间和自己程序所用时间。

func main() {
	// 打开CSV文件
	f, err := os.Open("transaction.1min.csv")
	if err != nil {
		log.Fatal("Unable to read input file example.csv", err)
	}
	defer f.Close()
	// 创建CSV阅读器
	csvReader := csv.NewReader(f)

	// 读取第一行，以获取列名称
	csvReader.Read()

	windMap := make(map[string]*stock)
	for {
		record, err := csvReader.Read()
		if err != nil {
			if err == csv.ErrFieldCount {
				continue
			}
			break // 文件结束或者发生其它错误
		}
		if record[9] == "C" {
			continue
		}
		volume, ok := windMap[record[0]]
		if !ok {
			volume = &stock{&stockTradeType{}, &stockTradeType{}}
		}

		v, err := strconv.Atoi(record[5])
		if err != nil {
			log.Println(err.Error())
		}
		if record[7] == "B" {
			volume.B.volume += v
			volume.B.times++
		} else if record[7] == "S" {
			volume.S.volume += v
			volume.S.times++
		}
		windMap[record[0]] = volume
	}
	fmt.Printf("WindCode,TotalCount,TotalVolume,BuyCount,BuyVolume,SellCount,SellVolume")
	for k, nums := range windMap {
		fmt.Println(k, nums.S.times+nums.B.times, nums.S.volume+nums.B.volume, nums.B.times, nums.B.volume, nums.S.times, nums.S.volume)
	}
}

type stock struct {
	B *stockTradeType
	S *stockTradeType
}

type stockTradeType struct {
	times  int
	volume int
}
