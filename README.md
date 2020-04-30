## 项目说明
Golang学习笔记，代码整理
## 目录说明
* ### crawler-distributed 
  分布式爬虫，利用Go在并发性上的天然优势实现爬虫任务的分发和调度完成并发需求；使用rpc分离并独立单机版中的并发任务，实现分布式爬虫。

* ### crawler 
  单任务爬虫，应用广度优先算法框架，嵌入数据爬取，信息提取等逻辑实现基本爬虫任务。

* ### goroutine 
  使用goroutine实现的聊天服务器，迷宫算法等。

* ### loadbalance 
  一致性Hash负载均衡的简单实现

* ### logcollection 
  日志收集系统，使用tail读取收集日志，kafka进行分发，beego实现后台管理页面，etcd对服务器ip、目录配置切换

* ### rpc 
  使用json实现rpc调用

* ### seckill 
  商城秒杀系统，秒杀接入层到逻辑层通过redis队列(LPUSH)来通信，完成秒杀逻辑后返回(RPOP)结果

* ### video_server 
  视频点播系统，go语言原生template实现web后端，集成阿里云SDK实现Cloud Native
