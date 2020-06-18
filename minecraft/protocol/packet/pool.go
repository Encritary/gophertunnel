package packet

import (
	"reflect"
)

// Register registers a function that returns a packet for a specific ID. Packets with this ID coming in from
// connections will resolve to the packet returned by the function passed.
func Register(id uint32, pk func() Packet) {
	registeredPackets[id] = pk
}

// registeredPackets holds packets registered by the user.
var registeredPackets = map[uint32]func() Packet{}

// Pool is a map holding packets indexed by a packet ID.
type Pool map[uint32]Packet

// NewPool returns a new pool with all supported packets sent. Packets may be retrieved from it simply by
// indexing it with the packet ID.
func NewPool() Pool {
	p := Pool{
		IDLogin:                      &Login{},
		IDPlayStatus:                 &PlayStatus{},
		IDServerToClientHandshake:    &ServerToClientHandshake{},
		IDClientToServerHandshake:    &ClientToServerHandshake{},
		IDDisconnect:                 &Disconnect{},
		IDResourcePacksInfo:          &ResourcePacksInfo{},
		IDResourcePackStack:          &ResourcePackStack{},
		IDResourcePackClientResponse: &ResourcePackClientResponse{},
		IDText:                       &Text{},
		IDSetTime:                    &SetTime{},
		IDStartGame:                  &StartGame{},
		IDAddPlayer:                  &AddPlayer{},
		IDAddActor:                   &AddActor{},
		IDRemoveActor:                &RemoveActor{},
		IDAddItemActor:               &AddItemActor{},
		// ---
		IDTakeItemActor:     &TakeItemActor{},
		IDMoveActorAbsolute: &MoveActorAbsolute{},
		IDMovePlayer:        &MovePlayer{},
		IDRiderJump:         &RiderJump{},
		IDUpdateBlock:       &UpdateBlock{},
		IDAddPainting:       &AddPainting{},
		IDTickSync:          &TickSync{},
		// ---
		IDLevelEvent:                  &LevelEvent{},
		IDBlockEvent:                  &BlockEvent{},
		IDActorEvent:                  &ActorEvent{},
		IDMobEffect:                   &MobEffect{},
		IDUpdateAttributes:            &UpdateAttributes{},
		IDInventoryTransaction:        &InventoryTransaction{},
		IDMobEquipment:                &MobEquipment{},
		IDMobArmourEquipment:          &MobArmourEquipment{},
		IDInteract:                    &Interact{},
		IDBlockPickRequest:            &BlockPickRequest{},
		IDActorPickRequest:            &ActorPickRequest{},
		IDPlayerAction:                &PlayerAction{},
		IDActorFall:                   &ActorFall{},
		IDHurtArmour:                  &HurtArmour{},
		IDSetActorData:                &SetActorData{},
		IDSetActorMotion:              &SetActorMotion{},
		IDSetActorLink:                &SetActorLink{},
		IDSetHealth:                   &SetHealth{},
		IDSetSpawnPosition:            &SetSpawnPosition{},
		IDAnimate:                     &Animate{},
		IDRespawn:                     &Respawn{},
		IDContainerOpen:               &ContainerOpen{},
		IDContainerClose:              &ContainerClose{},
		IDPlayerHotBar:                &PlayerHotBar{},
		IDInventoryContent:            &InventoryContent{},
		IDInventorySlot:               &InventorySlot{},
		IDContainerSetData:            &ContainerSetData{},
		IDCraftingData:                &CraftingData{},
		IDCraftingEvent:               &CraftingEvent{},
		IDGUIDataPickItem:             &GUIDataPickItem{},
		IDAdventureSettings:           &AdventureSettings{},
		IDBlockActorData:              &BlockActorData{},
		IDPlayerInput:                 &PlayerInput{},
		IDLevelChunk:                  &LevelChunk{},
		IDSetCommandsEnabled:          &SetCommandsEnabled{},
		IDSetDifficulty:               &SetDifficulty{},
		IDChangeDimension:             &ChangeDimension{},
		IDSetPlayerGameType:           &SetPlayerGameType{},
		IDPlayerList:                  &PlayerList{},
		IDSimpleEvent:                 &SimpleEvent{},
		IDEvent:                       &Event{},
		IDSpawnExperienceOrb:          &SpawnExperienceOrb{},
		IDClientBoundMapItemData:      &ClientBoundMapItemData{},
		IDMapInfoRequest:              &MapInfoRequest{},
		IDRequestChunkRadius:          &RequestChunkRadius{},
		IDChunkRadiusUpdated:          &ChunkRadiusUpdated{},
		IDItemFrameDropItem:           &ItemFrameDropItem{},
		IDGameRulesChanged:            &GameRulesChanged{},
		IDCamera:                      &Camera{},
		IDBossEvent:                   &BossEvent{},
		IDShowCredits:                 &ShowCredits{},
		IDAvailableCommands:           &AvailableCommands{},
		IDCommandRequest:              &CommandRequest{},
		IDCommandBlockUpdate:          &CommandBlockUpdate{},
		IDCommandOutput:               &CommandOutput{},
		IDUpdateTrade:                 &UpdateTrade{},
		IDUpdateEquip:                 &UpdateEquip{},
		IDResourcePackDataInfo:        &ResourcePackDataInfo{},
		IDResourcePackChunkData:       &ResourcePackChunkData{},
		IDResourcePackChunkRequest:    &ResourcePackChunkRequest{},
		IDTransfer:                    &Transfer{},
		IDPlaySound:                   &PlaySound{},
		IDStopSound:                   &StopSound{},
		IDSetTitle:                    &SetTitle{},
		IDAddBehaviourTree:            &AddBehaviourTree{},
		IDStructureBlockUpdate:        &StructureBlockUpdate{},
		IDShowStoreOffer:              &ShowStoreOffer{},
		IDPurchaseReceipt:             &PurchaseReceipt{},
		IDPlayerSkin:                  &PlayerSkin{},
		IDSubClientLogin:              &SubClientLogin{},
		IDAutomationClientConnect:     &AutomationClientConnect{},
		IDSetLastHurtBy:               &SetLastHurtBy{},
		IDBookEdit:                    &BookEdit{},
		IDNPCRequest:                  &NPCRequest{},
		IDPhotoTransfer:               &PhotoTransfer{},
		IDModalFormRequest:            &ModalFormRequest{},
		IDModalFormResponse:           &ModalFormResponse{},
		IDServerSettingsRequest:       &ServerSettingsRequest{},
		IDServerSettingsResponse:      &ServerSettingsResponse{},
		IDShowProfile:                 &ShowProfile{},
		IDSetDefaultGameType:          &SetDefaultGameType{},
		IDRemoveObjective:             &RemoveObjective{},
		IDSetDisplayObjective:         &SetDisplayObjective{},
		IDSetScore:                    &SetScore{},
		IDLabTable:                    &LabTable{},
		IDUpdateBlockSynced:           &UpdateBlockSynced{},
		IDMoveActorDelta:              &MoveActorDelta{},
		IDSetScoreboardIdentity:       &SetScoreboardIdentity{},
		IDSetLocalPlayerAsInitialised: &SetLocalPlayerAsInitialised{},
		IDUpdateSoftEnum:              &UpdateSoftEnum{},
		IDNetworkStackLatency:         &NetworkStackLatency{},
		// ---
		IDScriptCustomEvent:           &ScriptCustomEvent{},
		IDSpawnParticleEffect:         &SpawnParticleEffect{},
		IDAvailableActorIdentifiers:   &AvailableActorIdentifiers{},
		IDNetworkChunkPublisherUpdate: &NetworkChunkPublisherUpdate{},
		IDBiomeDefinitionList:         &BiomeDefinitionList{},
		IDLevelSoundEvent:             &LevelSoundEvent{},
		IDLevelEventGeneric:           &LevelEventGeneric{},
		IDLecternUpdate:               &LecternUpdate{},
		// ---
		IDAddEntity:                         &AddEntity{},
		IDRemoveEntity:                      &RemoveEntity{},
		IDClientCacheStatus:                 &ClientCacheStatus{},
		IDOnScreenTextureAnimation:          &OnScreenTextureAnimation{},
		IDMapCreateLockedCopy:               &MapCreateLockedCopy{},
		IDStructureTemplateDataRequest:      &StructureTemplateDataExportResponse{},
		IDStructureTemplateDataResponse:     &StructureTemplateDataExportResponse{},
		IDUpdateBlockProperties:             &UpdateBlockProperties{},
		IDClientCacheBlobStatus:             &ClientCacheBlobStatus{},
		IDClientCacheMissResponse:           &ClientCacheMissResponse{},
		IDEducationSettings:                 &EducationSettings{},
		IDEmote:                             &Emote{},
		IDMultiPlayerSettings:               &MultiPlayerSettings{},
		IDSettingsCommand:                   &SettingsCommand{},
		IDAnvilDamage:                       &AnvilDamage{},
		IDCompletedUsingItem:                &CompletedUsingItem{},
		IDNetworkSettings:                   &NetworkSettings{},
		IDPlayerAuthInput:                   &PlayerAuthInput{},
		IDCreativeContent:                   &CreativeContent{},
		IDPlayerEnchantOptions:              &PlayerEnchantOptions{},
		IDItemStackRequest:                  &ItemStackRequest{},
		IDItemStackResponse:                 &ItemStackResponse{},
		IDPlayerArmourDamage:                &PlayerArmourDamage{},
		IDCodeBuilder:                       &CodeBuilder{},
		IDUpdatePlayerGameType:              &UpdatePlayerGameType{},
		IDEmoteList:                         &EmoteList{},
		IDPositionTrackingDBServerBroadcast: &PositionTrackingDBServerBroadcast{},
		IDPositionTrackingDBClientRequest:   &PositionTrackingDBClientRequest{},
		IDDebugInfo:                         &DebugInfo{},
		IDPacketViolationWarning:            &PacketViolationWarning{},
	}
	for id, pk := range registeredPackets {
		p[id] = pk()
	}
	return p
}

// PacketsByName is a map holding a function to create a new packet for each packet registered in Pool. These
// functions are indexed using the exact packet name they return.
var PacketsByName = map[string]func() Packet{}

func init() {
	for _, packet := range NewPool() {
		pk := packet
		PacketsByName[reflect.TypeOf(pk).Elem().Name()] = func() Packet {
			return reflect.New(reflect.TypeOf(pk).Elem()).Interface().(Packet)
		}
	}
}
