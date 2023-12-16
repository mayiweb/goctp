/**
 * 结构体定义
 */

package main

import (
    "gitee.com/mayiweb/goctp"
)

// Ctp 行情 spi 回调函数
type FtdcMdSpi struct{
    CtpClient
}

// Ctp 交易 spi 回调函数
type FtdcTradeSpi struct{
    CtpClient
}

// Ctp 客户端 行情、交易模块 全局变量
type CtpClient struct {

    // 行情模块 api
    MdApi goctp.CThostFtdcMdApi

    // 交易模块 api
    TradeApi goctp.CThostFtdcTraderApi

    BrokerID    string
    InvestorID  string
    Password    string

    // 客户端认证
    AppID       string
    AuthCode    string

    // 当前交易日期
    TradingDay string

    // 当前交易月份
    TradeMonth string

    // 行情请求编号
    MdRequestId int

    // 交易请求编号
    TradeRequestId int

    // 交易系统是否已经初始化了
    IsTradeInit bool

    // 交易程序是否初始化完成（自动完成如下动作：交易账号登陆、结算单确认、查询合约、查询资金账户、查询用户报单、查询用户持仓 后算完成）
    IsTradeInitFinish bool

    // 交易程序是否已登录过
    IsTradeLogin bool

    // 行情系统是否已经初始化了
    IsMdInit bool

    // 行情程序是否已登录过
    IsMdLogin bool
}

// 深度行情
type TickStruct struct {
    // 交易日
    TradingDay string
    // 合约代码
    InstrumentID string
    // 合约名称
    InstrumentName string
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
    // 申买价二
    BidPrice2 float64
    // 申买量二
    BidVolume2 int
    // 申卖价二
    AskPrice2 float64
    // 申卖量二
    AskVolume2 int
    // 申买价三
    BidPrice3 float64
    // 申买量三
    BidVolume3 int
    // 申卖价三
    AskPrice3 float64
    // 申卖量三
    AskVolume3 int
    // 申买价四
    BidPrice4 float64
    // 申买量四
    BidVolume4 int
    // 申卖价四
    AskPrice4 float64
    // 申卖量四
    AskVolume4 int
    // 申买价五
    BidPrice5 float64
    // 申买量五
    BidVolume5 int
    // 申卖价五
    AskPrice5 float64
    // 申卖量五
    AskVolume5 int
    // 当日均价
    AveragePrice float64
    // 当日均价（可直接使用的价格）
    NewAveragePrice float64
}

// 用户登录成功
type LoginStruct struct {
    // Map 数据 Key 键
    MapKey string
    // 交易日
    TradingDay string
    // 登录成功时间
    LoginTime string
    // 经纪公司代码
    BrokerID string
    // 用户代码
    UserID string
    // 交易系统名称
    SystemName string
    // 前置编号
    FrontID int
    // 会话编号
    SessionID int
    // 最大报单引用
    MaxOrderRef string
    // 上期所时间
    SHFETime string
    // 大商所时间
    DCETime string
    // 郑商所时间
    CZCETime string
    // 中金所时间
    FFEXTime string
    // 能源中心时间
    INETime string
}

// 交易所合约详情
type InstrumentStruct struct {
    // Map 数据 Key 键
    MapKey string
    // 合约代码
    InstrumentID string
    // 交易所代码
    ExchangeID string
    // 合约名称
    InstrumentName string
    // 合约代码
    InstrumentCode string
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

    // 交易所名称
    ExchangeName string
    // 产品类型标题
    ProductClassTitle string
    // 行情波动一点价格
    MovePrice float64
    // 是否主力合约
    IsMainInstrument bool
}

