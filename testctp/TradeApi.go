package testctp

import (
    "gitee.com/mayiweb/goctp"
    "gitee.com/mayiweb/goctp/safe"
    "time"
    "strconv"
)

var (
    // 原始指针 Ctp.TradeApi = new(testctp.TradeApi)
    OldThis *TradeApi

    // 交易所合约详情列表 InstrumentStruct
    MapInstruments safe.Map

    // 资金账户信息 AccountStruct
    MapAccounts safe.Map

    // 报单列表（已成交、未成交、撤单等状态）的列表数据 OrderStruct
    MapOrders safe.Map

    // 报单成交列表
    MapTrades safe.Map

    // 持仓列表 PositionStruct
    MapPositions safe.Map

    // 持仓明细列表 PositionDetailStruct
    MapPositionDetails safe.Map

    // OrderRef 对应的 tick 数据
    MapOrderRefTick safe.Map
)

type TradeApi struct {
    goctp.CThostFtdcTraderApi
    TradeSpi goctp.CThostFtdcTraderSpi

    SessionID int

    BrokerID string

    InvestorID string

    // 报单编号
    OrderSysID int

    // 成交编号
    TradeID int

    TradingDay int

    // 1手保证金
    Margin float64

    // 1手手续费
    Commission float64
}

// 响应信息
func (p *TradeApi) GetRspInfo() goctp.CThostFtdcRspInfoField {
    pRspInfo := goctp.NewCThostFtdcRspInfoField()
    pRspInfo.SetErrorID(0)
    pRspInfo.SetErrorMsg("ok")

    return pRspInfo
}

// 原始指针 Ctp.TradeApi = new(testctp.TradeApi)
func (p *TradeApi) OldThis() *TradeApi {
    return OldThis
}

// 获取报单编号
func (p *TradeApi) GetOrderSysID() string {
    p.OrderSysID += 1
    return strconv.Itoa(p.OrderSysID)
}

// 获取成交编号
func (p *TradeApi) GetTradeID() string {
    p.TradeID += 1
    return strconv.Itoa(p.TradeID)
}

// 设置资金账户信息
func (p *TradeApi) SetAccount(BrokerID string, InvestorID string) {

    var sAccount AccountStruct

    // 初始资金
    InitialFunding := 10000000.00

    sAccount.BrokerID         = BrokerID
    sAccount.AccountID        = InvestorID
    sAccount.PreBalance       = InitialFunding
    sAccount.Deposit          = 0.00
    sAccount.Withdraw         = 0.00
    sAccount.FrozenMargin     = 0.00
    sAccount.FrozenCash       = 0.00
    sAccount.FrozenCommission = 0.00
    sAccount.CurrMargin       = 0.00
    sAccount.CashIn           = 0.00
    sAccount.Commission       = 0.00
    sAccount.CloseProfit      = 0.00
    sAccount.PositionProfit   = 0.00
    sAccount.Balance          = InitialFunding
    sAccount.Available        = InitialFunding
    sAccount.WithdrawQuota    = InitialFunding
    sAccount.ExchangeMargin   = 0.00
    sAccount.MapKey           = Sprintf("%v_%v", sAccount.BrokerID, sAccount.AccountID)

    MapAccounts.Set(sAccount.MapKey, sAccount)
}

// 重置数据
func  (p *TradeApi) Reset() {

    MapAccounts.Clear()
    MapOrders.Clear()
    MapTrades.Clear()
    MapPositions.Clear()
    MapPositionDetails.Clear()
    MapOrderRefTick.Clear()

    if p.OldThis().BrokerID == "" {
        p.OldThis().BrokerID = "8888"
    }

    if p.OldThis().InvestorID == "" {
        p.OldThis().InvestorID = "10001"
    }

    p.OldThis().SessionID  = 888888
    p.OldThis().OrderSysID = 1000000
    p.OldThis().TradeID    = 5000000
    p.OldThis().Margin     = 1000
    p.OldThis().Commission = 1
    p.OldThis().TradingDay = StringToInt(time.Now().Format("20060102"))

    p.OldThis().SetAccount(p.OldThis().BrokerID, p.OldThis().InvestorID)
}

// 获取API的版本信息
func (p *TradeApi) GetApiVersion() string {
    return "v6.6.5 Testing TradeApi"
}

// 删除接口对象本身
func (p *TradeApi) Release() {

}

// 初始化
func (p *TradeApi) Init() {

    // 将初始化的指针记录下来
    OldThis = p

    go p.TradeSpi.OnFrontConnected()

    p.Reset()
}

// 等待接口线程结束运行
func (p *TradeApi) Join() int {
    // 线程不能退出
    for {
        time.Sleep(time.Millisecond * time.Duration(1000))
    }

    return 0
}

// 获取当前交易日
func (p *TradeApi) GetTradingDay() string {
    return IntToString(p.OldThis().TradingDay)
}

// 注册前置机网络地址
func (p *TradeApi) RegisterFront(pszFrontAddress string) {

}

// 注册名字服务器网络地址
func (p *TradeApi) RegisterNameServer(pszNsAddress string) {

}

// 注册名字服务器用户信息
func (p *TradeApi) RegisterFensUserInfo(pFensUserInfo goctp.CThostFtdcFensUserInfoField) {

}

// 注册回调接口
func (p *TradeApi) RegisterSpi(pSpi goctp.CThostFtdcTraderSpi) {
    p.TradeSpi = pSpi
}

// 订阅私有流。
func (p *TradeApi) SubscribePrivateTopic(nResumeType goctp.THOST_TE_RESUME_TYPE) {

}

// 订阅公共流。
func (p *TradeApi) SubscribePublicTopic(nResumeType goctp.THOST_TE_RESUME_TYPE) {

}

// 客户端认证请求
func (p *TradeApi) ReqAuthenticate(pReqAuthenticateField goctp.CThostFtdcReqAuthenticateField, nRequestID int) int {

    pResult := goctp.NewCThostFtdcRspAuthenticateField()

    pResult = nil

    p.TradeSpi.OnRspAuthenticate(pResult, p.GetRspInfo(), nRequestID, true)

    return 0
}

// 注册用户终端信息，用于中继服务器多连接模式
// 需要在终端认证成功后，用户登录前调用该接口
func (p *TradeApi) RegisterUserSystemInfo(pUserSystemInfo goctp.CThostFtdcUserSystemInfoField) int {
    return 0
}

