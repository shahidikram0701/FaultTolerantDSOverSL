package address

import (
	"fmt"

	"github.com/spf13/viper"
)

type LocalDataAddr struct {
	numReplica int32
	basePort   uint16
}

func NewLocalDataAddr(numReplica int32, basePort uint16) *LocalDataAddr {
	return &LocalDataAddr{numReplica, basePort}
}

func (s *LocalDataAddr) UpdateBasePort(basePort uint16) {
	s.basePort = basePort
}

func (s *LocalDataAddr) Get(sid, rid int32) string {
	port := s.basePort + uint16(sid*s.numReplica+rid)
	ipAddr := string(viper.GetString("data-ip-address"))
	return fmt.Sprintf("%v:%v", ipAddr, port)
}
