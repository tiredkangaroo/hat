package database

import "unsafe"

func b2s(b []byte) string {
	return unsafe.String(&b[0], len(b))
}
