go run overdb.go -- tm 0 0 &
go run overdb.go -- tm 0 1 &
go run overdb.go -- tm 0 2 &

go run overdb.go -- kv 0 0 &
go run overdb.go -- kv 0 1 &
go run overdb.go -- kv 0 2 &

go run overdb.go -- kv 1 0 &
go run overdb.go -- kv 1 1 &
go run overdb.go -- kv 1 2 &
