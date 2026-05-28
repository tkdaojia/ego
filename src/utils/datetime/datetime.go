package datetime

import (
	"github.com/spf13/cast"
	"strings"
	"time"
)

type DateTime struct{}

func (datetime *DateTime) RFC3339ToYmdHis(tstr interface{}) string {
	if t, ok := tstr.(time.Time); ok {
		return t.Format("2006-01-02 15:04:05")
	}
	return ""
}

// DateToTime 日期字符串转时间戳（带安全防御，解析失败返回 0）
func (datetime *DateTime) DateToTime(str string) int64 {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0
	}

	t, err := time.ParseInLocation("2006-01-02 15:04:05", str, time.Local)
	if err != nil {
		t, err = time.ParseInLocation("2006-01-02", str, time.Local)
		if err != nil {
			return 0
		}
	}
	return t.Unix()
}

// DateToTimeMin 日期天始时间戳（如: 2026-05-28 -> 2026-05-28 00:00:00）
func (datetime *DateTime) DateToTimeMin(str string) int64 {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0
	}
	// 如果前端传了标准日期，直接拼上零点进行解析
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str+" 00:00:00", time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// DateToTimeMax 日期结束时间戳（如: 2026-05-28 -> 2026-05-28 23:59:59）
func (datetime *DateTime) DateToTimeMax(str string) int64 {
	str = strings.TrimSpace(str)
	if str == "" {
		return 0
	}
	t, err := time.ParseInLocation("2006-01-02 15:04:05", str+" 23:59:59", time.Local)
	if err != nil {
		return 0
	}
	return t.Unix()
}

// GetTimestampMills 获取当前秒级以下的毫秒级时间戳（现代 Go 最佳实践）
func (datetime *DateTime) GetTimestampMills() int64 {
	return time.Now().UnixMilli()
}

// CoreTimestampToDate 内部通用转换核心
func (datetime *DateTime) CoreTimestampToDate(any any, layout string) string {
	i := cast.ToInt64(any)
	if i < 1 {
		return ""
	}
	return time.Unix(i, 0).Format(layout)
}

// TimestampToDate 秒级时间戳转 -> "2006-01-02 15:04:05"
func (datetime *DateTime) TimestampToDate(any any) string {
	return datetime.CoreTimestampToDate(any, "2006-01-02 15:04:05")
}

// TimestampToDateMin 秒级时间戳转 -> "2006-01-02"
func (datetime *DateTime) TimestampToDateMin(any any) string {
	return datetime.CoreTimestampToDate(any, "2006-01-02")
}

// TimestampToDateYear 秒级时间戳转 -> "2006"
func (datetime *DateTime) TimestampToDateYear(any any) string {
	return datetime.CoreTimestampToDate(any, "2006")
}

func (datetime *DateTime) TimestampToDateInt(any any) string {
	timeLayout := "20060102" //转化所需模板
	i := cast.ToInt64(any)
	return time.Unix(i, 0).Format(timeLayout)
}
