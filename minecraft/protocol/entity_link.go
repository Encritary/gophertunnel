package protocol

import (
	"bytes"
	"encoding/binary"
)

const (
	// EntityLinkRemove is set to remove the link between two entities.
	EntityLinkRemove = iota
	// EntityLinkRider is set for entities that have control over the entity they're riding, such as in a
	// minecart.
	EntityLinkRider
	// EntityLinkPassenger is set for entities being a passenger of a vehicle they enter, such as the back
	// sit of a boat.
	EntityLinkPassenger
)

// EntityLink is a link between two entities, typically being one entity riding another.
type EntityLink struct {
	// RiddenEntityUniqueID is the entity unique ID of the entity that is being ridden. For a player sitting
	// in a boat, this is the unique ID of the boat.
	RiddenEntityUniqueID int64
	// RiderEntityUniqueID is the entity unique ID of the entity that is riding. For a player sitting in a
	// boat, this is the unique ID of the player.
	RiderEntityUniqueID int64
	// Type is one of the types above. It specifies the way the entity is linked to another entity.
	Type byte
	// Immediate is set to immediately dismount an entity from another. This should be set when the mount of
	// an entity is killed.
	Immediate bool
	// RiderInitiated specifies if the link was created by the rider, for example the player starting to ride
	// a horse by itself. This is generally true in vanilla environment for players.
	RiderInitiated bool
}

// EntityLinkAction reads a single entity link (action) from buffer src.
func EntityLinkAction(src *bytes.Buffer, x *EntityLink) error {
	return chainErr(
		Varint64(src, &x.RiddenEntityUniqueID),
		Varint64(src, &x.RiderEntityUniqueID),
		binary.Read(src, binary.LittleEndian, &x.Type),
		binary.Read(src, binary.LittleEndian, &x.Immediate),
		binary.Read(src, binary.LittleEndian, &x.RiderInitiated),
	)
}

// EntityLinks reads a list of entity links from buffer src that are currently active.
func EntityLinks(src *bytes.Buffer, x *[]EntityLink) error {
	var count uint32
	if err := Varuint32(src, &count); err != nil {
		return wrap(err)
	}
	if count > lowerLimit {
		return LimitHitError{Limit: lowerLimit, Type: "entity link"}
	}
	*x = make([]EntityLink, count)
	for i := uint32(0); i < count; i++ {
		if err := EntityLinkAction(src, &(*x)[i]); err != nil {
			return wrap(err)
		}
	}
	return nil
}

// WriteEntityLinkAction writes a single entity link x to buffer dst.
func WriteEntityLinkAction(dst *bytes.Buffer, x EntityLink) error {
	return chainErr(
		WriteVarint64(dst, x.RiddenEntityUniqueID),
		WriteVarint64(dst, x.RiderEntityUniqueID),
		binary.Write(dst, binary.LittleEndian, x.Type),
		binary.Write(dst, binary.LittleEndian, x.Immediate),
		binary.Write(dst, binary.LittleEndian, x.RiderInitiated),
	)
}

// WriteEntityLinks writes a list of entity links currently active to buffer dst.
func WriteEntityLinks(dst *bytes.Buffer, x []EntityLink) error {
	if err := WriteVaruint32(dst, uint32(len(x))); err != nil {
		return wrap(err)
	}
	for _, link := range x {
		if err := WriteEntityLinkAction(dst, link); err != nil {
			return wrap(err)
		}
	}
	return nil
}
