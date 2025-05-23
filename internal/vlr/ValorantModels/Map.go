package ValorantModels

type ValorantMap struct {
	ID    uint   `json:"id" gorm:"primaryKey"`
	Title string `json:"title"`
}
