package models

type BaseModel struct {
	// CreatedAt time.Time  `json:"createdAt"`
	// UpdatedAt time.Time  `json:"updatedAt"`
	// DeletedAt *time.Time `json:"deletedAt"`
	CreateTime string `gorm:"size:128" json:"createTime"`
	UpdateTime string `gorm:"size:128" json:"updateTime"`
}
