## Building on Linux

Dependencies:

  * go >= 1.10
  * redis-server >= 5.0

```
git clone https://github.com/duyk16/notifying-service.git
cd notifying-service
```
Copy config
```
cp config.example.json config.json
```
Install package
```
go get -v ./...
```
Build and run
```
go build -o notifying-service main.go
chmod +x notifying-service
./notifying-service
```

## Config

```
{
    "threads": 1, // CPU usages
    "name": "Secure Rest API",

    "auth_timeout": 5, // Time waiting for socket auth of socket connections

    "database": {
        "mongo_url": "mongodb://localhost:27017",
        "mongo_db": "notifying-service",
        "redis_url": "localhost:6379",
        "redis_password": "",
        "redis_pool_size": 10,
        "redis_db": 0,
        "redis_channel": "TEST"
    }
}
```

## Example

### Using REST API

| | |
|-|-|
|Endpoint|Body|
|/message|{type: "ALL", message: "test message"}|
|/message|{type: "Normal", users: ["userId1", "userId2"], message: "test message"}|

### Using Redis PubSub

Message example
| | |
|-|-|
|Description|Message|
|Send message to All users|{type: "ALL", message: "test message"}|
|Send message to Group users|{type: "Normal", users: ["userId1", "userId2"], message: "test message"}|