// 上报用户终端信息，用于中继服务器操作员登录模式
// 操作员登录后，可以多次调用该接口上报客户信息
func (p *TradeApi) SubmitUserSystemInfo(pUserSystemInfo goctp.CThostFtdcUserSystemInfoField) int {
    return 0
}

// 用户登录请求
func (p *TradeApi) ReqUserLogin(pReqUserLoginField goctp.CThostFtdcReqUserLoginField, nRequestID int) int {

    BrokerID   := pReqUserLoginField.GetBrokerID()
    InvestorID := pReqUserLoginField.GetUserID()

    if BrokerID != "" {
        p.BrokerID = BrokerID
    }

    if InvestorID != "" {
        p.InvestorID = InvestorID
    }

    // 设置资金账户
    p.SetAccount(p.BrokerID, p.InvestorID)

    pResult := goctp.NewCThostFtdcRspUserLoginField()
    pResult.SetUserID(p.InvestorID)
    pResult.SetBrokerID(p.BrokerID)
    pResult.SetTradingDay(GetCurrentDate())
    pResult.SetSessionID(p.SessionID)
    pResult.SetLoginTime(time.Now().Format("15:04:05"))

    p.TradeSpi.OnRspUserLogin(pResult, p.GetRspInfo(), nRequestID, true)

    return 0
}

// 登出请求
func (p *TradeApi) ReqUserLogout(pUserLogout goctp.CThostFtdcUserLogoutField, nRequestID int) int {
    return 0
}

// 用户口令更新请求
func (p *TradeApi) ReqUserPasswordUpdate(pUserPasswordUpdate goctp.CThostFtdcUserPasswordUpdateField, nRequestID int) int {
    return 0
}

// 资金账户口令更新请求
func (p *TradeApi) ReqTradingAccountPasswordUpdate(pTradingAccountPasswordUpdate goctp.CThostFtdcTradingAccountPasswordUpdateField, nRequestID int) int {
    return 0
}

// 查询用户当前支持的认证模式
func (p *TradeApi) ReqUserAuthMethod(pReqUserAuthMethod goctp.CThostFtdcReqUserAuthMethodField, nRequestID int) int {
    return 0
}

// 用户发出获取图形验证码请求
func (p *TradeApi) ReqGenUserCaptcha(pReqGenUserCaptcha goctp.CThostFtdcReqGenUserCaptchaField, nRequestID int) int {
    return 0
}

// 用户发出获取短信验证码请求
func (p *TradeApi) ReqGenUserText(pReqGenUserText goctp.CThostFtdcReqGenUserTextField, nRequestID int) int {
    return 0
}

// 用户发出带有图片验证码的登陆请求
func (p *TradeApi) ReqUserLoginWithCaptcha(pReqUserLoginWithCaptcha goctp.CThostFtdcReqUserLoginWithCaptchaField, nRequestID int) int {
    return 0
}

// 用户发出带有短信验证码的登陆请求
func (p *TradeApi) ReqUserLoginWithText(pReqUserLoginWithText goctp.CThostFtdcReqUserLoginWithTextField, nRequestID int) int {
    return 0
}

// 用户发出带有动态口令的登陆请求
func (p *TradeApi) ReqUserLoginWithOTP(pReqUserLoginWithOTP goctp.CThostFtdcReqUserLoginWithOTPField, nRequestID int) int {
    return 0
}

