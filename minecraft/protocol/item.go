package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/sandertv/gophertunnel/minecraft/nbt"
)

// ItemStack represents an item instance/stack over network. It has a network ID and a metadata value that
// define its type.
type ItemStack struct {
	ItemType
	// Count is the count of items that the item stack holds.
	Count int16
	// NBTData is a map that is serialised to its NBT representation when sent in a packet.
	NBTData map[string]interface{}
	// CanBePlacedOn is a list of block identifiers like 'minecraft:stone' which the item, if it is an item
	// that can be placed, can be placed on top of.
	CanBePlacedOn []string
	// CanBreak is a list of block identifiers like 'minecraft:dirt' that the item is able to break.
	CanBreak []string
}

// ItemType represents a consistent combination of network ID and metadata value of an item. It cannot usually
// be changed unless a new item is obtained.
type ItemType struct {
	// NetworkID is the numerical network ID of the item. This is sometimes a positive ID, and sometimes a
	// negative ID, depending on what item it concerns.
	NetworkID int32
	// MetadataValue is the metadata value of the item. For some items, this is the damage value, whereas for
	// other items it is simply an identifier of a variant of the item.
	MetadataValue int16
}

// Item reads an item stack from buffer src and stores it into item stack x.
func Item(src *bytes.Buffer, x *ItemStack) error {
	x.NBTData = make(map[string]interface{})
	if err := Varint32(src, &x.NetworkID); err != nil {
		return wrap(err)
	}
	if x.NetworkID == 0 {
		// The item was air, so there is no more data we should read for the item instance. After all, air
		// items aren't really anything.
		x.MetadataValue = 0
		x.Count = 0
		x.CanBePlacedOn, x.CanBreak = nil, nil
		return nil
	}
	var auxValue int32
	if err := Varint32(src, &auxValue); err != nil {
		return wrap(err)
	}
	x.MetadataValue = int16(auxValue >> 8)
	x.Count = int16(auxValue & 0xff)

	var legacyNBTLength int16
	if err := binary.Read(src, binary.LittleEndian, &legacyNBTLength); err != nil {
		return wrap(err)
	}
	if legacyNBTLength != 0 {
		if legacyNBTLength == -1 {
			var nbtCount byte
			if err := binary.Read(src, binary.LittleEndian, &nbtCount); err != nil {
				return wrap(err)
			}
			if nbtCount != 1 {
				// The NBT count seems to be always 1, so we return an error if it is not, just so we know there can
				// be more than one.
				return fmt.Errorf("%v: expected NBT count to be 1, got %v", callFrame(), nbtCount)
			}
			decoder := nbt.NewDecoder(src)
			for i := byte(0); i < nbtCount; i++ {
				if err := decoder.Decode(&x.NBTData); err != nil {
					return fmt.Errorf("%v: error decoding item NBT: %v", callFrame(), err)
				}
			}
		} else {
			if legacyNBTLength < 0 {
				return fmt.Errorf("%v: invalid NBT count %v", callFrame(), legacyNBTLength)
			}
			nbtData := src.Next(int(legacyNBTLength))
			if err := nbt.UnmarshalEncoding(nbtData, &x.NBTData, nbt.LittleEndian); err != nil {
				return fmt.Errorf("%v: error decoding item NBT: %v", callFrame(), err)
			}
		}
	}

	var count int32
	if err := Varint32(src, &count); err != nil {
		return wrap(err)
	}
	if count < 0 {
		return NegativeCountError{Type: "item can be placed on"}
	}
	if count > higherLimit {
		return LimitHitError{Limit: higherLimit, Type: "item can be placed on"}
	}
	x.CanBePlacedOn = make([]string, count)
	for i := int32(0); i < count; i++ {
		if err := String(src, &x.CanBePlacedOn[i]); err != nil {
			return wrap(err)
		}
	}
	if err := Varint32(src, &count); err != nil {
		return wrap(err)
	}
	if count < 0 {
		return NegativeCountError{Type: "item can break"}
	}
	if count > higherLimit {
		return LimitHitError{Limit: higherLimit, Type: "item can break"}
	}
	x.CanBreak = make([]string, count)
	for i := int32(0); i < count; i++ {
		if err := String(src, &x.CanBreak[i]); err != nil {
			return wrap(err)
		}
	}
	const shieldID = 513
	if x.NetworkID == shieldID {
		var blockingTick int64
		if err := Varint64(src, &blockingTick); err != nil {
			return wrap(err)
		}
	}
	return nil
}

// WriteItem writes an item stack x to buffer dst.
func WriteItem(dst *bytes.Buffer, x ItemStack) error {
	if err := WriteVarint32(dst, x.NetworkID); err != nil {
		return wrap(err)
	}
	if x.NetworkID == 0 {
		// The item was air, so there's no more data to follow. Return immediately.
		return nil
	}
	if err := WriteVarint32(dst, int32(x.MetadataValue<<8)|int32(x.Count)); err != nil {
		return wrap(err)
	}
	if len(x.NBTData) != 0 {
		// Write a fixed -1, which used to be the NBT length.
		if err := binary.Write(dst, binary.LittleEndian, int16(-1)); err != nil {
			return wrap(err)
		}
		// NBT Count, which is always one in our case.
		if err := binary.Write(dst, binary.LittleEndian, byte(1)); err != nil {
			return wrap(err)
		}
		b, err := nbt.Marshal(x.NBTData)
		if err != nil {
			return fmt.Errorf("%v: error writing NBT: %v", callFrame(), err)
		}
		_, _ = dst.Write(b)
	} else {
		// There's no NBT to write, so we just write 0 and continue with the rest.
		if err := binary.Write(dst, binary.LittleEndian, int16(0)); err != nil {
			return wrap(err)
		}
	}
	if err := WriteVarint32(dst, int32(len(x.CanBePlacedOn))); err != nil {
		return wrap(err)
	}
	for _, block := range x.CanBePlacedOn {
		if err := WriteString(dst, block); err != nil {
			return wrap(err)
		}
	}
	if err := WriteVarint32(dst, int32(len(x.CanBreak))); err != nil {
		return wrap(err)
	}
	for _, block := range x.CanBreak {
		if err := WriteString(dst, block); err != nil {
			return wrap(err)
		}
	}
	const shieldID = 513
	if x.NetworkID == shieldID {
		var blockingTick int64
		if err := WriteVarint64(dst, blockingTick); err != nil {
			return wrap(err)
		}
	}
	return nil
}

// RecipeIngredient reads an ItemStack x as a recipe ingredient from Buffer src.
func RecipeIngredient(src *bytes.Buffer, x *ItemStack) error {
	if err := Varint32(src, &x.NetworkID); err != nil {
		return nil
	}
	if x.NetworkID == 0 {
		return nil
	}
	var meta, count int32
	if err := Varint32(src, &meta); err != nil {
		return err
	}
	x.MetadataValue = int16(meta)
	if err := Varint32(src, &count); err != nil {
		return err
	}
	if count < 0 {
		return NegativeCountError{Type: "recipe ingredient count"}
	}
	x.Count = int16(count)
	return nil
}

// WriteRecipeIngredient writes an ItemStack x as a recipe ingredient to Buffer dst.
func WriteRecipeIngredient(dst *bytes.Buffer, x ItemStack) error {
	if err := WriteVarint32(dst, x.NetworkID); err != nil {
		return err
	}
	if x.NetworkID == 0 {
		return nil
	}
	return chainErr(
		WriteVarint32(dst, int32(x.MetadataValue)),
		WriteVarint32(dst, int32(x.Count)),
	)
}
