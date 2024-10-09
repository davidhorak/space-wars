package game

var uuid int64 = 0

func NewUUID() int64 {
	uuid++
	return uuid
}

func GetUUID() int64 {
	return uuid
}

func ResetUUID() {
	uuid = 0
}

func SetUUID(id int64) {
	uuid = id
}
