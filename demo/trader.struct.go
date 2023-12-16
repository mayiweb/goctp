/**
 * 交易模块 结构体数据转换
 */

package main

import (
    "gitee.com/mayiweb/goctp"
)

// 获得登录成功结构体数据
func GetLoginStruct(pRspUserLogin goctp.CThostFtdcRspUserLoginField) LoginStruct {

    var sLogin LoginStruct

    sLogin.TradingDay  = pRspUserLogin.GetTradingDay()
    sLogin.LoginTime   = pRspUserLogin.GetLoginTime()
    sLogin.BrokerID    = pRspUserLogin.GetBrokerID()
    sLogin.UserID      = pRspUserLogin.GetUserID()
    sLogin.SystemName  = pRspUserLogin.GetSystemName()
    sLogin.FrontID     = pRspUserLogin.GetFrontID()
    sLogin.SessionID   = pRspUserLogin.GetSessionID()
    sLogin.MaxOrderRef = pRspUserLogin.GetMaxOrderRef()
    sLogin.SHFETime    = pRspUserLogin.GetSHFETime()
    sLogin.DCETime     = pRspUserLogin.GetDCETime()
    sLogin.CZCETime    = pRspUserLogin.GetCZCETime()
    sLogin.FFEXTime    = pRspUserLogin.GetFFEXTime()
    sLogin.INETime     = pRspUserLogin.GetINETime()
    sLogin.MapKey      = Sprintf("%v_%v", sLogin.BrokerID, sLogin.UserID)

    MapLogin.Set(sLogin.MapKey, sLogin)

    return sLogin
}

// 获得合约结构体数据
func GetInstrumentStruct(pInstrument goctp.CThostFtdcInstrumentField) InstrumentStruct {

    var sInstrument InstrumentStruct

    sInstrument.InstrumentID           = pInstrument.GetInstrumentID()
    sInstrument.InstrumentName         = TrimSpace(GbkToUtf8(pInstrument.GetInstrumentName()))
    sInstrument.ExchangeID             = pInstrument.GetExchangeID()
    sInstrument.ExchangeName           = GetExchangeName(pInstrument.GetExchangeID())
    sInstrument.InstrumentCode         = GetInstrumentCode(pInstrument.GetInstrumentID())
    sInstrument.ExchangeInstID         = pInstrument.GetExchangeInstID()
    sInstrument.ProductID              = pInstrument.GetProductID()
    sInstrument.ProductClass           = string(pInstrument.GetProductClass())
    sInstrument.ProductClassTitle      = GetProductClassTitle(string(pInstrument.GetProductClass()))
    sInstrument.DeliveryYear           = pInstrument.GetDeliveryYear()
    sInstrument.DeliveryMonth          = pInstrument.GetDeliveryMonth()
    sInstrument.MaxMarketOrderVolume   = pInstrument.GetMaxMarketOrderVolume()
    sInstrument.MinMarketOrderVolume   = pInstrument.GetMinMarketOrderVolume()
    sInstrument.MaxLimitOrderVolume    = pInstrument.GetMaxLimitOrderVolume()
    sInstrument.MinLimitOrderVolume    = pInstrument.GetMinLimitOrderVolume()
    sInstrument.VolumeMultiple         = pInstrument.GetVolumeMultiple()
    sInstrument.PriceTick              = pInstrument.GetPriceTick()
    sInstrument.CreateDate             = pInstrument.GetCreateDate()
    sInstrument.OpenDate               = pInstrument.GetOpenDate()
    sInstrument.ExpireDate             = pInstrument.GetExpireDate()
    sInstrument.StartDelivDate         = pInstrument.GetStartDelivDate()
    sInstrument.EndDelivDate           = pInstrument.GetEndDelivDate()
    sInstrument.InstLifePhase          = string(pInstrument.GetInstLifePhase())
    sInstrument.IsTrading              = pInstrument.GetIsTrading()
    sInstrument.PositionType           = string(pInstrument.GetPositionType())
    sInstrument.PositionDateType       = string(pInstrument.GetPositionDateType())
    sInstrument.MaxMarginSideAlgorithm = string(pInstrument.GetMaxMarginSideAlgorithm())
    sInstrument.UnderlyingInstrID      = pInstrument.GetUnderlyingInstrID()
    sInstrument.OptionsType            = string(pInstrument.GetOptionsType())
    sInstrument.StrikePrice            = pInstrument.GetStrikePrice()
    sInstrument.UnderlyingMultiple     = pInstrument.GetUnderlyingMultiple()
    sInstrument.CombinationType        = string(pInstrument.GetCombinationType())

    // 将合约K键转换成大写
    sInstrument.MapKey                 = ToUpper(pInstrument.GetInstrumentID())

    MapInstruments.Set(sInstrument.MapKey, sInstrument)

    return sInstrument
}

