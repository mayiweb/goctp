/**
 * 交易模块
 */

package main

import (
    "gitee.com/mayiweb/goctp"
    "regexp"
    "strings"
)

// 获得交易请求编号
func GetTradeRequestId() int {
    Ctp.TradeRequestId += 1
    return Ctp.TradeRequestId
}

// 发送请求日志（仅查询类的函数需要调用）
func (p *FtdcTradeSpi) ReqMsg(Msg string) {

    defer CheckPanic()

    // 交易程序未初始化完成时，执行查询类的函数需要有1.5秒间隔
    if !Ctp.IsTradeInitFinish {
        Sleep(1500)
    }

    Println("")
    LogPrintln(Msg)
}

// 打印错误
func (p *FtdcTradeSpi) PrintError(Msg string, nRequestID int) {

    defer CheckPanic()

    LogPrintln(Msg)
}

// 交易系统错误通知
func (p *FtdcTradeSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int) bool {

    defer CheckPanic()

    // 容错处理 pRspInfo ，部分响应函数中，pRspInfo 为空
    if IsNullPointer(pRspInfo) {

        return false

    } else {

        // 如果ErrorID != 0, 说明收到了错误的响应
        bResult := (pRspInfo.GetErrorID() != 0)
        if bResult {

            Msg := Sprintf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), TrimSpace(GbkToUtf8(pRspInfo.GetErrorMsg())))

            // 输出日志
            LogPrintf(Msg)

            // 如果初始化未完成则退出程序，防止卡在某个流程
            if !Ctp.IsTradeInitFinish {
                Exit()
            }
        }

        return bResult
    }
}


// 当客户端与交易后台通信连接断开时，该方法被调用。当发生这个情况后，API会自动重新连接，客户端可不做处理。
// 服务器已断线，该函数也会被调用。【api 会自动初始化程序，并重新登录】
func (p *FtdcTradeSpi) OnFrontDisconnected(nReason int) {

    defer CheckPanic()

    if Ctp.IsTradeInit {

        Ctp.IsTradeInit       = false
        Ctp.IsTradeLogin      = false
        Ctp.IsTradeInitFinish = false

        LogPrintln("交易服务器已断线，尝试重新连接中...")
    }
}

// 当客户端与交易后台建立起通信连接时（还未登录前），该方法被调用。
func (p *FtdcTradeSpi) OnFrontConnected() {

    defer CheckPanic()

    Str := "-------------------------------------------------------------------------------------------------\n" +
        "- 交易系统初始化成功，Api 版本：" + goctp.CThostFtdcTraderApiGetApiVersion() + "\n" +
        "-------------------------------------------------------------------------------------------------"
    Println(Str)

    Ctp.IsTradeInit = true

    // 填写了 AppID 与 AuthCode 则进行客户端认证
    if Ctp.AppID != "" && Ctp.AuthCode != "" {
        p.ReqAuthenticate()
    } else {
        p.ReqUserLogin()
    }
}

// 客户端认证
func (p *FtdcTradeSpi) ReqAuthenticate() {

    defer CheckPanic()

    LogPrintln("客户端认证中...")

    req := goctp.NewCThostFtdcReqAuthenticateField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetUserID(Ctp.InvestorID)
    req.SetAppID(Ctp.AppID)
    req.SetAuthCode(Ctp.AuthCode)

    iResult := Ctp.TradeApi.ReqAuthenticate(req, GetTradeRequestId())

    if iResult != 0 {
        ReqFailMsg("发送客户端认证请求失败！", iResult)
    }
}

// 客户端认证响应
func (p *FtdcTradeSpi) OnRspAuthenticate(pRspAuthenticateField goctp.CThostFtdcRspAuthenticateField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if bIsLast && !p.IsErrorRspInfo(pRspInfo, nRequestID) {

        LogPrintln("客户端认证成功")

        p.ReqUserLogin()
    }
}

// 用户登录请求
func (p *FtdcTradeSpi) ReqUserLogin() {

    defer CheckPanic()

    Sleep(1000)

    LogPrintln("交易系统登录中...")

    req := goctp.NewCThostFtdcReqUserLoginField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetUserID(Ctp.InvestorID)
    req.SetPassword(Ctp.Password)

    iResult := Ctp.TradeApi.ReqUserLogin(req, GetTradeRequestId())

    if iResult != 0 {
        ReqFailMsg("发送交易系统登录请求失败！", iResult)
    }
}

