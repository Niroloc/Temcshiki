package mailings

import (
	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/tasks"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

func ApplyTasks(bot *tg.Bot, data *data.Data, tsks []tasks.Task) {
	for _, task := range tsks {
		if _, exists := task.Stages()[data.GetStage()]; !exists {
			continue
		}
		task.Apply(bot, data)
	}
}
