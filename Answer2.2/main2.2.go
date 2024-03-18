package main

import (
	"fmt"
	"log"

	"github.com/shopspring/decimal"
	"gonum.org/v1/hdf5"
)

const ffname = "transaction.h5"

func main() {

	sf, err := hdf5.OpenFile(ffname, hdf5.F_ACC_RDONLY)
	if err != nil {
		log.Fatalf("fail to open the file: %s", err)
	}
	defer sf.Close()
	sourceDatas := getSourceData(sf)
	group, err := sf.OpenGroup("DB")
	if err != nil {
		log.Fatal(err.Error())
	}
	// 创建或打开HDF5文件
	f, err := hdf5.CreateFile("depth.h5", hdf5.F_ACC_TRUNC)
	if err != nil {
		log.Fatalf("Failed to create file: %s", err)
	}
	defer f.Close()

	// 创建一个主组
	mainGroup, err := f.CreateGroup("DB")
	if err != nil {
		log.Fatalf("Failed to create group: %s", err)
	}
	defer mainGroup.Close()

	for _, sId := range sourceDatas.stockIds {
		stockId := fmt.Sprint(sId)
		for len(stockId) < 6 {
			stockId = "0" + stockId
		}
		data := getStockById(group, stockId)
		data.calculateData()

		func() {
			subGroup, err := mainGroup.CreateGroup(stockId)
			if err != nil {
				log.Fatalf("Failed to create subgroup: %s", err)
			}
			defer subGroup.Close()

			// 在子组中创建不同类型的数据集
			dsetType := []*hdf5.Datatype{hdf5.T_NATIVE_UINT32, hdf5.T_NATIVE_DOUBLE, hdf5.T_NATIVE_DOUBLE, hdf5.T_NATIVE_DOUBLE, hdf5.T_NATIVE_DOUBLE, hdf5.T_NATIVE_DOUBLE, hdf5.T_NATIVE_INT32}
			dataset := []string{"T", "Open", "High", "Low", "Close", "Amount", "Volume"}
			dims := len(data.arr)

			for j, dt := range dsetType {
				func() {
					dsetName := dataset[j]
					dspace, err := hdf5.CreateSimpleDataspace([]uint{uint(dims)}, nil)
					if err != nil {
						log.Fatalf("Failed to create dataspace: %s", err)
					}
					defer dspace.Close()

					dset, err := subGroup.CreateDataset(dsetName, dt, dspace)
					if err != nil {
						log.Fatalf("Failed to create dataset: %s", err)
					}
					defer dset.Close()

					// 根据数据集类型写入数据
					switch j {
					case 0:
						dd, _, _ := data.collectData("T", dims)
						err = dset.Write(&dd)
					case 1:
						_, dd, _ := data.collectData("Open", dims)
						err = dset.Write(&dd)
					case 2:
						_, dd, _ := data.collectData("High", dims)
						err = dset.Write(&dd)
					case 3:
						_, dd, _ := data.collectData("Low", dims)
						err = dset.Write(&dd)
					case 4:
						_, dd, _ := data.collectData("Close", dims)
						err = dset.Write(&dd)
					case 5:
						_, dd, _ := data.collectData("Amount", dims)
						err = dset.Write(&dd)
					case 6:
						_, _, dd := data.collectData("Volume", dims)
						err = dset.Write(&dd)
					}
					if err != nil {
						log.Fatalf("Failed to write data to dataset: %s", err.Error())
					}
				}()

			}
		}()

	}
	return
}

