package ValorantModels

type ValorantAgents struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}
