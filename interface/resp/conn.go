package resp

// redis  客户端的连接信息

type Connection interface {
	Write([]byte) error
	// used for multi database
	GetDBIndex() int
	SelectDB(int)
}
