package models

// User 用户模型
type User struct {
	BaseModel
	Username string `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password string `gorm:"type:varchar(255);not null" json:"-"`
	Email    string `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Nickname string `gorm:"type:varchar(50)" json:"nickname"`
	Avatar   string `gorm:"type:varchar(255)" json:"avatar"`
	Role     string `gorm:"type:varchar(20);default:'user'" json:"role"` // admin, user
	Status   int    `gorm:"default:1" json:"status"`                     // 1:正常 0:禁用
	Articles []Article `gorm:"foreignKey:AuthorID" json:"articles,omitempty"`
	Comments []Comment `gorm:"foreignKey:UserID" json:"comments,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
