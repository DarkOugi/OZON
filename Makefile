create_db:
	@echo "Наполняем бд мок данными..."
	@goose -dir ./db/migrations postgres "postgres://$(dbUsername):$(dbPassword)@$(dbHost):$(dbPort)/$(dbName)?sslmode=disable" up

create_proto_file:
	@echo "Генерим grpc код"
	@protoc \
      --proto_path=pkg/grpc \
      --go_out=pkg/grpc \
      --go-grpc_out=pkg/grpc \
      pkg/grpc/server.proto


