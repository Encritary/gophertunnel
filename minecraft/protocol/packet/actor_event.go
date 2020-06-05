package packet

import (
	"bytes"
	"encoding/binary"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

const (
	_ = iota
	ActorEventJump
	ActorEventHurt
	ActorEventDeath
	ActorEventStartAttack
	ActorEventStopAttack
	ActorEventTameFail
	ActorEventTameSucceed
	ActorEventShakeDry
	ActorEventUseItem
	ActorEventEatGrass
	ActorEventFishHookBubble
	ActorEventFishHookPosition
	ActorEventFishHook
	ActorEventFishHookTease
	ActorEventSquidInkCloud
	ActorEventZombieVillagerCure
	ActorEventPlayAmbientSound
	ActorEventRespawn
	ActorEventIronGolemOfferFlower
	ActorEventIronGolemWithdrawFlower
	ActorEventLookingForPartner
	ActorEventHappyVillager
	ActorEventAngryVillager
	ActorEventWitchSpell
	ActorEventFirework
	ActorEventFoundPartner
	ActorEventSilverfishSpawn
	ActorEventGuardianAttack
	ActorEventWitchDrinkPotion
	ActorEventWitchThrowPotion
	ActorEventMinecartTNTPrimeFuse
	ActorEventCreeperPrimeFuse
	ActorEventAirSupplyExpired
	ActorEventPlayerAddXPLevels
	ActorEventElderGuardianCurse
	ActorEventAgentArmSwing
	ActorEventEnderDragonDeath
	ActorEventDustParticles
	ActorEventArrowShake
)

const (
	ActorEventEatingItem = iota + 57
	_
	_
	ActorEventBabyAnimalFeed
	ActorEventDeathSmokeCloud
	ActorEventCompleteTrade
	ActorEventRemoveLeash
	ActorEventLlamaCaravanUpdated
	ActorEventConsumeToken
	ActorEventPlayerCheckTreasureHunterAchievement
	ActorEventEntitySpawn
	ActorEventDragonBreath
	ActorEventItemEntityMerge
	ActorEventStartSwimming
	ActorEventBalloonPop
	ActorEventTreasureHunt
	ActorEventSummonAgent
	ActorEventCrossbowCharged
)

// ActorEvent is sent by the server when a particular event happens that has to do with an entity. Some of
// these events are entity-specific, for example a wolf shaking itself dry, but others are used for each
// entity, such as dying.
type ActorEvent struct {
	// EntityRuntimeID is the runtime ID of the entity. The runtime ID is unique for each world session, and
	// entities are generally identified in packets using this runtime ID.
	EntityRuntimeID uint64
	// EventType is the ID of the event to be called. It is one of the constants that can be found above.
	EventType byte
	// EventData is optional data associated with a particular event. The data has a different function for
	// different events, however most events don't use this field at all.
	EventData int32
}

// ID ...
func (*ActorEvent) ID() uint32 {
	return IDActorEvent
}

// Marshal ...
func (pk *ActorEvent) Marshal(buf *bytes.Buffer) {
	_ = protocol.WriteVaruint64(buf, pk.EntityRuntimeID)
	_ = binary.Write(buf, binary.LittleEndian, pk.EventType)
	_ = protocol.WriteVarint32(buf, pk.EventData)
}

// Unmarshal ...
func (pk *ActorEvent) Unmarshal(buf *bytes.Buffer) error {
	return chainErr(
		protocol.Varuint64(buf, &pk.EntityRuntimeID),
		binary.Read(buf, binary.LittleEndian, &pk.EventType),
		protocol.Varint32(buf, &pk.EventData),
	)
}