// 用户登录响应
func (p *FtdcTradeSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if !p.IsErrorRspInfo(pRspInfo, nRequestID) {

        if bIsLast {

            Ctp.IsTradeLogin = true

            // 登录行情账号
            MdSpi.ReqUserLogin()

            LogPrintf("交易账号已登录，当前交易日：%v\n", Ctp.TradeApi.GetTradingDay())

            p.ReqSettlementInfoConfirm()
        }
    }
}

// 投资者结算单确认
func (p *FtdcTradeSpi) ReqSettlementInfoConfirm() int {

    defer CheckPanic()

    p.ReqMsg("结算单确认中...")

    req := goctp.NewCThostFtdcSettlementInfoConfirmField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TradeApi.ReqSettlementInfoConfirm(req, GetTradeRequestId())

    if iResult != 0 {
        ReqFailMsg("确认结算单失败！", iResult)
    }

    return iResult
}

// 发送投资者结算单确认响应
func (p *FtdcTradeSpi) OnRspSettlementInfoConfirm(pSettlementInfoConfirm goctp.CThostFtdcSettlementInfoConfirmField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if !p.IsErrorRspInfo(pRspInfo, nRequestID) {

        if bIsLast {

            LogPrintln("结算单确认成功")

            // 请求查询合约
            p.ReqQryInstrument()
        }
    }
}

// 请求查询合约
func (p *FtdcTradeSpi) ReqQryInstrument() int {

    defer CheckPanic()

    p.ReqMsg("查询合约中...")

    req := goctp.NewCThostFtdcQryInstrumentField()
    req.SetInstrumentID("")

    iResult := Ctp.TradeApi.ReqQryInstrument(req, GetTradeRequestId())

    if iResult != 0 {
        ReqFailMsg("查询合约失败！", iResult)
    }

    return iResult
}

// 请求查询合约响应
func (p *FtdcTradeSpi) OnRspQryInstrument(pInstrument goctp.CThostFtdcInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if !p.IsErrorRspInfo(pRspInfo, nRequestID) {

        if !IsNullPointer(pInstrument) {
            GetInstrumentStruct(pInstrument)
        }

        if bIsLast {

            LogPrintf("获得合约记录 %v 条\n", MapInstruments.Size())

            if !Ctp.IsTradeInitFinish {
                // 请求查询资金账户
                p.ReqQryTradingAccount()
            }
        }
    }
}

// 请求查询资金账户
func (p *FtdcTradeSpi) ReqQryTradingAccount() int {

    defer CheckPanic()

    p.ReqMsg("查询资金账户中...")

    req := goctp.NewCThostFtdcQryTradingAccountField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TradeApi.ReqQryTradingAccount(req, GetTradeRequestId())

    if iResult != 0 {
        ReqFailMsg("查询资金账户失败！", iResult)
    }

    return iResult
}

// 请求查询资金账户响应
func (p *FtdcTradeSpi) OnRspQryTradingAccount(pTradingAccount goctp.CThostFtdcTradingAccountField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if !p.IsErrorRspInfo(pRspInfo, nRequestID) {

        if !IsNullPointer(pTradingAccount) {
            GetAccountStruct(pTradingAccount)
        }

        if (bIsLast) {

            mAccount := GetAccount()

            Str := "-------------------------------------------------------------------------------------------------\n" +
                "- 资金账号：%v\n" +
                "- 期初资金：%.2f \n" +
                "- 动态权益：%.2f \n" +
                "- 可用资金：%.2f \n" +
                "- 持仓盈亏：%.2f \n" +
                "- 平仓盈亏：%.2f \n" +
                "- 手续费  ：%.2f \n" +
                "-------------------------------------------------------------------------------------------------"

            Println(Sprintf(Str, mAccount.AccountID, mAccount.PreBalance, mAccount.Balance, mAccount.Available, mAccount.PositionProfit, mAccount.CloseProfit, mAccount.Commission))

            if !Ctp.IsTradeInitFinish {
                // 请求查询投资者报单
                p.ReqQryOrder()
            }
        }
    }
}

// 请求查询报单（委托单）
func (p *FtdcTradeSpi) ReqQryOrder() int {

    defer CheckPanic()

    p.ReqMsg("查询报单中...")

    req := goctp.NewCThostFtdcQryOrderField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TradeApi.ReqQryOrder(req, GetTradeRequestId())

    if iResult != 0 {
        ReqFailMsg("查询报单失败！", iResult)
    }

    return iResult
}