// 报单录入请求
func (p *TradeApi) ReqOrderInsert(pInputOrder goctp.CThostFtdcInputOrderField, nRequestID int) int {

    sTick := TickStruct{}

    TradingDay := Sprintf("%v", p.TradingDay)
    UpdateTime := Sprintf("%v", time.Now().Format("15:04:05"))

    // 报单价格（如果传递了 OrderRefTick 使用对手价，否则使用设置的价格）
    LimitPrice := pInputOrder.GetLimitPrice()

    mOrderRefTickVal, mOrderRefTickOk := MapOrderRefTick.Get(StringToInt(pInputOrder.GetOrderRef()))
    if mOrderRefTickOk {
        sTick = mOrderRefTickVal.(TickStruct)

        if sTick.TradingDay != "" {
            TradingDay = sTick.TradingDay
        }

        if sTick.UpdateTime != "" {
            UpdateTime = sTick.UpdateTime
        }

        // 如果传递了 OrderRefTick 使用对手价，否则使用设置的价格
        if string(pInputOrder.GetDirection()) == "0" {
            if sTick.AskPrice1 > 0 {
                LimitPrice = sTick.AskPrice1
            }
        } else {
            if sTick.BidPrice1 > 0 {
                LimitPrice = sTick.BidPrice1
            }
        }
    }

    // 可平数量
    AvailableVolume := 0

    mPositionDetails := MapPositionDetails.GetAll()

    // 平仓时要反过来取
    CloseDirection := "1"
    if string(pInputOrder.GetDirection()) == "1" {
        CloseDirection = "0"
    }

    // 可平仓数量验证
    if string(pInputOrder.GetCombOffsetFlag()) != "0" {

        for _, v := range mPositionDetails {

            val := v.(PositionDetailStruct)

            // 同一个合约
            if val.InstrumentID == pInputOrder.GetInstrumentID() && val.Direction == CloseDirection && val.Volume > 0 {
                AvailableVolume += val.Volume
            }
        }

        // 可平仓数量不足
        if AvailableVolume < pInputOrder.GetVolumeTotalOriginal() {

            pRspInfo := goctp.NewCThostFtdcRspInfoField()
            pRspInfo.SetErrorID(30)
            // 可平仓数量不足  这边要写 gbk 的中文编码，
            // pRspInfo.SetErrorMsg("CTP:平仓量超过持仓量")
            pRspInfo.SetErrorMsg("CTP: Closing position exceeds holding position")

            pOrderAction := goctp.NewCThostFtdcOrderActionField()
            pOrderAction.SetBrokerID(pInputOrder.GetBrokerID())
            pOrderAction.SetInvestorID(pInputOrder.GetInvestorID())
            pOrderAction.SetInstrumentID(pInputOrder.GetInstrumentID())
            pOrderAction.SetExchangeID(pInputOrder.GetExchangeID())
            pOrderAction.SetOrderRef(pInputOrder.GetOrderRef())
            pOrderAction.SetOrderSysID(p.GetOrderSysID())
            pOrderAction.SetSessionID(p.SessionID)

            // 报单操作错误回报
            p.TradeSpi.OnErrRtnOrderAction(pOrderAction, pRspInfo)

            return 0
        }
    }

    var sOrder OrderStruct

    sOrder.BrokerID            = pInputOrder.GetBrokerID()
    sOrder.InvestorID          = pInputOrder.GetInvestorID()
    sOrder.InstrumentID        = pInputOrder.GetInstrumentID()
    sOrder.ExchangeID          = pInputOrder.GetExchangeID()
    sOrder.OrderRef            = pInputOrder.GetOrderRef()
    sOrder.SessionID           = p.SessionID
    sOrder.InsertTime          = UpdateTime
    sOrder.InsertDate          = TradingDay
    sOrder.OrderSysID          = p.GetOrderSysID()
    sOrder.LimitPrice          = LimitPrice
    sOrder.Volume              = pInputOrder.GetVolumeTotalOriginal()
    sOrder.VolumeTraded        = sOrder.Volume
    sOrder.VolumeTotal         = sOrder.Volume - sOrder.VolumeTraded
    sOrder.Direction           = string(pInputOrder.GetDirection())
    sOrder.CombOffsetFlag      = string(pInputOrder.GetCombOffsetFlag())
    sOrder.CombHedgeFlag       = string(pInputOrder.GetCombHedgeFlag())
    sOrder.OrderStatus         = "0"
    sOrder.MapKey              = Sprintf("%v_%v", sOrder.InstrumentID, sOrder.OrderSysID)

    MapOrders.Set(sOrder.MapKey, sOrder)


    // 发送报单成交通知
    pOrder := goctp.NewCThostFtdcOrderField()
    pOrder.SetBrokerID(sOrder.BrokerID)
    pOrder.SetInvestorID(sOrder.InvestorID)
    pOrder.SetInstrumentID(sOrder.InstrumentID)
    pOrder.SetExchangeID(sOrder.ExchangeID)
    pOrder.SetFrontID(sOrder.FrontID)
    pOrder.SetOrderRef(sOrder.OrderRef)
    pOrder.SetSessionID(sOrder.SessionID)
    pOrder.SetInsertTime(UpdateTime)
    pOrder.SetInsertDate(sOrder.InsertDate)
    pOrder.SetOrderSysID(sOrder.OrderSysID)
    pOrder.SetLimitPrice(sOrder.LimitPrice)
    pOrder.SetVolumeTotalOriginal(sOrder.Volume)
    pOrder.SetVolumeTraded(sOrder.VolumeTraded)
    pOrder.SetVolumeTotal(sOrder.VolumeTotal)
    pOrder.SetDirection(pInputOrder.GetDirection())
    pOrder.SetCombOffsetFlag(sOrder.CombOffsetFlag)
    pOrder.SetCombHedgeFlag(sOrder.CombHedgeFlag)
    pOrder.SetOrderStatus(goctp.THOST_FTDC_OST_AllTraded)
    pOrder.SetStatusMsg(sOrder.StatusMsg)



    var sTrade TradeStruct

    sTrade.BrokerID        = sOrder.BrokerID
    sTrade.InvestorID      = sOrder.InvestorID
    sTrade.OrderRef        = sOrder.OrderRef
    sTrade.UserID          = sOrder.InvestorID
    sTrade.ExchangeID      = sOrder.ExchangeID
    sTrade.TradeID         = p.GetTradeID()
    sTrade.Direction       = sOrder.Direction
    sTrade.OrderSysID      = sOrder.OrderSysID
    sTrade.OffsetFlag      = string(pInputOrder.GetCombOffsetFlag())
    sTrade.HedgeFlag       = string(pInputOrder.GetCombHedgeFlag())
    sTrade.Price           = sOrder.LimitPrice
    sTrade.Volume          = sOrder.Volume
    sTrade.TradeDate       = TradingDay
    sTrade.TradeTime       = UpdateTime
    sTrade.TradeType       = string(goctp.THOST_FTDC_TRDT_Common)
    sTrade.TradingDay      = TradingDay
    sTrade.TradeSource     = string(goctp.THOST_FTDC_TSRC_NORMAL)
    sTrade.InstrumentID    = sOrder.InstrumentID
    sTrade.MapKey          = Sprintf("%v_%v", sTrade.TradeID, sTrade.InstrumentID)

    MapTrades.Set(sTrade.MapKey, sTrade)

    // 发送成交通知
    pTrade := goctp.NewCThostFtdcTradeField()
    pTrade.SetBrokerID(sTrade.BrokerID)
    pTrade.SetInvestorID(sTrade.InvestorID)
    pTrade.SetInstrumentID(sTrade.InstrumentID)
    pTrade.SetOrderRef(sTrade.OrderRef)
    pTrade.SetUserID(sTrade.UserID)
    pTrade.SetExchangeID(sTrade.ExchangeID)
    pTrade.SetTradeID(sTrade.TradeID)
    pTrade.SetDirection(pInputOrder.GetDirection())
    pTrade.SetOrderSysID(sTrade.OrderSysID)
    pTrade.SetOffsetFlag([]byte(sOrder.CombOffsetFlag)[0])
    pTrade.SetHedgeFlag([]byte(sOrder.CombHedgeFlag)[0])
    pTrade.SetPrice(sTrade.Price)
    pTrade.SetVolume(sTrade.Volume)
    pTrade.SetTradeDate(sTrade.TradeDate)
    pTrade.SetTradeTime(sTrade.TradeTime)
    pTrade.SetTradeType(goctp.THOST_FTDC_TRDT_Common)
    pTrade.SetTradingDay(TradingDay)
    pTrade.SetTradeSource(goctp.THOST_FTDC_TSRC_NORMAL)



    var sPosition PositionDetailStruct

    sPosition.BrokerID           = sTrade.BrokerID
    sPosition.InvestorID         = sTrade.InvestorID
    sPosition.InstrumentID       = sTrade.InstrumentID
    sPosition.HedgeFlag          = sTrade.HedgeFlag
    sPosition.Direction          = sTrade.Direction
    sPosition.OpenDate           = TradingDay
    sPosition.TradeID            = sTrade.TradeID
    sPosition.Volume             = sTrade.Volume
    sPosition.OpenPrice          = sTrade.Price
    sPosition.TradingDay         = TradingDay
    sPosition.ExchangeID         = sTrade.ExchangeID
    sPosition.CloseProfitByTrade = 0.00
    sPosition.MapKey             = Sprintf("%v", sTrade.TradeID)


    // 平仓状态（默认为成功，如果未成功则不发送成交通知）
    ClosePositionStatus := true

    // 开仓
    if sTrade.OffsetFlag == "0" {

        sPosition.CloseVolume = 0

        // 记录持仓明细数据
        MapPositionDetails.Set(sPosition.MapKey, sPosition)

    // 平仓
    } else {

        // 平仓数量
        PositionDetailVolume := sTrade.Volume

        for _, v := range mPositionDetails {

            val := v.(PositionDetailStruct)

            // 同一个合约
            if val.InstrumentID == sTrade.InstrumentID && val.Direction == CloseDirection && PositionDetailVolume > 0 {

                if PositionDetailVolume <= val.Volume {
                    val.Volume             -= PositionDetailVolume
                    val.CloseVolume        += PositionDetailVolume
                    val.CloseProfitByTrade += GetPositionProfit(val.InstrumentID, val.OpenPrice, sTrade.Price, PositionDetailVolume, val.Direction)
                    PositionDetailVolume   = 0
                } else {
                    val.Volume              = 0
                    val.CloseVolume        += val.Volume
                    val.CloseProfitByTrade += GetPositionProfit(val.InstrumentID, val.OpenPrice, sTrade.Price, val.Volume, val.Direction)
                    PositionDetailVolume   -= val.Volume
                }

                // 更新持仓数据
                MapPositionDetails.Set(val.MapKey, val)
            }
        }

        // 有未执行平仓的数据
        if PositionDetailVolume != 0 {
            ClosePositionStatus = false
        }
    }

    if ClosePositionStatus {

    }

        p.TradeSpi.OnRtnOrder(pOrder)

        p.TradeSpi.OnRtnTrade(pTrade)

    return 0
}

