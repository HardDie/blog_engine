package utils

func GetPagination(limit, page int32) (newLimit, offset int64) {
	if page > 0 {
		offset = int64(page - 1)
	}
	if limit > 0 {
		newLimit = int64(limit)
	}
	offset = offset * newLimit
	return
}
