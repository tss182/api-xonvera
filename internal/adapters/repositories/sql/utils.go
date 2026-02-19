package repositoriesSql

func GetTotalPage(totalData uint64, limit uint8) uint8 {
	totalPages := totalData / uint64(limit)
	if totalData%uint64(limit) != 0 {
		totalPages++
	}
	return uint8(totalPages)
}