// 预埋单录入请求
func (p *TradeApi) ReqParkedOrderInsert(pParkedOrder goctp.CThostFtdcParkedOrderField, nRequestID int) int {
    return 0
}

// 预埋撤单录入请求
func (p *TradeApi) ReqParkedOrderAction(pParkedOrderAction goctp.CThostFtdcParkedOrderActionField, nRequestID int) int {
    return 0
}

// 报单操作请求
func (p *TradeApi) ReqOrderAction(pInputOrderAction goctp.CThostFtdcInputOrderActionField, nRequestID int) int {

    mKey := Sprintf("%v_%v", pInputOrderAction.GetInstrumentID(), pInputOrderAction.GetOrderSysID())

    mOrderVal, mTestingOrderOk := MapOrders.Get(mKey)
    if !mTestingOrderOk {
        return 0
    }

    mOrder := mOrderVal.(OrderStruct)
    if mOrder.OrderStatus == "3" {
        mOrder.OrderStatus = "5"

        MapOrders.Set(mOrder.MapKey, mOrder)

        // 发送报单撤单通知
        pOrder := goctp.NewCThostFtdcOrderField()
        pOrder.SetBrokerID(mOrder.BrokerID)
        pOrder.SetInvestorID(mOrder.InvestorID)
        pOrder.SetInstrumentID(mOrder.InstrumentID)
        pOrder.SetExchangeID(mOrder.ExchangeID)
        pOrder.SetFrontID(mOrder.FrontID)
        pOrder.SetOrderRef(mOrder.OrderRef)
        pOrder.SetSessionID(mOrder.SessionID)
        pOrder.SetInsertTime(mOrder.InsertTime)
        pOrder.SetInsertDate(mOrder.InsertDate)
        pOrder.SetOrderSysID(mOrder.OrderSysID)
        pOrder.SetLimitPrice(mOrder.LimitPrice)
        pOrder.SetVolumeTotalOriginal(mOrder.Volume)
        pOrder.SetVolumeTraded(mOrder.VolumeTraded)
        pOrder.SetVolumeTotal(mOrder.VolumeTotal)
        pOrder.SetDirection([]byte(mOrder.Direction)[0])
        pOrder.SetCombOffsetFlag(mOrder.CombOffsetFlag)
        pOrder.SetCombHedgeFlag(mOrder.CombHedgeFlag)
        pOrder.SetOrderStatus(goctp.THOST_FTDC_OST_Canceled)
        pOrder.SetStatusMsg(mOrder.StatusMsg)

        p.TradeSpi.OnRtnOrder(pOrder)
    }

    return 0
}

// 查询最大报单数量请求
func (p *TradeApi) ReqQryMaxOrderVolume(pQryMaxOrderVolume goctp.CThostFtdcQryMaxOrderVolumeField, nRequestID int) int {
    return 0
}

// 投资者结算结果确认
func (p *TradeApi) ReqSettlementInfoConfirm(pSettlementInfoConfirm goctp.CThostFtdcSettlementInfoConfirmField, nRequestID int) int {

    pResult := goctp.NewCThostFtdcSettlementInfoConfirmField()
    pResult.SetBrokerID(p.BrokerID)
    pResult.SetInvestorID(p.InvestorID)
    pResult.SetConfirmDate(GetCurrentDate())

    p.TradeSpi.OnRspSettlementInfoConfirm(pResult, p.GetRspInfo(), nRequestID, true)

    return 0
}

// 请求删除预埋单
func (p *TradeApi) ReqRemoveParkedOrder(pRemoveParkedOrder goctp.CThostFtdcRemoveParkedOrderField, nRequestID int) int {
    return 0
}

// 请求删除预埋撤单
func (p *TradeApi) ReqRemoveParkedOrderAction(pRemoveParkedOrderAction goctp.CThostFtdcRemoveParkedOrderActionField, nRequestID int) int {
    return 0
}

// 执行宣告录入请求
func (p *TradeApi) ReqExecOrderInsert(pInputExecOrder goctp.CThostFtdcInputExecOrderField, nRequestID int) int {
    return 0
}

// 执行宣告操作请求
func (p *TradeApi) ReqExecOrderAction(pInputExecOrderAction goctp.CThostFtdcInputExecOrderActionField, nRequestID int) int {
    return 0
}

// 询价录入请求
func (p *TradeApi) ReqForQuoteInsert(pInputForQuote goctp.CThostFtdcInputForQuoteField, nRequestID int) int {
    return 0
}

// 报价录入请求
func (p *TradeApi) ReqQuoteInsert(pInputQuote goctp.CThostFtdcInputQuoteField, nRequestID int) int {
    return 0
}

// 报价操作请求
func (p *TradeApi) ReqQuoteAction(pInputQuoteAction goctp.CThostFtdcInputQuoteActionField, nRequestID int) int {
    return 0
}

// 批量报单操作请求
func (p *TradeApi) ReqBatchOrderAction(pInputBatchOrderAction goctp.CThostFtdcInputBatchOrderActionField, nRequestID int) int {
    return 0
}

