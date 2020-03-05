package storage

func Init() {
	ConnectMongo()
	ConnectRedis()
}
