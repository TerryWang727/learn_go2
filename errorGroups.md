基于 errgroup 实现一个 http server 的启动和关闭 ，以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。

// http server
func main() {
  errgroup, context := .WithContext(context.Background())
  svr := http.NewServer()
  errgroup.Go(func() error {
  fmt.Println("http")
  go func() {
    <-context.Done()
    fmt.Println("http context done")
    svr.Shutdown(context.TODO())
  }()
  return svr.Start()
})


// Linux signal register handle
errgroup.Go(func() error {
   exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT}
   sig := make(chan os.Signal, len(exitSignals))
   signal.Notify(sig, exitSignals)
   for {
      select {
        case <-context.Done():
        fmt.Println("signal context done")
        return context.Err()
        case <-sig:
        // do something
        return nil
      }
   }
})

// inject error
errgroup.Go(func() error {
  fmt.Println("inject error ...")
  time.Sleep(time.Second)
  fmt.Println("inject finished")
  return errors.New("inject error")
})

err := errgroup.Wait() // first error return
fmt.Println(err)
