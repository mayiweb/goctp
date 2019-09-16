package main

import (
    "fmt"
    "os"
    "log"
    "strconv"
    "strings"
    "path/filepath"
    "github.com/axgle/mahonia"
)

// 交易所合约详情
type InstrumentInfoStruct struct {
    // 合约代码
    InstrumentID string
    // 交易所代码
    ExchangeID string
    // 合约名称
    InstrumentName string
    // 合约在交易所的代码
    ExchangeInstID string
    // 产品代码
    ProductID string
    // 产品类型
    ProductClass string
    // 交割年份
    DeliveryYear int
    // 交割月
    DeliveryMonth int
    // 市价单最大下单量
    MaxMarketOrderVolume int
    // 市价单最小下单量
    MinMarketOrderVolume int
    // 限价单最大下单量
    MaxLimitOrderVolume int
    // 限价单最小下单量
    MinLimitOrderVolume int
    // 合约数量乘数
    VolumeMultiple int
    // 最小变动价位
    PriceTick float64
    // 创建日
    CreateDate string
    // 上市日
    OpenDate string
    // 到期日
    ExpireDate string
    // 开始交割日
    StartDelivDate string
    // 结束交割日
    EndDelivDate string
    // 合约生命周期状态
    InstLifePhase string
    // 当前是否交易
    IsTrading int
    // 持仓类型
    PositionType string
    // 持仓日期类型
    PositionDateType string
    // 多头保证金率
    LongMarginRatio float64
    // 空头保证金率
    ShortMarginRatio float64
    // 是否使用大额单边保证金算法
    MaxMarginSideAlgorithm string
    // 基础商品代码
    UnderlyingInstrID string
    // 执行价
    StrikePrice float64
    // 期权类型
    OptionsType string
    // 合约基础商品乘数
    UnderlyingMultiple float64
    // 组合类型
    CombinationType string
}

// 资金账户信息
type AccountInfoStruct struct {
    // 更新 map 时用到
    MapKey string
    // 经纪公司代码
    BrokerID string
    // 投资者帐号
    AccountID string
    // 上次质押金额
    PreMortgage float64
    // 上次信用额度
    PreCredit float64
    // 上次存款额
    PreDeposit float64
    // 上次结算准备金
    PreBalance float64
    // 上次占用的保证金
    PreMargin float64
    // 利息基数
    InterestBase float64
    // 利息收入
    Interest float64
    // 入金金额
    Deposit float64
    // 出金金额
    Withdraw float64
    // 冻结的保证金
    FrozenMargin float64
    // 冻结的资金
    FrozenCash float64
    // 冻结的手续费
    FrozenCommission float64
    // 当前保证金总额
    CurrMargin float64
    // 资金差额
    CashIn float64
    // 手续费
    Commission float64
    // 平仓盈亏
    CloseProfit float64
    // 持仓盈亏
    PositionProfit float64
    // 期货结算准备金
    Balance float64
    // 可用资金
    Available float64
    // 可取资金
    WithdrawQuota float64
    // 基本准备金
    Reserve float64
    // 交易日
    TradingDay string
    // 结算编号
    SettlementID int
    // 信用额度
    Credit float64
    // 质押金额
    Mortgage float64
    // 交易所保证金
    ExchangeMargin float64
    // 投资者交割保证金
    DeliveryMargin float64
    // 交易所交割保证金
    ExchangeDeliveryMargin float64
    // 保底期货结算准备金
    ReserveBalance float64
    // 币种代码
    CurrencyID string
    // 上次货币质入金额
    PreFundMortgageIn float64
    // 上次货币质出金额
    PreFundMortgageOut float64
    // 货币质入金额
    FundMortgageIn float64
    // 货币质出金额
    FundMortgageOut float64
    // 货币质押余额
    FundMortgageAvailable float64
    // 可质押货币金额
    MortgageableFund float64
    // 特殊产品占用保证金
    SpecProductMargin float64
    // 特殊产品冻结保证金
    SpecProductFrozenMargin float64
    // 特殊产品手续费
    SpecProductCommission float64
    // 特殊产品冻结手续费
    SpecProductFrozenCommission float64
    // 特殊产品持仓盈亏
    SpecProductPositionProfit float64
    // 特殊产品平仓盈亏
    SpecProductCloseProfit float64
    // 根据持仓盈亏算法计算的特殊产品持仓盈亏
    SpecProductPositionProfitByAlg float64
    // 特殊产品交易所保证金
    SpecProductExchangeMargin float64
    // 业务类型
    BizType string
    // 延时换汇冻结金额
    FrozenSwap float64
    // 剩余换汇额度
    RemainSwap float64
}

