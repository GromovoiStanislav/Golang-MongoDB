## Example fiber web-server with MongoDB

setup:

```
go mod init main

go get go.mongodb.org/mongo-driver/mongo
go get github.com/bxcodec/faker/v3 
go get github.com/gofiber/fiber/v2  
go get github.com/gofiber/fiber/v2/middleware/cors
```

start:

```
go run ./server
go run ./client
```
