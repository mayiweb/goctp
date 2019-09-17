package main

import (
    "github.com/mayiweb/goctp"
    "fmt"
    "log"
    "time"
)

// 获得交易请求编号
func GetTraderRequestId() int {
    Ctp.TraderRequestId += 1
    return Ctp.TraderRequestId
}

// 当客户端与交易后台通信连接断开时，该方法被调用。当发生这个情况后，API会自动重新连接，客户端可不做处理。
// 服务器已断线，该函数也会被调用。【api 会自动初始化程序，并重新登陆】
func (p *FtdcTraderSpi) OnFrontDisconnected(nReason int) {

    Ctp.IsTraderLogin      = false
    Ctp.IsTraderInit       = false
    Ctp.IsTraderInitFinish = false

    log.Println("交易服务器已断线，尝试重新连接中...")
}

// 发送请求日志（仅查询类的函数需要调用）
func (p *FtdcTraderSpi) ReqMsg(Msg string) {

    // 交易程序未初始化完成时，执行查询类的函数需要有1.5秒间隔
    if !Ctp.IsTraderInitFinish {
        time.Sleep(time.Millisecond * 1500)
    }

    fmt.Println("")
    log.Println(Msg)
}

// 当客户端与交易后台建立起通信连接时（还未登录前），该方法被调用。
func (p *FtdcTraderSpi) OnFrontConnected() {

    TraderStr := "=================================================================================================\n" +
                 "= 交易模块初始化成功，API 版本：" + goctp.CThostFtdcTraderApiGetApiVersion() + "\n" +
                 "================================================================================================="
    fmt.Println(TraderStr)

    Ctp.IsTraderInit = true

    // 填写了 AppID 与 AuthCode 则进行客户端认证
    if AppID != "" && AuthCode != "" {

        p.ReqAuthenticate()
    } else {

        MdSpi.ReqUserLogin()

        p.ReqUserLogin()
    }
}

// 客户端认证
func (p *FtdcTraderSpi) ReqAuthenticate() {

    log.Println("客户端认证中...")

    req := goctp.NewCThostFtdcReqAuthenticateField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetUserID(Ctp.InvestorID)
    req.SetAppID(Ctp.AppID)
    req.SetAuthCode(Ctp.AuthCode)

    iResult := Ctp.TraderApi.ReqAuthenticate(req, GetTraderRequestId())

    if iResult != 0 {
        ReqFailMsg("发送客户端认证请求失败！", iResult)
    }
}

// 客户端认证响应
func (p *FtdcTraderSpi) OnRspAuthenticate(pRspAuthenticateField goctp.CThostFtdcRspAuthenticateField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

        log.Println("客户端认证成功！")

        MdSpi.ReqUserLogin()

        p.ReqUserLogin()
    }
}

// 用户登录请求
func (p *FtdcTraderSpi) ReqUserLogin() {

    time.Sleep(time.Second * 1)

    log.Println("交易系统账号登陆中...")

    req := goctp.NewCThostFtdcReqUserLoginField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetUserID(Ctp.InvestorID)
    req.SetPassword(Ctp.Password)

    iResult := Ctp.TraderApi.ReqUserLogin(req, GetTraderRequestId())

    if iResult != 0 {
        ReqFailMsg("发送用户登录请求失败！", iResult)
    }
}

func (p *FtdcTraderSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

        Ctp.IsTraderLogin = true

        log.Printf("交易账号已登录，当前交易日：%v\n", Ctp.TraderApi.GetTradingDay())

        p.ReqSettlementInfoConfirm()
    }
}

// 投资者结算单确认
func (p *FtdcTraderSpi) ReqSettlementInfoConfirm() int {

    p.ReqMsg("投资者结算单确认中...")

    req := goctp.NewCThostFtdcSettlementInfoConfirmField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TraderApi.ReqSettlementInfoConfirm(req, GetTraderRequestId())

    if iResult != 0 {
        ReqFailMsg("确认投资者结算单失败！", iResult)
    }

    return iResult
}

// 发送投资者结算单确认响应
func (p *FtdcTraderSpi) OnRspSettlementInfoConfirm(pSettlementInfoConfirm goctp.CThostFtdcSettlementInfoConfirmField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    if bIsLast && !p.IsErrorRspInfo(pRspInfo) {
        log.Println("投资者结算单确认成功")

        p.ReqQryInstrument()
    }
}

// 请求查询合约
func (p *FtdcTraderSpi) ReqQryInstrument() int {

    p.ReqMsg("查询合约中...")

    req := goctp.NewCThostFtdcQryInstrumentField()
    req.SetInstrumentID("")

    iResult := Ctp.TraderApi.ReqQryInstrument(req, GetTraderRequestId())

    if iResult != 0 {
        ReqFailMsg("查询合约失败！", iResult)
    }

    return iResult
}

