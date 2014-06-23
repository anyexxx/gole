gole
====

go+redis+mongoDB

初学乍练，用go做了个游戏服务器，搭了个架子，拍砖的欢迎···

目前已实现的特点有：

1.无服务器区分，全游戏大世界架构，前端有负载均衡，一般的云服务都有此服务提供

2.玩家数据按账号区分在不同的redis实例上，对redis的读写与单组服务器相同，保证效率

3.采用了MongoDB集群，对全游戏活动排名等活动有高效存储保证

4.目前支持pb和json数据格式以及socket,websocket连接，切换方便


运行：
设置好GOPATH，将src的文件夹替换到GOPATH下的就可以了

src/me.qqtu.game为服务器逻辑，其他为第三方包
src/me.qqtu.game/compile下的都是要编译的文件，用go install即可

整体还是很简单，主要是希望有大神看下这样的结构有没啥隐患，或者对redis和mongodb集群的使用有没啥建议，另外有go的高手指导下就更好了···谢谢啊 :)

联系：ruipengliu@live.cn
