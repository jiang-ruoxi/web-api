package web

type Lexicon struct {
	ID         uint   `gorm:"primaryKey;comment:主键ID" json:"id"`                                                                                       // 主键ID
	Word       string `gorm:"column:word;type:varchar(500);not null;default:'';comment:标题" json:"word" validate:"required" message:"标题必填"`             // 标题
	Mark       string `gorm:"column:mark;type:varchar(500);not null;default:'';comment:标题" json:"mark" validate:"required" message:"标题必填"`             // 标题
	Annotation string `gorm:"column:annotation;type:varchar(500);not null;default:'';comment:标题" json:"annotation" validate:"required" message:"标题必填"` // 标题
	Type       uint   `gorm:"column:type;type:int unsigned;default:0;not null;comment:类型ID" json:"type" validate:"required" message:"类型必填"`            // 类型ID
}
