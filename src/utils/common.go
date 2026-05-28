package utils

import (
	"context"
	"ego/src/boot/global"
	model "ego/src/model/basic"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math/rand"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func LibGetUid(c *gin.Context) int {
	Uid_t, _ := c.Get("User_Id")
	return cast.ToInt(Uid_t)
}

func LibGetUname(c *gin.Context) string {
	Username, _ := c.Get("User_Name")
	return cast.ToString(Username)
}

func LibGetUrole(c *gin.Context) []string {
	User_Role, _ := c.Get("User_Role")
	return cast.ToStringSlice(User_Role)
}

// LogSqlErr 统一记录 SQL 错误日志（自动带 uid / ip / 请求路由）
func LogSqlErr(c *gin.Context, msg string, err error, keys ...zap.Field) {
	// 1. 获取用户ID
	uid := LibGetUid(c)

	// 2. 获取客户端真实IP
	ip := c.ClientIP()

	// 3. 获取当前请求的路由地址（完整URI）
	uri := c.Request.URL.RequestURI()

	// 4. 固定日志字段（每次自动带上）
	fields := []zap.Field{
		zap.Int("uid", uid),
		zap.String("ip", ip),
		zap.String("uri", uri),
		zap.Error(err),
	}

	// 5. 追加自定义字段
	fields = append(fields, keys...)

	// 6. 输出错误日志
	global.C_LOG.Error(msg, fields...)
}

func GetParameter(ctx context.Context, pname string) map[int]interface{} {
	db := GetDB(ctx)

	var one model.SysDict
	db.Where("pname = ?", pname).Find(&one)

	var parval []model.SysDictValue
	db.Where("dict_id = ?", one.ID).Order("id ASC").Find(&parval)

	var arr = make(map[int]interface{})
	for _, value := range parval {
		k := value.Keyid
		arr[k] = value.Defaultval
	}
	return arr
}

func GetParameterStr(ctx context.Context, pname string) map[string]interface{} {
	db := GetDB(ctx)

	var one model.SysDict
	db.Where("pname = ?", pname).Find(&one)

	var parval []model.SysDictValue
	db.Where("dict_id = ?", one.ID).Order("id ASC").Find(&parval)

	var arr = make(map[string]interface{})
	for _, value := range parval {
		k := value.Keystr
		arr[k] = value.Defaultval
	}
	return arr
}

// SetMaptoById 将 []map[string]string 转换为以 id 为 key 的嵌套 map
func SetMaptoById(list []map[string]string) map[string]map[string]string {
	if len(list) == 0 {
		return make(map[string]map[string]string)
	}

	// 预先分配好空间
	result := make(map[string]map[string]string, len(list))
	for _, value := range list {
		if value == nil {
			continue
		}
		// 转换 id，如果不存在或为空，则继续
		idStr := cast.ToString(value["id"])
		result[idStr] = value
	}
	return result
}

// CapitalizeFirst 将字符串首字母转换为大写（高性能安全版）
func CapitalizeFirst(s string) string {
	if s == "" {
		return ""
	}

	var sb strings.Builder
	sb.Grow(len(s))

	for i, r := range s {
		if i == 0 {
			sb.WriteRune(unicode.ToUpper(r))
		} else {
			sb.WriteString(s[i:]) // 剩余部分直接追加，防止循环损耗
			break
		}
	}
	return sb.String()
}

// GetTimestamp 获取当前秒级时间戳
func GetTimestamp() int64 {
	return time.Now().Unix()
}

// RandN 生成 [min, max] 之间的随机整数
func RandN(min, max int) int {
	// 1. 边界安全防御：如果范围逆序，自动纠正，防止 Intn(n <= 0) 导致进程 Panic
	if min > max {
		min, max = max, min
	}
	// 2. 计算区间长度
	diff := max - min + 1
	// 🎯 Go 1.20+ 优化：直接使用 rand.Intn，无需任何 Seed 操作。
	return rand.Intn(diff) + min
}

// GetInterfaceMapByKey 从动态 map 中安全读取 key（防崩溃、高性能版）
func GetInterfaceMapByKey(value any, key any) any {
	if value == nil {
		return ""
	}
	val := reflect.ValueOf(value)
	if val.Kind() != reflect.Map {
		return ""
	}
	mapVal := val.MapIndex(reflect.ValueOf(key))
	if !mapVal.IsValid() {
		return ""
	}

	return mapVal.Interface()
}

