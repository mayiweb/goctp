package main

import (
    "gitee.com/mayiweb/goctp"
    "gitee.com/mayiweb/goctp/safe"
    "os"
    "fmt"
    "log"
)

var (
    // ctp 句柄及配置项
    Ctp CtpClient

    // ctp 队列 句柄
    Queue *CtpQueue

    // 行情模块函数 句柄
    MdSpi FtdcMdSpi

    // 交易模块函数 句柄
    TradeSpi FtdcTradeSpi

    // 交易用户登录信息
    MapLogin safe.Map

    // 资金账户信息 AccountStruct
    MapAccounts safe.Map

    // 交易所合约详情列表 InstrumentStruct
    MapInstruments safe.Map

    // 报单列表（已成交、未成交、撤单等状态）的列表数据 OrderStruct
    MapOrders safe.Map

    // 报单成交列表
    MapTrades safe.Map

    // 持仓列表 PositionStruct
    MapPositions safe.Map

    // 持仓明细列表 PositionDetailStruct
    MapPositionDetails safe.Map

    // 报单有效期任务
    MapOrderInsertTask safe.Map

    // 合约行情Tick数据（最新的一条tick）
    MapInstrumentTicks safe.Map

    // InstrumentID 转 InstrumentCode
    MapInstrumentCodes safe.Map

    // ctp 服务器，及交易账号
    MdFront []string
    TradeFront []string

    BrokerID string
    InvestorID string
    Password string

    // 客户端认证
    AppID string
    AuthCode string

    // ctp 流文件，绝对路径
    StreamFile string = GetCurrentDirectory() + "/StreamFile/"

    // tick 文件保存的目录（路径后面要带斜杠）
    TickDataDirectory string = GetCurrentDirectory() + "/csv/"

    // 买卖方向：买
    Buy byte  = '0'

    // 买卖方向：卖
    Sell byte = '1'

    // 运行模式（prod 生产，test 标准环境测试，dev 24小时测试）
    RunMode string
)

// 设置交易账号
func SetTradeAccount() {

    switch RunMode {

        // 生产环境
        case "prod":
            MdFront     = []string{}
            TradeFront  = []string{}
            BrokerID    = ""
            InvestorID  = ""
            Password    = ""
            AppID       = ""
            AuthCode    = ""

        // 测试环境 simnow (与实际生产环境保持一致)
        case "test":
            MdFront     = []string{"tcp://180.168.146.187:10211", "tcp://180.168.146.187:10212"}
            TradeFront  = []string{"tcp://180.168.146.187:10201", "tcp://180.168.146.187:10202"}
            BrokerID    = "9999"
            InvestorID  = ""
            Password    = ""
            AppID       = ""
            AuthCode    = ""

        // 7*24 服务器，交易日，16：00～次日09：00；非交易日，16：00～次日15：00
        case "dev":
            MdFront     = []string{"tcp://180.168.146.187:10131"}
            TradeFront  = []string{"tcp://180.168.146.187:10130"}
            BrokerID    = "9999"
            InvestorID  = ""
            Password    = ""
            AppID       = ""
            AuthCode    = ""

        default:
            Println("该模式未设置交易账号信息")
            os.Exit(1)
    }
}

func init() {
    // 全局 行情、交易 函数句柄
    MdSpi    = FtdcMdSpi{}
    TradeSpi = FtdcTradeSpi{}

    // 全局队列句柄
    Queue = &CtpQueue{}

    // 运行模式【运行程序时带上参数可设置】
    if len(os.Args) < 2 {
        RunMode = "test"
    } else {
        RunMode = os.Args[1]
    }

    // 需要检查的目录
    CheckDirectory := []string{StreamFile, TickDataDirectory}
    for _, val := range CheckDirectory {
        // 检查目录是否存在
        fileExists, _ := PathExists(val)
        if !fileExists {
            err := os.Mkdir(val, os.ModePerm)
            if err != nil {
               LogPrintln("创建目录失败，请检查是否有操作权限 Err:", err)
            }
        }
    }

    // 启动队列查询任务
    go RunQueryTask()
}

func main() {

    // 设置交易账号
    SetTradeAccount()

    log.Println("启动交易程序")

    // 检查流文件目录是否存在
    fileExists, _ := PathExists(StreamFile)
    if !fileExists {
        err := os.Mkdir(StreamFile, os.ModePerm)
        if err != nil {
           fmt.Println("创建目录失败，请检查是否有操作权限")
        }
    }

    Ctp = CtpClient {
        MdApi: goctp.CThostFtdcMdApiCreateFtdcMdApi(StreamFile),
        TradeApi: goctp.CThostFtdcTraderApiCreateFtdcTraderApi(StreamFile),
        BrokerID: BrokerID,
        InvestorID: InvestorID,
        Password: Password,
        AppID: AppID,
        AuthCode: AuthCode,
        MdRequestId: 1,
        TradeRequestId: 1,
        IsTradeInit: false,
        IsTradeInitFinish: false,
        IsMdLogin: false,
        IsTradeLogin: false,
    }

    Ctp.MdApi.RegisterSpi(goctp.NewDirectorCThostFtdcMdSpi(&FtdcMdSpi{Ctp}))

    for _, val := range MdFront {
        Ctp.MdApi.RegisterFront(val)
    }
    Ctp.MdApi.Init()


    Ctp.TradeApi.RegisterSpi(goctp.NewDirectorCThostFtdcTraderSpi(&FtdcTradeSpi{Ctp}))

    for _, val := range TradeFront {
        Ctp.TradeApi.RegisterFront(val)
    }

    Ctp.TradeApi.SubscribePublicTopic(goctp.THOST_TERT_QUICK);
    Ctp.TradeApi.SubscribePrivateTopic(goctp.THOST_TERT_QUICK);
    Ctp.TradeApi.Init()
    Ctp.TradeApi.Join()

    // .Join() 如果后面有其它需要处理的功能可以不写，但必须保证程序不能退出，Join 就是保证程序不退出的
}