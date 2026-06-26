package durableobjects

//easjson:json
type Map = map[string]any

//easyjson:json
type SqlStorageCursorProp struct {
	Value Map
	Done  bool
}
