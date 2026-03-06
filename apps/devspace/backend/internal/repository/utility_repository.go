package repository

import "gorm.io/gorm"

// Удаляет все записи в таблице с полем, значение которого val
func DeleteAll(field, tableName string, val any, db *gorm.DB) error {
	res := db.Table(tableName).Delete(nil, field+" = ?", val)
	return res.Error
}
