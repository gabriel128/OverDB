package servers


type ServersConfiguration struct {
	TransactionManager []int
	Kvstores map[int]([]int)
}

var ServersConfig ServersConfiguration = ServersConfiguration{
	TransactionManager: []int{8000, 8001, 8002},
	Kvstores: map[int][]int{
		0: {8003, 8004, 8005},
		1: {8006, 8007, 8008},
	},
}