// 期权自对冲录入请求
func (p *TradeApi) ReqOptionSelfCloseInsert(pInputOptionSelfClose goctp.CThostFtdcInputOptionSelfCloseField, nRequestID int) int {
    return 0
}

// 期权自对冲操作请求
func (p *TradeApi) ReqOptionSelfCloseAction(pInputOptionSelfCloseAction goctp.CThostFtdcInputOptionSelfCloseActionField, nRequestID int) int {
    return 0
}

// 申请组合录入请求
func (p *TradeApi) ReqCombActionInsert(pInputCombAction goctp.CThostFtdcInputCombActionField, nRequestID int) int {
    return 0
}

// 请求查询报单
func (p *TradeApi) ReqQryOrder(pQryOrder goctp.CThostFtdcQryOrderField, nRequestID int) int {

    Size := MapOrders.Size()

    if Size == 0 {

        pResult := goctp.NewCThostFtdcOrderField()

        pResult = nil

        p.TradeSpi.OnRspQryOrder(pResult, p.GetRspInfo(), nRequestID, true)

        return 0
    }

    count   := 0
    bIsLast := false

    mTestingOrders := MapOrders.GetAll()
    for _, v := range mTestingOrders {

        val := v.(OrderStruct)

        pResult := goctp.NewCThostFtdcOrderField()
        pResult.SetBrokerID(val.BrokerID)
        pResult.SetInvestorID(val.InvestorID)
        pResult.SetInstrumentID(val.InstrumentID)
        pResult.SetExchangeID(val.ExchangeID)
        pResult.SetFrontID(val.FrontID)
        pResult.SetOrderRef(val.OrderRef)
        pResult.SetSessionID(val.SessionID)
        pResult.SetInsertTime(val.InsertTime)
        pResult.SetInsertDate(val.InsertDate)
        pResult.SetOrderSysID(val.OrderSysID)
        pResult.SetLimitPrice(val.LimitPrice)
        pResult.SetVolumeTotalOriginal(val.Volume)
        pResult.SetVolumeTraded(val.VolumeTraded)
        pResult.SetVolumeTotal(val.VolumeTotal)
        pResult.SetDirection([]byte(val.Direction)[0])
        pResult.SetCombOffsetFlag(val.CombOffsetFlag)
        pResult.SetCombHedgeFlag(val.CombHedgeFlag)
        pResult.SetOrderStatus([]byte(val.OrderStatus)[0])
        pResult.SetStatusMsg(val.StatusMsg)

        count += 1

        if count == Size {
            bIsLast = true
        }

        p.TradeSpi.OnRspQryOrder(pResult, p.GetRspInfo(), nRequestID, bIsLast)
    }

    return 0
}

// 请求查询成交
func (p *TradeApi) ReqQryTrade(pQryTrade goctp.CThostFtdcQryTradeField, nRequestID int) int {

    Size := MapTrades.Size()

    if Size == 0 {

        pResult := goctp.NewCThostFtdcTradeField()

        pResult = nil

        p.TradeSpi.OnRspQryTrade(pResult, p.GetRspInfo(), nRequestID, true)

        return 0
    }

    count   := 0
    bIsLast := false

    mTestingTrades := MapTrades.GetAll()
    for _, v := range mTestingTrades {

        val := v.(TradeStruct)

        pResult := goctp.NewCThostFtdcTradeField()
        pResult.SetBrokerID(val.BrokerID)
        pResult.SetInvestorID(val.InvestorID)
        pResult.SetInstrumentID(val.InstrumentID)
        pResult.SetOrderRef(val.OrderRef)
        pResult.SetUserID(val.UserID)
        pResult.SetExchangeID(val.ExchangeID)
        pResult.SetTradeID(val.TradeID)
        pResult.SetDirection([]byte(val.Direction)[0])
        pResult.SetOrderSysID(val.OrderSysID)
        pResult.SetOffsetFlag([]byte(val.OffsetFlag)[0])
        pResult.SetHedgeFlag([]byte(val.HedgeFlag)[0])
        pResult.SetPrice(val.Price)
        pResult.SetVolume(val.Volume)
        pResult.SetTradeDate(val.TradeDate)
        pResult.SetTradeTime(val.TradeTime)
        pResult.SetTradeType([]byte(val.TradeType)[0])
        pResult.SetTradingDay(val.TradingDay)
        pResult.SetTradeSource([]byte(val.TradeSource)[0])

        count += 1

        if count == Size {
            bIsLast = true
        }

        p.TradeSpi.OnRspQryTrade(pResult, p.GetRspInfo(), nRequestID, bIsLast)
    }

    return 0
}

// 请求查询投资者持仓
func (p *TradeApi) ReqQryInvestorPosition(pQryInvestorPosition goctp.CThostFtdcQryInvestorPositionField, nRequestID int) int {

    Size := MapTrades.Size()

    if Size == 0 {

        pResult := goctp.NewCThostFtdcInvestorPositionField()

        pResult = nil

        p.TradeSpi.OnRspQryInvestorPosition(pResult, p.GetRspInfo(), nRequestID, true)

        return 0
    }

    mPositionDetails := MapPositionDetails.GetAll()
    for _, v := range mPositionDetails {
        val := v.(PositionDetailStruct)

        if val.Volume > 0 {

            Direction := "3"
            if val.Direction == "0" {
                Direction = "2"
            }

            mInstrument, _ := GetInstrumentInfo(val.InstrumentID)

            sPosition := PositionStruct{}

            sPosition.BrokerID       = val.BrokerID
            sPosition.InvestorID     = val.InvestorID
            sPosition.InstrumentID   = val.InstrumentID
            sPosition.Direction      = Direction
            sPosition.HedgeFlag      = val.HedgeFlag
            sPosition.PositionDate   = "1"
            sPosition.Volume         = val.Volume
            sPosition.ShortVolume    = 0
            sPosition.PositionProfit = 0
            sPosition.TradingDay     = val.TradingDay
            sPosition.OpenPrice      = val.OpenPrice * float64(val.Volume) * float64(mInstrument.VolumeMultiple)
            sPosition.ExchangeID     = val.ExchangeID

            if sPosition.TradingDay != IntToString(p.TradingDay) {
                sPosition.PositionDate = "2"
            }

            sPosition.MapKey = sPosition.InstrumentID + "_" + sPosition.Direction + "_" + sPosition.PositionDate

            MapPositions.Set(sPosition.MapKey, sPosition)
        }
    }


    Size     = MapPositions.Size()
    count   := 0
    bIsLast := false

    mPositions := MapPositions.GetAll()
    for _, v := range mPositions {
        val := v.(PositionStruct)

        pResult := goctp.NewCThostFtdcInvestorPositionField()
        pResult.SetBrokerID(val.BrokerID)
        pResult.SetInvestorID(val.InvestorID)
        pResult.SetInstrumentID(val.InstrumentID)
        pResult.SetHedgeFlag([]byte(val.HedgeFlag)[0])
        pResult.SetPosiDirection([]byte(val.Direction)[0])
        pResult.SetPositionDate([]byte(val.PositionDate)[0])
        pResult.SetPosition(val.Volume)
        pResult.SetYdPosition(val.YdPosition)
        pResult.SetTodayPosition(val.TodayPosition)
        pResult.SetLongFrozen(val.LongFrozen)
        pResult.SetShortFrozen(val.ShortFrozen)
        pResult.SetOpenCost(val.OpenPrice)
        pResult.SetCloseVolume(val.CloseVolume)
        pResult.SetTradingDay(val.TradingDay)
        pResult.SetOpenVolume(val.OpenVolume)
        pResult.SetCloseVolume(val.CloseVolume)
        pResult.SetExchangeID(val.ExchangeID)

        count += 1

        if count == Size {
            bIsLast = true
        }

        p.TradeSpi.OnRspQryInvestorPosition(pResult, p.GetRspInfo(), nRequestID, bIsLast)
    }

    return 0
}

