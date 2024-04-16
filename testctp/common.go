/**
 * 公共函数
 */

package testctp

import (
    "fmt"
    "os"
    "log"
    "time"
    "sort"
    "strconv"
    "strings"
    "runtime"
    "runtime/debug"
    "encoding/json"
)

// 是否空指针
func IsNullPointer(p interface{}) bool {

    if p == nil {
        return true
    }

    pv := Sprintf("%v", p)
    if pv == "0" {
        return true
    }

    return false
}

// 退出程序
func Exit() {
    os.Exit(1)
}

// 获取版本信息
func GetVersion() string {
    return runtime.Version()
}

// 获取字符串信息
func GetString(str interface{}) string {
    if str == nil {
        return ""
    }

    return Sprintf("%v", str)
}

// 获得当前日期
func GetCurrentDate() string {
    return string(time.Now().Format("20060102"))
}

// 获得当前时间
func GetCurrentTime() string {
    return string(time.Now().Format("2006-01-02 15:04:05"))
}

// 获得当前时间戳
func GetCurrentUnix() int {
    return int(time.Now().Unix())
}

// 日期转时间（例： 2020-06-29 08:00:00 转为 08:00:00）
func DateToTime(date string) string {
    loc, _      := time.LoadLocation("Local")
    theTime1, _ := time.ParseInLocation("2006-01-02 15:04:05", date, loc)

    return time.Unix(theTime1.Unix(), 0).Format("15:04:05")
}


// 获取日期中的月份（例： 20200629 转为 6）
func GetDateMonth(date string) string {
    loc, _      := time.LoadLocation("Local")
    theTime1, _ := time.ParseInLocation("20060102", date, loc)

    return time.Unix(theTime1.Unix(), 0).Format("1")
}

/**
 * 时间转时间戳
 *
 * 例：TimeToUnix("09:30:00")
 */
func TimeToUnix(strTime string) int {

    // 如果就写了个时间没写日期，视为当天日期
    if len(strTime) == 8 {
        strTime = time.Now().Format("2006-01-02") + " " + strTime
    }

    loc, _      := time.LoadLocation("Local")    // 获取时区
    theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", strTime, loc)

    return int(theTime.Unix())
}

/**
 * 时间戳转时间
 *
 * 例：UnixToTime(1686580152)
 */
func UnixToTime(TimeStamp int) string {
    return time.Unix(int64(TimeStamp), 0).Format("15:04:05")
}

/**
 * 时间戳转日期
 *
 * 例：UnixToTime(1686580152, "2006-01-02 15:04:05")
 */
func UnixToDateTime(TimeStamp int, Format string) string {
    return time.Unix(int64(TimeStamp), 0).Format(Format)
}

/**
 * 时间转时间戳（注意格式）
 *
 * 例：DateTimeToUnix("2006-01-02 15:04:05")
 */
func DateTimeToUnix(strTime string) int {

    loc, _     := time.LoadLocation("Local")    // 获取时区
    theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", strTime, loc)

    return int(theTime.Unix())
}

/**
 * 日期转时间戳（注意格式）
 *
 * 例：DateToUnix("2006-01-02")
 */
func DateToUnix(strTime string) int {

    loc, _     := time.LoadLocation("Local")    // 获取时区
    theTime, _ := time.ParseInLocation("2006-01-02", strTime, loc)

    return int(theTime.Unix())
}

// 请求日志
func ReqMsg(Msg string) {
    log.Println(Msg)
}

// 请求 api 出现错误
func ReqFailMsg(Msg string, iResult int) string {
    str := Sprintf("%v [%d: %s]", Msg, iResult, iResultMsg(iResult))

    log.Println(str)

    return str
}

// 请求失败的错误码对应消息
func iResultMsg(iResult int) (string) {

    msg := ""

    switch iResult {
        case 0:
            msg = "成功";
            break;

        case -1:
            msg = "请检查账号是否已经登录";
            break;

        case -2:
            msg = "未处理请求超过许可数";
            break;

        case -3:
            msg = "每秒发送请求数超过许可数";
            break;

        default:
            msg = "未知错误";
            break;
    }

    return msg;
}

// 打印错误（存在错误的时候）
func PrintErr(err error) {
    if err != nil {
        Println(err)
    }
}

// 检查错误，有就抛出 panic
func CheckErr(err error) {
    if err != nil {
        panic(err)
    }
}

