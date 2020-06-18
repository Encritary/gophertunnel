package protocol

import (
	"bytes"
	"encoding/binary"
)

// EnchantmentOption represents a single option in the enchantment table for a single item.
type EnchantmentOption struct {
	// Cost is the cost of the option. This is the amount of XP levels required to select this enchantment
	// option.
	Cost uint32
	// Enchantments holds the enchantments that will be applied to the item when this option is clicked.
	Enchantments ItemEnchantments
	// Name is a name that will be translated to the 'Standard Galactic Alphabet'
	// (https://minecraft.gamepedia.com/Enchanting_Table#Standard_Galactic_Alphabet) client-side. The names
	// generally have no meaning, such as:
	// 'animal imbue range galvanize '
	// 'bless inside creature shrink '
	// 'elder free of inside '
	Name string
	// RecipeNetworkID is a unique network ID for this enchantment option. When enchanting, the client
	// will submit this network ID in a ItemStackRequest packet with the CraftRecipe action, so that the
	// server knows which enchantment was selected.
	// Note that this ID should still be unique with other actual recipes. It's recommended to start counting
	// for enchantment network IDs from the counter used for producing network IDs for the normal recipes.
	RecipeNetworkID uint32
}

// WriteEnchantOption writes an EnchantmentOption x to Buffer dst.
func WriteEnchantOption(dst *bytes.Buffer, x EnchantmentOption) error {
	return chainErr(
		WriteVaruint32(dst, x.Cost),
		WriteItemEnchants(dst, x.Enchantments),
		WriteString(dst, x.Name),
		WriteVaruint32(dst, x.RecipeNetworkID),
	)
}

// EnchantOption reads an EnchantmentOption x from Buffer src.
func EnchantOption(src *bytes.Buffer, x *EnchantmentOption) error {
	return chainErr(
		Varuint32(src, &x.Cost),
		ItemEnchants(src, &x.Enchantments),
		String(src, &x.Name),
		Varuint32(src, &x.RecipeNetworkID),
	)
}

const (
	EnchantmentSlotNone           = 0
	EnchantmentSlotAll            = 0xffff
	EnchantmentSlotArmour         = EnchantmentSlotHelmet | EnchantmentSlotChestplate | EnchantmentSlotLeggings | EnchantmentSlotBoots
	EnchantmentSlotHelmet         = 0x1
	EnchantmentSlotChestplate     = 0x2
	EnchantmentSlotLeggings       = 0x4
	EnchantmentSlotBoots          = 0x8
	EnchantmentSlotSword          = 0x10
	EnchantmentSlotBow            = 0x20
	EnchantmentSlotToolOther      = EnchantmentSlotHoe | EnchantmentSlotShears | EnchantmentSlotFlintAndSteel
	EnchantmentSlotHoe            = 0x40
	EnchantmentSlotShears         = 0x80
	EnchantmentSlotFlintAndSteel  = 0x100
	EnchantmentSlotDig            = EnchantmentSlotAxe | EnchantmentSlotPickaxe | EnchantmentSlotShovel
	EnchantmentSlotAxe            = 0x200
	EnchantmentSlotPickaxe        = 0x400
	EnchantmentSlotShovel         = 0x800
	EnchantmentSlotFishingRod     = 0x1000
	EnchantmentSlotCarrotOnAStick = 0x2000
	EnchantmentSlotElytra         = 0x4000
	EnchantmentSlotTrident        = 0x8000
)

// ItemEnchantments holds information on the enchantments that are applied to an item when a specific button
// is clicked in the enchantment table.
type ItemEnchantments struct {
	// Slot is the enchantment slot of the item that was put into the enchantment table, for which the
	// following enchantments will apply.
	// The possible slots can be found above.
	Slot int32
	// Enchantments is an array of 3 slices of enchantment instances. Each array represents enchantments that
	// will be added to the item with a different activation type. The arrays in which enchantments are sent
	// by the vanilla server are as follows:
	// slice 1 { protection, fire protection, feather falling, blast protection, projectile protection,
	//           thorns, respiration, depth strider, aqua affinity, frost walker, soul speed }
	// slice 2 { sharpness, smite, bane of arthropods, fire aspect, looting, silk touch, unbreaking, fortune,
	//           flame, luck of the sea, impaling }
	// slice 3 { knockback, efficiency, power, punch, infinity, lure, mending, curse of binding,
	//           curse of vanishing, riptide, loyalty, channeling, multishot, piercing, quick charge }
	// The first slice holds armour enchantments, the differences between the slice 2 and slice 3 are more
	// vaguely defined.
	Enchantments [3][]EnchantmentInstance
}

// WriteItemEnchants writes an ItemEnchantments x to Buffer dst.
func WriteItemEnchants(dst *bytes.Buffer, x ItemEnchantments) error {
	if err := binary.Write(dst, binary.LittleEndian, x.Slot); err != nil {
		return err
	}
	for _, enchantments := range x.Enchantments {
		if err := WriteVaruint32(dst, uint32(len(enchantments))); err != nil {
			return err
		}
		for _, enchantment := range enchantments {
			if err := WriteEnchant(dst, enchantment); err != nil {
				return err
			}
		}
	}
	return nil
}

// ItemEnchants reads an ItemEnchantments x from Buffer src.
func ItemEnchants(src *bytes.Buffer, x *ItemEnchantments) error {
	if err := binary.Read(src, binary.LittleEndian, &x.Slot); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		var l uint32
		if err := Varuint32(src, &l); err != nil {
			return err
		}
		x.Enchantments[i] = make([]EnchantmentInstance, l)
		for j := uint32(0); j < l; j++ {
			if err := Enchant(src, &x.Enchantments[i][j]); err != nil {
				return err
			}
		}
	}
	return nil
}

// EnchantmentInstance represents a single enchantment instance with the type of the enchantment and its
// level.
type EnchantmentInstance struct {
	Type  byte
	Level byte
}

// WriteEnchant writes an EnchantmentInstance x to Buffer dst.
func WriteEnchant(dst *bytes.Buffer, x EnchantmentInstance) error {
	dst.WriteByte(x.Type)
	dst.WriteByte(x.Level)
	return nil
}

// Enchant reads an EnchantmentInstance x from Buffer src.
func Enchant(src *bytes.Buffer, x *EnchantmentInstance) error {
	return chainErr(
		binary.Read(src, binary.LittleEndian, &x.Type),
		binary.Read(src, binary.LittleEndian, &x.Level),
	)
}
