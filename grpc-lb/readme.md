# 验证grpc-lb行为的demo

## 结论
客户端使用默认值的lb策略pick first时, 会自动找到第一个可用的ip.
有一台服务器重启了, 如果在触发重连之前就已经重启完毕, 那么客户端还是会连到这台服务器.
如果没重启, 那么会找到下一台可用的服务器(TODO: 怎么定义可用? tcp能通?)

如果是round robin, 服务器重启了, 客户端一定是连接到下一台服务器的

## 其它知识
1. 由于服务器重启导致双向stream断开, 此时重连的逻辑需要重新发起一个双向strea的grpc请求, 在原来的stream去recv, 是一直等不到消息的
1. 如果客户端同个连接的stream数量达到了服务端的限制, 此时客户端新的stream会block住, 不会自动切换到另外一个服务器
1. 一个http2连接, 是可以有多个通往不同endpoint的stream的.
1. 服务端返回unavailable, 客户端是不会自动重试其他ip的, 因为grpc内部仅判断能建立http2连接

## 实验操作
1. docker-compose up来启动双向stream
1. 通过日志观察客户端连接了哪个container
1. 用docker stop container_name, docker restart container_name来模拟服务器下线, 重启