// 请求查询资金账户
func (p *TradeApi) ReqQryTradingAccount(pQryTradingAccount goctp.CThostFtdcQryTradingAccountField, nRequestID int) int {

    pResult := goctp.NewCThostFtdcTradingAccountField()

    mKey := Sprintf("%v_%v", pQryTradingAccount.GetBrokerID(), pQryTradingAccount.GetInvestorID())

    Margin         := 0.00
    ExchMargin     := 0.00
    CloseProfit    := 0.00
    Commission     := 0.00

    mPositionDetails := MapPositionDetails.GetAll()
    for _, v := range mPositionDetails {

        val := v.(PositionDetailStruct)

        Margin      += float64(val.Volume) * p.Margin
        ExchMargin  += float64(val.Volume) * p.Margin
        CloseProfit += val.CloseProfitByTrade
    }

    mTestingTrades := MapTrades.GetAll()
    for _, v := range mTestingTrades {
        val := v.(TradeStruct)

        Commission += float64(val.Volume) * p.Commission
    }

    mVal, mOk := MapAccounts.Get(mKey)
    if mOk {
        val := mVal.(AccountStruct)

        // 动态权益
        Balance := val.PreBalance - Commission + CloseProfit

        // 可用资金
        Available := val.PreBalance - Commission - Margin + CloseProfit

        pResult.SetBrokerID(val.BrokerID)
        pResult.SetAccountID(val.AccountID)
        pResult.SetPreBalance(val.PreBalance)
        pResult.SetDeposit(val.Deposit)
        pResult.SetWithdraw(val.Withdraw)
        pResult.SetFrozenMargin(val.FrozenMargin)
        pResult.SetFrozenCash(val.FrozenCash)
        pResult.SetFrozenCommission(val.FrozenCommission)
        pResult.SetCurrMargin(Margin)
        pResult.SetCashIn(val.CashIn)
        pResult.SetCommission(Commission)
        pResult.SetCloseProfit(CloseProfit)
        pResult.SetPositionProfit(val.PositionProfit)
        pResult.SetBalance(Balance)
        pResult.SetAvailable(Available)
        pResult.SetWithdrawQuota(Available)
        pResult.SetExchangeMargin(ExchMargin)
        pResult.SetTradingDay(GetCurrentDate())
    }

    p.TradeSpi.OnRspQryTradingAccount(pResult, p.GetRspInfo(), nRequestID, true)

    return 0
}

// 请求查询投资者
func (p *TradeApi) ReqQryInvestor(pQryInvestor goctp.CThostFtdcQryInvestorField, nRequestID int) int {
    return 0
}

// 请求查询交易编码
func (p *TradeApi) ReqQryTradingCode(pQryTradingCode goctp.CThostFtdcQryTradingCodeField, nRequestID int) int {
    return 0
}

// 请求查询合约保证金率
func (p *TradeApi) ReqQryInstrumentMarginRate(pQryInstrumentMarginRate goctp.CThostFtdcQryInstrumentMarginRateField, nRequestID int) int {
    return 0
}

// 请求查询合约手续费率
func (p *TradeApi) ReqQryInstrumentCommissionRate(pQryInstrumentCommissionRate goctp.CThostFtdcQryInstrumentCommissionRateField, nRequestID int) int {
    return 0
}

// 请求查询交易所
func (p *TradeApi) ReqQryExchange(pQryExchange goctp.CThostFtdcQryExchangeField, nRequestID int) int {
    return 0
}

// 请求查询产品
func (p *TradeApi) ReqQryProduct(pQryProduct goctp.CThostFtdcQryProductField, nRequestID int) int {
    return 0
}

// 请求查询合约
func (p *TradeApi) ReqQryInstrument(pQryInstrument goctp.CThostFtdcQryInstrumentField, nRequestID int) int {

    pResult := goctp.NewCThostFtdcInstrumentField()

    pResult = nil

    p.TradeSpi.OnRspQryInstrument(pResult, p.GetRspInfo(), nRequestID, true)

    return 0
}

// 请求查询行情
func (p *TradeApi) ReqQryDepthMarketData(pQryDepthMarketData goctp.CThostFtdcQryDepthMarketDataField, nRequestID int) int {
    return 0
}

// 请求查询交易员报盘机
func (p *TradeApi) ReqQryTraderOffer(pQryTraderOffer goctp.CThostFtdcQryTraderOfferField, nRequestID int) int {
    return 0
}

// 请求查询投资者结算结果
func (p *TradeApi) ReqQrySettlementInfo(pQrySettlementInfo goctp.CThostFtdcQrySettlementInfoField, nRequestID int) int {
    return 0
}

// 请求查询转帐银行
func (p *TradeApi) ReqQryTransferBank(pQryTransferBank goctp.CThostFtdcQryTransferBankField, nRequestID int) int {
    return 0
}

