package warmer

type Warmer interface {
	Process(url string)
	Refresh(<-chan string)
	ResetWriter()
	StartWriter()
}
