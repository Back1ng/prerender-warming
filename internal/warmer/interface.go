package warmer

type Warmer interface {
	Process(url string)
	Add(url string)
	Refresh()
}
