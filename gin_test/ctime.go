package main

import (
	"fmt"
	"time"
)

func main() {
	millisecondTimestamp := int64(1691682513091) // 假设这是一个毫秒时间戳示例

	// 将毫秒时间戳除以1000得到秒级时间戳
	secondTimestamp := millisecondTimestamp / 1000

	// 将秒级时间戳转换为time.Time类型
	timestampTime := time.Unix(secondTimestamp, 0)

	// 载入中国时区
	chinaTimezone, _ := time.LoadLocation("Asia/Shanghai")

	// 将时间调整为中国时区
	timestampTimeInChina := timestampTime.In(chinaTimezone)

	// 格式化为日期时间字符串
	formattedDate := timestampTimeInChina.Format("2006-01-02 15:04:05") // 格式可以根据需求进行调整

	fmt.Println("Formatted Date in China Timezone:", formattedDate)
}
