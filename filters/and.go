package filters

type And struct {
	source Column
}

func NewAnd() BooleanGrid {
	return And{}
}

func (boolean And) AddClause(parent BooleanGrid) BooleanGrid {
	return boolean
}
