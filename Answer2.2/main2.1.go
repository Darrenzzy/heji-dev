package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"

	"gonum.org/v1/hdf5"
)

const fname = "transaction.h5"

func main() {
	// 打开文件
	f, err := hdf5.OpenFile(fname, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Fatalf("fail to open the file: %s", err)
	}
	defer f.Close()

	dset, err := f.OpenDataset("StockPool")
	dataspace := dset.Space()
	defer dataspace.Close()

	// 获取数量
	dims, _, err := dataspace.SimpleExtentDims()
	if err != nil {
		log.Fatalln(err)
	}
	stockIds := make([]int32, dims[0])
	err = dset.Read(&stockIds)
	if err != nil {
		log.Fatal(err.Error())
	}

	group, err := f.OpenGroup("DB")
	if err != nil {
		log.Fatal(err.Error())
	}
	// 准备CSV文件
	csvFile, err := os.Create("max_drawdowns.csv")
	if err != nil {
		log.Fatal("Failed to create csv file:", err)
	}
	defer csvFile.Close()
	writer := csv.NewWriter(csvFile)
	defer writer.Flush()

	// 写入CSV头部
	headers := []string{"StockID", "StartTimestamp", "EndTimestamp", "MaxDrawdown"}
	if err := writer.Write(headers); err != nil {
		log.Fatal("Error writing headers to csv file:", err)
	}
	for _, sId := range stockIds {
		func() {
			stockId := fmt.Sprint(sId)
			for len(stockId) < 6 {
				stockId = "0" + stockId
			}
			aa, err := group.OpenGroup(stockId)
			if err != nil {
				log.Fatal(err.Error())
			}
			prices3, err := aa.OpenDataset("3")
			cancelFlagColumn, err := aa.OpenDataset("7")
			timestampColumn, err := aa.OpenDataset("T")
			if err != nil {
				log.Fatal(err.Error())
			}
			dataspace = prices3.Space()
			defer dataspace.Close()

			// 获取数量
			dims, _, err = dataspace.SimpleExtentDims()
			if err != nil {
				log.Fatalln(err)
			}
			// 当日所有价格
			prices := make([]float64, dims[0])
			// 当日所有订单状态
			cancelFlags := make([]int8, dims[0])
			// 当日所有时间戳
			timestampColumns := make([]uint32, dims[0])
			err = cancelFlagColumn.Read(&cancelFlags)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = timestampColumn.Read(&timestampColumns)
			if err != nil {
				log.Fatal(err.Error())
			}
			err = prices3.Read(&prices)
			if err != nil {
				log.Fatal(err.Error())
			}

			// 最大位置
			maxPlace := 0
			// 最小位置
			minPlace := 0
			maxPrice := float64(0)
			drawdown := float64(0)
			maxDrawdown := float64(0)
			for i, price := range prices {
				// 跳过无订单的数据
				if cancelFlags[i] != 70 {
					continue
				}
				if price <= maxPrice {
					drawdown = maxPrice - price
					if maxDrawdown < drawdown {
						maxPlace = i
						maxDrawdown = max(maxDrawdown, drawdown)
					}

				} else {
					maxPrice = price
					minPlace = i
				}
			}
			writer.Write([]string{
				fmt.Sprint(stockId),
				fmt.Sprint(timestampColumns[minPlace]),
				fmt.Sprint(timestampColumns[maxPlace]),
				fmt.Sprintf("%.2f", maxDrawdown)})
		}()

	}

}
