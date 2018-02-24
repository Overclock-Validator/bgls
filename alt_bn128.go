// Copyright (C) 2018 Authors
// distributed under Apache 2.0 license

package bgls

import (
	"math/big"

	"github.com/dchest/blake2b"
	"github.com/ethereum/go-ethereum/crypto/bn256"
	gosha3 "github.com/ethereum/go-ethereum/crypto/sha3"
	"golang.org/x/crypto/sha3"
)

//curve specific constants
var altbnB = big.NewInt(3)
var altbnQ, _ = new(big.Int).SetString("21888242871839275222246405745257275088696311157297823662689037894645226208583", 10)

//precomputed Z = (-1 + sqrt(-3))/2 in Fq
var altbnZ, _ = new(big.Int).SetString("2203960485148121921418603742825762020974279258880205651966", 10)

//precomputed sqrt(-3) in Fq
var altbnSqrtn3, _ = new(big.Int).SetString("4407920970296243842837207485651524041948558517760411303933", 10)

// Note that the cofactor in this curve is just 1

// AltbnSha3 Hashes a message to a point on Altbn128 using SHA3 and try and increment
// The return value is the x,y affine coordinate pair.
func AltbnSha3(message []byte) (p1, p2 *big.Int) {
	p1, p2 = hash64(message, sha3.Sum512, altbnQ, altbnXToYSquared)
	return
}

// AltbnKeccak3 Hashes a message to a point on Altbn128 using Keccak3 and try and increment
// Keccak3 is only for compatability with Ethereum hashing.
// The return value is the x,y affine coordinate pair.
func AltbnKeccak3(message []byte) (p1, p2 *big.Int) {
	p1, p2 = hash32(message, EthereumSum256, altbnQ, altbnXToYSquared)
	return
}

// AltbnBlake2b Hashes a message to a point on Altbn128 using Blake2b and try and increment
// The return value is the x,y affine coordinate pair.
func AltbnBlake2b(message []byte) (p1, p2 *big.Int) {
	p1, p2 = hash64(message, blake2b.Sum512, altbnQ, altbnXToYSquared)
	return
}

// AltbnKang12 Hashes a message to a point on Altbn128 using Blake2b and try and increment
// The return value is the x,y affine coordinate pair.
func AltbnKang12(message []byte) (p1, p2 *big.Int) {
	p1, p2 = hash64(message, kang12_64, altbnQ, altbnXToYSquared)
	return
}

func altbnXToYSquared(x *big.Int) *big.Int {
	result := new(big.Int)
	result.Exp(x, three, altbnQ)
	result.Add(result, altbnB)
	return result
}

// AltbnMkG1Point copies points into []byte and unmarshals to get around curvePoint not being exported
// This is copied from bn256.G1.Marshal (modified)
func AltbnMkG1Point(x, y *big.Int) (*bn256.G1, bool) {
	xBytes, yBytes := x.Bytes(), y.Bytes()
	ret := make([]byte, 64)
	copy(ret[32-len(xBytes):], xBytes)
	copy(ret[64-len(yBytes):], yBytes)
	return new(bn256.G1).Unmarshal(ret)
}

// AltbnHashToCurve Hashes a message to a point on Altbn128 using Keccak3 and try and increment
// This is for compatability with Ethereum hashing.
// The return value is the altbn_128 library's internel representation for points.
func AltbnHashToCurve(message []byte) *bn256.G1 {
	x, y := AltbnKeccak3(message)
	p, _ := AltbnMkG1Point(x, y)
	return p
}

// AltbnG1ToCoord takes a point in G1 of Altbn_128, and returns its affine coordinates
func AltbnG1ToCoord(pt *bn256.G1) (x, y *big.Int) {
	Bytestream := pt.Marshal()
	xBytes, yBytes := Bytestream[:32], Bytestream[32:64]
	x = new(big.Int).SetBytes(xBytes)
	y = new(big.Int).SetBytes(yBytes)
	return
}

// AltbnG2ToCoord takes a point in G2 of Altbn_128, and returns its affine coordinates
func AltbnG2ToCoord(pt *bn256.G2) (xx, xy, yx, yy *big.Int) {
	Bytestream := pt.Marshal()
	xxBytes, xyBytes, yxBytes, yyBytes := Bytestream[:32], Bytestream[32:64], Bytestream[64:96], Bytestream[96:128]
	xx = new(big.Int).SetBytes(xxBytes)
	xy = new(big.Int).SetBytes(xyBytes)
	yx = new(big.Int).SetBytes(yxBytes)
	yy = new(big.Int).SetBytes(yyBytes)
	return
}

// EthereumSum256 returns the Keccak3-256 digest of the data. This is because Ethereum
// uses a non-standard hashing algo.
func EthereumSum256(data []byte) (digest [32]byte) {
	h := gosha3.NewKeccak256()
	h.Write(data)
	h.Sum(digest[:0])
	return
}