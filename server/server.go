package server

func Init() {
	router := SetupRouter()
	router.Run(":80")
}
