package mailings

import "github.com/Niroloc/Temcshiki/v2/src/context"

func ApplyTasks(context *context.Context) {
	for _, task := range context.GetTasks() {
		if _, exists := task.Stages()[context.GetStage()]; !exists {
			continue
		}
		task.Apply(context)
	}
}
