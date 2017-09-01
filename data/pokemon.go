package data

type Pokemon struct {
	Name       string
	Abilities  []string
	Stats      []PokemonStat
	Types      []Type
	DamageType []DamageType
}

type PokemonStat struct {
	Name       string
	Base       int
	BattleOnly bool
}