// 请求查询报单响应
func (p *FtdcTradeSpi) OnRspQryOrder(pOrder goctp.CThostFtdcOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if !p.IsErrorRspInfo(pRspInfo, nRequestID) {

        if !IsNullPointer(pOrder) {

            // 只记录有报单编号的报单数据
            if pOrder.GetOrderSysID() != "" {

                // 获得报单结构体数据
                GetOrderStruct(pOrder)
            }
        }

        if bIsLast {

            Str := "-------------------------------------------------------------------------------------------------\n"

            MapOrderNoTradeSize := 0

            mOrders := MapOrders.GetAll()
            for _, v := range mOrders {
                val := v.(OrderStruct)

                // 输出 未成交、部分成交 的报单
                if val.OrderStatus == string(goctp.THOST_FTDC_OST_NoTradeQueueing) || val.OrderStatus == string(goctp.THOST_FTDC_OST_PartTradedQueueing) {
                    MapOrderNoTradeSize += 1
                    Str += Sprintf("- 合约：%v %v:%v 数量：%v 价格：%v 报单编号：%v (%v)\n", StrAfterSpace(val.InstrumentID, 16), val.DirectionTitle, StrAfterSpace(val.CombOffsetFlagTitle, 11), StrAfterSpace(IntToString(val.Volume), 6), StrAfterSpace(Float64ToString(val.LimitPrice), 12), TrimSpace(val.OrderSysID), val.OrderStatusTitle)
                }
            }

            Str += Sprintf("- 共有报单记录 %v 条，未成交 %v 条\n", MapOrders.Size(), MapOrderNoTradeSize)
            Str += "-------------------------------------------------------------------------------------------------"
            Println(Str)

            if !Ctp.IsTradeInitFinish {
                // 请求查询成交
                p.ReqQryTrade()
            }
        }
    }
}

// 请求查询成交
func (p *FtdcTradeSpi) ReqQryTrade() int {

    defer CheckPanic()

    p.ReqMsg("查询成交中...")

    req := goctp.NewCThostFtdcQryTradeField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TradeApi.ReqQryTrade(req, GetTradeRequestId())

    if iResult != 0 {
        ReqFailMsg("查询成交失败！", iResult)
    }

    return iResult
}

// 请求查询成交响应
func (p *FtdcTradeSpi) OnRspQryTrade(pTrade goctp.CThostFtdcTradeField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if !p.IsErrorRspInfo(pRspInfo, nRequestID) {

        if !IsNullPointer(pTrade) {
            GetTradeStruct(pTrade)
        }

        if bIsLast {

            Str := "-------------------------------------------------------------------------------------------------\n"
            Str += Sprintf("- 共有成交记录 %v 条\n", MapTrades.Size())
            Str += "-------------------------------------------------------------------------------------------------"
            Println(Str)

            if !Ctp.IsTradeInitFinish {
                // 请求查询投资者持仓汇总
                // p.ReqQryInvestorPosition()

                // 请求查询投资者持仓明细
                p.ReqQryInvestorPosition()
            }
        }
    }
}

// 请求查询投资者持仓（汇总）
func (p *FtdcTradeSpi) ReqQryInvestorPosition() int {

    defer CheckPanic()

    p.ReqMsg("查询持仓汇总中...")

    req := goctp.NewCThostFtdcQryInvestorPositionField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TradeApi.ReqQryInvestorPosition(req, GetTradeRequestId())

    if iResult != 0 {
        ReqFailMsg("查询持仓汇总失败！", iResult)
    }

    return iResult
}

