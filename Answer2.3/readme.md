# 逻辑和分析

### 运行说明

1. 运行监控程序
    ```./monitor```
2. 检查alice是否已经运行  
   ```ps aux |grep alice```
3. 如果alice没有运行，运行alice  
   ```./alice```
4. 启动bob程序  
   ```./bob```
5. 键入查询："query 601156.SH"
6. 启动多个bob程序
7. 退出bob程序 
   ```byb```
8. 退出alice程序, kill alice进程  
   ```kill -9 <alice进程号>```
9. 检查alice是否被重启
10. 退出监控程序  
   ```exit```
11. 检查日志文件  
   ```cat ./monitor.log
      cat ./alice.log
  ```