// 请求查询合约响应
func (p *FtdcTraderSpi) OnRspQryInstrument(pInstrument goctp.CThostFtdcInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    if !p.IsErrorRspInfo(pRspInfo) {

        var mInstrumentInfo InstrumentInfoStruct

        var mapKey string = pInstrument.GetInstrumentID()

        mInstrumentInfo.InstrumentID           = pInstrument.GetInstrumentID();
        mInstrumentInfo.ExchangeID             = pInstrument.GetExchangeID();
        mInstrumentInfo.InstrumentName         = ConvertToString(pInstrument.GetInstrumentName(), "gbk", "utf-8")
        mInstrumentInfo.ExchangeInstID         = pInstrument.GetExchangeInstID();
        mInstrumentInfo.ProductID              = pInstrument.GetProductID();
        mInstrumentInfo.ProductClass           = string(pInstrument.GetProductClass());
        mInstrumentInfo.DeliveryYear           = pInstrument.GetDeliveryYear();
        mInstrumentInfo.DeliveryMonth          = pInstrument.GetDeliveryMonth();
        mInstrumentInfo.MaxMarketOrderVolume   = pInstrument.GetMaxMarketOrderVolume();
        mInstrumentInfo.MinMarketOrderVolume   = pInstrument.GetMinMarketOrderVolume();
        mInstrumentInfo.MaxLimitOrderVolume    = pInstrument.GetMaxLimitOrderVolume();
        mInstrumentInfo.MinLimitOrderVolume    = pInstrument.GetMinLimitOrderVolume();
        mInstrumentInfo.VolumeMultiple         = pInstrument.GetVolumeMultiple();
        mInstrumentInfo.PriceTick              = pInstrument.GetPriceTick();
        mInstrumentInfo.CreateDate             = pInstrument.GetCreateDate();
        mInstrumentInfo.OpenDate               = pInstrument.GetOpenDate();
        mInstrumentInfo.ExpireDate             = pInstrument.GetExpireDate();
        mInstrumentInfo.StartDelivDate         = pInstrument.GetStartDelivDate();
        mInstrumentInfo.EndDelivDate           = pInstrument.GetEndDelivDate();
        mInstrumentInfo.InstLifePhase          = string(pInstrument.GetInstLifePhase());
        mInstrumentInfo.IsTrading              = pInstrument.GetIsTrading();
        mInstrumentInfo.PositionType           = string(pInstrument.GetPositionType());
        mInstrumentInfo.PositionDateType       = string(pInstrument.GetPositionDateType());
        mInstrumentInfo.LongMarginRatio        = pInstrument.GetLongMarginRatio();
        mInstrumentInfo.ShortMarginRatio       = pInstrument.GetShortMarginRatio();
        mInstrumentInfo.MaxMarginSideAlgorithm = string(pInstrument.GetMaxMarginSideAlgorithm());
        mInstrumentInfo.UnderlyingInstrID      = pInstrument.GetUnderlyingInstrID();
        mInstrumentInfo.StrikePrice            = pInstrument.GetStrikePrice();
        mInstrumentInfo.OptionsType            = string(pInstrument.GetOptionsType());
        mInstrumentInfo.UnderlyingMultiple     = pInstrument.GetUnderlyingMultiple();
        mInstrumentInfo.CombinationType        = string(pInstrument.GetCombinationType());

        MapInstrumentInfos.Store(mapKey, mInstrumentInfo)

        if bIsLast {

            MapInstrumentInfoSize := 0

            MapInstrumentInfos.Range(func(k, v interface{}) bool {
                MapInstrumentInfoSize += 1
                return true
            })

            log.Printf("获得合约记录 %v 条\n", MapInstrumentInfoSize)

            if !Ctp.IsTraderInitFinish {
                // 请求查询资金账户
                p.ReqQryTradingAccount()
            }
        }
    }
}

// 请求查询资金账户
func (p *FtdcTraderSpi) ReqQryTradingAccount() int {

    p.ReqMsg("查询资金账户中...")

    req := goctp.NewCThostFtdcQryTradingAccountField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TraderApi.ReqQryTradingAccount(req, GetTraderRequestId())

    if iResult != 0 {
        ReqFailMsg("查询资金账户失败！", iResult)
    }

    return iResult
}

