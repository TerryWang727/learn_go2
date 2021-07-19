1、使用 redis benchmark 工具, 测试 10 20 50 100 200 1k 5k 字节 value 大小，redis get set 性能。

使用方法

Usage: redis-benchmark [-h <host>] [-p <port>] [-c <clients>] [-n <requests]> [-k <boolean>]

 -h <hostname>      Server hostname (default 127.0.0.1)
 -p <port>          Server port (default 6379)
 -s <socket>        Server socket (overrides host and port)
 -a <password>      Password for Redis Auth
 -c <clients>       Number of parallel connections (default 50)
 -n <requests>      Total number of requests (default 100000)
 -d <size>          Data size of SET/GET value in bytes (default 2)
 -dbnum <db>        SELECT the specified db number (default 0)
 -k <boolean>       1=keep alive 0=reconnect (default 1)
 -r <keyspacelen>   Use random keys for SET/GET/INCR, random values for SADD
  Using this option the benchmark will expand the string __rand_int__
  inside an argument with a 12 digits number in the specified range
  from 0 to keyspacelen-1. The substitution changes every time a command
  is executed. Default tests use this to hit random keys in the
  specified range.
 -P <numreq>        Pipeline <numreq> requests. Default 1 (no pipeline).
 -q                 Quiet. Just show query/sec values
 --csv              Output in CSV format
 -l                 Loop. Run the tests forever
 -t <tests>         Only run the comma separated list of tests. The test
                    names are the same as the ones produced as output.
 -I                 Idle mode. Just open N idle connections and wait.
  
  
 [root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 20
====== SET ======
  100000 requests completed in 1.48 seconds
  50 parallel clients
  20 bytes payload
  keep alive: 1

99.73% <= 1 milliseconds
99.85% <= 2 milliseconds
99.88% <= 3 milliseconds
99.90% <= 6 milliseconds
99.95% <= 11 milliseconds
100.00% <= 11 milliseconds
67430.88 requests per second

====== GET ======
  100000 requests completed in 1.19 seconds
  50 parallel clients
  20 bytes payload
  keep alive: 1

99.78% <= 1 milliseconds
99.92% <= 4 milliseconds
99.97% <= 5 milliseconds
100.00% <= 5 milliseconds
84317.03 requests per second
  
  
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 50
====== SET ======
  100000 requests completed in 1.00 seconds
  50 parallel clients
  50 bytes payload
  keep alive: 1

99.90% <= 1 milliseconds
100.00% <= 1 milliseconds
100200.40 requests per second

====== GET ======
  100000 requests completed in 1.06 seconds
  50 parallel clients
  50 bytes payload
  keep alive: 1

99.96% <= 1 milliseconds
100.00% <= 1 milliseconds
94696.97 requests per second
  
  
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 100
====== SET ======
  100000 requests completed in 1.02 seconds
  50 parallel clients
  100 bytes payload
  keep alive: 1

99.96% <= 1 milliseconds
100.00% <= 1 milliseconds
98425.20 requests per second

====== GET ======
  100000 requests completed in 1.01 seconds
  50 parallel clients
  100 bytes payload
  keep alive: 1

99.93% <= 1 milliseconds
100.00% <= 1 milliseconds
98522.17 requests per second
  
  
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 200
====== SET ======
  100000 requests completed in 1.04 seconds
  50 parallel clients
  200 bytes payload
  keep alive: 1

99.98% <= 1 milliseconds
100.00% <= 1 milliseconds
95877.28 requests per second

====== GET ======
  100000 requests completed in 1.02 seconds
  50 parallel clients
  200 bytes payload
  keep alive: 1

99.96% <= 1 milliseconds
100.00% <= 1 milliseconds
98328.42 requests per second
  
  
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 1000
====== SET ======
  100000 requests completed in 1.05 seconds
  50 parallel clients
  1000 bytes payload
  keep alive: 1

99.89% <= 1 milliseconds
99.95% <= 3 milliseconds
100.00% <= 3 milliseconds
95602.30 requests per second

====== GET ======
  100000 requests completed in 1.11 seconds
  50 parallel clients
  1000 bytes payload
  keep alive: 1

99.89% <= 1 milliseconds
100.00% <= 1 milliseconds
90497.73 requests per second
           
           
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 2000
====== SET ======
  100000 requests completed in 1.01 seconds
  50 parallel clients
  2000 bytes payload
  keep alive: 1

99.93% <= 1 milliseconds
100.00% <= 1 milliseconds
99206.34 requests per second

====== GET ======
  100000 requests completed in 1.10 seconds
  50 parallel clients
  2000 bytes payload
  keep alive: 1

99.89% <= 1 milliseconds
100.00% <= 1 milliseconds
90579.71 requests per second
           
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 4000
====== SET ======
  100000 requests completed in 1.08 seconds
  50 parallel clients
  4000 bytes payload
  keep alive: 1

99.95% <= 1 milliseconds
100.00% <= 1 milliseconds
92592.59 requests per second

====== GET ======
  100000 requests completed in 1.17 seconds
  50 parallel clients
  4000 bytes payload
  keep alive: 1

99.88% <= 1 milliseconds
100.00% <= 1 milliseconds
85397.09 requests per second

[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 5000
====== SET ======
  100000 requests completed in 1.15 seconds
  50 parallel clients
  5000 bytes payload
  keep alive: 1

99.85% <= 1 milliseconds
100.00% <= 1 milliseconds
86805.56 requests per second

====== GET ======
  100000 requests completed in 1.23 seconds
  50 parallel clients
  5000 bytes payload
  keep alive: 1

99.82% <= 1 milliseconds
100.00% <= 1 milliseconds
81499.59 requests per second
           
           
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 10000
====== SET ======
  100000 requests completed in 1.19 seconds
  50 parallel clients
  10000 bytes payload
  keep alive: 1

99.89% <= 1 milliseconds
100.00% <= 1 milliseconds
84033.61 requests per second

====== GET ======
  100000 requests completed in 1.32 seconds
  50 parallel clients
  10000 bytes payload
  keep alive: 1

99.75% <= 1 milliseconds
99.99% <= 2 milliseconds
100.00% <= 2 milliseconds
75757.57 requests per second
  
  
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 15000
====== SET ======
  100000 requests completed in 1.23 seconds
  50 parallel clients
  15000 bytes payload
  keep alive: 1

99.89% <= 1 milliseconds
100.00% <= 1 milliseconds
81037.28 requests per second

====== GET ======
  100000 requests completed in 1.44 seconds
  50 parallel clients
  15000 bytes payload
  keep alive: 1

99.43% <= 1 milliseconds
99.98% <= 2 milliseconds
100.00% <= 2 milliseconds
69252.08 requests per second
           
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 20000
====== SET ======
  100000 requests completed in 1.30 seconds
  50 parallel clients
  20000 bytes payload
  keep alive: 1

98.80% <= 1 milliseconds
100.00% <= 2 milliseconds
100.00% <= 2 milliseconds
76745.97 requests per second

====== GET ======
  100000 requests completed in 1.88 seconds
  50 parallel clients
  20000 bytes payload
  keep alive: 1

98.58% <= 1 milliseconds
99.94% <= 2 milliseconds
99.98% <= 5 milliseconds
100.00% <= 6 milliseconds
53219.80 requests per second
  
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 40000
====== SET ======
  100000 requests completed in 1.85 seconds
  50 parallel clients
  40000 bytes payload
  keep alive: 1

92.57% <= 1 milliseconds
99.90% <= 2 milliseconds
99.99% <= 3 milliseconds
100.00% <= 3 milliseconds
54171.18 requests per second

====== GET ======
  100000 requests completed in 2.56 seconds
  50 parallel clients
  40000 bytes payload
  keep alive: 1

95.31% <= 1 milliseconds
99.97% <= 2 milliseconds
100.00% <= 2 milliseconds
39108.33 requests per second
           
[root@test-redis]# redis-benchmark -h 127.0.0.1 -p 6379 -t set,get -d 50000
====== SET ======
  100000 requests completed in 2.18 seconds
  50 parallel clients
  50000 bytes payload
  keep alive: 1

45.04% <= 1 milliseconds
99.86% <= 2 milliseconds
99.95% <= 3 milliseconds
99.96% <= 4 milliseconds
100.00% <= 4 milliseconds
45955.88 requests per second

====== GET ======
  100000 requests completed in 3.68 seconds
  50 parallel clients
  50000 bytes payload
  keep alive: 1

59.49% <= 1 milliseconds
99.57% <= 2 milliseconds
99.82% <= 3 milliseconds
99.83% <= 5 milliseconds
99.87% <= 6 milliseconds
99.93% <= 8 milliseconds
99.95% <= 9 milliseconds
99.95% <= 11 milliseconds
99.99% <= 12 milliseconds
100.00% <= 12 milliseconds
27173.91 requests per second
  

2、写入一定量的 kv 数据, 根据数据大小 1w-50w 自己评估, 结合写入前后的 info memory 信息  , 分析上述不同 value 大小下，平均每个 key 的占用内存空间。
  
Redis字符串占用的内存，远比实际字符串的长度要大。
  
初始 used_memory都是873584 
 
Test1 value10字节 1w个key set后used_memory 1645704 平均每个 key 的占用内存空间 77.212
Test1 value100字节 1w个key set后used_memory 2604824 平均每个 key 的占用内存空间 173.124
Test1 value1000字节 1w个key set后used_memory 11725688 平均每个 key 的占用内存空间 1085.2104
Test1 value5000字节 1w个key set后used_memory 52686120 平均每个 key 的占用内存空间 5181.2536
