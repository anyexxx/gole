gole
====

golang gameserver

1.无服务器区分，全游戏大世界架构

2.游戏服务器做了负载均衡，增强服务器稳定性以及易于横向扩展

3.玩家数据按账号区分在不同的redis实例上，对redis的读写与单组服务器相同，保证效率

4.采用了MongoDB集群，对全游戏活动排名等活动有高效存储保证

5.目前支持pb和json数据格式以及socket,websocket连接，切换方便

