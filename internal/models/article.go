package models

// Article 文章模型
type Article struct {
	BaseModel
	Title       string    `gorm:"type:varchar(200);not null" json:"title"`
	Description string    `gorm:"type:varchar(500)" json:"description"`
	Content     string    `gorm:"type:longtext;not null" json:"content"`
	Cover       string    `gorm:"type:varchar(255)" json:"cover"`
	AuthorID    uint      `gorm:"not null;index" json:"author_id"`
	Author      User      `gorm:"foreignKey:AuthorID" json:"author,omitempty"`
	CategoryID  uint      `gorm:"index" json:"category_id"`
	Category    Category  `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags        []Tag     `gorm:"many2many:article_tags;" json:"tags,omitempty"`
	ViewCount   int       `gorm:"default:0" json:"view_count"`
	LikeCount   int       `gorm:"default:0" json:"like_count"`
	Status      int       `gorm:"default:1;index" json:"status"` // 1:已发布 0:草稿
	IsTop       bool      `gorm:"default:false" json:"is_top"`
	Comments    []Comment `gorm:"foreignKey:ArticleID" json:"comments,omitempty"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}