// 获得资金账户结构体数据
func GetAccountStruct(pTradingAccount goctp.CThostFtdcTradingAccountField) AccountStruct {

    var sAccount AccountStruct

    sAccount.BrokerID                       = pTradingAccount.GetBrokerID()
    sAccount.AccountID                      = pTradingAccount.GetAccountID()
    sAccount.PreMortgage                    = Decimal(pTradingAccount.GetPreMortgage(), 2)
    sAccount.PreCredit                      = Decimal(pTradingAccount.GetPreCredit(), 2)
    sAccount.PreDeposit                     = Decimal(pTradingAccount.GetPreDeposit(), 2)
    sAccount.PreBalance                     = Decimal(pTradingAccount.GetPreBalance(), 2)
    sAccount.PreMargin                      = Decimal(pTradingAccount.GetPreMargin(), 2)
    sAccount.InterestBase                   = Decimal(pTradingAccount.GetInterestBase(), 2)
    sAccount.Interest                       = Decimal(pTradingAccount.GetInterest(), 2)
    sAccount.Deposit                        = Decimal(pTradingAccount.GetDeposit(), 2)
    sAccount.Withdraw                       = Decimal(pTradingAccount.GetWithdraw(), 2)
    sAccount.FrozenMargin                   = Decimal(pTradingAccount.GetFrozenMargin(), 2)
    sAccount.FrozenCash                     = Decimal(pTradingAccount.GetFrozenCash(), 2)
    sAccount.FrozenCommission               = Decimal(pTradingAccount.GetFrozenCommission(), 2)
    sAccount.CurrMargin                     = Decimal(pTradingAccount.GetCurrMargin(), 2)
    sAccount.CashIn                         = Decimal(pTradingAccount.GetCashIn(), 2)
    sAccount.Commission                     = Decimal(pTradingAccount.GetCommission(), 2)
    sAccount.CloseProfit                    = Decimal(pTradingAccount.GetCloseProfit(), 2)
    sAccount.PositionProfit                 = Decimal(pTradingAccount.GetPositionProfit(), 2)
    sAccount.Balance                        = Decimal(pTradingAccount.GetBalance(), 2)
    sAccount.Available                      = Decimal(pTradingAccount.GetAvailable(), 2)
    sAccount.WithdrawQuota                  = Decimal(pTradingAccount.GetWithdrawQuota(), 2)
    sAccount.Reserve                        = Decimal(pTradingAccount.GetReserve(), 2)
    sAccount.TradingDay                     = pTradingAccount.GetTradingDay()
    sAccount.SettlementID                   = pTradingAccount.GetSettlementID()
    sAccount.Credit                         = Decimal(pTradingAccount.GetCredit(), 2)
    sAccount.Mortgage                       = Decimal(pTradingAccount.GetMortgage(), 2)
    sAccount.ExchangeMargin                 = Decimal(pTradingAccount.GetExchangeMargin(), 2)
    sAccount.DeliveryMargin                 = Decimal(pTradingAccount.GetDeliveryMargin(), 2)
    sAccount.ExchangeDeliveryMargin         = Decimal(pTradingAccount.GetExchangeDeliveryMargin(), 2)
    sAccount.ReserveBalance                 = Decimal(pTradingAccount.GetReserveBalance(), 2)
    sAccount.CurrencyID                     = pTradingAccount.GetCurrencyID()
    sAccount.PreFundMortgageIn              = Decimal(pTradingAccount.GetPreFundMortgageIn(), 2)
    sAccount.PreFundMortgageOut             = Decimal(pTradingAccount.GetPreFundMortgageOut(), 2)
    sAccount.FundMortgageIn                 = Decimal(pTradingAccount.GetFundMortgageIn(), 2)
    sAccount.FundMortgageOut                = Decimal(pTradingAccount.GetFundMortgageOut(), 2)
    sAccount.FundMortgageAvailable          = Decimal(pTradingAccount.GetFundMortgageAvailable(), 2)
    sAccount.MortgageableFund               = Decimal(pTradingAccount.GetMortgageableFund(), 2)
    sAccount.SpecProductMargin              = Decimal(pTradingAccount.GetSpecProductMargin(), 2)
    sAccount.SpecProductFrozenMargin        = Decimal(pTradingAccount.GetSpecProductFrozenMargin(), 2)
    sAccount.SpecProductCommission          = Decimal(pTradingAccount.GetSpecProductCommission(), 2)
    sAccount.SpecProductFrozenCommission    = Decimal(pTradingAccount.GetSpecProductFrozenCommission(), 2)
    sAccount.SpecProductPositionProfit      = Decimal(pTradingAccount.GetSpecProductPositionProfit(), 2)
    sAccount.SpecProductCloseProfit         = Decimal(pTradingAccount.GetSpecProductCloseProfit(), 2)
    sAccount.SpecProductPositionProfitByAlg = Decimal(pTradingAccount.GetSpecProductPositionProfitByAlg(), 2)
    sAccount.SpecProductExchangeMargin      = Decimal(pTradingAccount.GetSpecProductExchangeMargin(), 2)
    sAccount.BizType                        = string(pTradingAccount.GetBizType())
    sAccount.FrozenSwap                     = Decimal(pTradingAccount.GetFrozenSwap(), 2)
    sAccount.RemainSwap                     = Decimal(pTradingAccount.GetRemainSwap(), 2)
    sAccount.MapKey                         = Sprintf("%v_%v", sAccount.BrokerID, sAccount.AccountID)

    MapAccounts.Set(sAccount.MapKey, sAccount)

    return sAccount
}

