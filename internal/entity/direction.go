package entity

type Direction string

const (
	DirectionInbound  Direction = "inbound"
	DirectionOutbound Direction = "outbound"
	DirectionBoth     Direction = "both"
)

func (d Direction) String() string {
	return string(d)
}
