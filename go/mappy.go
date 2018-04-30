package main

// #include "minimap.h"
// #cgo CFLAGS: -I${SRCDIR}/..
// #cgo LDFLAGS: -L${SRCDIR}/.. -lminimap2 -lz -pthread
import "C"
import (
	"unsafe"
	"fmt"
	"os"
)

func align(idxFile, read string, n_threads int) {
	c_idxopt := C.mm_idxopt_t{}
	c_map_opt := C.mm_mapopt_t{}

	C.mm_set_opt(
		nil,
		(*C.mm_idxopt_t)(unsafe.Pointer(&c_idxopt)),
		(*C.mm_mapopt_t)(unsafe.Pointer(&c_map_opt)),
	)

	c_map_opt.flag |= 4                      // always perform alignment
	c_idxopt.batch_size = 0x7fffffffffffffff // always build a uni-part index

	r := C.mm_idx_reader_open(
		C.CString(idxFile),
		(*C.mm_idxopt_t)(unsafe.Pointer(&c_idxopt)),
		nil,
	)

	if r == nil {
		fmt.Println("Error openning index file")
		os.Exit(1)
	}

	c_idx := C.mm_idx_reader_read(
		(*C.mm_idx_reader_t)(unsafe.Pointer(r)),
		C.int(n_threads),
	)

	C.mm_idx_reader_close(r)

	C.mm_mapopt_update(
		(*C.mm_mapopt_t)(unsafe.Pointer(&c_map_opt)),
		(*C.mm_idx_t)(unsafe.Pointer(c_idx)),
	)

	C.mm_idx_index_name(
		(*C.mm_idx_t)(unsafe.Pointer(c_idx)),
	)

	n_regs := 0
	buf := C.mm_tbuf_init()

	// map single end
	c_mm_reg1_t := C.mm_map(
		(*C.mm_idx_t)(unsafe.Pointer(c_idx)),
		C.int(len(read)),
		C.CString(read),
		(*C.int)(unsafe.Pointer(&n_regs)),
		(*C.mm_tbuf_t)(unsafe.Pointer(buf)),
		(*C.mm_mapopt_t)(unsafe.Pointer(&c_map_opt)),
		nil,
	)

	if c_mm_reg1_t != nil {
		fmt.Println(*c_mm_reg1_t)
	} else {
		fmt.Println("Failed to be aligned")
	}

	// clean up
	C.mm_tbuf_destroy((*C.mm_tbuf_t)(unsafe.Pointer(buf)))
	C.mm_idx_destroy((*C.mm_idx_t)(unsafe.Pointer(c_idx)))
}

func main() {
	align(
		"/Users/admin_bo/src/minimap2/go/dm3.mmi",
		"CAATCTTCCGGCCAGCCAATCGAGCGGCCAAATCTGGCGGGC"+
			"AAGTCGTGGTATTATGGAGCCATCACCCGCAGCCAGTGCGACAC"+
			"GGTGCTCAACGGCCACGGGCACGATGGCGACTTCCTCATCAGAG",
		3,
	)
}
