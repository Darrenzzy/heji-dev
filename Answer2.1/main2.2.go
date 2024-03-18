package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"
)

// 1.将上交所的股票和深交所的标的 (WindCode) 分割成两个文件分别存储，且只保留列 WindCode,Time, Price, Volume, Turnover, BSFlag, FunctinCode.
// 2. 统计上交所和深交所分别出现了多少个不同的标的 (WindCode)
// 3. 加总上海和深圳交易所分别的总成交金额 (需要使用 FunctionCode, Turnover), 以及主买主卖 (BSFlag) 分别的成交金额
// 4. 统计每只股票的撤单次数 (FunctionCode)
// 5. 每只股票成交的次数和总成交量 (FunctinCode, Volume), 分主卖主买 (BSFlag) 分别统计的次数和
// 成交量 (Volume)
// 6. 每只股票的最高成交价格和最低成交价格，以及最新成交价格 (Price).
// 最后比较 awk 计算需要的时间和自己程序所用时间。

func main() {
	// 打开CSV文件
	f, err := os.Open("transaction.1min.csv")
	if err != nil {
		log.Fatal("Unable to read input file example.csv", err)
	}
	defer f.Close()
	shFile, err := os.Create("sh-transaction.csv")
	if err != nil {
		fmt.Println("open file is failed, err: ", err)
	}
	szFile, err := os.Create("sz-transaction.csv")
	if err != nil {
		fmt.Println("open file is failed, err: ", err)
	}
	// 延迟关闭
	defer szFile.Close()
	defer shFile.Close()
	// 创建CSV阅读器
	csvReader := csv.NewReader(f)
	shMap, szMap := make(map[string]struct{}), make(map[string]struct{})
	for {
		recored, err := csvReader.Read()
		if err != nil {
			if err == csv.ErrFieldCount {
				// 如果是字段数量不匹配的错误，可以选择忽略，继续处理下一行
				continue
			}
			break // 文件结束或者发生其它错误
		}
		if strings.HasSuffix(recored[0], "SH") {
			shMap[recored[0]] = struct{}{}
		} else if strings.HasSuffix(recored[0], "SZ") {
			szMap[recored[0]] = struct{}{}
		}
	}

	fmt.Println("SHWindCodeNums,SZWindCodeNums")
	fmt.Printf("%d,%d\n", len(shMap), len(szMap))
}
