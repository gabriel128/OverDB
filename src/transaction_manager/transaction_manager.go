package transaction_manager

import "overdb/src/rm_server"

type TransactionManager struct {

}

func (tm TransactionManager) Create() rm_server.RMServer {
	return tm
}
