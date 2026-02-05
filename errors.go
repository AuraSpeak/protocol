package protocol

import "errors"

var ErrInvalidMagic = errors.New("invalid magic")
var ErrInvalidVersion = errors.New("invalid version")
var ErrInvalidPacketType = errors.New("invalid packet type")
var ErrDataTooShort = errors.New("data too short")
