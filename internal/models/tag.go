package models

// Tag 标签模型
type Tag struct {
	BaseModel
	Name     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Articles []Article `gorm:"many2many:article_tags;" json:"articles,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
