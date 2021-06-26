限流的常用方式
1、计数器
基本实现：
type LimitRate struct {
   rate  int           //阀值
   begin time.Time     //计数开始时间
   cycle time.Duration //计数周期
   count int           //收到的请求数
   lock  sync.Mutex    //锁
}

func (limit *LimitRate) Allow() bool {
   limit.lock.Lock()
   defer limit.lock.Unlock()

   // 判断收到请求数是否达到阀值
   if limit.count == limit.rate-1 {
      now := time.Now()
      // 达到阀值后，判断是否是请求周期内
      if now.Sub(limit.begin) >= limit.cycle {
         limit.Reset(now)
         return true
      }
      return false
   } else {
      limit.count++
      return true
   }
}

func (limit *LimitRate) Set(rate int, cycle time.Duration) {
   limit.rate = rate
   limit.begin = time.Now()
   limit.cycle = cycle
   limit.count = 0
}

func (limit *LimitRate) Reset(begin time.Time) {
   limit.begin = begin
   limit.count = 0
}




2、滑动窗口
kratos框架里circuit breaker用循环列表保存time slot对象的实现，这个实现的好处是不用频繁的创建和销毁time slot对象
基本实现：
type timeSlot struct {
    timestamp time.Time // 这个timeSlot的时间起点
    count     int       // 落在这个timeSlot内的请求数
}

// 统计整个时间窗口中已经发生的请求次数
func countReq(win []*timeSlot) int {
    var count int
    for _, ts := range win {
        count += ts.count
    }
    return count
}

type SlidingWindowLimiter struct {
    mu           sync.Mutex    // 互斥锁保护其他字段
    SlotDuration time.Duration // time slot的长度
    WinDuration  time.Duration // sliding window的长度
    numSlots     int           // window内最多有多少个slot
    windows      []*timeSlot
    maxReq       int // 大窗口时间内允许的最大请求数
}

func NewSliding(slotDuration time.Duration, winDuration time.Duration, maxReq int) *SlidingWindowLimiter {
    return &SlidingWindowLimiter{
        SlotDuration: slotDuration,
        WinDuration:  winDuration,
        numSlots:     int(winDuration / slotDuration),
        maxReq:       maxReq,
    }
}


func (l *SlidingWindowLimiter) validate() bool {
    l.mu.Lock()
    defer l.mu.Unlock()


    now := time.Now()
    // 已经过期的time slot移出时间窗
    timeoutOffset := -1
    for i, ts := range l.windows {
        if ts.timestamp.Add(l.WinDuration).After(now) {
            break
        }
        timeoutOffset = i
    }
    if timeoutOffset > -1 {
        l.windows = l.windows[timeoutOffset+1:]
    }

    // 判断请求是否超限
    var result bool
    if countReq(l.windows) < l.maxReq {
        result = true
    }

    // 记录这次的请求数
    var lastSlot *timeSlot
    if len(l.windows) > 0 {
        lastSlot = l.windows[len(l.windows)-1]
        if lastSlot.timestamp.Add(l.SlotDuration).Before(now) {
            // 如果当前时间已经超过这个时间插槽的跨度，那么新建一个时间插槽
            lastSlot = &timeSlot{timestamp: now, count: 1}
            l.windows = append(l.windows, lastSlot)
        } else {
            lastSlot.count++
        }
    } else {
        lastSlot = &timeSlot{timestamp: now, count: 1}
        l.windows = append(l.windows, lastSlot)
    }


    return result
}

滑动窗口算法将一个大的时间窗口分成多个小窗口，每次大窗口向后滑动一个小窗口，并保证大的窗口内流量不会超出最大值，这种实现比固定窗口的流量曲线更加平滑。

普通时间窗口有一个问题，比如窗口期内请求的上限是100，假设有100个请求集中在前1s的后100ms，100个请求集中在后1s的前100ms，其实在这200ms内就已经请求超限了，但是由于时间窗每经过1s就会重置计数，就无法识别到这种请求超限。


3、令牌桶
算法思想
令牌桶是反向的"漏桶"，它是以恒定的速度往木桶里加入令牌，木桶满了则不再加入令牌。服务收到请求时尝试从木桶中取出一个令牌，如果能够得到令牌则继续执行后续的业务逻辑。如果没有得到令牌，直接返回访问频率超限的错误码或页面等，不继续执行后续的业务逻辑。

