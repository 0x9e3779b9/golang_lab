package models

type Task struct {
	ID         int64       `json:"id"`
	EnNo       string      `json:"en_no"`
	Owner      string      `json:"owner"`
	Cases      string      `json:"case_ids"`
	PlanFinish string      `json:"plan_finish"`
	State      int         `json:"state"`
	Opt        *TaskOption `json:"opt"`
}

type TaskCase struct {
	ID     int64       `json:"id"`
	TaskID int64       `json:"task_id"`
	CaseID int64       `json:"case_id"`
	Opt    *CaseOption `json:"opt"`
}

type Record struct {
	ID         int64       `json:"id"`
	TaskCaseID int64       `json:"task_case_id"`
	TaskOpt    *TaskOption `json:"task_opt"`
	CaseOpt    *CaseOption `json:"case_opt"`
	StartTime  string      `json:"start_time"`
	EndTime    string      `json:"end_time"`
}

type TaskOption struct {
	MeterConf map[string]string `json:"meter"`
	DUTConf   map[string]string `json:"dut"`
}

type CaseOption struct {
	Param    map[string]string `json:"param"`
	ConfPath string            `json:"conf_path"`
}

func NewRecord(tc *TaskCase) *Record {
	return &Record{
		TaskCaseID: tc.ID,
		CaseOpt:    tc.Opt,
		StartTime:  "now",
		EndTime:    "tom",
	}
}

func (r *Record) Result(byt []byte) {
}

func InsertRecord(r *Record) {

}