// 请求查询资金账户响应
func (p *FtdcTraderSpi) OnRspQryTradingAccount(pTradingAccount goctp.CThostFtdcTradingAccountField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    if bIsLast && !p.IsErrorRspInfo(pRspInfo) {

        var mAccountInfo AccountInfoStruct

        mAccountInfo.MapKey = pTradingAccount.GetBrokerID() + "_" + pTradingAccount.GetAccountID()

        mAccountInfo.BrokerID                       = pTradingAccount.GetBrokerID()
        mAccountInfo.AccountID                      = pTradingAccount.GetAccountID()
        mAccountInfo.PreMortgage                    = Decimal(pTradingAccount.GetPreMortgage(), 2)
        mAccountInfo.PreCredit                      = Decimal(pTradingAccount.GetPreCredit(), 2)
        mAccountInfo.PreDeposit                     = Decimal(pTradingAccount.GetPreDeposit(), 2)
        mAccountInfo.PreBalance                     = Decimal(pTradingAccount.GetPreBalance(), 2)
        mAccountInfo.PreMargin                      = Decimal(pTradingAccount.GetPreMargin(), 2)
        mAccountInfo.InterestBase                   = Decimal(pTradingAccount.GetInterestBase(), 2)
        mAccountInfo.Interest                       = Decimal(pTradingAccount.GetInterest(), 2)
        mAccountInfo.Deposit                        = Decimal(pTradingAccount.GetDeposit(), 2)
        mAccountInfo.Withdraw                       = Decimal(pTradingAccount.GetWithdraw(), 2)
        mAccountInfo.FrozenMargin                   = Decimal(pTradingAccount.GetFrozenMargin(), 2)
        mAccountInfo.FrozenCash                     = Decimal(pTradingAccount.GetFrozenCash(), 2)
        mAccountInfo.FrozenCommission               = Decimal(pTradingAccount.GetFrozenCommission(), 2)
        mAccountInfo.CurrMargin                     = Decimal(pTradingAccount.GetCurrMargin(), 2)
        mAccountInfo.CashIn                         = Decimal(pTradingAccount.GetCashIn(), 2)
        mAccountInfo.Commission                     = Decimal(pTradingAccount.GetCommission(), 2)
        mAccountInfo.CloseProfit                    = Decimal(pTradingAccount.GetCloseProfit(), 2)
        mAccountInfo.PositionProfit                 = Decimal(pTradingAccount.GetPositionProfit(), 2)
        mAccountInfo.Balance                        = Decimal(pTradingAccount.GetBalance(), 2)
        mAccountInfo.Available                      = Decimal(pTradingAccount.GetAvailable(), 2)
        mAccountInfo.WithdrawQuota                  = Decimal(pTradingAccount.GetWithdrawQuota(), 2)
        mAccountInfo.Reserve                        = Decimal(pTradingAccount.GetReserve(), 2)
        mAccountInfo.TradingDay                     = pTradingAccount.GetTradingDay()
        mAccountInfo.SettlementID                   = pTradingAccount.GetSettlementID()
        mAccountInfo.Credit                         = Decimal(pTradingAccount.GetCredit(), 2)
        mAccountInfo.Mortgage                       = Decimal(pTradingAccount.GetMortgage(), 2)
        mAccountInfo.ExchangeMargin                 = Decimal(pTradingAccount.GetExchangeMargin(), 2)
        mAccountInfo.DeliveryMargin                 = Decimal(pTradingAccount.GetDeliveryMargin(), 2)
        mAccountInfo.ExchangeDeliveryMargin         = Decimal(pTradingAccount.GetExchangeDeliveryMargin(), 2)
        mAccountInfo.ReserveBalance                 = Decimal(pTradingAccount.GetReserveBalance(), 2)
        mAccountInfo.CurrencyID                     = pTradingAccount.GetCurrencyID()
        mAccountInfo.PreFundMortgageIn              = Decimal(pTradingAccount.GetPreFundMortgageIn(), 2)
        mAccountInfo.PreFundMortgageOut             = Decimal(pTradingAccount.GetPreFundMortgageOut(), 2)
        mAccountInfo.FundMortgageIn                 = Decimal(pTradingAccount.GetFundMortgageIn(), 2)
        mAccountInfo.FundMortgageOut                = Decimal(pTradingAccount.GetFundMortgageOut(), 2)
        mAccountInfo.FundMortgageAvailable          = Decimal(pTradingAccount.GetFundMortgageAvailable(), 2)
        mAccountInfo.MortgageableFund               = Decimal(pTradingAccount.GetMortgageableFund(), 2)
        mAccountInfo.SpecProductMargin              = Decimal(pTradingAccount.GetSpecProductMargin(), 2)
        mAccountInfo.SpecProductFrozenMargin        = Decimal(pTradingAccount.GetSpecProductFrozenMargin(), 2)
        mAccountInfo.SpecProductCommission          = Decimal(pTradingAccount.GetSpecProductCommission(), 2)
        mAccountInfo.SpecProductFrozenCommission    = Decimal(pTradingAccount.GetSpecProductFrozenCommission(), 2)
        mAccountInfo.SpecProductPositionProfit      = Decimal(pTradingAccount.GetSpecProductPositionProfit(), 2)
        mAccountInfo.SpecProductCloseProfit         = Decimal(pTradingAccount.GetSpecProductCloseProfit(), 2)
        mAccountInfo.SpecProductPositionProfitByAlg = Decimal(pTradingAccount.GetSpecProductPositionProfitByAlg(), 2)
        mAccountInfo.SpecProductExchangeMargin      = Decimal(pTradingAccount.GetSpecProductExchangeMargin(), 2)
        mAccountInfo.BizType                        = string(pTradingAccount.GetBizType())
        mAccountInfo.FrozenSwap                     = Decimal(pTradingAccount.GetFrozenSwap(), 2)
        mAccountInfo.RemainSwap                     = Decimal(pTradingAccount.GetRemainSwap(), 2)

        AccountInfoStr := "-------------------------------------------------------------------------------------------------\n" +
                          "- 公司代码：" + pTradingAccount.GetBrokerID() + "\n" +
                          "- 资金账号：" + pTradingAccount.GetAccountID() + "\n" +
                          "- 期初资金：" + Float64ToString(mAccountInfo.PreBalance) + "\n" +
                          "- 动态权益：" + Float64ToString(mAccountInfo.Balance) + "\n" +
                          "- 可用资金：" + Float64ToString(mAccountInfo.Available) + "\n" +
                          "- 持仓盈亏：" + Float64ToString(mAccountInfo.PositionProfit) + "\n" +
                          "- 平仓盈亏：" + Float64ToString(mAccountInfo.CloseProfit) + "\n" +
                          "- 手续费  ：" + Float64ToString(mAccountInfo.Commission) + "\n" +
                          "-------------------------------------------------------------------------------------------------"
        fmt.Println(AccountInfoStr)

        if !Ctp.IsTraderInitFinish {
            // 请求查询投资者报单（委托单）
            p.ReqQryOrder()
        }
    }
}

