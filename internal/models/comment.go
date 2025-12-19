package models

// Comment 评论模型
type Comment struct {
	BaseModel
	ArticleID uint   `gorm:"not null;index" json:"article_id"`
	Article   Article `gorm:"foreignKey:ArticleID" json:"article,omitempty"`
	UserID    uint   `gorm:"not null;index" json:"user_id"`
	User      User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Content   string `gorm:"type:text;not null" json:"content"`
	ParentID  *uint  `gorm:"index" json:"parent_id"` // 父评论ID，用于回复
	Status    int    `gorm:"default:1" json:"status"` // 1:正常 0:已删除
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