// 处理 panic
func CheckPanic() {
    if err := recover(); err != nil {

        ErrStr := "-------------------------------------------------------------------------------------------------\n" +
                "- 错误提示\n" +
                "- 出错时间：" + GetCurrentTime() + "\n" +
                "- 错误内容：" + fmt.Sprintf("%s", err) + "\n\n" +
                string(debug.Stack()) +
                "-------------------------------------------------------------------------------------------------"

        Println(ErrStr)
    }
}

// float64 保留几位小数点
func Decimal(f float64, n int) float64 {

    defer CheckPanic()

    if f == 1.7976931348623157e+308 {
        f = 0.0000
    }

    value, _ := strconv.ParseFloat(fmt.Sprintf("%." + strconv.Itoa(n) + "f", f), 64)
    return value
}

// int 转 string
func IntToString(i int) string {
    return strconv.Itoa(i)
}

// float64 转 string
func Float64ToString(f float64) string {
    return strconv.FormatFloat(f, 'f', 2, 64)
}

// string 转 float64
func StringToFloat64(str string) float64 {
    f64, _ := strconv.ParseFloat(str, 64)
    return f64
}

// string 转 int
func StringToInt(str string) int {
    num, _ := strconv.Atoi(str)
    return num
}

// string 转 int32
func StringToInt32(str string) int32 {
    num, err := strconv.ParseInt(str, 10, 32)
    if err != nil {
        fmt.Println(str, err)
    }
    return int32(num)
}

// int32 转 string
func Int32ToString(n int32) string {
    buf := [11]byte{}
    pos := len(buf)
    i := int64(n)
    signed := i < 0
    if signed {
        i = -i
    }
    for {
        pos--
        buf[pos], i = '0'+byte(i%10), i/10
        if i == 0 {
            if signed {
                pos--
                buf[pos] = '-'
            }
            return string(buf[pos:])
        }
    }
}

// 去掉左右两边空格
func TrimSpace(str string) string {
    return strings.TrimSpace(str)
}

// fmt.Println
func Println(a ...interface{}) (n int, err error) {
    return fmt.Println(a...)
}

// fmt.Sprintf
func Sprintf(format string, a ...interface{}) string {
    return fmt.Sprintf(format, a...)
}

// fmt.Printf
func Printf(format string, a ...interface{}) (n int, err error) {
    return fmt.Printf(format, a...)
}

/**
 * 字符串替换
 *
 * 例： result = StrReplace("文本内容", "被替换字符串", "替换为新的字符串")
 */
func StrReplace(content string, old string, newStr string) string {
    return strings.Replace(content, old, newStr, -1)
}

// JsonEncode
func JsonEncode(str interface{}) string {
    jsonBytes, err := json.Marshal(str)
    if err != nil {
        fmt.Println(err)
    }

    return string(jsonBytes)
}

// JsonDecode
func JsonDecode(str string, data interface{}) error {
    return json.Unmarshal([]byte(str), data)
}

/**
 * 日期格式化
 *
 * 例：DateFormat("2020-01-23", "2006-01-02", "20060102")
 *     DateFormat("2020-01-23 14:05:06", "2006-01-02 15:04:05", "20060102150405")
 */
func DateFormat(date string, inputFormat string, outputFormat string) string {
    theTime, _ := time.ParseInLocation(inputFormat, date, time.Local)
    return time.Unix(theTime.Unix(), 0).Format(outputFormat)
}

// 获得 Map key 键的 正序（sort.Ints）
func GetMapIntsKeys(m map[int]interface{}) []int {

    ints := make([]int, 0)

    for k, _ := range m {
        ints = append(ints, k)
    }

    sort.Ints(ints)

    return ints
}

// 获得 Map key 键的 倒序（sort.Sort(sort.Reverse(sort.IntSlice(s)))）
func GetMapIntrsKeys(m map[int]interface{}) []int {

    ints := make([]int, 0)

    for k, _ := range m {
        ints = append(ints, k)
    }

    sort.Sort(sort.Reverse(sort.IntSlice(ints)))

    return ints
}

// 休眠（毫秒，1000 毫秒等于 1 秒）
func Sleep(number int) {
    time.Sleep(time.Millisecond * time.Duration(number))
}

// 将字符串转换成大写
func ToUpper(str string) string {
    return strings.ToUpper(str)
}

// 将字符串转换成小写
func ToLower(str string) string {
    return strings.ToLower(str)
}