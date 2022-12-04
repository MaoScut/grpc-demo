# 验证grpc-lb行为的demo

## 结论
客户端使用默认值的lb策略pick first时, 会自动找到第一个可用的ip.
有一台服务器重启了, 如果在触发重连之前就已经重启完毕, 那么客户端还是会连到这台服务器.
如果没重启, 那么会找到下一台可用的服务器(TODO: 怎么定义可用? tcp能通?)

如果是round robin, 服务器重启了, 客户端一定是连接到下一台服务器的

## 其它知识
由于服务器重启导致双向stream断开, 此时重连的逻辑需要重新发起一个双向strea的grpc请求, 在原来的stream去recv, 是一直等不到消息的

## 实验操作
1. docker-compose up来启动双向stream
1. 通过日志观察客户端连接了哪个container
1. 用docker stop container_name, docker restart container_name来模拟服务器下线, 重启