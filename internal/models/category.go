package models

// Category 分类模型
type Category struct {
	BaseModel
	Name        string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string    `gorm:"type:varchar(255)" json:"description"`
	Sort        int       `gorm:"default:0" json:"sort"`
	Articles    []Article `gorm:"foreignKey:CategoryID" json:"articles,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}