// 请求查询投资者持仓（汇总）响应
func (p *FtdcTradeSpi) OnRspQryInvestorPosition(pPosition goctp.CThostFtdcInvestorPositionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if !p.IsErrorRspInfo(pRspInfo, nRequestID) {

        if !IsNullPointer(pPosition) {

            // 获得持仓结构体数据
            GetPositionStruct(pPosition)
        }

        if bIsLast {

            PositionStr := "-------------------------------------------------------------------------------------------------\n"

            // 输出当前持仓
            mPositions := MapPositions.GetAll()
            for _, v := range mPositions {

                val := v.(PositionStruct)

                if val.Volume != 0 {
                    PositionStr += Sprintf("- 合约：%v %v:%v 持仓：%v 均价：%v 盈亏：%v\n", StrAfterSpace(val.InstrumentID, 16), val.PositionDateTitle, StrAfterSpace(val.DirectionTitle, 9), StrAfterSpace(IntToString(val.Volume), 7), StrAfterSpace(Float64ToString(val.OpenPrice), 12), val.PositionProfit)
                }
            }

            PositionStr += "-------------------------------------------------------------------------------------------------"
            Println(PositionStr)


            if !Ctp.IsTradeInitFinish {

                // 交易程序初始化流程走完了
                Ctp.IsTradeInitFinish = true

                // 订阅行情
                Subscribe := []string{"rb2410", "sc2401"}
                MdSpi.SubscribeTick(Subscribe)
            }
        }
    }
}


// 请求查询投资者持仓明细
func (p *FtdcTradeSpi) ReqQryInvestorPositionDetail() int {

    defer CheckPanic()

    p.ReqMsg("查询持仓明细中...")

    req := goctp.NewCThostFtdcQryInvestorPositionDetailField()
    req.SetBrokerID(Ctp.BrokerID)
    req.SetInvestorID(Ctp.InvestorID)

    iResult := Ctp.TradeApi.ReqQryInvestorPositionDetail(req, GetTradeRequestId())

    if iResult != 0 {
        ReqFailMsg("查询持仓明细失败！", iResult)
    }

    return iResult
}

// 请求查询投资者持仓明细响应
func (p *FtdcTradeSpi) OnRspQryInvestorPositionDetail(pPositionDetail goctp.CThostFtdcInvestorPositionDetailField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if !p.IsErrorRspInfo(pRspInfo, nRequestID) {

        if !IsNullPointer(pPositionDetail) {

            // 获得持仓明细结构体数据
            GetPositionDetailStruct(pPositionDetail)
        }

        if bIsLast {

            PositionStr := "-------------------------------------------------------------------------------------------------\n"

            // 当前持仓明细
            mPositionDetails := MapPositionDetails.GetAll()
            for _, v := range mPositionDetails {

                val := v.(PositionDetailStruct)

                if val.Volume != 0 {
                    PositionStr += Sprintf("- 合约：%v %v:%v 数量：%v 价格：%v 盈亏：%v\n", StrAfterSpace(val.InstrumentID, 16), val.DirectionTitle, StrAfterSpace(val.PositionDateTitle, 14), StrAfterSpace(IntToString(val.Volume), 7), StrAfterSpace(Float64ToString(val.OpenPrice), 14), val.PositionProfitByTrade)
                }
            }

            PositionStr += "-------------------------------------------------------------------------------------------------"
            Println(PositionStr)


            if !Ctp.IsTradeInitFinish {

                // 交易程序初始化流程走完了
                // Ctp.IsTradeInitFinish = true

                // 订阅行情
                // Subscribe := []string{"rb2410"}
                // MdSpi.SubscribeMarketData(Subscribe)
            }
        }
    }
}


