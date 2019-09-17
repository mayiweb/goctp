# goctp
上海期货交易所 ctp 接口 Golang版 (for linux64)

## 环境
    install go
    install swig

## 构建
    go get -u github.com/mayiweb/goctp
    构建过程比较慢，需要多等一会儿
    通过 go get 下载会自动构建，也可以手动构建，进入 github.com/mayiweb/goctp 目录 执行 make install 即可

## 依赖
    ctp 默认使用 gbk 编码，需要用到开源库转换为 utf8
    go get -u github.com/axgle/mahonia

## 编译
    进入 github.com/mayiweb/goctp/demo 目录，或将该目录里面的文件全部复制出来组成一个新项目
    make build
    编译成功后会在当前目录生成 ctp 可执行文件（可以修改 Makefile 文件改变生成的文件名）
    ./ctp [运行模式|默认test]

## 部署发布
    生成 ctp 文件后，使用 ldd ctp 命令查看依赖关系，并将如下文件放在同一文件夹（注意需要有执行权限）:
        ctp
        libruntime,sync-atomic.so
        libthostmduserapi_se.so
        libthosttraderapi_se.so

    将部署目录路径写入 /etc/ld.so.conf 文件最后一行，并执行 /sbin/ldconfig 命令