特点：由于木桶内只要有令牌，请求就可以被处理，所以令牌桶算法可以支持突发流量。

同时由于往木桶添加令牌的速度是恒定的，且木桶的容量有上限，所以单位时间内处理的请求书也能够得到控制，起到限流的目的。假设加入令牌的速度为 1token/10ms，桶的容量为500，在请求比较的少的时候（小于每10毫秒1个请求）时，木桶可以先"攒"一些令牌（最多500个）。当有突发流量时，一下把木桶内的令牌取空，也就是有500个在并发执行的业务逻辑，之后要等每10ms补充一个新的令牌才能接收一个新的请求。

参数设置
木桶的容量 - 考虑业务逻辑的资源消耗和机器能承载并发处理多少业务逻辑。

生成令牌的速度 - 太慢的话起不到“攒”令牌应对突发流量的效果。

适用场景
适合电商抢购或者微博出现热点事件这种场景，因为在限流的同时可以应对一定的突发流量。如果采用漏桶那样的均匀速度处理请求的算法，在发生热点时间的时候，会造成大量的用户无法访问，对用户体验的损害比较大。

代码实现
type TokenBucket struct {
   rate         int64 //固定的token放入速率, r/s
   capacity     int64 //桶的容量
   tokens       int64 //桶中当前token数量
   lastTokenSec int64 //上次向桶中放令牌的时间的时间戳，单位为秒

   lock sync.Mutex
}

func (bucket *TokenBucket) Take() bool {
   bucket.lock.Lock()
   defer bucket.lock.Unlock()

   now := time.Now().Unix()
   bucket.tokens = bucket.tokens + (now-bucket.lastTokenSec)*bucket.rate // 先添加令牌
   if bucket.tokens > bucket.capacity {
      bucket.tokens = bucket.capacity
   }
   bucket.lastTokenSec = now
   if bucket.tokens > 0 {
      // 还有令牌，领取令牌
      bucket.tokens--
      return true
   } else {
      // 没有令牌,则拒绝
      return false
   }
}

func (bucket *TokenBucket) Init(rate, cap int64) {
   bucket.rate = rate
   bucket.capacity = cap
   bucket.tokens = 0
   bucket.lastTokenSec = time.Now().Unix()
}


4、漏桶
算法思想
漏桶算法是首先想象有一个木桶，桶的容量是固定的。当有请求到来时先放到木桶中，处理请求的worker以固定的速度从木桶中取出请求进行相应。如果木桶已经满了，直接返回请求频率超限的错误码或者页面。
适用场景
漏桶算法是流量最均匀的限流实现方式，一般用于流量“整形”。例如保护数据库的限流，先把对数据库的访问加入到木桶中，worker再以db能够承受的qps从木桶中取出请求，去访问数据库。

存在的问题
木桶流入请求的速率是不固定的，但是流出的速率是恒定的。这样的话能保护系统资源不被打满，但是面对突发流量时会有大量请求失败，不适合电商抢购和微博出现热点事件等场景的限流。

代码实现
// 漏桶
// 一个固定大小的桶，请求按照固定的速率流出
// 如果桶是空的，不需要流出请求
// 请求数大于桶的容量，则抛弃多余请求

type LeakyBucket struct {
   rate       float64    // 每秒固定流出速率
   capacity   float64    // 桶的容量
   water      float64    // 当前桶中请求量
   lastLeakMs int64      // 桶上次漏水微秒数
   lock       sync.Mutex // 锁
}

func (leaky *LeakyBucket) Allow() bool {
   leaky.lock.Lock()
   defer leaky.lock.Unlock()

   now := time.Now().UnixNano() / 1e6
   // 计算剩余水量,两次执行时间中需要漏掉的水
   leakyWater := leaky.water - (float64(now-leaky.lastLeakMs) * leaky.rate / 1000)
   leaky.water = math.Max(0, leakyWater)
   leaky.lastLeakMs = now
   if leaky.water+1 <= leaky.capacity {
      leaky.water++
      return true
   } else {
      return false
   }
}