// 获得报单结构体数据
func GetOrderStruct(pOrder goctp.CThostFtdcOrderField) OrderStruct {

    var sOrder OrderStruct

    sOrder.BrokerID            = pOrder.GetBrokerID()
    sOrder.InvestorID          = pOrder.GetInvestorID()
    sOrder.InstrumentID        = pOrder.GetInstrumentID()
    sOrder.InstrumentName      = GetInstrumentName(pOrder.GetInstrumentID())
    sOrder.ExchangeID          = pOrder.GetExchangeID()
    sOrder.FrontID             = pOrder.GetFrontID()
    sOrder.OrderRef            = pOrder.GetOrderRef()
    sOrder.SessionID           = pOrder.GetSessionID()
    sOrder.InsertTime          = pOrder.GetInsertTime()
    sOrder.InsertDate          = pOrder.GetInsertDate()
    sOrder.OrderSysID          = pOrder.GetOrderSysID()
    sOrder.LimitPrice          = pOrder.GetLimitPrice()
    sOrder.Volume              = pOrder.GetVolumeTotalOriginal()
    sOrder.VolumeTraded        = pOrder.GetVolumeTraded()
    sOrder.VolumeTotal         = pOrder.GetVolumeTotal()
    sOrder.Direction           = string(pOrder.GetDirection())
    sOrder.CombOffsetFlag      = string(pOrder.GetCombOffsetFlag())
    sOrder.CombHedgeFlag       = string(pOrder.GetCombHedgeFlag())
    sOrder.OrderStatus         = string(pOrder.GetOrderStatus())
    sOrder.StatusMsg           = GbkToUtf8(pOrder.GetStatusMsg())
    sOrder.DirectionTitle      = GetDirectionTitle(sOrder.Direction)
    sOrder.OrderStatusTitle    = GetOrderStatusTitle(sOrder.OrderStatus)
    sOrder.CombOffsetFlagTitle = GetOffsetFlagTitle(sOrder.CombOffsetFlag)
    sOrder.MapKey              = Sprintf("%v_%v", sOrder.InstrumentID, TrimSpace(sOrder.OrderSysID))

    // 只记录有报单编号的报单数据
    if sOrder.OrderSysID != "" {
        // 记录报单数据
        MapOrders.Set(sOrder.MapKey, sOrder)
    }

    return sOrder
}

