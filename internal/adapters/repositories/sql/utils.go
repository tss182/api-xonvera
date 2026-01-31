package repositoriesSql

func GetTotalPage(totalData int64, limit uint8) uint8 {
	totalPages := totalData / int64(limit)
	if totalData%int64(limit) != 0 {
		totalPages++
	}
	return uint8(totalPages)
}