func (leaky *LeakyBucket) Set(rate, capacity float64) {
   leaky.rate = rate
   leaky.capacity = capacity
   leaky.water = 0
   leaky.lastLeakMs = time.Now().UnixNano() / 1e6
}


分布式限流
-使用redis限流
--单个大流量接口, 使用redis容易产生热点
--pre-request模式对性能有一定影响, 高频的网络往返
--从获取单个quota升级成批量quota - 异步批量获取quota可以大幅度减少redis的请求频次
--但申请的配额需要手动设定静态值, 缺乏灵活    利用历史窗口的数据自动修改quota的请求数量




熔断
熔断 - 客户端限流
为了限制操作的持续时间
当某个用户超过资源配额时, 后端任务会快速拒绝请求, 返回配额不足的错误, 但是拒绝回复仍会消耗一定资源, 有可能因为不断拒绝请求而导致过载
max(0, (requests - K * accepts) / (requests + 1))

降级
丢弃不重要的请求, 提供一个降级的服务, 对某几个服务可进行空回复
基于cpu、错误的降级回复, 回复一些mock值
进入降级时, 不反悔一个复杂数据, 而是从一些缓存中捞取或者直接空回复
降级一般在bff或者gate way层做, 防止缓存污染
降级在意外流量或者意外负载时候触发
降级在不定时演练, 保证功能可用
降级的本质: 提供有损服务
ui模块化, 非核心模块降级
BFF层聚合API, 模块降级
页面上一次缓存副本
默认值、热门推荐等。
流量拦截+定期数据缓存(过期副本策略)
页面降级、延迟服务、写/读降级、缓存降级(local cache)
抛异常、返回约定协议、Mock数据、Fallback处理
Case Study
客户端解析协议失败, APP奔溃
客户端部分协议不兼容, 导致页面失败
local cache 数据源缓存, 发版失效+依赖接口故障, 引起的白屏 - 备份remote cache避免白屏
没有playbook, 导致的MTTR上升


重试
当请求返回错误(例: 配额不足、超时、内部错误等), 对于backend部分节点过载的情况下, 倾向于立刻重试, 但是需要留意重试带来的流量放大
限制重试次数(内网中一般不超过两次)和基于重试分布的策略(重试比例10%)
随机化、指数型增长的重试周期: exponential ackoff + jitter
client侧记录重试次数直方图, 传递到server, 进行分布判定, 交由server判定拒绝
只应该在失败这层重试, 当重试仍然失败, 全局约定错误码"过载, 无需重试", 避免级联重试
Case Study
Nginx upstream retry过大, 导致服务雪崩
业务不幂等, 导致的重试, 数据重复 - (写请求不重试)
多层重试传递, 放大流量引起雪崩


负载均衡
某个服务的负载会完全均匀的分发给所有后端任务, 最忙和最不忙的节点永远消耗同样数量的CPU
均衡的流量分发
可靠的识别异常节点
scale-out, 增加同质节点扩容
减少错误, 提高可用性 (N+2冗余)
backend之间的load差异比较大
每个请求的处理成本不同
物理机环境的差异
服务器很难强同质性
存在内存资源争用(内存缓存、带宽、IO等)
性能因素
FullGC
JVM JIT
JSQ(最闲轮训)
缺乏服务端全局视图, 目标: 需要综合考虑 负载+可用性
the choice-of-2
选择backend: CPU, client: health, inflight, latency 作为指标, 使用一个简单的线性方程进行打分
对新启动的节点使用常量惩罚值(penalty), 以及使用探针方式最小化放量, 进行预热
打分较低的节点, 避免进入"永久黑名单"而无法恢复, 使用统计衰减的方式, 让节点指标逐渐恢复到初始状态
指标计算结合moving average, 使用时间衰减, 计算vt = v(t-1) * b + at * (1-b), b为若干次幂的倒数即: Math.Exp((-span) / 600ms)
最佳实践
变更管理
出问题先恢复可用代码
避免过载
过载保护、流量调度等
依赖管理
任何依赖都可能故障, 做chaos monkey testing, 注入故障测试
优雅降级
有损服务, 避免核心链路依赖故障
重试退避
退让算法, 冻结时间, API retry detail控制策略
超时控制
进程内 + 服务间 超时控制
极限压测 + 故障演练
扩容 + 重启 + 消除有害流量