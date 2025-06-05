package nower

import "time"

//go:generate go run go.uber.org/mock/mockgen -source=contract.go -destination=mocks/contract_mock.go -package=nower Nower
type Nower interface {
	Now() time.Time
}