// 请求查询投资者持仓明细
func (p *TradeApi) ReqQryInvestorPositionDetail(pQryInvestorPositionDetail goctp.CThostFtdcQryInvestorPositionDetailField, nRequestID int) int {

    Size := MapPositionDetails.Size()

    if Size == 0 {
        pResult := goctp.NewCThostFtdcInvestorPositionDetailField()

        pResult = nil

        p.TradeSpi.OnRspQryInvestorPositionDetail(pResult, p.GetRspInfo(), nRequestID, true)

        return 0
    }

    count   := 0
    bIsLast := false

    mPositionDetails := MapPositionDetails.GetAll()
    for _, v := range mPositionDetails {

        val := v.(PositionDetailStruct)

        pResult := goctp.NewCThostFtdcInvestorPositionDetailField()
        pResult.SetBrokerID(val.BrokerID)
        pResult.SetInvestorID(val.InvestorID)
        pResult.SetInstrumentID(val.InstrumentID)
        pResult.SetHedgeFlag([]byte(val.HedgeFlag)[0])
        pResult.SetDirection([]byte(val.Direction)[0])
        pResult.SetOpenDate(val.OpenDate)
        pResult.SetTradeID(val.TradeID)
        pResult.SetVolume(val.Volume)
        pResult.SetOpenPrice(val.OpenPrice)
        pResult.SetTradingDay(val.TradingDay)
        pResult.SetExchangeID(val.ExchangeID)
        pResult.SetCloseVolume(val.CloseVolume)
        pResult.SetCloseProfitByTrade(val.CloseProfitByTrade)

        count += 1

        if count == Size {
            bIsLast = true
        }

        p.TradeSpi.OnRspQryInvestorPositionDetail(pResult, p.GetRspInfo(), nRequestID, bIsLast)
    }

    return 0
}

// 请求查询客户通知
func (p *TradeApi) ReqQryNotice(pQryNotice goctp.CThostFtdcQryNoticeField, nRequestID int) int {
    return 0
}

// 请求查询结算信息确认
func (p *TradeApi) ReqQrySettlementInfoConfirm(pQrySettlementInfoConfirm goctp.CThostFtdcQrySettlementInfoConfirmField, nRequestID int) int {
    return 0
}

// 请求查询投资者持仓明细
func (p *TradeApi) ReqQryInvestorPositionCombineDetail(pQryInvestorPositionCombineDetail goctp.CThostFtdcQryInvestorPositionCombineDetailField, nRequestID int) int {
    return 0
}

// 请求查询保证金监管系统经纪公司资金账户密钥
func (p *TradeApi) ReqQryCFMMCTradingAccountKey(pQryCFMMCTradingAccountKey goctp.CThostFtdcQryCFMMCTradingAccountKeyField, nRequestID int) int {
    return 0
}

// 请求查询仓单折抵信息
func (p *TradeApi) ReqQryEWarrantOffset(pQryEWarrantOffset goctp.CThostFtdcQryEWarrantOffsetField, nRequestID int) int {
    return 0
}

// 请求查询投资者品种/跨品种保证金
func (p *TradeApi) ReqQryInvestorProductGroupMargin(pQryInvestorProductGroupMargin goctp.CThostFtdcQryInvestorProductGroupMarginField, nRequestID int) int {
    return 0
}

// 请求查询交易所保证金率
func (p *TradeApi) ReqQryExchangeMarginRate(pQryExchangeMarginRate goctp.CThostFtdcQryExchangeMarginRateField, nRequestID int) int {
    return 0
}

// 请求查询交易所调整保证金率
func (p *TradeApi) ReqQryExchangeMarginRateAdjust(pQryExchangeMarginRateAdjust goctp.CThostFtdcQryExchangeMarginRateAdjustField, nRequestID int) int {
    return 0
}

// 请求查询汇率
func (p *TradeApi) ReqQryExchangeRate(pQryExchangeRate goctp.CThostFtdcQryExchangeRateField, nRequestID int) int {
    return 0
}

// 请求查询二级代理操作员银期权限
func (p *TradeApi) ReqQrySecAgentACIDMap(pQrySecAgentACIDMap goctp.CThostFtdcQrySecAgentACIDMapField, nRequestID int) int {
    return 0
}

// 请求查询产品报价汇率
func (p *TradeApi) ReqQryProductExchRate(pQryProductExchRate goctp.CThostFtdcQryProductExchRateField, nRequestID int) int {
    return 0
}

// 请求查询产品组
func (p *TradeApi) ReqQryProductGroup(pQryProductGroup goctp.CThostFtdcQryProductGroupField, nRequestID int) int {
    return 0
}

// 请求查询做市商合约手续费率
func (p *TradeApi) ReqQryMMInstrumentCommissionRate(pQryMMInstrumentCommissionRate goctp.CThostFtdcQryMMInstrumentCommissionRateField, nRequestID int) int {
    return 0
}

// 请求查询做市商期权合约手续费
func (p *TradeApi) ReqQryMMOptionInstrCommRate(pQryMMOptionInstrCommRate goctp.CThostFtdcQryMMOptionInstrCommRateField, nRequestID int) int {
    return 0
}

// 请求查询报单手续费
func (p *TradeApi) ReqQryInstrumentOrderCommRate(pQryInstrumentOrderCommRate goctp.CThostFtdcQryInstrumentOrderCommRateField, nRequestID int) int {
    return 0
}

// 请求查询资金账户
func (p *TradeApi) ReqQrySecAgentTradingAccount(pQryTradingAccount goctp.CThostFtdcQryTradingAccountField, nRequestID int) int {
    return 0
}

// 请求查询二级代理商资金校验模式
func (p *TradeApi) ReqQrySecAgentCheckMode(pQrySecAgentCheckMode goctp.CThostFtdcQrySecAgentCheckModeField, nRequestID int) int {
    return 0
}

// 请求查询二级代理商信息
func (p *TradeApi) ReqQrySecAgentTradeInfo(pQrySecAgentTradeInfo goctp.CThostFtdcQrySecAgentTradeInfoField, nRequestID int) int {
    return 0
}

// 请求查询期权交易成本
func (p *TradeApi) ReqQryOptionInstrTradeCost(pQryOptionInstrTradeCost goctp.CThostFtdcQryOptionInstrTradeCostField, nRequestID int) int {
    return 0
}

// 请求查询期权合约手续费
func (p *TradeApi) ReqQryOptionInstrCommRate(pQryOptionInstrCommRate goctp.CThostFtdcQryOptionInstrCommRateField, nRequestID int) int {
    return 0
}

// 请求查询执行宣告
func (p *TradeApi) ReqQryExecOrder(pQryExecOrder goctp.CThostFtdcQryExecOrderField, nRequestID int) int {
    return 0
}

