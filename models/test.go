package models

import (
	orm "dh-passwd/global"
	_ "time"
)

// 定义你的数据表字段
type Test struct {
	Id   int    `json:"id" gorm:"primary_key;AUTO_INCREMENT"` //ID
	Name string `json:"name" gorm:"name"`                     //名称
}

func (Test) TableName() string {
	return "test"
}

// 创建数据的方法
func (e *Test) Create() (Test, error) {
	var doc Test
	result := orm.Eloquent.Table(e.TableName()).Create(&e)
	if result.Error != nil {
		err := result.Error
		return doc, err
	}
	doc = *e
	return doc, nil
}

// 获取你的单个数据
func (e *Test) Get() (Test, error) {
	var doc Test
	table := orm.Eloquent.Table(e.TableName())

	if e.Id != 0 {
		table = table.Where("id = ?", e.Id)
	}

	if err := table.First(&doc).Error; err != nil {
		return doc, err
	}
	return doc, nil
}

// 获取你的批量数据
func (e *Test) GetList() ([]Test, int, error) {
	var doc []Test
	var count int

	table := orm.Eloquent.Select("*").Table(e.TableName())
	if e.Id != 0 {
		table = table.Where("id = ?", e.Id)
	}

	//table.Where("`update_time` IS NULL").Count(&count)
	if err := table.Find(&doc).Error; err != nil {
		return nil, count, err
	}
	return doc, count, nil
}

// 更新你的数据，通过ID
func (e *Test) Update(id int) (update Test, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).First(&update).Error; err != nil {
		return
	}
	//参数1:是要修改的数据
	//参数2:是修改的数据
	if err = orm.Eloquent.Table(e.TableName()).Model(&update).Updates(&e).Error; err != nil {
		return
	}
	return
}

// 删除你的数据，哦你给过ID
func (e *Test) Delete(id int) (success bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id = ?", id).Delete(&Test{}).Error; err != nil {
		success = false
		return
	}
	success = true
	return
}

//批量删除你的数据，根据ID列表
func (e *Test) BatchDelete(id []int) (Result bool, err error) {
	if err = orm.Eloquent.Table(e.TableName()).Where("id in (?)", id).Delete(&Test{}).Error; err != nil {
		return
	}
	Result = true
	return
}
