package v1

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=dto/handler/config.yaml dto/handler/dto.yaml
//go:generate go run github.com/bufbuild/buf/cmd/buf generate proto/pvz.proto --template proto/buf.gen.yaml
//go:generate go run github.com/bufbuild/buf/cmd/buf generate proto/pvz.proto --template proto/buf.gen.doc.yaml
