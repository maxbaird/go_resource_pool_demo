# go_resource_pool_demo

Simple example of GO resource pool implementation

## How to run

After cloning, first build, and then run.

1. Build the program

```sh
cd go_thread_pool_demo/cmd
go build
```

2. Run the program to do profiling *without* resource pool

```sh
./cmd -memprofile mem_no_pool.prof -cpuprofile cpu_no_pool.prof
```

3. Run the program to do profiling *with* resoruce pool

```sh
./cmd -memprofile mem_pool.prof -cpuprofile cpu_pool.prof -wp
```

4. Get profile info as immages

```sh
go tool pprof -png cpu_no_pool.prof
go tool pprof -png mem_no_pool.prof
go tool pprof -png cpu_pool.prof
go tool pprof -png mem_pool.prof
```
