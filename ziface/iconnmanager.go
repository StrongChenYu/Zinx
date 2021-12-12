package ziface

type IConnManager interface {
	// 增删查
	Add(connection IConnection)
	Delete(connection IConnection)
	Get(uint32) (IConnection, error)

	// how many
	Len() int

	// clear all
	ClearAll()
}
