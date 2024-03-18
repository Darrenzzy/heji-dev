# heji-dev

### 运行说明

- 题目1：
```bash
  time ./main2.1
./main2.1  0.75s user 0.08s system 94% cpu 0.880 total
```

```bash
  time awk -F, 'BEGIN { OFS="," }
{
  if (NR == 1) { # 处理文件头
    print $1, $3, $5, $6, $7, $8, $10 > "SH_Stocks.csv"
    print $1, $3, $5, $6, $7, $8, $10 > "SZ_Stocks.csv"
  } else {
    # 将价格从整数转换为实数
    real_price = $5 / 10000
    # 根据WindCode后缀将数据写入对应的文件
    if (index($1, ".SH") != 0) {
      print $1, $3, real_price, $6, $7, $8, $10 > "SH_Stocks.csv"
    } else if (index($1, ".SZ") != 0) {
      print $1, $3, real_price, $6, $7, $8, $10 > "SZ_Stocks.csv"
    }
  }
}' transaction.1min.csv
awk -F,  transaction.1min.csv  4.55s user 2.51s system 99% cpu 7.119 total
```

- 题目2：
```bash
  time ./main2.2
  ./main2.2  0.60s user 0.04s system 102% cpu 0.618 total  
```

```bash

time awk -F, 'NR > 1 { # 跳过表头
  if (index($1, ".SH") != 0) {
    sh_stocks[$1] # 为上交所的标的创建数组元素
  } else if (index($1, ".SZ") != 0) {
    sz_stocks[$1] # 为深交所的标的创建数组元素
  }
}
END {
  
  sh_count = 0
  for (w in sh_stocks) {
    sh_count++
  }
  print "上交所的不同标的数量:", sh_count
  sz_count = 0
  for (w in sz_stocks) {
    sz_count++
  }
  print "SHWindCodeNums,SZWindCodeNums"
  print sh_count,sz_count
}' transaction.1min.csv

上交所的不同标的数量: 2364
SHWindCodeNums,SZWindCodeNums
2364 3055
awk -F,  transaction.1min.csv  3.16s user 0.03s system 99% cpu 3.192 total  
```

- 题目3：
```bash
  time ./main2.3
  ./main2.3  0.58s user 0.04s system 65% cpu 0.947 total
```

```bash
  
 time  awk -F, '
$2 != "TradingDay" && $10 != "C" {
    turnover = $7
    if (index($1, ".SH") > 0) {
        if ($8 == "B") {
            shbuy += turnover
        } else if ($8 == "S") {
            shsell += turnover
        }
    } else if (index($1, ".SZ") > 0) {
        if ($8 == "B") {
            szbuy += turnover
        } else if ($8 == "S") {
            szsell += turnover
        }
    }
}
END {
    print "SHTotalAmount,SZTotalAmount,SHTotalBuyAmount,SZTotalBuyAmount,SHTotalSellAmount,SZTotalSellAmount"
       printf "%d, %d, %d, %d, %d, %d \n",shbuy+shsell,szbuy+szsell,shbuy,szbuy,shsell,szsell

}' transaction.1min.csv


awk -F,  transaction.1min.csv  3.92s user 0.05s system 98% cpu 4.018 total
```


- 题目4：
```bash
  time ./main2.4
  ./main2.4  0.57s user 0.04s system 101% cpu 0.596 total  
```

```bash  
    time awk -F, '
    BEGIN {
        print "WindCode, Cancel_Count"
    }
    $10 == "C" {
        cancel_count[$1]++
    }
    END {
        for (code in cancel_count) {
            printf "%s, %d\n", code, cancel_count[code]
        }
    }' transaction.1min.csv


  awk -F,  transaction.1min.csv  3.10s user 0.04s system 99% cpu 3.150 total
```

- 题目5：
```bash
  time ./main2.6
  ./main2.5  0.72s user 0.06s system 116% cpu 0.670 total  
```

```bash
  
time awk -F, '
BEGIN {
    print "WindCode, BuyTransactionTimes,BuyTotalTurnover,SellTransactionTimes,SellTotalTurnover"
}
{
    if (NR > 1 && $10 != "C") {
        total_volume[$1] += $6
        total_count[$1]++
        if ($8 == "B") {
            buy_volume[$1] += $6
            buy_count[$1]++
        } else if ($8 == "S") {
            sell_volume[$1] += $6
            sell_count[$1]++
        }
    }
}
END {
    for (code in total_volume) {
        printf "%s, %d, %d, %d, %d, %d, %d\n",
               code,
               total_count[code],
               total_volume[code],
               buy_count[code],
               buy_volume[code],
               sell_count[code],
               sell_volume[code]

    }
}' transaction.1min.csv


awk -F,  transaction.1min.csv  4.05s user 0.04s system 98% cpu 4.146 total 
```

- 题目6：
```bash
  time ./main2.6
  ./main2.6  0.63s user 0.04s system 98% cpu 0.678 total  
```

```bash
  awk -F, '
BEGIN {
    OFS = ",";
    print "WindCode","HighestPrice","LowestPrice","LastPrice";
}
(NR > 1) && ($10 != "C") {
    realPrice = $5 / 10000;
    if (!(($1 in highest) && ($1 in lowest) && ($1 in last))) {
        highest[$1] = realPrice;
        lowest[$1] = realPrice;
        last[$1] = realPrice;
    } else {
        if (realPrice > highest[$1]) {
            highest[$1] = realPrice;
        }
        if (realPrice < lowest[$1]) {
            lowest[$1] = realPrice;
        }
        last[$1] = realPrice;
    }
}
END {
    for (code in highest) {
        print code, highest[code], lowest[code], last[code];
    }
}
' transaction.1min.csv

awk -F,  transaction.1min.csv  3.93s user 0.04s system 99% cpu 3.983 total  
```
