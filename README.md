# 饭否消息备份小工具
这是一个用于备份你的饭否消息到本地数据库的小工具，它使用Go语言编写。

## 初次使用
第一次使用时需要指定你的consumer key和consumer secret，你可以在[这里](http://fanfou.com/apps)申请一个自己的应用。
获取你的key和secret后运行一下命令即可。
```
FanfouStatusBakup -key 你的key -secret 你的secret
```
初次运行后，目录下会自动生成一个配置文件，保留这个配置文件可以让你不必每次运行时都手工指定key和secret。

## 关于本地存储
这个小工具将会把你的消息存到一个本地的sqlite数据库中。
