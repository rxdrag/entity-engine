package table

import "rxdrag.com/entify/model/meta"

type Column struct {
	meta.AttributeMeta
	PartialId bool
	Key       bool
}
