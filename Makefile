proto:
	protoc --go_out=. --go-grpc_out=. auth/auth.proto
	protoc --go_out=. --go-grpc_out=. roles/roles.proto

up:
	sql-migrate new up

down:
	sql-migrate down