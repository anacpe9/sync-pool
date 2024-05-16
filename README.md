# anacpe9/sync-pool


## Usage

```bash
go get github.com/anacpe9/sync-pool
```

### Basic
```go
package main

import (
	syncpool "github.com/anacpe9/sync-pool"
)

func main() {
	type Foo struct{ x int }
	type Bar struct{ y int }

	fooPool_01 := syncpool.GetPool[Foo]()
	fooPool_02 := syncpool.GetPool[Foo]()
	barPool_01 := syncpool.GetPool[Bar]()

	// ...

	foo_01 := fooPool_01.Get()
	foo_02 := fooPool_02.Get()
	bar_01 := barPool_01.Get()

	fooPool_01.Put(foo_01)
	fooPool_01.Put(foo_02)

	barPool_01.Put(bar_01)
}
```

### Fiber example

- [Go Fiber with DTO validate as middleware -- Full example](./example/fiber/main.go)

```go
func NewDTOMiddleware[T interface{ Reset() }]() *DTOMiddleware[T] {
	pool := syncpool.GetPool[T]() // when call GetPoo[T], it's call new (T) every time
	return &DTOMiddleware[T]{
		dtoPool: pool,
	}
}
```
```go
func (mid *DTOMiddleware[T]) DTOValidate() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		resMsg := responseHTTPPool.Get()
		dtoBody := mid.dtoPool.Get()
		defer func() {
			dto := *dtoBody
			dto.Reset()
			mid.dtoPool.Put(&dto)

			resMsg.Reset()
			responseHTTPPool.Put(resMsg)
		}()

		// // for body
		// ctx.BodyParser(dtoBody);
		// // ...
		//
		// validate DTO
		// errors := ValidateStruct(*dtoBody)
		// // ...
		// //

		// share dto to next chain
		ctx.Locals("body", *dtoBody)
		return ctx.Next()
	}
}
```

#### Test invalid payload

```bash
curl -i -s --location 'http://localhost:3000/otp' \
  --header 'Content-Type: application/json' \
  --data '{}'
```
```text
HTTP/1.1 400 Bad Request
Date: Thu, 16 May 2024 14:29:50 GMT
Content-Type: application/json
Content-Length: 358

{"error":"'TargetId' has a value of '' which does not satisfy 'required'.","errors":["'TargetId' has a value of '' which does not satisfy 'required'.","'CommandId' has a value of '' which does not satisfy 'required'.","'Timestamp' has a value of '' which does not satisfy 'required'.","'OTP' has a value of '' which does not satisfy 'required'."],"code":400}
```

#### Test success

```bash
curl -i -s --location 'http://localhost:3000/otp' \
  --header 'Content-Type: application/json' \
  --data '{
      "targetId": "123",
      "commandId": "456",
      "timestamp": "202405151730",
      "otp": "123"
  }'
```
```text
HTTP/1.1 201 Created
Date: Thu, 16 May 2024 14:33:48 GMT
Content-Type: application/json
Content-Length: 42

{"message":"Success","code":201,"ok":true}
```
