/**
 * 行情模块
 */

package main

/*
#include <stdio.h>
#include <stdlib.h>
*/
import "C"

import (
    "gitee.com/mayiweb/goctp"
    "os"
    "unsafe"
)

// 获得行情请求编号
func (p *FtdcMdSpi) GetMdRequestId() int {
    Ctp.MdRequestId += 1
    return Ctp.MdRequestId
}

// 发送请求日志（仅查询类的函数需要调用）
func (p *FtdcMdSpi) ReqMsg(Msg string) {

    // 交易程序未初始化完成时，执行查询类的函数需要有1.5秒间隔
    if !Ctp.IsTradeInitFinish {
        // Sleep(1500)
    }

    Println("")
    LogPrintln(Msg)
}

// 行情系统错误通知
func (p *FtdcMdSpi) IsErrorRspInfo(pRspInfo goctp.CThostFtdcRspInfoField) bool {

    // 容错处理 pRspInfo ，部分响应函数中，pRspInfo 为空
    if IsNullPointer(pRspInfo) {

        return false

    } else {

        // 如果ErrorID != 0, 说明收到了错误的响应
        bResult := (pRspInfo.GetErrorID() != 0)
        if bResult {
            LogPrintf("ErrorID=%v ErrorMsg=%v\n", pRspInfo.GetErrorID(), GbkToUtf8(pRspInfo.GetErrorMsg()))

            // 如果初始化未完成则退出程序，防止卡在某个流程
            if !Ctp.IsTradeInitFinish {
                Exit()
            }
        }

        return bResult
    }
}

// 获取API的版本信息
func (p *FtdcMdSpi) GetApiVersion() string {
    if RunMode == BackTestingMode {
        return TestCtpMdApi.GetApiVersion()
    } else {
        return goctp.CThostFtdcTraderApiGetApiVersion()
    }
}

// 当客户端与交易后台通信连接断开时，该方法被调用。当发生这个情况后，API会自动重新连接，客户端可不做处理。
// 服务器已断线，该函数也会被调用。【api 会自动初始化程序，并重新登录】
func (p *FtdcMdSpi) OnFrontDisconnected(nReason int) {

    if Ctp.IsMdInit {

        Ctp.IsMdInit  = false
        Ctp.IsMdLogin = false
        LogPrintln("行情服务器已断线，尝试重新连接中...")
    }
}

// 当客户端与交易后台建立起通信连接时（还未登录前），该方法被调用。
func (p *FtdcMdSpi) OnFrontConnected() {

    MdStr := "-------------------------------------------------------------------------------------------------\n" +
        "- 行情系统初始化成功，Api 版本：" + p.GetApiVersion() + "\n" +
        "-------------------------------------------------------------------------------------------------"
    Println(MdStr)

    Ctp.IsMdInit = true

    // 登录（如果交易系统已登录，行情模块才初始化完成则进行登录）
    if Ctp.IsTradeLogin {
        p.ReqUserLogin()
    }
}

// 行情用户登录
func (p *FtdcMdSpi) ReqUserLogin() {

    // 只有未登录时才可以进行登录操作
    if !Ctp.IsMdInit || Ctp.IsMdLogin {
        return
    }

    LogPrintln("行情系统登录中...")

    req := goctp.NewCThostFtdcReqUserLoginField()
    req.SetBrokerID(BrokerID)
    req.SetUserID(InvestorID)
    req.SetPassword(Password)

    iResult := Ctp.MdApi.ReqUserLogin(req, p.GetMdRequestId())

    if iResult != 0 {
        ReqFailMsg("发送行情系统登录请求失败！", iResult)
    }
}

// 登录请求响应
func (p *FtdcMdSpi) OnRspUserLogin(pRspUserLogin goctp.CThostFtdcRspUserLoginField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {

    if bIsLast && !p.IsErrorRspInfo(pRspInfo) {
        Ctp.IsMdLogin = true
        LogPrintf("行情系统已登录，当前交易日：%v\n", Ctp.MdApi.GetTradingDay())
    }
}

// 订阅行情
func (p *FtdcMdSpi) SubscribeTick(InstrumentID []string) int {

    if len(InstrumentID) == 0 {
        LogPrintln("没有指定需要订阅的行情数据")
        return 0
    }

    Println("")
    LogPrintln("订阅行情数据中...")

    // 本次需要订阅的合约
    args := make([]*C.char, 0)
    for _, v := range InstrumentID {
        char := C.CString(v)
        defer C.free(unsafe.Pointer(char))
        args = append(args, char)
    }

    iResult := Ctp.MdApi.SubscribeMarketData((*string)(unsafe.Pointer(&args[0])), len(InstrumentID))

    if iResult != 0 {
        ReqFailMsg("发送订阅行情请求失败！", iResult)
        return iResult
    }

    return iResult
}

// 订阅行情响应
func (p *FtdcMdSpi) OnRspSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
    if !p.IsErrorRspInfo(pRspInfo) {
        LogPrintf("订阅合约 %v 行情数据成功！\n", pSpecificInstrument.GetInstrumentID())
    }
}

