package interfaces

type EnvService interface {
	GetListenAddr() string
	GetTLSCertPath() string
	GetTLSKeyPath() string
}
