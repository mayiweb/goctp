/**
 * 公共函数
 */

package main

import (
    "os"
    "fmt"
    "log"
    "time"
    "strconv"
    "strings"
    "runtime/debug"
    "path/filepath"
    "github.com/axgle/mahonia"
)

// 获得交易所名称
func GetExchangeName(ExchangeID string) string {

    title := ""

    switch ExchangeID {
        case "SHFE":
            title = "上海期货交易所"

        case "CZCE":
            title = "郑州商品交易所"

        case "DCE":
            title = "大连商品交易所"

        case "GFEX":
            title = "广州期货交易所"

        case "CFFEX":
            title = "中国金融期货交易所"

        case "INE":
            title = "上海国际能源交易中心"

        default:
            title = "未知"
    }

    return title
}

// 获得报单多空方向
func GetDirectionTitle(Direction string) string {

    title := ""

    switch Direction {
        case "0":
            title = "买"

        case "1":
            title = "卖"

        default:
            title = "未知"
    }

    return title
}

// 获得持仓多空方向
func GetPosiDirectionTitle(Direction string) string {

    title := ""

    switch Direction {
        case "1":
            title = "净"

        case "2":
            title = "买"

        case "3":
            title = "卖"

        default:
            title = "未知"
    }

    return title
}

// 获得报单状态
func GetOrderStatusTitle(OrderStatus string) string {

    title := ""

    switch OrderStatus {
        case "0":
            title = "已成交"

        // 部分成交还在队列中
        case "1":
            title = "部分成交"

        // 部分成交不在队列中
        case "2":
            title = "部分成交"

        case "3":
            title = "未成交"

        // 未成交不在队列中
        case "4":
            title = "未成交"

        case "5":
            title = "已撤单"

        case "a":
            title = "未知"

        case "b":
            title = "尚未触发"

        case "c":
            title = "已触发"

        default:
            title = "未知状态"
    }

    return title
}

// 获得开平标志
func GetOffsetFlagTitle(OrderStatus string) string {

    title := ""

    switch OrderStatus {
        case "0":
            title = "开仓"

        case "1":
            title = "平仓"

        case "2":
            title = "强平"

        case "3":
            title = "平今"

        case "4":
            title = "平昨"

        case "5":
            title = "强减"

        case "6":
            title = "本地强平"

        default:
            title = "未知"
    }

    return title
}

// 获得投机套保标志
func GetHedgeFlagTitle(HedgeFlag string) string {

    title := ""

    switch HedgeFlag {

        case "1":
            title = "投机"

        case "2":
            title = "套利"

        case "3":
            title = "套保"

        case "5":
            title = "做市商"

        case "6":
            title = "第一腿投机第二腿套保"

        case "7":
            title = "第一腿套保第二腿投机"

        default:
            title = "未知"
    }

    return title
}

// 获得持仓日期类型
func GetPositionDateTitle(PositionDate string) string {

    title := ""

    switch PositionDate {

        case "1":
            title = "今仓"

        case "2":
            title = "昨仓"

        default:
            title = "未知"
    }

    return title
}

// 获得产品类型
func GetProductClassTitle(ProductClass string) string {

    title := ""

    switch ProductClass {

        case "1":
            title = "期货"

        case "2":
            title = "期货期权"

        case "3":
            title = "组合"

        case "4":
            title = "即期"

        case "5":
            title = "期转现"

        case "6":
            title = "现货期权"

        case "7":
            title = "TAS合约"

        case "I":
            title = "金属指数"

        default:
            title = "未知"
    }

    return title
}

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

// 请求日志
func ReqMsg(Msg string) {
    log.Println(Msg)
}

// 请求 api 出现错误
func ReqFailMsg(Msg string, iResult int) {
    log.Printf("%v [%d: %s]\n", Msg, iResult, iResultMsg(iResult))
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

// 判断文件或目录是否存在
func PathExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil {
        return true, nil
    }

    if os.IsNotExist(err) {
        return false, nil
    }

    return false, err
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

// log.Println
func LogPrintln(a ...interface{}) {
    log.Println(a...)
}

// log.Printf
func LogPrintf(format string, a ...interface{}) {
    log.Printf(format, a...)
}

func StrInArray(str string, arr []string) bool {
    for _, v := range arr {
        if str == v {
            return true
        }
    }
    return false
}

/**
 * 字符串后面补位空格
 *
 * 例： ConvertToString("abc", 6)  // 字符串长6位，不足的补空格
 */
func StrAfterSpace(str string, length int) string {
    result := str

    for i := len(str); i < length; i++ {
        result += " "
    }

    return result
}

// 编码转换（gbk 转 utf8）
func GbkToUtf8(text string) string {

    srcCoder := mahonia.NewDecoder("gbk")
    tagCoder := mahonia.NewDecoder("utf-8")
    srcResult := srcCoder.ConvertString(text)

    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

    result := string(cdata)

    return result
}

/**
 * 字符串替换
 *
 * 例： result = StrReplace("文本内容", "被替换字符串", "替换为新的字符串")
 */
func StrReplace(content string, old string, newStr string) string {
    return strings.Replace(content, old, newStr, -1)
}

// 获取当前路径
func GetCurrentDirectory() string {
    dir, err := filepath.Abs(filepath.Dir(os.Args[0]))  //返回绝对路径  filepath.Dir(os.Args[0])去除最后一个元素的路径
    if err != nil {
        log.Fatal(err)
    }
    return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

/**
 * 编码转换
 *
 * 例： result = ConvertToString(text, "gbk", "utf-8")
 */
func ConvertToString(text string, srcCode string, tagCode string) string {

    srcCoder := mahonia.NewDecoder(srcCode)
    tagCoder := mahonia.NewDecoder(tagCode)
    srcResult := srcCoder.ConvertString(text)

    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

    result := string(cdata)

    return result
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