package main

type Todo struct {
	//gorm.Model
	Id              string     `json:"id"`
	Name            NullString `json:"name"`
	Description     NullString `json:"description"`
	Priority        NullInt64  `json:"priority"`
	Due_Date        NullTime   `json:"due"`
	Completed       NullBool   `json:"completed"`
	Completion_Date NullTime   `json:"completion_date"`
}

// func (t *Todo) MarshalJSON() ([]byte, error) {
// 	return json.Marshal(&struct {
// 		Id              string `json:"id"`
// 		Name            string `json:"name"`
// 		Description     string `json:"description"`
// 		Priority        int64  `json:"priority"`
// 		Due_date        string `json:"due_date"`
// 		Completed       bool   `json:"completed"`
// 		Completion_date string `json:"completion_date"`
// 	}{
// 		Id:              t.Id,
// 		Name:            t.Name.String,
// 		Description:     t.Description.String,
// 		Priority:        t.Priority.Int64,
// 		Due_date:        getCorrectDate(t.Completion_Date),
// 		Completed:       t.Completed.Bool,
// 		Completion_date: getCorrectDate(t.Completion_Date),
// 	})
// }

// func getCorrectDate(t NullTime) string {
// 	if t.Valid {
// 		return t.Time.String()
// 	}
// 	return "null"
// }
