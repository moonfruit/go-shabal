package shabal

/*
#include "sph_shabal.h"

typedef void (*sph_shabal_init)(void *cc);
static void shabal_init(sph_shabal_init init, void *cc) {
	init(cc);
}

typedef void (*sph_shabal)(void *cc, const void *data, size_t len);
static void shabal(sph_shabal update, void *cc, const void *data, size_t len) {
	update(cc, data, len);
}

typedef void (*sph_shabal_close)(void *cc, void *dst);
static void shabal_close(sph_shabal_close close, void *cc, void *dst) {
	close(cc, dst);
}
*/
import "C"
import (
	"hash"
	"strconv"
	"unsafe"
)

type hashFunc struct {
	init   C.sph_shabal_init
	update C.sph_shabal
	close  C.sph_shabal_close
}

var shabal192 = hashFunc{
	C.sph_shabal_init(C.sph_shabal192_init),
	C.sph_shabal(C.sph_shabal192),
	C.sph_shabal_close(C.sph_shabal192_close),
}
var shabal224 = hashFunc{
	C.sph_shabal_init(C.sph_shabal224_init),
	C.sph_shabal(C.sph_shabal224),
	C.sph_shabal_close(C.sph_shabal224_close),
}
var shabal256 = hashFunc{
	C.sph_shabal_init(C.sph_shabal256_init),
	C.sph_shabal(C.sph_shabal256),
	C.sph_shabal_close(C.sph_shabal256_close),
}
var shabal384 = hashFunc{
	C.sph_shabal_init(C.sph_shabal384_init),
	C.sph_shabal(C.sph_shabal384),
	C.sph_shabal_close(C.sph_shabal384_close),
}
var shabal512 = hashFunc{
	C.sph_shabal_init(C.sph_shabal512_init),
	C.sph_shabal(C.sph_shabal512),
	C.sph_shabal_close(C.sph_shabal512_close),
}

func getHashFunc(size int) *hashFunc {
	switch size {
	case 192:
		return &shabal192
	case 224:
		return &shabal224
	case 256:
		return &shabal256
	case 384:
		return &shabal384
	case 512:
		return &shabal512
	default:
		panic("Unsupported size " + strconv.Itoa(size))
	}
}

type shabal struct {
	context C.sph_shabal_context
	size    int
	*hashFunc
}

func (s *shabal) Write(b []byte) (n int, err error) {
	pointer, size := bytesToC(b)
	C.shabal(s.hashFunc.update, unsafe.Pointer(&s.context), pointer, size)
	return len(b), nil
}

func (s *shabal) Sum(b []byte) []byte {
	ret := make([]byte, s.Size())
	C.shabal_close(s.hashFunc.close, unsafe.Pointer(&s.context), unsafe.Pointer(&ret[0]))
	return append(b, ret...)
}

func (s *shabal) Reset() {
	C.shabal_init(s.hashFunc.init, unsafe.Pointer(&s.context))
}

func (s *shabal) Size() int {
	return s.size / 8
}

func (*shabal) BlockSize() int {
	return 64
}

func newShabal(size int) hash.Hash {
	s := &shabal{size: size, hashFunc: getHashFunc(size)}
	s.Reset()
	return s
}

func NewShabal192() hash.Hash {
	return newShabal(192)
}

func NewShabal224() hash.Hash {
	return newShabal(224)
}

func NewShabal256() hash.Hash {
	return newShabal(256)
}

func NewShabal384() hash.Hash {
	return newShabal(384)
}

func NewShabal512() hash.Hash {
	return newShabal(512)
}

func bytesToC(buf []byte) (unsafe.Pointer, C.size_t) {
	length := len(buf)
	if length == 0 {
		return nil, 0
	}
	return unsafe.Pointer(&buf[0]), C.size_t(length)
}
