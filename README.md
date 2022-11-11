## grace Overview
A package for graceful exit of Go programs

## Install
```go
go get github.com/petew9527/grace
```
**Note:** grace uses [Go Modules](https://github.com/golang/go/wiki/Modules) to manage dependencies.

## Examples:
```go
grace.Wait(grace.WithOutTime(time.Second*3), grace.WithHandlers(func() error {return nil}))
```