// 报单列表（已成交、未成交、撤单等状态）的列表数据
type OrderListStruct struct {
    // 更新 map 时用到
    MapKey string
    // 经纪公司代码
    BrokerID string
    // 投资者代码
    InvestorID string
    // 合约代码
    InstrumentID string
    // 合约名称
    InstrumentName string
    // 交易所代码
    ExchangeID string
    // 前置编号
    FrontID int
    // 会话编号
    SessionID int
    // 报单编号
    OrderSysID string
    // 委托时间
    InsertTime string
    // 报单引用
    OrderRef string
    // 买卖方向
    Direction string
    // 组合开平标志
    CombOffsetFlag string
    // 组合投机套保标志
    CombHedgeFlag string
    // 价格
    LimitPrice float64
    // 数量
    Volume int
    // 状态信息
    StatusMsg string
    // 报单状态
    OrderStatus string
    // 买卖方向，中文
    DirectionTitle string
    // 报单状态，中文
    OrderStatusTitle string
    // 投机套保标志
    CombHedgeFlagTitle string
    // 开平标志，中文
    CombOffsetFlagTitle string
}

// 持仓列表
type InvestorPositionStruct struct {
    // 更新 map 时用到
    MapKey string
    // 经纪公司代码
    BrokerID string
    // 投资者帐号
    InvestorID string
    // 合约代码
    InstrumentID string
    // 合约名称
    InstrumentName string
    // 交易所代码
    ExchangeID string
    // 投机套保标志
    HedgeFlag string
    // 投机套保标志标题
    HedgeFlagTitle string
    // 持仓日期类型（1：今日持仓，历史持仓）
    PositionDate string
    // 持仓日期类型标题
    PositionDateTitle string
    // 持仓多空方向
    PosiDirection string
    // 持仓多空方向标题
    PosiDirectionTitle string
    // 开仓成本
    OpenCost float64
    // 持仓成本
    PositionCost float64
    // 手续费
    Commission float64
    // 总持仓
    Position int
    // 上日持仓
    YdPosition int
    // 今日持仓
    TodayPosition int
    // 冻结的持仓量
    ShortVolume int
    // 多头冻结
    LongFrozen int
    // 空头冻结
    ShortFrozen int
    // 开仓量
    OpenVolume int
    // 平仓量
    CloseVolume int
    // 平仓盈亏
    CloseProfit float64
    // 持仓盈亏
    PositionProfit float64
    // 上次结算价
    PreSettlementPrice float64
    // 本次结算价
    SettlementPrice float64
    // 结算编号
    SettlementID int

    // 止损价
    StopPrice float64
    // 持仓时段最高价
    MaxTickPrice float64
    // 移动止损价
    MovingStopPrice float64
    // 移动止损触发开启价差（仅盈利有效）
    MovingPriceDiff float64
    // 移动止损回退平仓价差（最高价回退到价差就平）
    MovingStopPriceDiff float64
}

