一、微服务架构
现在BFF是跟路由合在一起的agent做路由和数据处理和转发 由lua开发 后续会转go
baseserver 也是强依赖mysql和redis做数据缓存和数据持久化处理
前后的结构优化xmind文件展示


优化之前存在的问题  
1、大部分基类未作提取 服务之间高耦合  
2、针对外部依赖跟引用优化目录结构
3、根据功能点针对抽象
4、将中间件提取
5、针对rpc使用提取修改单独提取出来
6、依赖注入不明确

同时做了上述优化后，针对goroutine的生命周期处理
之前大部分的goroutine野生没有管理
go func(){
	// 业务逻辑
}()
完全没有管理这个 goroutine 的生命周期，如果代码里面造成 panic 还使得整个程序崩溃。

现在把 errgroup 包里的代码拷贝处理修改，管理 golang 的生命周期：

var g errgroup.Group
g.Go()
g.Go()
g.Wait()
errgroup 可以使用 context 的方式管理 goroutine 声明周期，同时适用 defer revocer 捕获 panic ，防止意外情况发生，大大提升了代码的可靠性。

总结
经过本期学期对于go有了初步认识 错误处理 存储处理 怎样工程化最佳实践的等等  还要继续学习  道阻且长 坚持才是胜利  最后恭喜发财
