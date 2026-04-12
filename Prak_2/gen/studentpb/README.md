# Generated code — studentpb

Эта директория содержит сгенерированный Go-код из `proto/student.proto`.

## Генерация

Выполните из корня проекта:

```bash
protoc --proto_path=proto --go_out=. --go-grpc_out=. proto/student.proto
```

После этого здесь появятся файлы:
- `student.pb.go` — структуры данных (сообщения protobuf)
- `student_grpc.pb.go` — интерфейсы и заглушки gRPC-сервиса

## Требования

- `protoc` (Protocol Buffers compiler)
- `protoc-gen-go`: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`
- `protoc-gen-go-grpc`: `go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest`

Убедитесь, что `$GOPATH/bin` добавлен в `PATH`.