// 深度行情
type MarketDataStruct struct {
    // 交易日
    TradingDay string
    // 合约代码
    InstrumentID string
    // 合约code【合约代码字母部分，非官方字段】
    InstrumentCode string
    // 交易所代码
    ExchangeID string
    // 最新价
    LastPrice float64
    // 上次结算价
    PreSettlementPrice float64
    // 昨收盘
    PreClosePrice float64
    // 昨持仓量
    PreOpenInterest float64
    // 今开盘
    OpenPrice float64
    // 最高价
    HighestPrice float64
    // 最低价
    LowestPrice float64
    // 数量
    Volume int
    // 成交金额
    Turnover float64
    // 持仓量
    OpenInterest float64
    // 今收盘
    ClosePrice float64
    // 本次结算价
    SettlementPrice float64
    // 涨停板价
    UpperLimitPrice float64
    // 跌停板价
    LowerLimitPrice float64
    // 最后修改时间
    UpdateTime string
    // 最后修改毫秒
    UpdateMillisec int
    // 申买价一
    BidPrice1 float64
    // 申买量一
    BidVolume1 int
    // 申卖价一
    AskPrice1 float64
    // 申卖量一
    AskVolume1 int
    // 当日均价
    AveragePrice float64
}

// 昨日收盘数据
type LastTimeMarketStruct struct {
    // id
    id int
    // 交易日期
    TradingDay string
    // 合约代码
    InstrumentID string
    // 昨收盘
    PreClosePrice float64
    // 最新价
    LastPrice float64
    // 今开盘
    OpenPrice float64
    // 今收盘
    ClosePrice float64
    // 最低价
    LowestPrice float64
    // 最高价
    HighestPrice float64
    // 跌停板价
    LowerLimitPrice float64
    // 涨停板价
    UpperLimitPrice float64
    // 持仓量
    OpenInterest int
    // 成交数量
    Volume int
    // 创建时间
    CreateTime string
    // 涨跌方向
    Direction string
    // 昨收与今收的价差
    PriceDiff float64
}

// 输入报单
type InputOrderStruct struct {
    // 合约代码
    InstrumentID string
    // 买卖方向【0：买，1：卖】
    Direction byte
    // 价格
    Price float64
    // 数量
    Volume int
    // 组合开平标志【平仓可以设置】
    CombOffsetFlag byte
}

// 获得报单多空方向
func GetDirectionTitle(Direction string) string {
    var title string

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
func GetPosiDirectionTitle(PosiDirection string) string {

    title := ""

    switch PosiDirection {
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

        case "1":
            title = "部分成交还在队列中"

        case "2":
            title = "部分成交不在队列中"

        case "3":
            title = "未成交"

        case "4":
            title = "未成交不在队列中"

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

// 请求日志
func ReqMsg(Msg string) {
    log.Println(Msg)
}

// 请求 api 出现错误
func ReqFailMsg(Msg string, iResult int) {
    fmt.Printf("%v [%d: %s]\n", Msg, iResult, iResultMsg(iResult));
}

// 请求失败的错误码对应消息
func iResultMsg(iResult int) (string) {

    msg := ""

    switch (iResult) {
        case 0:
            msg = "成功";
            break;

        case -1:
            msg = "请检查账号是否已经登陆";
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

// 检查错误，有就抛出
func CheckErr(err error) {
    if err != nil {
        panic(err)
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
    value, _ := strconv.ParseFloat(fmt.Sprintf("%." + strconv.Itoa(n) + "f", f), 64)
    return value
}

// int 转 string
func IntToString (i int) string {
    return strconv.Itoa(i)
}

// float64 转 string
func Float64ToString (f float64) string {
    return strconv.FormatFloat(f, 'f', 2, 64)
}

// string 转 float64
func StringToFloat64 (str string) float64 {
    f64, _ := strconv.ParseFloat(str, 64)
    return f64
}

// string 转 int
func StringToInt (str string) int {
    num, _ := strconv.Atoi(str)
    return num
}

// 去掉左右两边空格
func TrimSpace (str string) string {
    return strings.TrimSpace(str)
}

// fmt.Println
func Println(a ...interface{}) (n int, err error) {
    return fmt.Println(a...)
}


func StrInArray(str string, arr []string) bool {
    for _, v := range arr {
        if str == v {
            return true
        }
    }
    return false
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

    srcCoder  := mahonia.NewDecoder(srcCode)
    tagCoder  := mahonia.NewDecoder(tagCode)
    srcResult := srcCoder.ConvertString(text)

    _, cdata, _ := tagCoder.Translate([]byte(srcResult), true)

    result := string(cdata)

    return result
}