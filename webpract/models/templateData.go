package models

type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	CSRFtoken string
	Flash     string
	Warning   string
	Error     string
}