// 请求查询询价
func (p *TradeApi) ReqQryForQuote(pQryForQuote goctp.CThostFtdcQryForQuoteField, nRequestID int) int {
    return 0
}

// 请求查询报价
func (p *TradeApi) ReqQryQuote(pQryQuote goctp.CThostFtdcQryQuoteField, nRequestID int) int {
    return 0
}

// 请求查询期权自对冲
func (p *TradeApi) ReqQryOptionSelfClose(pQryOptionSelfClose goctp.CThostFtdcQryOptionSelfCloseField, nRequestID int) int {
    return 0
}

// 请求查询投资单元
func (p *TradeApi) ReqQryInvestUnit(pQryInvestUnit goctp.CThostFtdcQryInvestUnitField, nRequestID int) int {
    return 0
}

// 请求查询组合合约安全系数
func (p *TradeApi) ReqQryCombInstrumentGuard(pQryCombInstrumentGuard goctp.CThostFtdcQryCombInstrumentGuardField, nRequestID int) int {
    return 0
}

// 请求查询申请组合
func (p *TradeApi) ReqQryCombAction(pQryCombAction goctp.CThostFtdcQryCombActionField, nRequestID int) int {
    return 0
}

// 请求查询转帐流水
func (p *TradeApi) ReqQryTransferSerial(pQryTransferSerial goctp.CThostFtdcQryTransferSerialField, nRequestID int) int {
    return 0
}

// 请求查询银期签约关系
func (p *TradeApi) ReqQryAccountregister(pQryAccountregister goctp.CThostFtdcQryAccountregisterField, nRequestID int) int {
    return 0
}

// 请求查询签约银行
func (p *TradeApi) ReqQryContractBank(pQryContractBank goctp.CThostFtdcQryContractBankField, nRequestID int) int {
    return 0
}

// 请求查询预埋单
func (p *TradeApi) ReqQryParkedOrder(pQryParkedOrder goctp.CThostFtdcQryParkedOrderField, nRequestID int) int {
    return 0
}

// 请求查询预埋撤单
func (p *TradeApi) ReqQryParkedOrderAction(pQryParkedOrderAction goctp.CThostFtdcQryParkedOrderActionField, nRequestID int) int {
    return 0
}

// 请求查询交易通知
func (p *TradeApi) ReqQryTradingNotice(pQryTradingNotice goctp.CThostFtdcQryTradingNoticeField, nRequestID int) int {
    return 0
}

// 请求查询经纪公司交易参数
func (p *TradeApi) ReqQryBrokerTradingParams(pQryBrokerTradingParams goctp.CThostFtdcQryBrokerTradingParamsField, nRequestID int) int {
    return 0
}

// 请求查询经纪公司交易算法
func (p *TradeApi) ReqQryBrokerTradingAlgos(pQryBrokerTradingAlgos goctp.CThostFtdcQryBrokerTradingAlgosField, nRequestID int) int {
    return 0
}

// 请求查询监控中心用户令牌
func (p *TradeApi) ReqQueryCFMMCTradingAccountToken(pQueryCFMMCTradingAccountToken goctp.CThostFtdcQueryCFMMCTradingAccountTokenField, nRequestID int) int {
    return 0
}

// 期货发起银行资金转期货请求
func (p *TradeApi) ReqFromBankToFutureByFuture(pReqTransfer goctp.CThostFtdcReqTransferField, nRequestID int) int {
    return 0
}

// 期货发起期货资金转银行请求
func (p *TradeApi) ReqFromFutureToBankByFuture(pReqTransfer goctp.CThostFtdcReqTransferField, nRequestID int) int {
    return 0
}

// 期货发起查询银行余额请求
func (p *TradeApi) ReqQueryBankAccountMoneyByFuture(pReqQueryAccount goctp.CThostFtdcReqQueryAccountField, nRequestID int) int {
    return 0
}

// 请求查询分类合约
func (p *TradeApi) ReqQryClassifiedInstrument(pQryClassifiedInstrument goctp.CThostFtdcQryClassifiedInstrumentField, nRequestID int) int {
    return 0
}

// 请求组合优惠比例
func (p *TradeApi) ReqQryCombPromotionParam(pQryCombPromotionParam goctp.CThostFtdcQryCombPromotionParamField, nRequestID int) int {
    return 0
}

// 投资者风险结算持仓查询
func (p *TradeApi) ReqQryRiskSettleInvstPosition(pQryRiskSettleInvstPosition goctp.CThostFtdcQryRiskSettleInvstPositionField, nRequestID int) int {
    return 0
}

// 风险结算产品查询
func (p *TradeApi) ReqQryRiskSettleProductStatus(pQryRiskSettleProductStatus goctp.CThostFtdcQryRiskSettleProductStatusField, nRequestID int) int {
    return 0
}

// 设置报单引用的 tick 数据
func (p *TradeApi) SetOrderRefTick(iRequestID int, pTick TickStruct) {

    MapOrderRefTick.Set(iRequestID, pTick)
}

// 设置合约信息
func (p *TradeApi) SetInstrument(sInstrument InstrumentStruct) {

    MapInstruments.Set(ToUpper(sInstrument.InstrumentID), sInstrument)
}

// 设置当前交易日期（格式：20230101）
func (p *TradeApi) SetTradingDay(TradingDay int) {

    p.OldThis().TradingDay = TradingDay
}

// 获得合约详情信息
func GetInstrumentInfo(InstrumentID string) (InstrumentStruct, bool) {

    // 将合约转换成大写
    InstrumentID = ToUpper(InstrumentID)

    if v, ok := MapInstruments.Get(InstrumentID); ok {
        return v.(InstrumentStruct), true
    }

    var mInstrument InstrumentStruct

    return mInstrument, false
}

/**
 * 计算盈亏
 *
 * @param   InstrumentID  string  合约
 * @param   OpenPrice     float64 开仓价格
 * @param   LastPrice     float64 最新价|平仓价格
 * @param   Number        int     数量
 * @param   Direction     string  持仓方向[0|2：买，1|3：卖]
 */
func GetPositionProfit(InstrumentID string, OpenPrice float64, LastPrice float64, Number int, Direction string) float64 {

    mInstrument, mInstrumentOk := GetInstrumentInfo(InstrumentID)
    if !mInstrumentOk {
        return 0.00
    }

    if Direction == "2" || Direction == "0" {
        return ((LastPrice - OpenPrice) * float64(mInstrument.VolumeMultiple)) * float64(Number)
    } else {
        return ((OpenPrice - LastPrice) * float64(mInstrument.VolumeMultiple)) * float64(Number)
    }
}