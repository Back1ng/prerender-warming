package warmer

type Warmer interface {
	Process(url string)
	Refresh(<-chan string, *int)
	ResetWriter()
	StartWriter()
}
