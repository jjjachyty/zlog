package zlog

import "time"

//DailyRollingLogs func 每隔24小时执行一次任务
func DailyRollingLogs(f func()) {
	go func() {
		for {

			now := time.Now()
			// 计算下一个零点
			next := now.Add(time.Hour * 24)
			//next := now.Add(time.Second * 10)
			next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())

			t := time.NewTimer(next.Sub(now))
			//t := time.NewTimer(time.Second * 10)
			<-t.C
			f()
		}
	}()
}
