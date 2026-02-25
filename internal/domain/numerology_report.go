package domain

type ShortSection struct {
	Number      string `json:"number"`
	Description string `json:"description"`
}

type NumerologyReport struct {
	PersonalityOutlook   []string     `json:"personality_outlook"`
	LifePathNumber       ShortSection `json:"life_path_number"`
	DestinyNumber        ShortSection `json:"destiny_number"`
	CareerPrediction     []string     `json:"career_prediction"`
	WealthPrediction     []string     `json:"wealth_prediction"`
	MarriageRelationship []string     `json:"marriage_relationship"`
	HealthPrediction     []string     `json:"health_prediction"`
	Remedies             []string     `json:"remedies"`
}
