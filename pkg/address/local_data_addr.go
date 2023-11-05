package address

import (
	"fmt"
	"log"

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
	port := s.basePort
	//ipAddr := string(viper.GetString(fmt.Sprintf("data-ip-address-%d-%d", sid, rid)))

	ipAddr := viper.GetStringSlice(fmt.Sprintf("data-ip-address-%d", sid))[rid]

	log.Printf("The data addr is %v", ipAddr)
	return fmt.Sprintf("%v:%v", ipAddr, port)
}
