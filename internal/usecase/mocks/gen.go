package mocks

//go:generate mockgen -destination=repository_mock.go -package=mocks messenger/internal/usecase Repo
//go:generate mockgen -destination=message_broker.go -package=mocks messenger/internal/usecase MessageBroker
//go:generate mockgen -destination=token_issuer_mock.go -package=mocks messenger/internal/usecase TokenIssuer
