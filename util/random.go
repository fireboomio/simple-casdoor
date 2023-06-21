package util

import "github.com/thanhpk/randstr"

func GenerateClientId() string {
	return randstr.Hex(10)
}

func GenerateClientSecret() string {
	return randstr.Hex(20)
}