// 资金账户信息
type AccountStruct struct {
    // Map 数据 Key 键
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

// 持仓列表
type PositionStruct struct {
    // key
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
    // 持仓多空方向（接口字段 PosiDirection）
    Direction string
    // 持仓多空方向标题
    DirectionTitle string
    // 开仓成本（接口字段 OpenCost）
    OpenPrice float64
    // 持仓成本
    PositionCost float64
    // 手续费
    Commission float64
    // 总持仓（接口字段 Position）
    Volume int
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

    // 最新价
    LastPrice float64
    // 策略ID
    StrategyId int
    // 策略代码
    StrategyCode string
    // 策略名称
    StrategyName string
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

// 持仓明细
type PositionDetailStruct struct {
    // key
    MapKey string
    // 经纪公司代码
    BrokerID string
    // 投资者代码
    InvestorID string
    // 合约代码
    InstrumentID string
    // 合约名称
    InstrumentName string
    // 投机套保标志
    HedgeFlag string
    // 投机套保标志标题
    HedgeFlagTitle string
    // 买卖
    Direction string
    // 买卖标题
    DirectionTitle string
    // 开仓日期
    OpenDate string
    // 成交编号
    TradeID string
    // 数量
    Volume int
    // Position int
    // 开仓价
    OpenPrice float64
    // 交易日
    TradingDay string
    // 结算价
    SettlementPrice float64
    // 交易所代码
    ExchangeID string
    // 逐日盯市平仓盈亏
    CloseProfitByDate float64
    // 逐笔对冲平仓盈亏
    CloseProfitByTrade float64
    // 逐日盯市持仓盈亏
    PositionProfitByDate float64
    // 逐笔对冲持仓盈亏
    PositionProfitByTrade float64
    // 平仓量
    CloseVolume int
    // 平仓金额
    CloseAmount float64
    // 先开先平剩余数量（DCE）
    TimeFirstVolume int

    // 最新价
    LastPrice float64
    // 盈亏
    Profit float64
    // 持仓日期类型（1：今日持仓，历史持仓）
    PositionDate string
    // 持仓日期类型标题
    PositionDateTitle string
    // 冻结的持仓量
    ShortVolume int
    // 策略ID
    StrategyId int
    // 策略代码
    StrategyCode string
    // 策略名称
    StrategyName string
    // 止损价
    StopPrice float64
    // 持仓时段最高价（多头为最高价，空头为最低价）
    MaxTickPrice float64
    // 移动止损价
    MovingStopPrice float64
    // 移动止损触发开启价差（仅盈利有效）
    MovingPriceDiff float64
    // 移动止损回退平仓价差（最高价回退到价差就平）
    MovingStopPriceDiff float64
}

// 报单列表（已成交、未成交、撤单等状态）的列表数据
type OrderStruct struct {
    // Map 数据 Key 键
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
    // 报单日期
    InsertDate string
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
    // 成交数量
    VolumeTraded int
    // 剩余数量
    VolumeTotal int
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

    // 盈亏
    Profit float64
    // 最新价
    LastPrice float64
}

// 成交
type TradeStruct struct {
    // key
    MapKey string
    // 经纪公司代码
    BrokerID string
    // 投资者代码
    InvestorID string
    // 报单引用
    OrderRef string
    // 用户代码
    UserID string
    // 交易所代码
    ExchangeID string
    // 成交编号
    TradeID string
    // 买卖方向
    Direction string
    // 买卖方标题
    DirectionTitle string
    // 报单编号
    OrderSysID string
    // 会员代码
    ParticipantID string
    // 客户代码
    ClientID string
    // 开平标志
    OffsetFlag string
    // 开平标志标题
    OffsetFlagTitle string
    // 投机套保标志
    HedgeFlag string
    // 投机套保标志标题
    HedgeFlagTitle string
    // 价格
    Price float64
    // 数量
    Volume int
    // 成交时期
    TradeDate string
    // 成交时间
    TradeTime string
    // 成交类型
    TradeType string
    // 成交价来源
    PriceSource string
    // 交易所交易员代码
    TraderID string
    // 本地报单编号
    OrderLocalID string
    // 结算会员编号
    ClearingPartID string
    // 业务单元
    BusinessUnit string
    // 序号
    SequenceNo int
    // 交易日
    TradingDay string
    // 结算编号
    SettlementID int
    // 经纪公司报单编号
    BrokerOrderSeq int
    // 成交来源
    TradeSource string
    // 投资单元代码
    InvestUnitID string
    // 合约代码
    InstrumentID string
    // 合约在交易所的代码
    ExchangeInstID string

    // 合约名称
    InstrumentName string
    // 盈亏
    Profit float64
    // 手续费
    Commission float64
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
    // 报单引用
    OrderRef int
    // 组合开平标志【平仓可以设置】
    CombOffsetFlag byte
    // 开仓策略代码
    StrategyCode string
    // 平仓策略代码
    CloseStrategyCode string
    // 策略消息
    StrategyMsg string
    // 有效秒数【0表示不限制，秒倒计时，到0撤单】
    Second int
    // 时间【回测模式用到了】
    UpdateTime string
    // 持仓均价（平仓时用来计算盈亏）
    OpenPrice float64
    // 股票价（标的价，期权会用到）
    StockPrice float64
    // 持仓key（平仓时会用到）
    PositionKey string
}