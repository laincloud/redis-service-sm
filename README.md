# redis-service-sm

[![MIT license](https://img.shields.io/github/license/mashape/apistatus.svg)](https://opensource.org/licenses/MIT)

## 基本介绍

redis-service-sm 是 lain layer 2 的应用，主要提供key-value storage及缓存服务。

基于redis2.6之后的版本，redis单主多从的自动化运维管理

对于redis 主从的介绍可以见[官方文档](http://redis.io/topics/replication)

## 主要功能

```
1. 智能构建：自动初始化一个一主一备的redis 集群，并每分钟检查一次集群状态。
2. 智能切换：通过三节点的redis sentinel集群监控redis server集群状态，并在master节点故障时自动主从切换。
3. 智能代理：redis-service-sm通过proxy与master保持连接，客户端无需关注redis主从切换。
4. 集群监控：将redis server实时运行参数信息发送至监控系统。
5. 数据备份：在集群配置了backupd的情况下，可以通过backupd实现全量物理备份以及增量备份，这两者都是数据恢复的基础。
```

## License
 
redis-service-sm 遵循[MIT](https://github.com/laincloud/redis-service-sm/blob/master/LICENSE)开源协议。
