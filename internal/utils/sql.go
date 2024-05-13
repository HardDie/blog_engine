package utils

import (
	"database/sql"
	"strings"
)

func PrepareStringToLike(s string) string {
	if s == "" {
		return s
	}
	arr := strings.Split(s, " ")
	if len(arr) == 1 {
		return "%" + s + "%"
	}
	return "%" + strings.Join(arr, "%") + "%"
}
func NewSqlString(val *string) sql.NullString {
	if val == nil {
		return sql.NullString{}
	}
	return sql.NullString{
		String: *val,
		Valid:  true,
	}
}
func NewSqlInt64(val *int64) sql.NullInt64 {
	if val == nil {
		return sql.NullInt64{}
	}
	return sql.NullInt64{
		Int64: *val,
		Valid: true,
	}
}
func SqlStringToString(val sql.NullString) *string {
	if !val.Valid {
		return nil
	}
	return &val.String
}
