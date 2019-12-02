package task

import (
	"time"
)

// 启动定时任务
func Run(d time.Duration, f func()) {
	// 创建Ticker
	ticker := time.NewTicker(d)
	// 启动任务
	go loop(ticker, f)
}

// 定时任务循环部分
func loop(ticker *time.Ticker, f func()) {
	for range ticker.C {
		f()
	}
}
