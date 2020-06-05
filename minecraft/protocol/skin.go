package protocol

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

// Skin represents the skin of a player as sent over network. The skin holds a texture and a model, and
// optional animations which may be present when the skin is created using persona or bought from the
// marketplace.
type Skin struct {
	// SkinID is a unique ID produced for the skin, for example 'c18e65aa-7b21-4637-9b63-8ad63622ef01_Alex'
	// for the default Alex skin.
	SkinID string
	// SkinResourcePatch is a JSON encoded object holding some fields that point to the geometry that the
	// skin has.
	// The JSON object that this holds specifies the way that the geometry of animations and the default skin
	// of the player are combined.
	SkinResourcePatch []byte
	// SkinImageWidth and SkinImageHeight hold the dimensions of the skin image. Note that these are not the
	// dimensions in bytes, but in pixels.
	SkinImageWidth, SkinImageHeight uint32
	// SkinData is a byte slice of SkinImageWidth * SkinImageHeight bytes. It is an RGBA ordered byte
	// representation of the skin pixels.
	SkinData []byte
	// Animations is a list of all animations that the skin has.
	Animations []SkinAnimation
	// CapeImageWidth and CapeImageHeight hold the dimensions of the cape image. Note that these are not the
	// dimensions in bytes, but in pixels.
	CapeImageWidth, CapeImageHeight uint32
	// CapeData is a byte slice of 64*32*4 bytes. It is a RGBA ordered byte representation of the cape
	// colours, much like the SkinData.
	CapeData []byte
	// SkinGeometry is a JSON encoded structure of the geometry data of a skin, containing properties
	// such as bones, uv, pivot etc.
	SkinGeometry []byte
	// TODO: Find out what value AnimationData holds and when it does hold something.
	AnimationData []byte
	// PremiumSkin specifies if this is a skin that was purchased from the marketplace.
	PremiumSkin bool
	// PersonaSkin specifies if this is a skin that was created using the in-game skin creator.
	PersonaSkin bool
	// PersonaCapeOnClassicSkin specifies if the skin had a Persona cape (in-game skin creator cape) equipped
	// on a classic skin.
	PersonaCapeOnClassicSkin bool
	// CapeID is a unique identifier that identifies the cape. It usually holds a UUID in it.
	CapeID string
	// FullSkinID is an ID that represents the skin in full. The actual functionality is unknown: The client
	// does not seem to send a value for this.
	FullSkinID string
	// SkinColour is a hex representation (including #) of the base colour of the skin. An example of the
	// colour sent here is '#b37b62'.
	SkinColour string
	// ArmSize is the size of the arms of the player's model. This is either 'wide' (generally for male skins)
	// or 'slim' (generally for female skins).
	ArmSize string
	// PersonaPieces is a list of all persona pieces that the skin is composed of.
	PersonaPieces []PersonaPiece
	// PieceTintColours is a list of specific tint colours for (some of) the persona pieces found in the list
	// above.
	PieceTintColours []PersonaPieceTintColour
	// Trusted specifies if the skin is 'trusted'. No code should rely on this field, as any proxy or client
	// can easily change it.
	Trusted bool
}

// WriteSerialisedSkin writes a Skin x to Buffer dst. WriteSerialisedSkin panics if the fields of the skin
// have invalid values, usually indicating that the dimensions of the skin images are incorrect.
func WriteSerialisedSkin(dst *bytes.Buffer, x Skin) error {
	if err := x.validate(); err != nil {
		panic(err)
	}
	if err := chainErr(
		WriteString(dst, x.SkinID),
		WriteByteSlice(dst, x.SkinResourcePatch),
		binary.Write(dst, binary.LittleEndian, x.SkinImageWidth),
		binary.Write(dst, binary.LittleEndian, x.SkinImageHeight),
		WriteByteSlice(dst, x.SkinData),
		binary.Write(dst, binary.LittleEndian, uint32(len(x.Animations))),
	); err != nil {
		return err
	}
	for _, anim := range x.Animations {
		if err := WriteAnimation(dst, anim); err != nil {
			return err
		}
	}
	if err := chainErr(
		binary.Write(dst, binary.LittleEndian, x.CapeImageWidth),
		binary.Write(dst, binary.LittleEndian, x.CapeImageHeight),
		WriteByteSlice(dst, x.CapeData),
		WriteByteSlice(dst, x.SkinGeometry),
		WriteByteSlice(dst, x.AnimationData),
		binary.Write(dst, binary.LittleEndian, x.PremiumSkin),
		binary.Write(dst, binary.LittleEndian, x.PersonaSkin),
		binary.Write(dst, binary.LittleEndian, x.PersonaCapeOnClassicSkin),
		WriteString(dst, x.CapeID),
		WriteString(dst, x.FullSkinID),
		WriteString(dst, x.ArmSize),
		WriteString(dst, x.SkinColour),
		binary.Write(dst, binary.LittleEndian, uint32(len(x.PersonaPieces))),
	); err != nil {
		return err
	}
	for _, piece := range x.PersonaPieces {
		if err := WriteSkinPiece(dst, piece); err != nil {
			return err
		}
	}
	if err := binary.Write(dst, binary.LittleEndian, uint32(len(x.PieceTintColours))); err != nil {
		return err
	}
	for _, tint := range x.PieceTintColours {
		if err := WriteSkinPieceTint(dst, tint); err != nil {
			return err
		}
	}
	return nil
}