// 获得报单成交结构体数据
func GetTradeStruct(pTrade goctp.CThostFtdcTradeField) TradeStruct {

    var sTrade TradeStruct

    sTrade.BrokerID        = pTrade.GetBrokerID()
    sTrade.InvestorID      = pTrade.GetInvestorID()
    sTrade.OrderRef        = TrimSpace(pTrade.GetOrderRef())
    sTrade.UserID          = pTrade.GetUserID()
    sTrade.ExchangeID      = pTrade.GetExchangeID()
    sTrade.TradeID         = TrimSpace(pTrade.GetTradeID())
    sTrade.Direction       = string(pTrade.GetDirection())
    sTrade.OrderSysID      = TrimSpace(pTrade.GetOrderSysID())
    sTrade.ParticipantID   = pTrade.GetParticipantID()
    sTrade.ClientID        = pTrade.GetClientID()
    sTrade.OffsetFlag      = string(pTrade.GetOffsetFlag())
    sTrade.HedgeFlag       = string(pTrade.GetHedgeFlag())
    sTrade.Price           = Decimal(pTrade.GetPrice(), 4)
    sTrade.Volume          = pTrade.GetVolume()
    sTrade.TradeDate       = pTrade.GetTradeDate()
    sTrade.TradeTime       = pTrade.GetTradeTime()
    sTrade.TradeType       = string(pTrade.GetTradeType())
    sTrade.PriceSource     = string(pTrade.GetPriceSource())
    sTrade.TraderID        = pTrade.GetTraderID()
    sTrade.OrderLocalID    = pTrade.GetOrderLocalID()
    sTrade.ClearingPartID  = pTrade.GetClearingPartID()
    sTrade.BusinessUnit    = pTrade.GetBusinessUnit()
    sTrade.SequenceNo      = pTrade.GetSequenceNo()
    sTrade.TradingDay      = pTrade.GetTradingDay()
    sTrade.SettlementID    = pTrade.GetSettlementID()
    sTrade.BrokerOrderSeq  = pTrade.GetBrokerOrderSeq()
    sTrade.TradeSource     = string(pTrade.GetTradeSource())
    sTrade.InvestUnitID    = pTrade.GetInvestUnitID()
    sTrade.InstrumentID    = pTrade.GetInstrumentID()
    sTrade.ExchangeInstID  = pTrade.GetExchangeInstID()
    sTrade.DirectionTitle  = GetDirectionTitle(sTrade.Direction)
    sTrade.HedgeFlagTitle  = GetHedgeFlagTitle(sTrade.HedgeFlag)
    sTrade.OffsetFlagTitle = GetOffsetFlagTitle(sTrade.OffsetFlag)
    sTrade.InstrumentName  = GetInstrumentName(sTrade.InstrumentID)

    sTrade.MapKey = Sprintf("%v_%v", sTrade.TradeID, sTrade.InstrumentID)

    // 将报单成交数据记录下来
    MapTrades.Set(sTrade.MapKey, sTrade)

    return sTrade
}

// 获得持仓结构体数据
func GetPositionStruct(pPosition goctp.CThostFtdcInvestorPositionField) PositionStruct {

    defer CheckPanic()

    var sPosition PositionStruct

    // 检查合约详情是否存在
    mInstrument, _ := GetInstrumentInfo(pPosition.GetInstrumentID())

    // 合约乘数
    var VolumeMultiple int = mInstrument.VolumeMultiple

    // 开仓成本
    var OpenPrice float64 = pPosition.GetOpenCost() / float64(pPosition.GetPosition() * VolumeMultiple)

    sPosition.BrokerID          = pPosition.GetBrokerID()
    sPosition.InvestorID        = pPosition.GetInvestorID()
    sPosition.InstrumentID      = pPosition.GetInstrumentID()
    sPosition.InstrumentName    = mInstrument.InstrumentName
    sPosition.Direction         = string(pPosition.GetPosiDirection())
    sPosition.DirectionTitle    = GetPosiDirectionTitle(sPosition.Direction)
    sPosition.HedgeFlag         = string(pPosition.GetHedgeFlag())
    sPosition.HedgeFlagTitle    = GetHedgeFlagTitle(sPosition.HedgeFlag)
    sPosition.PositionDate      = string(pPosition.GetPositionDate())
    sPosition.PositionDateTitle = GetPositionDateTitle(sPosition.PositionDate)
    sPosition.Volume            = pPosition.GetPosition()
    sPosition.YdPosition        = pPosition.GetYdPosition()
    sPosition.TodayPosition     = pPosition.GetTodayPosition()
    sPosition.LongFrozen        = pPosition.GetLongFrozen()
    sPosition.ShortFrozen       = pPosition.GetShortFrozen()

    // 冻结的持仓量（多空并成一个字段）
    if sPosition.Direction == string(goctp.THOST_FTDC_PD_Long) {
        // 多头冻结的持仓量
        sPosition.ShortVolume = pPosition.GetShortFrozen()
    } else {
        // 空头冻结的持仓量
        sPosition.ShortVolume = pPosition.GetLongFrozen()
    }

    sPosition.OpenVolume         = pPosition.GetOpenVolume()
    sPosition.CloseVolume        = pPosition.GetCloseVolume()
    sPosition.PositionCost       = Decimal(pPosition.GetPositionCost(), 4)
    sPosition.Commission         = Decimal(pPosition.GetCommission(), 4)
    sPosition.CloseProfit        = pPosition.GetCloseProfit()
    sPosition.PositionProfit     = Decimal(pPosition.GetPositionProfit(), 4)
    sPosition.PreSettlementPrice = pPosition.GetPreSettlementPrice()
    sPosition.SettlementPrice    = Decimal(pPosition.GetSettlementPrice(), 4)
    sPosition.SettlementID       = pPosition.GetSettlementID()
    sPosition.OpenPrice          = Decimal(OpenPrice, 2)
    sPosition.ExchangeID         = pPosition.GetExchangeID()

    // 持仓key（合约+"_"+持仓方向+"_"+今昨仓类型）
    sPosition.MapKey = Sprintf("%v_%v_%v", sPosition.InstrumentID, sPosition.Direction, sPosition.PositionDate)

    // 记录持仓数据
    MapPositions.Set(sPosition.MapKey, sPosition)

    return sPosition
}

