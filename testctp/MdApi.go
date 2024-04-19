package testctp

import (
    "gitee.com/mayiweb/goctp"
)

type MdApi struct {
    goctp.CThostFtdcMdApi
    MdSpi goctp.CThostFtdcMdSpi
}

// 响应信息
func (p *MdApi) GetRspInfo() goctp.CThostFtdcRspInfoField {
    pRspInfo := goctp.NewCThostFtdcRspInfoField()
    pRspInfo.SetErrorID(0)
    pRspInfo.SetErrorMsg("ok")

    return pRspInfo
}

// 获取API的版本信息
func (p *MdApi) GetApiVersion() string {
    return "v6.6.5 Testing MdApi"
}

// 删除接口对象本身
func (p *MdApi) Release() {

}

// 初始化
func (p *MdApi) Init() {
    p.MdSpi.OnFrontConnected()
}

// 等待接口线程结束运行
func (p *MdApi) Join() int {
    return 0
}

// 获取当前交易日
func (p *MdApi) GetTradingDay() string {
    return ""
}

// 注册回调接口
func (p *MdApi) RegisterSpi(pSpi goctp.CThostFtdcMdSpi) {
    p.MdSpi = pSpi
}

// 注册前置机网络地址
func (p *MdApi) RegisterFront(pszFrontAddress string) {

}

// 订阅行情
func (p *MdApi) SubscribeMarketData(ppInstrumentID *string, nRequestID int) int {
    return 0
}

// 退订行情
func (p *MdApi) UnSubscribeMarketData(ppInstrumentID *string, nRequestID int) int {
    return 0
}

// 订阅询价。
func (p *MdApi) SubscribeForQuoteRsp(ppInstrumentID *string, nRequestID int) int {
    return 0
}

// 退订询价
func (p *MdApi) UnSubscribeForQuoteRsp(ppInstrumentID *string, nRequestID int) int {
    return 0
}

// 用户登录请求
func (p *MdApi) ReqUserLogin(pReqUserLoginField goctp.CThostFtdcReqUserLoginField, nRequestID int) int {

    pResult := goctp.NewCThostFtdcRspUserLoginField()
    pResult  = nil

    p.MdSpi.OnRspUserLogin(pResult, p.GetRspInfo(), nRequestID, true)

    return 0
}

// 登出请求
func (p *MdApi) ReqUserLogout(pUserLogout goctp.CThostFtdcUserLogoutField, nRequestID int) int {

    pResult := goctp.NewCThostFtdcUserLogoutField()
    pResult  = nil

    p.MdSpi.OnRspUserLogout(pResult, p.GetRspInfo(), nRequestID, true)

    return 0
}

// 请求查询组播合约
func (p *MdApi) ReqQryMulticastInstrument(pQryMulticastInstrument goctp.CThostFtdcQryMulticastInstrumentField, nRequestID int) int {
    return 0
}