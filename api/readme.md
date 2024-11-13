# zerocms

打造可能是go-zero下最好用的快速开发平台

## 推荐环境
go-zero 1.7.3+  
goctl 1.5.5+  
mysql5.7+  
redis latest 

## 快速开始

### admin api模块
```shell
# 从admin下进入api定义目录
cd api
# 执行goctl api代码生成
goctl api go --api admin.api --dir ..
```

```shell
# 从model下进入sql定义目录
cd sql
# 执行goctl model代码生成
 goctl model mysql ddl --src *.sql --dir .. --cache
```