// 报单通知（委托单）
func (p *FtdcTradeSpi) OnRtnOrder(pOrder goctp.CThostFtdcOrderField) {

    defer CheckPanic()

    if !IsNullPointer(pOrder) {

        // 报单状态
        OrderStatus := pOrder.GetOrderStatus()

        // 获得报单结构体数据
        sOrder := GetOrderStruct(pOrder)

        if sOrder.OrderSysID == "" {

            // 报单就自动撤单，且没有编号的 都视为报错
            if OrderStatus == goctp.THOST_FTDC_OST_Canceled {

                OrderErrorStr := "\n-------------------------------------------------------------------------------------------------\n" +
                    "- 报单错误  " + sOrder.InstrumentName + "（" + sOrder.InstrumentID + "）\n" +
                    "- 报单时间：" + sOrder.InsertTime + "      报单引用：" + sOrder.OrderRef + "\n" +
                    "- 报单方向：" + sOrder.DirectionTitle + "            报单价格：" + Float64ToString(sOrder.LimitPrice) + "\n" +
                    "- 报单开平：" + sOrder.CombOffsetFlagTitle + "          报单数量：" + IntToString(sOrder.Volume) + "\n" +
                    "- 错误消息：" + sOrder.StatusMsg + "\n" +
                    "-------------------------------------------------------------------------------------------------"
                Println(OrderErrorStr)

                // 更新持仓数据及资金账户
                AddQueryTask("ReqQryInvestorPosition")
                AddQueryTask("ReqQryTradingAccount")
            }

            return
        }

        // 未成交和撤单的报单（已成交的通知在 OnRtnTrade 函数中处理）
        if OrderStatus == goctp.THOST_FTDC_OST_NoTradeQueueing || OrderStatus == goctp.THOST_FTDC_OST_Canceled {


            OrderStr := "\n-------------------------------------------------------------------------------------------------\n" +
                "- 报单通知  " + sOrder.InstrumentName + "（" + sOrder.InstrumentID + "）\n" +
                "- 报单时间：" + sOrder.InsertTime + "      报单编号：" + StrAfterSpace(TrimSpace(sOrder.OrderSysID), 20) + "报单引用：" + TrimSpace(sOrder.OrderRef) + "\n" +
                "- 报单方向：" + sOrder.DirectionTitle + "            报单价格：" + Float64ToString(sOrder.LimitPrice) + "\n" +
                "- 报单开平：" + sOrder.CombOffsetFlagTitle + "          报单数量：" + IntToString(sOrder.Volume) + "\n" +
                "- 报单状态：" + sOrder.OrderStatusTitle + "        状态信息：" + sOrder.StatusMsg + "\n" +
                "-------------------------------------------------------------------------------------------------"
            Println(OrderStr)

            // 更新资金账户
            AddQueryTask("ReqQryTradingAccount")
        }
    }
}

// 成交通知（委托单在交易所成交了）
func (p *FtdcTradeSpi) OnRtnTrade(pTrade goctp.CThostFtdcTradeField) {

    defer CheckPanic()

    if !IsNullPointer(pTrade) {

        // 获得报单成交结构体数据
        sTrade := GetTradeStruct(pTrade)

        Str := "\n-------------------------------------------------------------------------------------------------\n" +
            "- 成交通知  " + sTrade.InstrumentName + "（" + sTrade.InstrumentID + "）\n" +
            "- 成交时间：" + sTrade.TradeTime + "      报单编号：" + StrAfterSpace(sTrade.OrderSysID, 20) + "成交编号：" + sTrade.TradeID + "\n" +
            "- 成交方向：" + sTrade.DirectionTitle + "            成交价格：" + Float64ToString(sTrade.Price) + "\n" +
            "- 成交开平：" + sTrade.OffsetFlagTitle + "          成交数量：" + IntToString(sTrade.Volume) + "\n" +
            "-------------------------------------------------------------------------------------------------"
        Println(Str)

        // 更新持仓数据及资金账户
        AddQueryTask("ReqQryInvestorPosition")
        AddQueryTask("ReqQryTradingAccount")
    }
}

// 开仓
func (p *FtdcTradeSpi) OrderOpen(Input InputOrderStruct) int {

    defer CheckPanic()

    // 报单引用
    OrderRef := 0

    // 限价单最大报单量
    MaxLimitOrderVolume := 20

    // 组合开平标志: 开仓
    Input.CombOffsetFlag = goctp.THOST_FTDC_OF_Open

    // 设置限价单最大报单量
    mInstrument, mInstrumentOk := GetInstrumentInfo(Input.InstrumentID)
    if mInstrumentOk {
        MaxLimitOrderVolume = mInstrument.MaxLimitOrderVolume
    }

    // 如果持仓数量小于1则直接使用交易所错误提示
    if Input.Volume < 1 {
        OrderRef = p.ReqOrderInsert(Input)
    }

    // 自动拆单（每次下单数不能超过数量限额）
    for volume := Input.Volume; volume > 0; volume -= MaxLimitOrderVolume {

        // 剩余数量大于 最大报单量，则使用最大报单量
        if volume > MaxLimitOrderVolume {
            Input.Volume = MaxLimitOrderVolume
        } else {
            Input.Volume = volume
        }

        OrderRef = p.ReqOrderInsert(Input)
    }

    return OrderRef
}

