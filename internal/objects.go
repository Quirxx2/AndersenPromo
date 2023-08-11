package promo

type Grade int

const (
	trainee Grade = iota + 1
	junior
	middle
	senior
)

var dGrades = map[Grade]string{
	trainee: "trainee",
	junior:  "junior",
	middle:  "middle",
	senior:  "senior",
}

var bGrades = map[string]Grade{
	"trainee": trainee,
	"junior":  junior,
	"middle":  middle,
	"senior":  senior,
}

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Position Grade  `json:"position"`
	Project  string `json:"project"`
}