func (S *sourceData) collectData(tp string, dims int) ([]uint32,
	[]float64,
	[]int32) {
	uu := make([]uint32, dims)
	ff := make([]float64, dims)
	ii := make([]int32, dims)
	for i, ohlc := range S.arr {
		switch tp {
		case "T":
			uu[i] = formatTimestamp(int(ohlc.Ts))
		case "Open":

			ff[i] = ohlc.Open.InexactFloat64()
		case "High":

			ff[i] = ohlc.High.InexactFloat64()
		case "Low":

			ff[i] = ohlc.Low.InexactFloat64()
		case "Close":
			ff[i] = ohlc.Close.InexactFloat64()
		case "Amount":

			ff[i] = ohlc.Amount.InexactFloat64()
		case "Volumeqq":
			ii[i] = ohlc.Volume
		}

	}
	return uu, ff, ii
}
func (S *sourceData) addNext(newData OHLC) {
	if len(S.arr) == 0 {
		S.arr = append(S.arr, newData)
		return
	}
	data := S.arr[len(S.arr)-1]
	if data.Ts != newData.Ts {
		S.arr = append(S.arr, newData)
		return
	}
	if data.High.Cmp(newData.High) < 0 {
		data.High = newData.High
	}
	// 过滤无效最小值
	if data.Low.Cmp(newData.Low) > 0 && !newData.Low.IsZero() {
		data.Low = newData.Low
	}
	if !newData.Close.IsZero() {
		data.Close = newData.Close
	}
	data.Volume += newData.Volume
	data.Amount = data.Amount.Add(newData.Amount)
	S.arr[len(S.arr)-1] = data
}
func (S *sourceData) calculateData() {
	// nextTimestamp := S.timestamps[0]
	// 下一次区间的结尾时间戳
	next := int(0)
	for i, t := range S.timestamps {
		ms := parseTimestamp(t)
		if next == 0 {
			next = ms + 3000
		}

		if ms > next {
			next = ms + 3000
		}

		vo := int32(S.volumes[i])
		// 防止无效成交量算入
		if S.prices[i] <= 0 {
			vo = 0
		}
		S.addNext(OHLC{
			Ts:     uint32(next),
			Open:   decimal.NewFromFloat(S.prices[i]),
			High:   decimal.NewFromFloat(S.prices[i]),
			Low:    decimal.NewFromFloat(S.prices[i]),
			Close:  decimal.NewFromFloat(S.prices[i]),
			Volume: vo,
			Amount: decimal.NewFromFloat(S.prices[i] * float64(S.volumes[i])),
		})

	}

}

// parseTimestamp 将HHMMSSsss格式的时间戳转换为总毫秒数
func parseTimestamp(ts uint32) int {
	hours := ts / 10000000
	minutes := (ts % 10000000) / 100000
	seconds := (ts % 100000) / 1000
	milliseconds := ts % 1000
	return int(hours)*3600000 + int(minutes)*60000 + int(seconds)*1000 + int(milliseconds)
}

// formatTimestamp 总毫秒数转换回HHMMSSsss格式的时间戳
func formatTimestamp(ms int) uint32 {
	hours := ms / 3600000
	minutes := (ms % 3600000) / 60000
	seconds := (ms % 60000) / 1000
	milliseconds := ms % 1000
	return uint32(hours)*10000000 + uint32(minutes)*100000 + uint32(seconds)*1000 + uint32(milliseconds)
}

type sourceData struct {
	stockIds   []int32
	turnover   []int32
	total      int64
	prices     []float64
	volumes    []int32
	amounts    []float64
	timestamps []uint32
	arr        []OHLC
}

// 定义一个结构体来保存OHLC数据
type OHLC struct {
	Ts     uint32
	Open   decimal.Decimal
	High   decimal.Decimal
	Low    decimal.Decimal
	Close  decimal.Decimal
	Volume int32
	Amount decimal.Decimal
}

func getStockById(group *hdf5.Group, stockId string) *sourceData {
	stockGroup, err := group.OpenGroup(stockId)
	if err != nil {
		log.Fatal(err.Error())
	}
	pricesColumn, err := stockGroup.OpenDataset("3")
	amountsColumn, err := stockGroup.OpenDataset("5")
	volumeColumn, err := stockGroup.OpenDataset("4")
	timestampColumn, err := stockGroup.OpenDataset("T")
	if err != nil {
		log.Fatal(err.Error())
	}
	dataspace := pricesColumn.Space()
	defer dataspace.Close()
	// 获取数量
	dims, _, err := dataspace.SimpleExtentDims()
	if err != nil {
		log.Fatalln(err)
	}
	// 当日所有价格
	prices := make([]float64, dims[0])
	err = pricesColumn.Read(&prices)
	if err != nil {
		log.Fatal(err.Error())
	}

	volumes := make([]int32, dims[0])
	err = volumeColumn.Read(&volumes)
	if err != nil {
		log.Fatal(err.Error())
	}
	amounts := make([]float64, dims[0])
	err = amountsColumn.Read(&amounts)
	if err != nil {
		log.Fatal(err.Error())
	}
	timestamps := make([]uint32, dims[0])
	err = timestampColumn.Read(&timestamps)
	if err != nil {
		log.Fatal(err.Error())
	}

	return &sourceData{
		prices:     prices,
		volumes:    volumes,
		amounts:    amounts,
		timestamps: timestamps,
	}
}

// 从HDF5文件中获取股票id 和总数
func getSourceData(ff *hdf5.File) *sourceData {
	dset, err := ff.OpenDataset("StockPool")
	if err != nil {
		log.Fatal(err)
	}
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

	return &sourceData{
		total:    int64(dims[0]),
		stockIds: stockIds,
	}

}
