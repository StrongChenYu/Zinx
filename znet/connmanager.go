package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	conn  map[uint32]ziface.IConnection
	mutex sync.RWMutex
}

func NewConnManager() ziface.IConnManager {
	connManager := &ConnManager{
		conn: make(map[uint32]ziface.IConnection),
	}
	return connManager
}

// add connection to connection manager
func (connManager *ConnManager) Add(connection ziface.IConnection) {
	connManager.mutex.Lock()
	defer connManager.mutex.Unlock()

	connManager.conn[connection.GetConnID()] = connection
	fmt.Printf("Add connection: %d to connection manager\n", connection.GetConnID())
}

// delete connection from connection manager
func (connManager *ConnManager) Delete(connection ziface.IConnection) {
	connManager.mutex.Lock()
	defer connManager.mutex.Unlock()

	delete(connManager.conn, connection.GetConnID())
	fmt.Printf("Delete connection: %d from connection manager\n", connection.GetConnID())
}

// get connection from connection manager
func (connManager *ConnManager) Get(cId uint32) (ziface.IConnection, error) {
	connManager.mutex.RLock()
	defer connManager.mutex.RUnlock()

	if conn, ok := connManager.conn[cId]; ok {
		return conn, nil
	}

	return nil, errors.New("connection is not existed\n!")
}

// return length of connection manager
func (connManager *ConnManager) Len() int {
	connManager.mutex.RLock()
	defer connManager.mutex.RUnlock()

	return len(connManager.conn)
}

// clear all connections in manager
func (connManager *ConnManager) ClearAll() {
	connManager.mutex.Lock()
	defer connManager.mutex.Unlock()

	for id, conn := range connManager.conn {
		conn.Stop()
		delete(connManager.conn, id)
	}

	fmt.Printf("Clear all connection\n")
}
