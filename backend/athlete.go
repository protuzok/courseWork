package backend

type Athlete struct {
	Id           int     `json:"id"`
	Name         string  `json:"name"`
	Surname      string  `json:"surname"`
	Run100m      float32 `json:"run_100m"`
	Run3km       float32 `json:"run_3km"`
	PressCnt     int     `json:"press_cnt"`
	JumpDistance float32 `json:"jump_distance"`
}
