package mailings

import "github.com/Niroloc/Temcshiki/v2/src/data"

func ApplyTasks(data *data.Data) {
	for _, task := range data.GetTasks() {
		if _, exists := task.Stages()[data.GetStage()]; !exists {
			continue
		}
		task.Apply(data)
	}
}