// 平仓
func (p *FtdcTradeSpi) OrderClose(Input InputOrderStruct) int {

    defer CheckPanic()

    // 报单引用
    OrderRef := 0

    // 限价单最大报单量
    MaxLimitOrderVolume := 20

    // 设置限价单最大报单量
    mInstrument, mInstrumentOk := GetInstrumentInfo(Input.InstrumentID)
    if mInstrumentOk {
        MaxLimitOrderVolume = mInstrument.MaxLimitOrderVolume
    }

    // 没有设置平仓类型（默认值）
    if Input.CombOffsetFlag == 0 {

        if mInstrument.ExchangeID == "SHFE" || mInstrument.ExchangeID == "INE" {
            // 上期所（默认平今仓）
            Input.CombOffsetFlag = goctp.THOST_FTDC_OF_CloseToday
        } else {
            // 非上期所，不用区分今昨仓，直接使用平仓即可
            Input.CombOffsetFlag = goctp.THOST_FTDC_OF_Close
        }
    }

    // 如果持仓数量小于1则直接使用交易所错误提示
    if Input.Volume < 1 {
        OrderRef = p.ReqOrderInsert(Input)
    }

    // 自动拆单（每次下单数不能超过数量限额）
    for volume := Input.Volume; volume > 0; volume -= MaxLimitOrderVolume {

        // 剩余数量大于 最大报单量，则使用最大报单量
        if volume > MaxLimitOrderVolume {
            Input.Volume = MaxLimitOrderVolume
        } else {
            Input.Volume = volume
        }

        OrderRef = p.ReqOrderInsert(Input)
    }

    return OrderRef
}

