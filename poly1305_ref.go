// Copyright (c) 2016 Andreas Auernhammer. All rights reserved.
// Use of this source code is governed by a license that can be
// found in the LICENSE file.

// +build !amd64

package poly1305

func initialize(r *[5]uint32, pad *[4]uint32, key *[32]byte) {
	r[0] = (uint32(key[0]) | uint32(key[1])<<8 | uint32(key[2])<<16 | uint32(key[3])<<24) & 0x3ffffff
	r[1] = ((uint32(key[3]) | uint32(key[4])<<8 | uint32(key[5])<<16 | uint32(key[6])<<24) >> 2) & 0x3ffff03
	r[2] = ((uint32(key[6]) | uint32(key[7])<<8 | uint32(key[8])<<16 | uint32(key[9])<<24) >> 4) & 0x3ffc0ff
	r[3] = ((uint32(key[9]) | uint32(key[10])<<8 | uint32(key[11])<<16 | uint32(key[12])<<24) >> 6) & 0x3f03fff
	r[4] = ((uint32(key[12]) | uint32(key[13])<<8 | uint32(key[14])<<16 | uint32(key[15])<<24) >> 8) & 0x00fffff

	pad[0] = (uint32(key[16]) | uint32(key[17])<<8 | uint32(key[18])<<16 | uint32(key[19])<<24)
	pad[1] = (uint32(key[20]) | uint32(key[21])<<8 | uint32(key[22])<<16 | uint32(key[23])<<24)
	pad[2] = (uint32(key[24]) | uint32(key[25])<<8 | uint32(key[26])<<16 | uint32(key[27])<<24)
	pad[3] = (uint32(key[28]) | uint32(key[29])<<8 | uint32(key[30])<<16 | uint32(key[31])<<24)
}

func core(msg []byte, flag uint32, h, r *[5]uint32) {
	h0, h1, h2, h3, h4 := h[0], h[1], h[2], h[3], h[4]
	r0, r1, r2, r3, r4 := uint64(r[0]), uint64(r[1]), uint64(r[2]), uint64(r[3]), uint64(r[4])
	s1, s2, s3, s4 := uint64(r[1]*5), uint64(r[2]*5), uint64(r[3]*5), uint64(r[4]*5)

	var d0, d1, d2, d3, d4 uint64
	for i := 0; i < len(msg); i += TagSize {
		// h += m
		h0 += (uint32(msg[i+0]) | uint32(msg[i+1])<<8 | uint32(msg[i+2])<<16 | uint32(msg[i+3])<<24) & 0x3ffffff
		h1 += ((uint32(msg[i+3]) | uint32(msg[i+4])<<8 | uint32(msg[i+5])<<16 | uint32(msg[i+6])<<24) >> 2) & 0x3ffffff
		h2 += ((uint32(msg[i+6]) | uint32(msg[i+7])<<8 | uint32(msg[i+8])<<16 | uint32(msg[i+9])<<24) >> 4) & 0x3ffffff
		h3 += ((uint32(msg[i+9]) | uint32(msg[i+10])<<8 | uint32(msg[i+11])<<16 | uint32(msg[i+12])<<24) >> 6) & 0x3ffffff
		h4 += ((uint32(msg[i+12]) | uint32(msg[i+13])<<8 | uint32(msg[i+14])<<16 | uint32(msg[i+15])<<24) >> 8) | flag

		// h *= r
		d0 = (uint64(h0) * r0) + (uint64(h1) * s4) + (uint64(h2) * s3) + (uint64(h3) * s2) + (uint64(h4) * s1)
		d1 = (d0 >> 26) + (uint64(h0) * r1) + (uint64(h1) * r0) + (uint64(h2) * s4) + (uint64(h3) * s3) + (uint64(h4) * s2)
		d2 = (d1 >> 26) + (uint64(h0) * r2) + (uint64(h1) * r1) + (uint64(h2) * r0) + (uint64(h3) * s4) + (uint64(h4) * s3)
		d3 = (d2 >> 26) + (uint64(h0) * r3) + (uint64(h1) * r2) + (uint64(h2) * r1) + (uint64(h3) * r0) + (uint64(h4) * s4)
		d4 = (d3 >> 26) + (uint64(h0) * r4) + (uint64(h1) * r3) + (uint64(h2) * r2) + (uint64(h3) * r1) + (uint64(h4) * r0)

		// h %= p
		h0 = uint32(d0) & 0x3ffffff
		h1 = uint32(d1) & 0x3ffffff
		h2 = uint32(d2) & 0x3ffffff
		h3 = uint32(d3) & 0x3ffffff
		h4 = uint32(d4) & 0x3ffffff

		h0 += uint32(d4>>26) * 5
		h1 += h0 >> 26
		h0 = h0 & 0x3ffffff
	}
	h[0], h[1], h[2], h[3], h[4] = h0, h1, h2, h3, h4
}

func extractHash(tag *[16]byte, h0, h1, h2, h3 uint32) {
	tag[0] = byte(h0)
	tag[1] = byte(h0 >> 8)
	tag[2] = byte(h0 >> 16)
	tag[3] = byte(h0 >> 24)
	tag[4] = byte(h1)
	tag[5] = byte(h1 >> 8)
	tag[6] = byte(h1 >> 16)
	tag[7] = byte(h1 >> 24)
	tag[8] = byte(h2)
	tag[9] = byte(h2 >> 8)
	tag[10] = byte(h2 >> 16)
	tag[11] = byte(h2 >> 24)
	tag[12] = byte(h3)
	tag[13] = byte(h3 >> 8)
	tag[14] = byte(h3 >> 16)
	tag[15] = byte(h3 >> 24)
}