// 退订行情
func (p *FtdcMdSpi) UnSubscribeTick(InstrumentID []string) int {

    if len(InstrumentID) == 0 {
        LogPrintln("没有指定需要退订的行情数据")
        return 0
    }

    Println("")
    LogPrintln("退订行情数据中...")

    args := make([]*C.char, 0)
    for _, v := range InstrumentID {
        char := C.CString(v)
        defer C.free(unsafe.Pointer(char))
        args = append(args, char)
    }

    iResult := Ctp.MdApi.UnSubscribeMarketData((*string)(unsafe.Pointer(&args[0])), len(InstrumentID))

    if iResult != 0 {
        ReqFailMsg("发送退订行情请求失败！", iResult)
    }

    return iResult
}

// 退订行情响应
func (p *FtdcMdSpi) OnRspUnSubMarketData(pSpecificInstrument goctp.CThostFtdcSpecificInstrumentField, pRspInfo goctp.CThostFtdcRspInfoField, nRequestID int, bIsLast bool) {
    if !p.IsErrorRspInfo(pRspInfo) {
        LogPrintln("取消订阅 %v 行情数据成功！\n", pSpecificInstrument.GetInstrumentID())
    }
}

// 深度行情响应
// 注意：五档行情需要将服务器放在期货公司内部网络才能接收
func (p *FtdcMdSpi) OnRtnDepthMarketData(pTick goctp.CThostFtdcDepthMarketDataField) {

    defer CheckPanic()

    if !IsNullPointer(pTick) {

        // 合约详情
        mInstrument, _ := GetInstrumentInfo(pTick.GetInstrumentID())

        var sTick TickStruct

        // 有的交易所夜盘交易日期还是当天的，没有切换到下一个交易日
        // sTick.TradingDay         = string(pTick.GetTradingDay())
        sTick.TradingDay         = Ctp.TradeApi.GetTradingDay()
        sTick.InstrumentID       = string(pTick.GetInstrumentID())
        sTick.InstrumentCode     = mInstrument.InstrumentCode
        sTick.InstrumentName     = mInstrument.InstrumentName
        sTick.ExchangeID         = string(pTick.GetExchangeID())
        sTick.LastPrice          = Decimal(pTick.GetLastPrice(), 4)
        sTick.PreSettlementPrice = Decimal(pTick.GetPreSettlementPrice(), 4)
        sTick.PreClosePrice      = Decimal(pTick.GetPreClosePrice(), 4)
        sTick.PreOpenInterest    = Decimal(pTick.GetPreOpenInterest(), 4)
        sTick.OpenPrice          = Decimal(pTick.GetOpenPrice(), 4)
        sTick.HighestPrice       = Decimal(pTick.GetHighestPrice(), 4)
        sTick.LowestPrice        = Decimal(pTick.GetLowestPrice(), 4)
        sTick.Volume             = pTick.GetVolume()
        sTick.Turnover           = Decimal(pTick.GetTurnover(), 4)
        sTick.OpenInterest       = Decimal(pTick.GetOpenInterest(), 4)
        sTick.ClosePrice         = Decimal(pTick.GetClosePrice(), 4)
        sTick.SettlementPrice    = Decimal(pTick.GetSettlementPrice(), 4)
        sTick.UpperLimitPrice    = Decimal(pTick.GetUpperLimitPrice(), 4)
        sTick.LowerLimitPrice    = Decimal(pTick.GetLowerLimitPrice(), 4)
        sTick.UpdateTime         = pTick.GetUpdateTime()
        sTick.UpdateMillisec     = pTick.GetUpdateMillisec()
        sTick.BidPrice1          = Decimal(pTick.GetBidPrice1(), 4)
        sTick.BidVolume1         = pTick.GetBidVolume1()
        sTick.AskPrice1          = Decimal(pTick.GetAskPrice1(), 4)
        sTick.AskVolume1         = pTick.GetAskVolume1()
        sTick.BidPrice2          = Decimal(pTick.GetBidPrice2(), 4)
        sTick.BidVolume2         = pTick.GetBidVolume2()
        sTick.AskPrice2          = Decimal(pTick.GetAskPrice2(), 4)
        sTick.AskVolume2         = pTick.GetAskVolume2()
        sTick.BidPrice3          = Decimal(pTick.GetBidPrice3(), 4)
        sTick.BidVolume3         = pTick.GetBidVolume3()
        sTick.AskPrice3          = Decimal(pTick.GetAskPrice3(), 4)
        sTick.AskVolume3         = pTick.GetAskVolume3()
        sTick.BidPrice4          = Decimal(pTick.GetBidPrice4(), 4)
        sTick.BidVolume4         = pTick.GetBidVolume4()
        sTick.AskPrice4          = Decimal(pTick.GetAskPrice4(), 4)
        sTick.AskVolume4         = pTick.GetAskVolume4()
        sTick.BidPrice5          = Decimal(pTick.GetBidPrice5(), 4)
        sTick.BidVolume5         = pTick.GetBidVolume5()
        sTick.AskPrice5          = Decimal(pTick.GetAskPrice5(), 4)
        sTick.AskVolume5         = pTick.GetAskVolume5()
        sTick.AveragePrice       = Decimal(pTick.GetAveragePrice(), 4)

        p.Tick(sTick)
    }
}

