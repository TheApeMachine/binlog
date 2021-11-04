package filters

type BooleanGrid interface {
	AddClause(BooleanGrid) BooleanGrid
}
