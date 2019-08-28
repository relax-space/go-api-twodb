# go-api-twodb template

You can connect many database at the same time

## prepare

install mysql
- start local mysql(port is 3308)
- create two database handerly named fruit,fruit2

### run test
```bash
$ cd $GOPATH/src/go-api-twodb
$ go test -p 1 -count 1 -v ./...
```

### run
```bash
$ cd $GOPATH/src/go-api-twodb
$ go run .
```


