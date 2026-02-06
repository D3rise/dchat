package env

import (
	"os"

	"github.com/D3rise/dchat/internal/interfaces"
)

type EnvKey string

func NewProcessEnvService() interfaces.EnvService {
	return &processEnvServiceImpl{}
}

type processEnvServiceImpl struct{}

func (s *processEnvServiceImpl) GetDatabaseDSN() string {
	return s.getValue(DB_DSN)
}

func (s *processEnvServiceImpl) GetTLSKeyPath() string {
	return s.getValue(TLS_KEY_PATH)
}

func (s *processEnvServiceImpl) GetTLSCertPath() string {
	return s.getValue(TLS_CERT_PATH)
}

func (s *processEnvServiceImpl) GetListenAddr() string {
	addr := s.getValue(LISTEN_ADDR)
	if addr == "" {
		return ":4000"
	}
	return addr
}

func (s *processEnvServiceImpl) getValue(key EnvKey) string {
	return os.Getenv(string(key))
}
