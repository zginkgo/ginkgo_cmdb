## 认证问题

### HTTP接口的认证(OpenAPI)

该接口会别哪些用户使用:
+ Web界面 ---> OpenAPI
+ 其他第三方服务, 监控系统(Prometheus), 监控资源发现

如何认证:
+ basic auth: (user/password) ---> user:
    + 存放的Header： Authorization
    + "Basic " +  user:password 的base64编码
    + header: Authorization = "Basic base64(user:password)"

用户的密码 在每次API调用的时候都会把明文传输, 甚至可以登录Web 去修改你的密码

系统 --> 接口 (编程用户: 程序使用)
+ 需要登录界面的用户
+ 编程用户

最好的方式 不直接使用用户的 user/password, 可以用户颁发一个代表自己身份的令牌(token)

即使用户发现自己的令牌被盗取了, 可以吊销该令牌

怎么基于Token做一个统一认证系统? 我们总不至于每个系统 都自己保存一个user表

### GRPC接口的认证