// 请求查询投资者报单（委托单）
func (p *FtdcTraderSpi) ReqQryOrder() int {

    p.ReqMsg("查询投资者报单中...")

    req := goctp.NewCThostFtdcQryOrderField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TraderApi.ReqQryOrder(req, GetTraderRequestId())

    if iResult != 0 {
        ReqFailMsg("查询投资者报单失败！", iResult)
    }

    return iResult
}

// 请求查询投资者报单响应
func (p *FtdcTraderSpi) OnRspQryOrder(pOrder goctp.CThostFtdcOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    if !p.IsErrorRspInfo(pRspInfo) {

        // 如果 没有数据 pOrder 等于0
        pOrderCode := fmt.Sprintf("%v", pOrder)

        // 只记录有报单编号的报单数据
        if pOrderCode != "0" && pOrder.GetOrderSysID() != "" {
            // 获得报单结构体数据
            mOrder := GetOrderListStruct(pOrder)

            // 报单列表数据 key 键
            mOrder.MapKey = pOrder.GetInstrumentID() + "_" + TrimSpace(pOrder.GetOrderSysID())

            // 记录报单数据
            MapOrderList.Store(mOrder.MapKey, mOrder)
        }

        if bIsLast {

            fmt.Println("-------------------------------------------------------------------------------------------------")

            MapOrderListSize    := 0
            MapOrderNoTradeSize := 0
            MapOrderList.Range(func(key, v interface{}) bool {

                val := v.(OrderListStruct)

                MapOrderListSize += 1

                // 输出 未成交、部分成交 的报单
                if val.OrderStatus == string(goctp.THOST_FTDC_OST_NoTradeQueueing) || val.OrderStatus == string(goctp.THOST_FTDC_OST_PartTradedQueueing) {
                    MapOrderNoTradeSize += 1
                    fmt.Printf("- 合约：%v   \t%v:%v   \t数量：%v   \t价格：%v   \t报单编号：%v (%v)\n", val.InstrumentID, val.DirectionTitle, val.CombOffsetFlagTitle, val.Volume, val.LimitPrice, TrimSpace(val.OrderSysID), val.OrderStatusTitle)
                }
                return true
            })

            fmt.Printf("- 共有报单记录 %v 条，未成交 %v 条（不含错单）\n", MapOrderListSize, MapOrderNoTradeSize)
            fmt.Println("-------------------------------------------------------------------------------------------------")

            if !Ctp.IsTraderInitFinish {
                // 请求查询投资者持仓（汇总）
                p.ReqQryInvestorPosition()
            }
        }
    }
}

// 请求查询投资者持仓（汇总）
func (p *FtdcTraderSpi) ReqQryInvestorPosition() int {

    p.ReqMsg("查询投资者持仓中...")

    req := goctp.NewCThostFtdcQryInvestorPositionField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TraderApi.ReqQryInvestorPosition(req, GetTraderRequestId())

    if iResult != 0 {
        ReqFailMsg("查询投资者持仓失败！", iResult)
    }

    fmt.Println("-------------------------------------------------------------------------------------------------")

    return iResult
}

// 请求查询投资者持仓（汇总）响应
func (p *FtdcTraderSpi) OnRspQryInvestorPosition(pInvestorPosition goctp.CThostFtdcInvestorPositionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    if !p.IsErrorRspInfo(pRspInfo) {

        // 没有数据 pInvestorPosition 会等于 0
        pInvestorPositionCode := fmt.Sprintf("%v", pInvestorPosition)

        if pInvestorPositionCode != "0" {

            // 获得持仓结构体数据
            mInvestorPosition := GetInvestorPositionStruct(pInvestorPosition)

            if mInvestorPosition.Position != 0 {
                fmt.Printf("- 合约：%v   \t%v:%v   \t总持仓：%v   \t持仓均价：%v   \t持仓盈亏：%v\n", mInvestorPosition.InstrumentID, mInvestorPosition.PositionDateTitle, mInvestorPosition.PosiDirectionTitle, mInvestorPosition.Position, mInvestorPosition.OpenCost, mInvestorPosition.PositionProfit);
            }
        }

        if bIsLast {

            fmt.Println("-------------------------------------------------------------------------------------------------")

            if !Ctp.IsTraderInitFinish {
                // 交易程序初始化流程走完了
                Ctp.IsTraderInitFinish = true

                // 订阅行情
                Subscribe := []string{"rb2001"}
                MdSpi.SubscribeMarketData(Subscribe)
            }
        }
    }
}


