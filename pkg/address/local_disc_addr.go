package address

import (
	"fmt"

	"github.com/spf13/viper"
)

type LocalDiscAddr struct {
	port uint16
}

func NewLocalDiscAddr(port uint16) *LocalDiscAddr {
	return &LocalDiscAddr{port}
}

func (s *LocalDiscAddr) UpdateAddr(port uint16) {
	s.port = port
}

func (s *LocalDiscAddr) Get() string {
	ipAddr := string(viper.GetString("disc-ip-address"))
	return fmt.Sprintf("%v:%v", ipAddr, s.port)
}
