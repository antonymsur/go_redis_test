# Redis,PostgreSQL Benchmark with <i><b>GO</b></i>

Benchmarking test has been performed in Linux Fedora 23 running in Virtual Machine with the following config. 

### CPU information(lscpu)
* Architecture:          x86_64
* CPU op-mode(s):        32-bit, 64-bit
* Byte Order:            Little Endian
* CPU(s):                1
* On-line CPU(s) list:   0
* Thread(s) per core:    1
* Core(s) per socket:    1
* Socket(s):             1
* NUMA node(s):          1
* Vendor ID:             GenuineIntel
* CPU family:            6
* Model:                 58
* Model name:            Intel(R) Core(TM) i7-3517U CPU @ 1.90GHz
* Stepping:              9
* CPU MHz:               2394.560
* BogoMIPS:              4789.12
* Hypervisor vendor:     KVM
* Virtualization type:   full
* L1d cache:             32K
* L1i cache:             32K
* L2 cache:              256K
* L3 cache:              4096K
* NUMA node0 CPU(s):     0

### RAM/Memory Info(free)
<pre>
          total      used       free   shared  buff/cache   available
Mem:     2025596    662860    191560    54040     1171176     1231804
Swap:    2097148       968   2096180</pre>

### Kernel,Linux,Architecture(uname -r)

<b>4.4.4-301.fc23.x86_64</b>


#### Test Command :

<b>go test -v -bench=. -count 1 -benchtime 1s  myapp</b><br>

<b>RESULT</b> with <i><b>PostgreSQL</b></i> as DB gorm as ORM utility

<pre>
PASS
BenchmarkAddUser         	     100	  10349840 ns/op
BenchmarkAddSubscriptions	     300	   6722208 ns/op
BenchmarkAddApps         	     200	   6439089 ns/op
ok  	myapp	6.363s</pre>

<b>RESULT</b> with <i><b>REDIS</b></i> as NoSQL DB and Caching with RDB and AOF (everysec)
<pre>
PASS
BenchmarkAddUser         	     500	   2634495 ns/op
BenchmarkAddSubscriptions	    1000	   1579358 ns/op
BenchmarkAddApps         	    2000	    858782 ns/op
ok  	myapp	8.116s</pre>





