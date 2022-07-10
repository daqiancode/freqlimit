# freqlimit
A tiny library to limit the request frequency based on redis

### Installation
```shell
go get github.com/daqiancode/freqlimit
```

### Example
```go
// add frequency limit for user/1
limit := freqlimit.NewFreqLimit(getRedisClient(), "user/1")
// add freq limit < 2 times/1 seconds
limit.AddLimit(1, 2)
// add freq limit < 100 times/60 seconds
limit.AddLimit(60, 100)

ok, err:= limit.Incr()
if err!=nil {
    return err
}
if !ok {
    return errors.New("exceed the request frequency")
}

```