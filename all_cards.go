package main

import "encoding/json"

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

func (s Suit) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
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
		return "food"
	case Resource:
		return "resource"
	case Animal:
		return "animal"
	case Leisure:
		return "leisure"
	case Governance:
		return "governance"
	case Market:
		return "market"
	case Crafting:
		return "crafting"
	default:
		return "error"
	}
}

func (a Archetype) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
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

func (e Effect) MarshalJSON() ([]byte, error) {
	return json.Marshal(e.String())
}

type Card struct {
	Name         string    `json:"name"`
	Description  string    `json:"description"`
	PurchaseCost int       `json:"purchaseCost"`
	DieTriggers  []int     `json:"dieTriggers"`
	Suit         Suit      `json:"suit"`
	Archetype    Archetype `json:"archetype"`
	Effect       Effect    `json:"effect"`
	EffectSource Archetype `json:"effectSource"`
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

		//TODO: To be continued...
	}
)
