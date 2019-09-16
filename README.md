# goctp
上海期货交易所 ctp 接口 Golang版 (for linux64)

## 环境
    install go
    install swig

## 构建
    go get -u github.com/mayiweb/goctp
    构建过程比较慢，需要多等一会儿。

## 依赖
    ctp 默认使用的 gbk 编码，需要用到开源库转换为 utf8
    go get -u github.com/axgle/mahonia

## 编译
    make build
    编译成功后会在当前目录生成 ctp 可执行文件。
    ./ctp [运行模式|默认testing]

## 部署发布
    可将 demo 目录下的全部文件复制出来组成一个新项目
    生成 ctp 文件后，使用 ldd ctp 命令查看依赖关系，并将如下文件放在同一文件夹（注意需要有执行权限）:
        ctp
        libruntime,sync-atomic.so
        libthostmduserapi_se.so
        libthosttraderapi_se.so

    将部署目录路径写入 /etc/ld.so.conf 文件最后一行，并执行 /sbin/ldconfig 命令