// SerialisedSkin reads a Skin x from Buffer src.
func SerialisedSkin(src *bytes.Buffer, x *Skin) error {
	var animationCount uint32
	var c uint32
	if err := chainErr(
		String(src, &x.SkinID),
		ByteSlice(src, &x.SkinResourcePatch),
		binary.Read(src, binary.LittleEndian, &x.SkinImageWidth),
		binary.Read(src, binary.LittleEndian, &x.SkinImageHeight),
		ByteSlice(src, &x.SkinData),
		binary.Read(src, binary.LittleEndian, &animationCount),
	); err != nil {
		return err
	}
	x.Animations = make([]SkinAnimation, animationCount)

	for i := uint32(0); i < animationCount; i++ {
		if err := Animation(src, &x.Animations[i]); err != nil {
			return err
		}
	}
	if err := chainErr(
		binary.Read(src, binary.LittleEndian, &x.CapeImageWidth),
		binary.Read(src, binary.LittleEndian, &x.CapeImageHeight),
		ByteSlice(src, &x.CapeData),
		ByteSlice(src, &x.SkinGeometry),
		ByteSlice(src, &x.AnimationData),
		binary.Read(src, binary.LittleEndian, &x.PremiumSkin),
		binary.Read(src, binary.LittleEndian, &x.PersonaSkin),
		binary.Read(src, binary.LittleEndian, &x.PersonaCapeOnClassicSkin),
		String(src, &x.CapeID),
		String(src, &x.FullSkinID),
		String(src, &x.ArmSize),
		String(src, &x.SkinColour),
		binary.Read(src, binary.LittleEndian, &c),
	); err != nil {
		return err
	}
	x.PersonaPieces = make([]PersonaPiece, c)
	for i := uint32(0); i < c; i++ {
		if err := SkinPiece(src, &x.PersonaPieces[i]); err != nil {
			return err
		}
	}
	if err := binary.Read(src, binary.LittleEndian, &c); err != nil {
		return err
	}
	x.PieceTintColours = make([]PersonaPieceTintColour, c)
	for i := uint32(0); i < c; i++ {
		if err := SkinPieceTint(src, &x.PieceTintColours[i]); err != nil {
			return err
		}
	}
	return x.validate()
}

// validate checks the skin and makes sure every one of its values are correct. It checks the image dimensions
// and makes sure they match the image size of the skin, cape and the skin's animations.
func (skin Skin) validate() error {
	if skin.SkinImageHeight*skin.SkinImageWidth*4 != uint32(len(skin.SkinData)) {
		return fmt.Errorf("expected size of skin is %vx%v (%v bytes total), but got %v bytes", skin.SkinImageWidth, skin.SkinImageHeight, skin.SkinImageHeight*skin.SkinImageWidth*4, len(skin.SkinData))
	}
	if skin.CapeImageHeight*skin.CapeImageWidth*4 != uint32(len(skin.CapeData)) {
		return fmt.Errorf("expected size of cape is %vx%v (%v bytes total), but got %v bytes", skin.CapeImageWidth, skin.CapeImageHeight, skin.CapeImageHeight*skin.CapeImageWidth*4, len(skin.CapeData))
	}
	for i, animation := range skin.Animations {
		if animation.ImageHeight*animation.ImageWidth*4 != uint32(len(animation.ImageData)) {
			return fmt.Errorf("expected size of animation %v is %vx%v (%v bytes total), but got %v bytes", i, animation.ImageWidth, animation.ImageHeight, animation.ImageHeight*animation.ImageWidth*4, len(animation.ImageData))
		}
	}
	return nil
}

const (
	SkinAnimationHead = iota + 1
	SkinAnimationBody32x32
	SkinAnimationBody128x128
)

// SkinAnimation represents an animation that may be added to a skin. The client plays the animation itself,
// without the server having to do so.
// The rate at which these animations play appears to be decided by the client.
type SkinAnimation struct {
	// ImageWidth and ImageHeight hold the dimensions of the animation image. Note that these are not the
	// dimensions in bytes, but in pixels.
	ImageWidth, ImageHeight uint32
	// ImageData is a byte slice of ImageWidth * ImageHeight bytes. It is an RGBA ordered byte representation
	// of the animation image pixels. The ImageData contains FrameCount images in it, which each represent one
	// stage of the animation. The actual part of the skin that this field holds depends on the AnimationType,
	// where SkinAnimationHead holds only the head and its hat, whereas the other animations hold the entire
	// body of the skin.
	ImageData []byte
	// AnimationType is the type of the animation, which is one of the types found above. The data that
	// ImageData contains depends on this type.
	AnimationType uint32
	// FrameCount is the amount of frames that the skin animation holds. The number of frames here is the
	// amount of images that may be found in the ImageData field.
	FrameCount float32
}

