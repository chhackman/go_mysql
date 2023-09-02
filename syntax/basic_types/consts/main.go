package consts

const External = "包外"
const internal = "包内"
const (
	a = 123
)

const (
	Init = iota
	Running
	Paused
	Stop

	StatusE = 100
	StatusF
)

const (
	DayA = iota
	DayB
	DayC
	DayD
	DayE
)

func main() {
	const a = 123
}