// 深度行情
func (p *FtdcMdSpi) Tick(pTick TickStruct) {

    // 未开盘的行情 金额太长，直接过滤 1.7976931348623157e+308
    if pTick.BidPrice1 == 1.7976931348623157e+308 || pTick.AskPrice1 == 1.7976931348623157e+308 || pTick.OpenPrice == 1.7976931348623157e+308 {
        return
    }

    // 新版本 6.3.15 会有一些为0的垃圾数据，此处过滤
    if pTick.Volume == 0 && pTick.BidPrice1 == 0.00 && pTick.AskPrice1 == 0.00 && pTick.BidVolume1 == 0 && pTick.AskVolume1 == 0 {
        return
    }

    // 注意：如果涨停，卖一价没有报价会为0
    if pTick.BidPrice1 == 0.00 {
        pTick.BidPrice1 = pTick.LastPrice
    }

    // 注意：如果跌停，买一价没有报价会为0
    if pTick.AskPrice1 == 0.00 {
        pTick.AskPrice1 = pTick.LastPrice
    }

    // 风险管理（止盈、止损）

    // 策略运行

    // 保存 tick 数据
    SaveTickData(pTick)
}

// 保存 tick 数据
func SaveTickData(pTick TickStruct) {

    defer CheckPanic()

    // tick 写到csv文件中（写tick 时，不要打开文件，不然会写不进去）

    // tick目录
    tickDirectory := Sprintf("%v%v", TickDataDirectory, ToUpper(pTick.InstrumentID))

    tickPath := Sprintf("%v/%v.%v.csv", tickDirectory, ToUpper(pTick.InstrumentID), pTick.TradingDay)

    // 检查Tick文件目录是否存在
    tickDirectoryExists, _ := PathExists(tickDirectory)
    if !tickDirectoryExists {
        err := os.Mkdir(tickDirectory, os.ModePerm)
        if err != nil {
           LogPrintln("创建目录失败，请检查是否有操作权限 Err:", err)
        }
    }

    // 检查Tick文件是否存在
    tickFileExists, _ := PathExists(tickPath)

    file, fileErr := os.OpenFile(tickPath, os.O_RDWR | os.O_APPEND | os.O_CREATE, 0666)
    defer file.Close()

    if fileErr != nil {
        LogPrintln("打开 Tick 文件失败 Err:", fileErr)
        return
    }

    // 如果tick 文件不存在，则在创建文件后写入表头
    if !tickFileExists {
        csvTitle := "UpdateTime,TradingDay,InstrumentID,LastPrice,OpenPrice,PreClosePrice,HighestPrice,LowestPrice,UpperLimitPrice,LowerLimitPrice,AveragePrice,Volume,OpenInterest,BidPrice1,BidVolume1,AskPrice1,AskVolume1,BidPrice2,BidVolume2,AskPrice2,AskVolume2,BidPrice3,BidVolume3,AskPrice3,AskVolume3,BidPrice4,BidVolume4,AskPrice4,AskVolume4,BidPrice5,BidVolume5,AskPrice5,AskVolume5\n"
        _, writeErr := file.WriteString(csvTitle)
        if writeErr != nil {
            LogPrintln("写入 Tick 表头失败 Err:", writeErr)
            return
        }
    }

    // tick 文本
    tickContent := Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", pTick.UpdateTime, pTick.TradingDay, pTick.InstrumentID, pTick.LastPrice, pTick.OpenPrice, pTick.PreClosePrice, pTick.HighestPrice, pTick.LowestPrice, pTick.UpperLimitPrice, pTick.LowerLimitPrice, pTick.AveragePrice, pTick.Volume, pTick.OpenInterest, pTick.BidPrice1, pTick.BidVolume1, pTick.AskPrice1, pTick.AskVolume1, pTick.BidPrice2, pTick.BidVolume2, pTick.AskPrice2, pTick.AskVolume2, pTick.BidPrice3, pTick.BidVolume3, pTick.AskPrice3, pTick.AskVolume3, pTick.BidPrice4, pTick.BidVolume4, pTick.AskPrice4, pTick.AskVolume4, pTick.BidPrice5, pTick.BidVolume5, pTick.AskPrice5, pTick.AskVolume5)

    _, writeErr := file.WriteString(tickContent)
    if writeErr != nil {
        LogPrintln("写入 Tick 数据失败 Err:", writeErr)
        return
    }
}