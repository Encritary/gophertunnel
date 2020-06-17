package packet

import (
	"bytes"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
)

// CreativeContent is a packet sent by the server to set the creative inventory's content for a player.
// Introduced in 1.16, this packet replaces the previous method - sending an InventoryContent packet with
// creative inventory window ID.
type CreativeContent struct {
	// Items is a list of the items that should be added to the creative inventory.
	Items []protocol.CreativeItem
}

// ID ...
func (*CreativeContent) ID() uint32 {
	return IDCreativeContent
}

// Marshal ...
func (pk *CreativeContent) Marshal(buf *bytes.Buffer) {
	_ = protocol.WriteVaruint32(buf, uint32(len(pk.Items)))
	for _, item := range pk.Items {
		_ = protocol.WriteCreativeEntry(buf, item)
	}
}

// Unmarshal ...
func (pk *CreativeContent) Unmarshal(buf *bytes.Buffer) error {
	var count uint32
	if err := protocol.Varuint32(buf, &count); err != nil {
		return err
	}
	pk.Items = make([]protocol.CreativeItem, count)
	for i := 0; i < int(count); i++ {
		if err := protocol.CreativeEntry(buf, &pk.Items[i]); err != nil {
			return err
		}
	}
	return nil
}
