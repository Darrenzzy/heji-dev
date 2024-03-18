// monitor.go
package main

import (
	"log"
	"os"
	"os/exec"
	"time"
)

const (
	aliceExecutable = "./alice" // Alice 程序的路径
	logFile         = "monitor.log"
	checkInterval   = 5 * time.Second // 检查 Alice 进程的间隔时间
)

func main() {
	// 设置日志文件
	setupLogging()

	// 开始监控 Alice 进程
	monitorAlice()
}

func setupLogging() {
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %s", err)
	}
	log.SetOutput(file)
}

func monitorAlice() {
	for {
		// 启动 Alice 进程
		cmd := exec.Command(aliceExecutable)
		err := cmd.Start()
		if err != nil {
			log.Printf("Unable to start Alice: %s", err)
			notify("Failed to start Alice")
		} else {
			log.Println("Alice started successfully")
		}

		// 等待 Alice 进程结束
		err = cmd.Wait()
		if err != nil {
			log.Printf("Alice exited with error: %s", err)
			notify("Alice exited unexpectedly")
		}

		// 等待一段时间后重新启动 Alice
		time.Sleep(checkInterval)
	}
}

func notify(message string) {
	// 这里可以加入发送通知的代码，比如发送邮件、短信或推送通知
	log.Printf("Notification: %s", message)
	// 例如：sendEmail(message) 或 sendSMS(message)
}
