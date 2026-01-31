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

func (s *processEnvServiceImpl) getValue(key EnvKey) string {
	return os.Getenv(string(key))
}

func (s *processEnvServiceImpl) GetListenAddr() string {
	addr := s.getValue(LISTEN_ADDR)
	if addr == "" {
		return ":4000"
	}
	return addr
}