// 报单录入
func (p *FtdcTradeSpi) ReqOrderInsert(Input InputOrderStruct) int {

    defer CheckPanic()

    iRequestID := GetTradeRequestId()

    // 设置报单引用
    Input.OrderRef = iRequestID

    mInstrument, mInstrumentOk := GetInstrumentInfo(Input.InstrumentID)
    if !mInstrumentOk {
        p.PrintError(Sprintf("合约 %v 不存在，禁止报单录入！", Input.InstrumentID), iRequestID)
        return iRequestID
    }

    if Input.Volume < 1 {
        p.PrintError(Sprintf("报单数量不能小于1，当前数值为 %v", Input.Volume), iRequestID)
        return iRequestID
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
    req.SetExchangeID(mInstrument.ExchangeID)
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

    iResult := Ctp.TradeApi.ReqOrderInsert(req, iRequestID)

    if iResult != 0 {
        ReqFailMsg("报单录入失败！", iResult)
        return 0
    }

    return iRequestID
}

// 撤消报单
func (p *FtdcTradeSpi) OrderCancel(InstrumentID string, OrderSysID string) int {

    defer CheckPanic()

    iRequestID := GetTradeRequestId()

    // 检查报单数据是否存在
    mOrder, mOrderOk := GetOrder(InstrumentID, OrderSysID)
    if !mOrderOk {
        p.PrintError(Sprintf("撤单失败：合约 %v 报单编号 %v 不存在！", InstrumentID, OrderSysID), iRequestID)
        return iRequestID
    }

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

    iResult := Ctp.TradeApi.ReqOrderAction(req, iRequestID)

    if iResult != 0 {
        ReqFailMsg("提交报单失败！", iResult)
        return 0
    }

    return iRequestID
}

// 报单出错响应（综合交易平台交易核心返回的包含错误信息的报单响应）
func (p *FtdcTradeSpi) OnRspOrderInsert(pInputOrder goctp.CThostFtdcInputOrderField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    defer CheckPanic()

    if !p.IsErrorRspInfo(pRspInfo, nRequestID) && !IsNullPointer(pInputOrder) {

        // 报单方向
        DirectionTitle := GetDirectionTitle(string(pInputOrder.GetDirection()))

        // 报单开平
        OffsetFlagTitle := GetOffsetFlagTitle(string(pInputOrder.GetCombOffsetFlag()))

        Str := "\n-------------------------------------------------------------------------------------------------\n" +
            "- 报单错误  " + GetInstrumentName(pInputOrder.GetInstrumentID()) + "（" + pInputOrder.GetInstrumentID() + "）\n" +
            "- 报单时间：" + DateToTime(GetCurrentTime()) + "      报单引用：" + pInputOrder.GetOrderRef() + "\n" +
            "- 报单方向：" + DirectionTitle + "            报单价格：" + Float64ToString(pInputOrder.GetLimitPrice()) + "\n" +
            "- 报单开平：" + OffsetFlagTitle + "          报单数量：" + IntToString(pInputOrder.GetVolumeTotalOriginal()) + "\n" +
            "- 错误代码：" + string(pRspInfo.GetErrorID()) + "    \t错误消息：" + GbkToUtf8(pRspInfo.GetErrorMsg()) + "\n" +
            "-------------------------------------------------------------------------------------------------"
        Println(Str)

        // 更新持仓数据及资金账户
        AddQueryTask("ReqQryInvestorPosition")
        AddQueryTask("ReqQryTradingAccount")

    }
}

// 期货发起银行资金转期货通知
func (p *FtdcTradeSpi) OnRtnFromBankToFutureByFuture(pTransfer goctp.CThostFtdcRspTransferField) {

    defer CheckPanic()

    if !IsNullPointer(pTransfer) {

        Str := "\n-------------------------------------------------------------------------------------------------\n" +
            "- 资金转账通知\n" +
            "- 交易时间：" + pTransfer.GetTradeTime() + " \t银行帐号：" + pTransfer.GetBankAccount() + "\n" +
            "- 转帐金额：" + Float64ToString(Decimal(pTransfer.GetTradeAmount(), 2)) + "   \t转帐方向：银行 -> 期货\n" +
            "- 转账消息：" + GbkToUtf8(pTransfer.GetErrorMsg()) + "\n" +
            "-------------------------------------------------------------------------------------------------"
        Println(Str)

        // 更新资金账户
        AddQueryTask("ReqQryTradingAccount")
    }
}

// 期货发起期货资金转银行通知
func (p *FtdcTradeSpi) OnRtnFromFutureToBankByFuture(pTransfer goctp.CThostFtdcRspTransferField) {

    defer CheckPanic()

    if !IsNullPointer(pTransfer) {

        Str := "\n-------------------------------------------------------------------------------------------------\n" +
            "- 资金转账通知\n" +
            "- 交易时间：" + pTransfer.GetTradeTime() + " \t银行帐号：" + pTransfer.GetBankAccount() + "\n" +
            "- 转帐金额：" + Float64ToString(Decimal(pTransfer.GetTradeAmount(), 2)) + "   \t转帐方向：期货 -> 银行\n" +
            "- 转账消息：" + GbkToUtf8(pTransfer.GetErrorMsg()) + "\n" +
            "-------------------------------------------------------------------------------------------------"
        Println(Str)

        // 更新资金账户
        AddQueryTask("ReqQryTradingAccount")
    }
}

// 银行发起银行资金转期货通知
func (p *FtdcTradeSpi) OnRtnFromBankToFutureByBank(pTransfer goctp.CThostFtdcRspTransferField) {

    defer CheckPanic()

    if !IsNullPointer(pTransfer) {

        Str := "\n-------------------------------------------------------------------------------------------------\n" +
            "- 资金转账通知\n" +
            "- 交易时间：" + pTransfer.GetTradeTime() + " \t银行帐号：" + pTransfer.GetBankAccount() + "\n" +
            "- 转帐金额：" + Float64ToString(Decimal(pTransfer.GetTradeAmount(), 2)) + "   \t转帐方向：银行 -> 期货\n" +
            "- 转账消息：" + GbkToUtf8(pTransfer.GetErrorMsg()) + "\n" +
            "-------------------------------------------------------------------------------------------------"
        Println(Str)

        // 更新资金账户
        AddQueryTask("ReqQryTradingAccount")
    }
}

// 银行发起期货资金转银行通知
func (p *FtdcTradeSpi) OnRtnFromFutureToBankByBank(pTransfer goctp.CThostFtdcRspTransferField) {

    defer CheckPanic()

    if !IsNullPointer(pTransfer) {

        Str := "\n-------------------------------------------------------------------------------------------------\n" +
            "- 资金转账通知\n" +
            "- 交易时间：" + pTransfer.GetTradeTime() + " \t银行帐号：" + pTransfer.GetBankAccount() + "\n" +
            "- 转帐金额：" + Float64ToString(Decimal(pTransfer.GetTradeAmount(), 2)) + "   \t转帐方向：期货 -> 银行\n" +
            "- 转账消息：" + GbkToUtf8(pTransfer.GetErrorMsg()) + "\n" +
            "-------------------------------------------------------------------------------------------------"
        Println(Str)

        // 更新资金账户
        AddQueryTask("ReqQryTradingAccount")
    }
}

// 错误应答
func (p *FtdcTradeSpi) OnRspError(pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
    p.IsErrorRspInfo(pRspInfo, nRequestID)
}

// 报单操作错误回报
func (p *FtdcTradeSpi) OnErrRtnOrderAction(pOrderAction goctp.CThostFtdcOrderActionField, pRspInfo goctp.CThostFtdcRspInfoField) {
    p.IsErrorRspInfo(pRspInfo, StringToInt(pOrderAction.GetOrderRef()))
}

// 报单操作请求响应（撤单失败会触发）
func (p *FtdcTradeSpi) OnRspOrderAction(pInputOrderAction goctp.CThostFtdcInputOrderActionField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
    p.IsErrorRspInfo(pRspInfo, nRequestID)
}

// 心跳超时警告。当长时间未收到报文时，该方法被调用。
func (p *FtdcTradeSpi) OnHeartBeatWarning(nTimeLapse int) {
    Println("心跳超时警告（OnHeartBeatWarning） nTimerLapse=", nTimeLapse)
}

// 处理持仓明细方向与持仓汇总方向不一致的问题（将明细改为与汇总一致）
func GetPositionDetailDirection(Direction string) string {

    result := ""

    switch Direction {
        case "0":
            result = "2"

        case "1":
            result = "3"

        default:
            result = Direction
    }

    return result
}

// 获得合约详情信息
func GetInstrumentInfo(InstrumentID string) (InstrumentStruct, bool) {

    defer CheckPanic()

    // 将合约转换成大写
    InstrumentID = ToUpper(InstrumentID)

    if v, ok := MapInstruments.Get(InstrumentID); ok {
        return v.(InstrumentStruct), true
    } else {

        var mInstrument InstrumentStruct
        return mInstrument, false
    }
}

// 获得合约名称
func GetInstrumentName(InstrumentID string) string {

    defer CheckPanic()

    mInstrument, mInstrumentOk := GetInstrumentInfo(InstrumentID)
    if !mInstrumentOk {
        return InstrumentID
    }

    return mInstrument.InstrumentName
}

// 获得合约 code（合约的字母部分）
func GetInstrumentCode(InstrumentID string) string {

    // 如果 map 有记录则使用已存在的记录
    mInstrumentCode, mInstrumentCodeOk := MapInstrumentCodes.Get(InstrumentID)
    if mInstrumentCodeOk {
        return mInstrumentCode.(string)

    } else {

        InstrumentCode := regexp.MustCompile("^[a-zA-Z]+").FindString(InstrumentID)
        InstrumentCode = strings.ToUpper(InstrumentCode)

        MapInstrumentCodes.Set(InstrumentID, InstrumentCode)

        return InstrumentCode
    }
}

// 获得期权对应的合约
func GetOptionInstrumentID(InstrumentID string) string {
    return regexp.MustCompile("^[a-zA-Z]+\\d+").FindString(InstrumentID)
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

    defer CheckPanic()

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

/**
 * 获得报单数据
 *
 * @param string InstrumentID
 * @param string OrderSysID
 */
func GetOrder(InstrumentID string, OrderSysID string) (OrderStruct, bool) {

    defer CheckPanic()

    var sOrder OrderStruct

    MapKey := Sprintf("%v_%v", InstrumentID, TrimSpace(OrderSysID))

    mOrder, mOrderOk := MapOrders.Get(MapKey)
    if !mOrderOk {
        return sOrder, false
    }

    return mOrder.(OrderStruct), true
}

/**
 * 获得成交数据
 *
 * @param string InstrumentID
 * @param string TradeID
 */
func GetTrade(InstrumentID string, TradeID string) (TradeStruct, bool) {

    defer CheckPanic()

    var sTrade TradeStruct

    MapKey := Sprintf("%v_%v", TradeID, InstrumentID)

    mTrade, mTradeOk := MapTrades.Get(MapKey)
    if !mTradeOk {
        return sTrade, false
    }

    return mTrade.(TradeStruct), true
}

// 获得资金账户信息
func GetAccount() AccountStruct {

    defer CheckPanic()

    var sAccount AccountStruct

    MapKey := Sprintf("%v_%v", Ctp.BrokerID, Ctp.InvestorID)

    mAccountVal, mAccountOk := MapAccounts.Get(MapKey)
    if !mAccountOk {
        return sAccount
    }

    return mAccountVal.(AccountStruct)
}

// 添加队列查询任务（写上对应的函数名字）
func AddQueryTask(funcName string) {

    defer CheckPanic()

    Queue.PushUnique(funcName)
}