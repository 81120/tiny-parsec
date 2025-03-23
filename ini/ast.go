package ini

type Ini struct {
	Sections []Section
}

type Section struct {
	Name    string
	Entries []Entry
}

type Entry struct {
	Key   string
	Value string
}
