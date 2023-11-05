package address

import (
	"fmt"

	"github.com/spf13/viper"
)

type LocalOrderAddr struct {
	port uint16
}

func NewLocalOrderAddr(port uint16) *LocalOrderAddr {
	return &LocalOrderAddr{port}
}

func (s *LocalOrderAddr) UpdateAddr(port uint16) {
	s.port = port
}

func (s *LocalOrderAddr) Get() string {

	ipAddr := viper.GetStringSlice(fmt.Sprintf("order-ip-address"))[0]

	return fmt.Sprintf("%v:%v", ipAddr, s.port)
}
