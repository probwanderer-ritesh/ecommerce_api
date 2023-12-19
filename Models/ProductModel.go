package Models

// import "github.com/jinzhu/gorm"

type Product struct {
	Id          uint    `json:"id"`
	Name        string  `json:"name"`
	Quantity    int     `json:"quantity"`
	Description string  `json:"description"`
	Price       float64 `json:"float"`
}

func (p *Product) TableName() string {
	return "product"
}
