package filters

type Column interface {
	Name() string
}

type GridColumn struct {
	name string
}

func NewGridColumn() Column {
	return GridColumn{}
}

func (col GridColumn) Name() string {
	return col.name
}