// 报单通知（委托单）
func (p *FtdcTraderSpi) OnRtnOrder(pOrder goctp.CThostFtdcOrderField) {

    // 报单编号
    OrderSysID := pOrder.GetOrderSysID()

    // 报单状态
    OrderStatus := pOrder.GetOrderStatus()

    // 获得报单结构体数据
    mOrder := GetOrderListStruct(pOrder)

    // 报单列表数据 key 键
    mOrder.MapKey = pOrder.GetInstrumentID() + "_" + TrimSpace(pOrder.GetOrderSysID())

    if OrderSysID == "" {

        // 报单就自动撤单，且没有编号的 都视为报错
        if OrderStatus == goctp.THOST_FTDC_OST_Canceled {

            OrderErrorStr := "-------------------------------------------------------------------------------------------------\n" +
                "- 报单出错了\n" +
                "- 报单合约：" + mOrder.InstrumentID + "\t报单引用：" + mOrder.OrderRef + "\n" +
                "- 报单方向：" + mOrder.DirectionTitle + "   \t报单价格：" + Float64ToString(mOrder.LimitPrice) + "\n" +
                "- 报单开平：" + mOrder.CombOffsetFlagTitle + " \t报单数量：" + IntToString(mOrder.Volume) + "\n" +
                "- 错误代码：-1   \t错误消息：" + mOrder.StatusMsg + "\n" +
                "-------------------------------------------------------------------------------------------------"
            fmt.Println(OrderErrorStr)
        }

        return
    }

    // 未成交和撤单的报单（已成交的通知在 OnRtnTrade 函数中处理）
    if OrderStatus == goctp.THOST_FTDC_OST_NoTradeQueueing || OrderStatus == goctp.THOST_FTDC_OST_Canceled {

        OrderStr := "-------------------------------------------------------------------------------------------------\n" +
                    "- 报单通知 " + mOrder.InsertTime + "\n" +
                    "- 报单合约：" + mOrder.InstrumentID + " \t报单编号：" + TrimSpace(mOrder.OrderSysID) + "\n" +
                    "- 报单方向：" + mOrder.DirectionTitle + "   \t报单价格：" + Float64ToString(mOrder.LimitPrice) + "\n" +
                    "- 报单开平：" + mOrder.CombOffsetFlagTitle + " \t报单数量：" + IntToString(mOrder.Volume) + "\n" +
                    "- 报单状态：" + mOrder.OrderStatusTitle + " \t状态信息：" + mOrder.StatusMsg + "\n" +
                    "-------------------------------------------------------------------------------------------------"
        fmt.Println(OrderStr)
    }

    // 将报单数据记录下来
    MapOrderList.Store(mOrder.MapKey, mOrder)
}

// 成交通知（委托单在交易所成交了）
func (p *FtdcTraderSpi) OnRtnTrade(pTrade goctp.CThostFtdcTradeField) {

    // 报单方向
    DirectionTitle := GetDirectionTitle(string(pTrade.GetDirection()))

    // 报单开平
    OffsetFlagTitle := GetOffsetFlagTitle(string(pTrade.GetOffsetFlag()))

    OrderStr := "-------------------------------------------------------------------------------------------------\n" +
                "- 成交通知 " + pTrade.GetTradeTime() + "\n" +
                "- 成交合约：" + pTrade.GetInstrumentID() + "\t成交编号：" + TrimSpace(pTrade.GetTradeID()) + " \t报单编号：" + TrimSpace(pTrade.GetOrderSysID()) + "\n" +
                "- 成交方向：" + DirectionTitle + "   \t成交价格：" + Float64ToString(pTrade.GetPrice()) + "\n" +
                "- 成交开平：" + OffsetFlagTitle + " \t成交数量：" + IntToString(pTrade.GetVolume()) + "\n" +
                "-------------------------------------------------------------------------------------------------"
    fmt.Println(OrderStr)
}

// 报单出错响应（综合交易平台交易核心返回的包含错误信息的报单响应）
func (p *FtdcTraderSpi) OnRspOrderInsert(pInputOrder goctp.CThostFtdcInputOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    // 报单方向
    DirectionTitle := GetDirectionTitle(string(pInputOrder.GetDirection()))

    // 报单开平
    OffsetFlagTitle := GetOffsetFlagTitle(string(pInputOrder.GetCombOffsetFlag()))

    OrderStr := "-------------------------------------------------------------------------------------------------\n" +
                "- 报单出错了\n" +
                "- 报单合约：" + pInputOrder.GetInstrumentID() + "\t报单引用：" + pInputOrder.GetOrderRef() + "\n" +
                "- 报单方向：" + DirectionTitle + "   \t报单价格：" + Float64ToString(pInputOrder.GetLimitPrice()) + "\n" +
                "- 报单开平：" + OffsetFlagTitle + " \t报单数量：" + IntToString(pInputOrder.GetVolumeTotalOriginal()) + "\n" +
                "- 错误代码：" + string(pRspInfo.GetErrorID()) + "    \t错误消息：" + ConvertToString(pRspInfo.GetErrorMsg(), "gbk", "utf-8") + "\n" +
                "-------------------------------------------------------------------------------------------------"
    fmt.Println(OrderStr)
}

