package channelserver

import (
	"erupe-ce/common/byteframe"
	cfg "erupe-ce/config"
	"erupe-ce/network/mhfpacket"

	"go.uber.org/zap"
)

// Scenario represents scenario counter data.
type Scenario struct {
	MainID uint32
	// 0 = Basic
	// 1 = Veteran
	// 3 = Other
	// 6 = Pallone
	// 7 = Diva
	CategoryID uint8
}

func handleMsgMhfInfoScenarioCounter(s *Session, p mhfpacket.MHFPacket) {
	pkt := p.(*mhfpacket.MsgMhfInfoScenarioCounter)
	scenarios, err := s.server.scenarioRepo.GetCounters()
	if err != nil {
		s.logger.Error("Failed to get scenario counter info from db", zap.Error(err))
		doAckBufSucceed(s, pkt.AckHandle, make([]byte, 1))
		return
	}

	maxScenarios := 128
	if s.server.erupeConfig.RealClientMode == cfg.ZZ {
		// 256 would wrap the count to zero.
		maxScenarios = 255
	}
	if len(scenarios) > maxScenarios {
		scenarios = scenarios[:maxScenarios]
	}

	bf := byteframe.NewByteFrame()
	bf.WriteUint8(uint8(len(scenarios)))
	for _, scenario := range scenarios {
		bf.WriteUint32(scenario.MainID)
		// If item exchange
		switch scenario.CategoryID {
		case 3, 6, 7:
			bf.WriteBool(true)
		default:
			bf.WriteBool(false)
		}
		bf.WriteUint8(scenario.CategoryID)
	}
	doAckBufSucceed(s, pkt.AckHandle, bf.Data())
}