// GetMapDataByMore 提取多键值并拼接（预分配优化版）
func GetMapDataByMore(arr map[string]any, s any) string {
	str := cast.ToString(s)
	if str == "" || arr == nil {
		return ""
	}

	tmp := strings.Split(str, ",")
	tmparr := make([]string, 0, len(tmp))
	for _, tv := range tmp {
		if val, exists := arr[tv]; exists {
			tmparr = append(tmparr, cast.ToString(val))
		} else {
			tmparr = append(tmparr, "")
		}
	}
	return strings.Join(tmparr, ",")
}

// GetMapDataByMoreArr 提取多键值并返回数组（预分配优化版）
func GetMapDataByMoreArr(arr map[string]any, s any) []string {
	str := cast.ToString(s)
	if str == "" || arr == nil {
		return []string{}
	}

	tmp := strings.Split(str, ",")
	tmparr := make([]string, 0, len(tmp))
	for _, tv := range tmp {
		if val, exists := arr[tv]; exists {
			tmparr = append(tmparr, cast.ToString(val))
		} else {
			tmparr = append(tmparr, "")
		}
	}
	return tmparr
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, os.ErrNotExist) { // 🎯 现代 Go 推荐的错误包装判定
		return false, nil
	}
	return false, err
}

func CreateDir(dirs ...string) error {
	for _, path := range dirs {
		if path == "" {
			continue
		}

		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			// 妥善记录错误日志
			if global.C_LOG != nil {
				global.C_LOG.Error("failed to create directory",
					zap.String("path", path),
					zap.Error(err),
				)
			}
			return err // 🎯 确保发生错误时能真正返回给上层
		}

		if global.C_LOG != nil {
			global.C_LOG.Debug("directory ready", zap.String("path", path))
		}
	}
	return nil
}

// 根据当前日期创建级联存储目录
func Mkdir(basePath string) (string, error) {
	now := time.Now()
	// 🎯 优化：不硬编码斜杠，通过多参数交由 filepath 适配 Windows/Linux 平台规范
	folderPath := filepath.Join(
		basePath,
		now.Format("2006"),
		now.Format("01"),
		now.Format("02"),
	)

	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return "", err
	}
	return folderPath, nil
}

// 格式化文件大小
func FormatFileSize(fileSize int64) string {
	const minBuf = 1024
	if fileSize < minBuf {
		return fmt.Sprintf("%.2f B", float64(fileSize))
	}

	const (
		_   = iota
		_KB = 1 << (10 * iota) // 1024
		_MB                    // 1048576
		_GB                    // 1073741824
		_TB                    // 1099511627776
		_PB                    // 1125899906842624 (注意：TB 上面是 PB)
	)

	fSize := float64(fileSize)
	switch {
	case fileSize < _MB:
		return fmt.Sprintf("%.2f KB", fSize/float64(_KB))
	case fileSize < _GB:
		return fmt.Sprintf("%.2f MB", fSize/float64(_MB))
	case fileSize < _TB:
		return fmt.Sprintf("%.2f GB", fSize/float64(_GB))
	case fileSize < _PB:
		return fmt.Sprintf("%.2f TB", fSize/float64(_TB))
	default:
		return fmt.Sprintf("%.2f PB", fSize/float64(_PB))
	}
}

func GetDB(ctx context.Context) *gorm.DB {
	return global.C_DB.WithContext(ctx)
}

func GetTabPre() string {
	m := global.C_CONFIG.System
	return m.Tabpre
}

func BoolInarray(a []string, s int) bool {
	d := cast.ToString(s)
	return slices.Contains(a, strings.ToLower(d))
}

func ReturnEgoRootUrl() string {
	return global.C_CONFIG.System.Host + cast.ToString(global.C_CONFIG.System.Addr) + "/"
}

func StringSliceToIntSlice(strSlice []string) []int {
	intSlice := make([]int, 0, len(strSlice))
	for _, s := range strSlice {
		// 忽略空字符串
		if s == "" {
			continue
		}
		// 转换为数字，失败则跳过
		num, err := strconv.Atoi(s)
		if err == nil {
			intSlice = append(intSlice, num)
		}
	}
	return intSlice
}