// 错误应答
func (p *FtdcTraderSpi) OnRspError(pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
    p.IsErrorRspInfo(pRspInfo)
}

// 报单操作错误回报
func (p *FtdcTraderSpi) OnErrRtnOrderAction(pOrderAction goctp.CThostFtdcOrderActionField, pRspInfo goctp.CThostFtdcRspInfoField) {
    p.IsErrorRspInfo(pRspInfo)
}

// 报单操作请求响应（撤单失败会触发）
func (p *FtdcTraderSpi) OnRspOrderAction(pInputOrderAction goctp.CThostFtdcInputOrderActionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
    p.IsErrorRspInfo(pRspInfo)
}

// 交易系统错误通知
func (p *FtdcTraderSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField) bool {

    rspInfo := fmt.Sprintf("%v", pRspInfo)

    // 容错处理 pRspInfo ，部分响应函数中，pRspInfo 为 0
    if rspInfo == "0" {
        return false

    } else {

        // 如果ErrorID != 0, 说明收到了错误的响应
        bResult := (pRspInfo.GetErrorID() != 0)
        if bResult {
            log.Printf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), ConvertToString(pRspInfo.GetErrorMsg(), "gbk", "utf-8"))
        }

        return bResult
    }
}

// 心跳超时警告。当长时间未收到报文时，该方法被调用。
func (p *FtdcTraderSpi) OnHeartBeatWarning(nTimeLapse int) {
    fmt.Println("心跳超时警告（OnHeartBeatWarning） nTimerLapse=", nTimeLapse)
}

// 开仓
func (p *FtdcTraderSpi) OrderOpen(Input InputOrderStruct) int {

    iRequestID := GetTraderRequestId()

    mInstrumentInfo, mapRes := GetInstrumentInfo(Input.InstrumentID)
    if !mapRes {
        fmt.Println("开仓失败，合约不存在！")
        return 0
    }

    req := goctp.NewCThostFtdcInputOrderField()

    // 经纪公司代码
    req.SetBrokerID(Ctp.BrokerID)
    // 投资者代码
    req.SetInvestorID(Ctp.InvestorID)
    // 合约代码
    req.SetInstrumentID(Input.InstrumentID)
    // 报单引用
    req.SetOrderRef(IntToString(iRequestID))
    // 买卖方向:买(THOST_FTDC_D_Buy),卖(THOST_FTDC_D_Sell)
    req.SetDirection(Input.Direction)
    // 交易所代码
    req.SetExchangeID(mInstrumentInfo.ExchangeID)
    // 组合开平标志: 开仓
    req.SetCombOffsetFlag(string(goctp.THOST_FTDC_OF_Open))
    // 组合投机套保标志: 投机
    req.SetCombHedgeFlag(string(goctp.THOST_FTDC_HF_Speculation))
    // 报单价格条件: 限价
    req.SetOrderPriceType(goctp.THOST_FTDC_OPT_LimitPrice)
    // 价格
    req.SetLimitPrice(Input.Price)
    // 数量
    req.SetVolumeTotalOriginal(Input.Volume)
    // 有效期类型: 当日有效
    req.SetTimeCondition(goctp.THOST_FTDC_TC_GFD)
    // 成交量类型: 任何数量
    req.SetVolumeCondition(goctp.THOST_FTDC_VC_AV)
    // 最小成交量
    req.SetMinVolume(1)
    // 触发条件: 立即
    req.SetContingentCondition(goctp.THOST_FTDC_CC_Immediately)
    // 强平原因: 非强平
    req.SetForceCloseReason(goctp.THOST_FTDC_FCC_NotForceClose)
    // 自动挂起标志: 否
    req.SetIsAutoSuspend(0)
    // 用户强评标志: 否
    req.SetUserForceClose(0)

    iResult := Ctp.TraderApi.ReqOrderInsert(req, iRequestID)

    if iResult != 0 {
        ReqFailMsg("提交报单失败！", iResult)
        return 0
    }

    return iRequestID
}

