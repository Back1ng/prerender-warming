package warmer

type Warmer interface {
	Process(<-chan string, chan<- struct{})
	Add(url string)
	Refresh()
}
