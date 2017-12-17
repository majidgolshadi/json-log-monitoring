package rest

var isJsonFlag bool

func init() {
	isJsonFlag = true
}

func SetNotJSON()  {
	isJsonFlag = false
}

func SetIsJSON()  {
	isJsonFlag = true
}

func getStatus()  bool {
	return isJsonFlag
}