// 平仓
func (p *FtdcTraderSpi) OrderClose(Input InputOrderStruct) int {

    iRequestID := GetTraderRequestId()

    mInstrumentInfo, mapRes := GetInstrumentInfo(Input.InstrumentID)
    if !mapRes {
        fmt.Println("平仓失败，合约不存在！")
        return 0
    }

    // 没有设置平仓类型
    if Input.CombOffsetFlag == 0 {

        if mInstrumentInfo.ExchangeID == "SHFE" {
            // 上期所（默认平今仓）
            Input.CombOffsetFlag = goctp.THOST_FTDC_OF_CloseToday
        } else {
            // 非上期所，不用区分今昨仓，直接使用平仓即可
            Input.CombOffsetFlag = goctp.THOST_FTDC_OF_Close
        }
    }

    req := goctp.NewCThostFtdcInputOrderField()

    // 经纪公司代码
    req.SetBrokerID(Ctp.BrokerID)
    // 投资者代码
    req.SetInvestorID(Ctp.InvestorID)
    // 合约代码
    req.SetInstrumentID(Input.InstrumentID)
    // 报单引用
    req.SetOrderRef(IntToString(iRequestID))
    // 买卖方向:买(THOST_FTDC_D_Buy),卖(THOST_FTDC_D_Sell)
    req.SetDirection(Input.Direction)
    // 交易所代码
    req.SetExchangeID(mInstrumentInfo.ExchangeID)
    // 组合开平标志: 平仓 (针对上期所，区分昨仓、今仓)
    req.SetCombOffsetFlag(string(Input.CombOffsetFlag))
    // 组合投机套保标志: 投机
    req.SetCombHedgeFlag(string(goctp.THOST_FTDC_HF_Speculation))
    // 报单价格条件: 限价
    req.SetOrderPriceType(goctp.THOST_FTDC_OPT_LimitPrice)
    // 价格
    req.SetLimitPrice(Input.Price)
    // 数量
    req.SetVolumeTotalOriginal(Input.Volume)
    // 有效期类型: 当日有效
    req.SetTimeCondition(goctp.THOST_FTDC_TC_GFD)
    // 成交量类型: 任何数量
    req.SetVolumeCondition(goctp.THOST_FTDC_VC_AV)
    // 最小成交量
    req.SetMinVolume(1)
    // 触发条件: 立即
    req.SetContingentCondition(goctp.THOST_FTDC_CC_Immediately)
    // 强平原因: 非强平
    req.SetForceCloseReason(goctp.THOST_FTDC_FCC_NotForceClose)
    // 自动挂起标志: 否
    req.SetIsAutoSuspend(0)
    // 用户强评标志: 否
    req.SetUserForceClose(0)

    iResult := Ctp.TraderApi.ReqOrderInsert(req, iRequestID)

    if iResult != 0 {
        ReqFailMsg("提交报单失败！", iResult)
        return 0
    }

    return iRequestID
}

// 撤消报单
func (p *FtdcTraderSpi) OrderCancel (InstrumentID string, OrderSysID string) int {

    iRequestID := GetTraderRequestId()

    mapKey := InstrumentID + "_" + OrderSysID

    // 检查报单列表数据是否存在
    mOrderV, mapRes := MapOrderList.Load(mapKey)
    if !mapRes {
        fmt.Printf("撤消报单失败：合约 %v 报单编号 %v 不存在！\n", InstrumentID, OrderSysID)
        return 0
    }

    mOrder := mOrderV.(OrderListStruct)

    req := goctp.NewCThostFtdcInputOrderActionField()

    // 经纪公司代码
    req.SetBrokerID(mOrder.BrokerID)
    // 投资者代码
    req.SetInvestorID(mOrder.InvestorID)
    // 合约代码
    req.SetInstrumentID(InstrumentID)
    // 报单引用
    req.SetOrderRef(mOrder.OrderRef)
    // 交易所代码
    req.SetExchangeID(mOrder.ExchangeID)
    // 前置编号
    req.SetFrontID(mOrder.FrontID)
    // 会话编号
    req.SetSessionID(mOrder.SessionID)
    // 报单编号
    req.SetOrderSysID(mOrder.OrderSysID)
    // 操作标志
    req.SetActionFlag(goctp.THOST_FTDC_AF_Delete)

    iResult := Ctp.TraderApi.ReqOrderAction(req, iRequestID)

    if iResult != 0 {
        ReqFailMsg("提交报单失败！", iResult)
        return 0
    }

    return iRequestID
}

// 获得报单结构体数据
func GetOrderListStruct(pOrder goctp.CThostFtdcOrderField) OrderListStruct {

    var mOrder OrderListStruct

    mOrder.BrokerID            = pOrder.GetBrokerID()
    mOrder.InvestorID          = pOrder.GetInvestorID()
    mOrder.InstrumentID        = pOrder.GetInstrumentID()
    mOrder.ExchangeID          = pOrder.GetExchangeID()
    mOrder.FrontID             = pOrder.GetFrontID()
    mOrder.OrderRef            = pOrder.GetOrderRef()
    mOrder.SessionID           = pOrder.GetSessionID()
    mOrder.InsertTime          = pOrder.GetInsertTime()
    mOrder.OrderSysID          = pOrder.GetOrderSysID()
    mOrder.LimitPrice          = pOrder.GetLimitPrice()
    mOrder.Volume              = pOrder.GetVolumeTotalOriginal()
    mOrder.Direction           = string(pOrder.GetDirection())
    mOrder.CombOffsetFlag      = string(pOrder.GetCombOffsetFlag())
    mOrder.CombHedgeFlag       = string(pOrder.GetCombHedgeFlag())
    mOrder.OrderStatus         = string(pOrder.GetOrderStatus())
    mOrder.StatusMsg           = ConvertToString(pOrder.GetStatusMsg(), "gbk", "utf-8")
    mOrder.DirectionTitle      = GetDirectionTitle(mOrder.Direction)
    mOrder.OrderStatusTitle    = GetOrderStatusTitle(mOrder.OrderStatus)
    mOrder.CombOffsetFlagTitle = GetOffsetFlagTitle(mOrder.CombOffsetFlag)

    return mOrder
}

