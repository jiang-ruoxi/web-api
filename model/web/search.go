package web

import (
	"gorm.io/plugin/soft_delete"
)

type Section struct {
	ID         uint                  `gorm:"primaryKey;comment:主键ID" json:"id"`                                                                             // 主键ID
	Name       string                `gorm:"column:name;type:varchar(500);default:'';not null;comment:栏目名称" json:"name" validate:"required" message:"名称必填"` // 栏目名称
	SecType    uint8                 `gorm:"column:sec_type;index;type:int unsigned;default:0;not null;comment:类型:1栏目 2类型" json:"sec_type"`                 // 类型:1栏目 2类型
	LibType    uint8                 `gorm:"column:lib_type;type:tinyint unsigned;default:0;not null;comment:库类型:1动态库 2法律库 3案例库 4成果库 5报告库" json:"lib_type"` // 库类型:1动态库 2法律库 3案例库 4成果库 5报告库
	CreateTime int64                 `gorm:"type:uint;not null;default:0;comment:创建时间" json:"create_time"`                                                  // 创建时间(系统自动生成)
	UpdateTime int64                 `gorm:"type:uint;not null;default:0;comment:最近更新" json:"update_time"`                                                  // 最近更新(系统自动生成)
	DeleteTime soft_delete.DeletedAt `gorm:"index;type:uint;not null;default:0;comment:删除时间" json:"-"`                                                      // 删除时间
}

func (Section) TableName() string {
	return "e_sections"
}
