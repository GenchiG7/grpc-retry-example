# grpc-retry-example
test retry behavior for the retry mechanism of gRPC

1. run test server
```bash
go run greeter_server/main.go 
```

2. run client
```bash
go run greeter_client/main.go
```

3. the stdout of server will show log with retry times, like
```bash
2022/12/27 23:58:39 Received: world
2022/12/27 23:58:40 Received: world
2022/12/27 23:58:41 Received: world
2022/12/27 23:58:42 Received: world
2022/12/27 23:58:44 Received: world
```

