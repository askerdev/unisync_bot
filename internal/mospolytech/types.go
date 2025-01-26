package mospolytech

type Auditory struct {
	Title string `json:"title"`
}

type Lecture struct {
	Subject    string     `json:"sbj"`
	Teacher    string     `json:"teacher"`
	DateFrom   string     `json:"df"`
	DateTo     string     `json:"dt"`
	Location   string     `json:"location"`
	Type       string     `json:"type"`
	Link       *string    `json:"e_link"`
	Auditories []Auditory `json:"auditories"`
}

type Grid map[string]map[string][]*Lecture

type Group struct {
	Grid Grid `json:"grid"`
}

type Contents map[string]Group

type SemesterSchedule struct {
	Content Contents `json:"contents"`
}
