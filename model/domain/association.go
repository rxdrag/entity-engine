package domain

import "rxdrag.com/entity-engine/model/meta"

type Association struct {
	Relation       *Relation
	OwnerClassUuid string
}

func (a *Association) Name() string {
	if a.Relation.Source.Uuid == a.OwnerClassUuid {
		return a.Relation.RoleOfTarget
	} else {
		return a.Relation.RoleOfSource
	}
}

func (a *Association) Owner() *Class {
	if a.Relation.Source.Uuid == a.OwnerClassUuid {
		return a.Relation.Source
	} else {
		return a.Relation.Target
	}
}

func (a *Association) TypeClass() *Class {
	if a.Relation.Source.Uuid == a.OwnerClassUuid {
		return a.Relation.Target
	} else {
		return a.Relation.Source
	}
}

func (a *Association) Description() string {
	if a.Relation.Source.Uuid == a.OwnerClassUuid {
		return a.Relation.DescriptionOnTarget
	} else {
		return a.Relation.DescriptionOnSource
	}
}

func (a *Association) isArray() bool {
	if a.Relation.Source.Uuid == a.OwnerClassUuid {
		return a.Relation.SourceMutiplicity == meta.ZERO_MANY
	} else {
		return a.Relation.TargetMultiplicity == meta.ZERO_MANY
	}
}
