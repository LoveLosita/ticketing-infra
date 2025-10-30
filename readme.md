# 1.目标与约束

下列是本票务系统的目标：

- 高并发秒杀型开售（峰值 10–50 万 QPS 入站，核心写路径受控在 5–10k TPS）。
- 不超卖；座位原子性预占与释放；支付成功后再出票。
- 多终端/多地区接入；队列防饥饿与公平性；严格防重入。
- 核验一次性二维码（离线可验，在线可撤销），此处为了方便改用密钥。生产环境中可以使用二维码识别程序将二维码转换成密钥后再传入。
- 合规存储 PII，审计支付回调与操作日志。

# 2.模块分布

以下为该rpc系统的各个模块：

- api-gateway（WAF/验证码/限流）
- queue-service（排队发号、资格校验）
- inventory-service（座位库存、预占、释放、分段压测）
- order-service（订单、支付对接、回调幂等）
- ticket-service（出票、二维码签名、状态机）
- checkin-service（检票、一次性验证、撤销/退票联动）
- risk-service（设备指纹、黑白名单、打分、行为速率）
- refund-waitlist-service（退票回库、候补补票）
- ops/obs（日志、指标、链路追踪、审计）

# 3.技术栈

- 编程语言：Go（主语言），thrift（用于rpc通信）
- 框架：Kitex（rpc），Hertz（在kitex客户端区域作为网关）
- 数据库：MySQL；主从 + 只读分离
- 缓存/队列：Redis，Kafka（异步审计、事件驱动）
- 搜索/分析：ClickHouse/Elastic（可选，审计与看板）
- 签名：KMS/HSM 存储私钥，Ed25519/Secp256k1
- 鉴权：JWT
- 负载均衡：Nginx
- 监控：Prometheus + Grafana，OpenTelemetry

# 4.业务核心链路

1.开售排队与发号

- 进入等待室：校验人机、设备指纹、黑名单；ZSET 排队。
- 发放 queue_token：基于队列位置与风控评分；签名 token 防伪造；TTL 短（如 2–3 分钟）。
- 进入购买页需携带 queue_token，网关与 queue-service 验签 + 幂等登记，限制并发窗口（如每秒放行 N 人）。

2.预占座（Hold）与释放

- 用户选座或系统自动分配：
  - Redis Lua 原子预占 seat，写入 seat_holds + 事务日志（DB 最终一致补偿）。
  - 预占 TTL 例如 120 秒；超时释放（Redis 过期 + 定时任务兜底）。
- 预占成功后生成临时报价快照（防止价格变化）。

3.下单与支付

- 创建订单（DB 事务）：
  - 幂等键 idempotency_key（前端/服务端生成，存 Orders.unique）
  - 校验持有的 hold_token 是否有效且归属当前用户
- 调用支付渠道，落库 payments(pending)。
- 支付回调
  - 校验签名，按 (channel,out_trade_no) 幂等更新为 succeeded/failed
  - 成功则推进订单状态到 paid 并发出票事件；失败则触发释放座位（如还在 held）

4.出票与二维码一次性核验

- 出票：为每个 seat 生成 Ticket，ticket_sn 全局唯一；签名 QR payload
  - QR 内容建议：base64(event_id, seat_id, ticket_sn, ts, nonce) + 签名
  - 私钥存 KMS/HSM，ticket-service 仅可签名，不可导出密钥
- 检票
  - 在线核验：checkin-service 验签payload -> 校验门禁 nonce 一次性 -> 落 checkin_logs 并原子设置已验
  - 离线容错：可用短时 CRL 缓存 + 有界时钟偏移校验；门禁定时同步撤销列表
- 一票一验：以 tickets.id 或 ticket_sn 建唯一 checkin_logs 索引保证原子性

5.退票回库与候补

- 退票申请 -> 审核 -> 改票状态 refunded -> 释放 seat：
  - seat.status: sold -> available（或由 waitlist 立即消费）
- 候补策略
  - waitlist ZSET（priority,joined_at），空位产生后按照策略通知/自动下单
  - 时效+限次，防刷补票

```
进入等待室 → 排队发号 → 验签放行
          ↓
      选座 / 预占座 (Hold)
          ↓
      创建订单 → 支付 → 回调
          ↓
         出票
          ↓
         检票验签
          ↓
      退票 → 回库 → 候补
```

```text
Browser/App
   │
   ▼
Nginx / Cloud LB  ──(HTTP)──►  Hertz 网关 x N（每个进程内嵌 Kitex client）
                                 │
                                 └─(RPC, client-side LB/发现)─►  Kitex Server: queue
                                                               Kitex Server: seat
                                                               Kitex Server: order
                                                               ...

```



# 5.常见状态码及其对应错误

| 状态码 | HTTP状态码 |            描述            | 原因 |        解决方案        |
| :----: | :--------: | :------------------------: | :--: | :--------------------: |
| 20000  |    200     |            成功            |  -   |           -            |
| 40001  |    400     |        用户名已存在        |      | 请使用全新的用户名注册 |
| 40002  |    400     |       传入的参数错误       |      |     请传入正确参数     |
| 40003  |    400     |         没找到用户         |      |   请传入正确的用户名   |
| 40004  |    400     |          密码错误          |      |                        |
| 40005  |    400     |     token签名方法无效      |      |                        |
| 40006  |    400     |         无效token          |      |                        |
| 40007  |    400     |          无效声明          |      |                        |
| 40008  |    400     |        无效令牌类型        |      |                        |
| 40009  |    400     | 更改密码时传入的旧密码错误 |      |                        |
|        |            |                            |      |                        |
|        |            |                            |      |                        |
|        |            |                            |      |                        |
|        |            |                            |      |                        |
|        |            |                            |      |                        |
|        |            |                            |      |                        |
|        |            |                            |      |                        |
|        |            |                            |      |                        |

# 6.接口文档

为了方便我的编写和测试，本次接口文档使用ApiFox托管，链接：https://mnxzotmuxs.apifox.cn

# 7.快速开始