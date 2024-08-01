package main

type Suit int

const (
	Red Suit = iota
	Green
	Blue
	Purple
)

func (s Suit) String() string {
	switch s {
	case Red:
		return "red"
	case Green:
		return "green"
	case Blue:
		return "blue"
	case Purple:
		return "purple"
	default:
		return "error"
	}
}

type Archetype int

const (
	None Archetype = iota
	Food
	Resource
	Animal
	Leisure
	Governance
	Market
	Crafting
)

func (a Archetype) String() string {
	switch a {
	case None:
		return "none"
	case Food:
		return "none"
	case Resource:
		return "none"
	case Animal:
		return "none"
	case Leisure:
		return "none"
	case Governance:
		return "none"
	case Market:
		return "none"
	case Crafting:
		return "none"
	default:
		return "error"
	}
}

type Effect int

const (
	Gain Effect = iota
	GainAny
	GainMultiplied
	Take
	All
	Choose
	Exchange
)

func (e Effect) String() string {
	switch e {
	case Gain:
		return "gain"
	case GainAny:
		return "gain (any)"
	case GainMultiplied:
		return "gain (multiplied)"
	case Take:
		return "take"
	case All:
		return "all"
	case Choose:
		return "choose"
	case Exchange:
		return "exchange"
	default:
		return "error"
	}
}

type Card struct {
	Name         string
	Description  string
	PurchaseCost int
	DieTriggers  []int
	Suit
	Archetype
	Effect
	EffectSource Archetype
}

type CardJson struct {
	Name        string `json:"name"`
	Description string `json:"description"`

	PurchaseCost int    `json:"purchaseCost"`
	DieTriggers  []int  `json:"dieTriggers"`
	Suit         string `json:"suit"`
	Archetype    string `json:"archetype"`
	Effect       string `json:"effect"`
	EffectSource string `json:"effectSource"`
}

var (
	AllCards = []Card{
		{
			Name:         "Wheat field",
			Description:  "Get 1 coin from the bank.",
			PurchaseCost: 1,
			DieTriggers: []int{
				1,
			},
			Suit:         Blue,
			Archetype:    Food,
			Effect:       GainAny,
			EffectSource: None,
		},
		{
			Name:         "Ranch",
			Description:  "Get 1 coin from the bank.",
			PurchaseCost: 1,
			DieTriggers: []int{
				2,
			},
			Suit:         Blue,
			Archetype:    Animal,
			Effect:       GainAny,
			EffectSource: None,
		},
		{
			Name:         "Bakery",
			Description:  "Get 1 coin from the bank.",
			PurchaseCost: 1,
			DieTriggers: []int{
				2,
				3,
			},
			Suit:         Green,
			Archetype:    Market,
			Effect:       Gain,
			EffectSource: None,
		},
		{
			Name:         "Café",
			Description:  "Take 1 coin from the active player.",
			PurchaseCost: 2,
			DieTriggers: []int{
				3,
			},
			Suit:         Red,
			Archetype:    Leisure,
			Effect:       Take,
			EffectSource: None,
		},
	}
)