// 获得持仓结构体数据
func GetInvestorPositionStruct (pInvestorPosition goctp.CThostFtdcInvestorPositionField) InvestorPositionStruct {

    var mInvestorPosition InvestorPositionStruct

    // 检查合约详情是否存在
    mInstrumentInfo, mapRes := GetInstrumentInfo(pInvestorPosition.GetInstrumentID())
    if !mapRes {
        fmt.Printf("合约 %v 不存在！\n", pInvestorPosition.GetInstrumentID())
        return mInvestorPosition
    }

    // 合约乘数
    var VolumeMultiple int = mInstrumentInfo.VolumeMultiple

    // 开仓成本
    var OpenCost float64 = pInvestorPosition.GetOpenCost() / float64(pInvestorPosition.GetPosition() * VolumeMultiple)

    mInvestorPosition.BrokerID           = pInvestorPosition.GetBrokerID()
    mInvestorPosition.InvestorID         = pInvestorPosition.GetInvestorID()
    mInvestorPosition.InstrumentID       = pInvestorPosition.GetInstrumentID()
    mInvestorPosition.InstrumentName     = mInstrumentInfo.InstrumentName
    mInvestorPosition.PosiDirection      = string(pInvestorPosition.GetPosiDirection())
    mInvestorPosition.PosiDirectionTitle = GetPosiDirectionTitle(mInvestorPosition.PosiDirection)
    mInvestorPosition.HedgeFlag          = string(pInvestorPosition.GetHedgeFlag())
    mInvestorPosition.HedgeFlagTitle     = GetHedgeFlagTitle(mInvestorPosition.HedgeFlag)
    mInvestorPosition.PositionDate       = string(pInvestorPosition.GetPositionDate())
    mInvestorPosition.PositionDateTitle  = GetPositionDateTitle(mInvestorPosition.PositionDate)
    mInvestorPosition.Position           = pInvestorPosition.GetPosition()
    mInvestorPosition.YdPosition         = pInvestorPosition.GetYdPosition()
    mInvestorPosition.TodayPosition      = pInvestorPosition.GetTodayPosition()
    mInvestorPosition.LongFrozen         = pInvestorPosition.GetLongFrozen()
    mInvestorPosition.ShortFrozen        = pInvestorPosition.GetShortFrozen()

    // 冻结的持仓量（多空并成一个字段）
    if mInvestorPosition.PosiDirection == string(goctp.THOST_FTDC_PD_Long) {
        // 多头冻结的持仓量
        mInvestorPosition.ShortVolume = pInvestorPosition.GetShortFrozen()
    } else {
        // 空头冻结的持仓量
        mInvestorPosition.ShortVolume = pInvestorPosition.GetLongFrozen()
    }

    mInvestorPosition.OpenVolume         = pInvestorPosition.GetOpenVolume()
    mInvestorPosition.CloseVolume        = pInvestorPosition.GetCloseVolume()
    mInvestorPosition.PositionCost       = Decimal(pInvestorPosition.GetPositionCost(), 2)
    mInvestorPosition.Commission         = Decimal(pInvestorPosition.GetCommission(), 2)
    mInvestorPosition.CloseProfit        = pInvestorPosition.GetCloseProfit()
    mInvestorPosition.PositionProfit     = Decimal(pInvestorPosition.GetPositionProfit(), 2)
    mInvestorPosition.PreSettlementPrice = pInvestorPosition.GetPreSettlementPrice()
    mInvestorPosition.SettlementPrice    = Decimal(pInvestorPosition.GetSettlementPrice(), 2)
    mInvestorPosition.SettlementID       = pInvestorPosition.GetSettlementID()
    mInvestorPosition.OpenCost           = Decimal(OpenCost, 2)
    mInvestorPosition.ExchangeID         = pInvestorPosition.GetExchangeID()

    return mInvestorPosition
}

// 获得合约详情信息
func GetInstrumentInfo(InstrumentID string) (InstrumentInfoStruct, bool) {
    if v, ok := MapInstrumentInfos.Load(InstrumentID); ok {
        return v.(InstrumentInfoStruct), true
    } else {
        var mInstrumentInfo InstrumentInfoStruct
        return mInstrumentInfo, false
    }
}

// 获得期货合约列表【只有期货，不含期权】
func GetFuturesList() []string {
    var InstrumentList []string
    MapInstrumentInfos.Range(func(k, v interface{}) bool {

        val := v.(InstrumentInfoStruct)

        // 类型为期货的合约
        if val.ProductClass == "1" {
            InstrumentList = append(InstrumentList, val.InstrumentID)
        }

        return true
    })

    return InstrumentList
}

/**
 * 计算盈亏
 *
 * @param   InstrumentID  string  合约
 * @param   OpenPrice     float64 开仓价格
 * @param   LastPrice     float64 最新价|平仓价格
 * @param   Number        int     数量
 * @param   PosiDirection string  持仓方向[2：买，3：卖]
 */
func GetPositionProfit(InstrumentID string, OpenPrice float64, LastPrice float64, Number int, PosiDirection string) float64 {

    InstrumentInfo, _ := GetInstrumentInfo(InstrumentID)
    if PosiDirection == "2" {
        return ((LastPrice - OpenPrice) * float64(InstrumentInfo.VolumeMultiple)) * float64(Number)
    } else {
        return ((OpenPrice - LastPrice) * float64(InstrumentInfo.VolumeMultiple)) * float64(Number)
    }
}