// 获得持仓明细结构体数据
func GetPositionDetailStruct(pPosition goctp.CThostFtdcInvestorPositionDetailField) PositionDetailStruct {

    defer CheckPanic()

    var sPosition PositionDetailStruct

    sPosition.BrokerID              = pPosition.GetBrokerID()
    sPosition.InvestorID            = pPosition.GetInvestorID()
    sPosition.InstrumentID          = pPosition.GetInstrumentID()
    sPosition.InstrumentName        = GetInstrumentName(pPosition.GetInstrumentID())
    sPosition.HedgeFlag             = string(pPosition.GetHedgeFlag())
    sPosition.HedgeFlagTitle        = GetHedgeFlagTitle(sPosition.HedgeFlag)
    sPosition.Direction             = GetPositionDetailDirection(string(pPosition.GetDirection()))
    sPosition.DirectionTitle        = GetPosiDirectionTitle(sPosition.Direction)
    sPosition.OpenDate              = pPosition.GetOpenDate()
    sPosition.TradeID               = TrimSpace(pPosition.GetTradeID())
    sPosition.Volume                = pPosition.GetVolume()
    sPosition.OpenPrice             = Decimal(pPosition.GetOpenPrice(), 4)
    sPosition.TradingDay            = pPosition.GetTradingDay()
    sPosition.SettlementPrice       = Decimal(pPosition.GetSettlementPrice(), 4)
    sPosition.ExchangeID            = pPosition.GetExchangeID()
    sPosition.CloseProfitByDate     = Decimal(pPosition.GetCloseProfitByDate(), 4)
    sPosition.CloseProfitByTrade    = Decimal(pPosition.GetCloseProfitByTrade(), 4)
    sPosition.PositionProfitByDate  = Decimal(pPosition.GetPositionProfitByDate(), 4)
    sPosition.PositionProfitByTrade = Decimal(pPosition.GetPositionProfitByTrade(), 4)
    sPosition.CloseVolume           = pPosition.GetCloseVolume()
    sPosition.CloseAmount           = Decimal(pPosition.GetCloseAmount(), 4)
    sPosition.TimeFirstVolume       = pPosition.GetTimeFirstVolume()
    sPosition.PositionDate          = "1"
    sPosition.PositionDateTitle     = "今仓"

    if sPosition.OpenDate != sPosition.TradingDay {
        sPosition.PositionDate      = "2"
        sPosition.PositionDateTitle = "昨仓"
    }

    // 持仓key（成交编号_合约_开仓日期）
    sPosition.MapKey = Sprintf("%v_%v_%v", sPosition.TradeID, sPosition.InstrumentID, sPosition.OpenDate)


    // 记录持仓数据
    MapPositionDetails.Set(sPosition.MapKey, sPosition)

    return sPosition
}