// WriteAnimation writes a SkinAnimation x to Buffer dst.
func WriteAnimation(dst *bytes.Buffer, x SkinAnimation) error {
	return chainErr(
		binary.Write(dst, binary.LittleEndian, x.ImageWidth),
		binary.Write(dst, binary.LittleEndian, x.ImageHeight),
		WriteByteSlice(dst, x.ImageData),
		binary.Write(dst, binary.LittleEndian, x.AnimationType),
		WriteFloat32(dst, x.FrameCount),
	)
}

// Animation reads a SkinAnimation x from Buffer src.
func Animation(src *bytes.Buffer, x *SkinAnimation) error {
	return chainErr(
		binary.Read(src, binary.LittleEndian, &x.ImageWidth),
		binary.Read(src, binary.LittleEndian, &x.ImageHeight),
		ByteSlice(src, &x.ImageData),
		binary.Read(src, binary.LittleEndian, &x.AnimationType),
		Float32(src, &x.FrameCount),
	)
}

// PersonaPiece represents a piece of a persona skin. All pieces are sent separately.
type PersonaPiece struct {
	// PieceId is a UUID that identifies the piece itself, which is unique for each separate piece.
	PieceID string
	// PieceType holds the type of the piece. Several types I was able to find immediately are listed below.
	// - persona_skeleton
	// - persona_body
	// - persona_skin
	// - persona_bottom
	// - persona_feet
	// - persona_top
	// - persona_mouth
	// - persona_hair
	// - persona_eyes
	// - persona_facial_hair
	PieceType string
	// PackID is a UUID that identifies the pack that the persona piece belongs to.
	PackID string
	// Default specifies if the piece is one of the default pieces. This is true when the piece is one of
	// those that a Steve or Alex skin have.
	Default bool
	// ProductID is a UUID that identifies the piece when it comes to purchases. It is empty for pieces that
	// have the 'Default' field set to true.
	ProductID string
}

// WriteSkinPiece writes a PersonaPiece x to Buffer dst.
func WriteSkinPiece(dst *bytes.Buffer, x PersonaPiece) error {
	return chainErr(
		WriteString(dst, x.PieceID),
		WriteString(dst, x.PieceType),
		WriteString(dst, x.PackID),
		binary.Write(dst, binary.LittleEndian, x.Default),
		WriteString(dst, x.ProductID),
	)
}

// SkinPiece reads a PersonaPiece x from Buffer src.
func SkinPiece(src *bytes.Buffer, x *PersonaPiece) error {
	return chainErr(
		String(src, &x.PieceID),
		String(src, &x.PieceType),
		String(src, &x.PackID),
		binary.Read(src, binary.LittleEndian, &x.Default),
		String(src, &x.ProductID),
	)
}

// PersonaPieceTintColour describes the tint colours of a specific piece of a persona skin.
type PersonaPieceTintColour struct {
	// PieceType is the type of the persona skin piece that this tint colour concerns. The piece type must
	// always be present in the persona pieces list, but not each piece type has a tint colour sent.
	// Pieces that do have a tint colour that I was able to find immediately are listed below.
	// - persona_mouth
	// - persona_eyes
	// - persona_hair
	PieceType string
	// Colours is a list four colours written in hex notation (note, that unlike the SkinColour field in
	// the ClientData struct, this is actually ARGB, not just RGB).
	// The colours refer to different parts of the skin piece. The 'persona_eyes' may have the following
	// colours: ["#ffa12722","#ff2f1f0f","#ff3aafd9","#0"]
	// The first hex colour represents the tint colour of the iris, the second hex colour represents the
	// eyebrows and the third represents the sclera. The fourth is #0 because there are only 3 parts of the
	// persona_eyes skin piece.
	Colours []string
}

// WriteSkinPieceTint writes a PersonaPieceTintColour x to Buffer dst.
func WriteSkinPieceTint(dst *bytes.Buffer, x PersonaPieceTintColour) error {
	if err := chainErr(
		WriteString(dst, x.PieceType),
		binary.Write(dst, binary.LittleEndian, uint32(len(x.Colours))),
	); err != nil {
		return err
	}
	for _, c := range x.Colours {
		if err := WriteString(dst, c); err != nil {
			return err
		}
	}
	return nil
}

// SkinPieceTint reads a PersonaPieceTintColour x from Buffer src.
func SkinPieceTint(src *bytes.Buffer, x *PersonaPieceTintColour) error {
	var c uint32
	if err := chainErr(
		String(src, &x.PieceType),
		binary.Read(src, binary.LittleEndian, &c),
	); err != nil {
		return err
	}
	x.Colours = make([]string, c)
	for i := uint32(0); i < c; i++ {
		if err := String(src, &x.Colours[i]); err != nil {
			return err
		}
	}
	return nil
}
