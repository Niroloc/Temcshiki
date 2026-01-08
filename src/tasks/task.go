package tasks

import (
	"reflect"
	"strings"
	"time"

	"github.com/Niroloc/Temcshiki/v2/src/data"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
	"github.com/Niroloc/Temcshiki/v2/src/tg"
)

type Task interface {
	Stages() map[data.Stage]struct{}
	Apply(*tg.Bot, *data.Data)
}

type Tasks struct {
	exportedData *data.Data
	bot          *tg.Bot
	tasks        []Task
	scheduled    []*time.Timer
	stageToTime  map[data.Stage]string
	logger       *logger.Logger
}

func InitTasks(exportedData *data.Data, bot *tg.Bot) *Tasks {
	return &Tasks{
		exportedData: exportedData,
		bot:          bot,
		tasks: []Task{
			NewChooseReminder(),
		},
		scheduled: []*time.Timer{},
		stageToTime: map[data.Stage]string{
			data.CHOOSING:    "12:00:00",
			data.VOTING:      "19:00:00",
			data.REMINDING:   "19:00:00",
			data.APPROVING:   "19:00:00",
			data.RESERVATING: "13:00:00",
			data.REVIEWING:   "19:00:00",
		},
		logger: logger.GetLogger(reflect.TypeFor[Tasks]()),
	}
}

func (this *Tasks) applyTasks() {
	for _, task := range this.tasks {
		if _, exists := task.Stages()[this.exportedData.GetStage()]; !exists {
			continue
		}
		task.Apply(this.bot, this.exportedData)
	}
	this.exportedData.NextStage()
}

func (this *Tasks) Loop() {
	for true {
		if len(this.scheduled) > 0 {
			time.Sleep(time.Minute * 5)
		}
		ts, exists := this.stageToTime[this.exportedData.GetStage()]
		if !exists {
			ts = "19:00:00"
		}
		date := this.exportedData.GetDb().GetNextTaskDate()
		taskTime, err := time.Parse(time.DateTime, strings.Join([]string{date, ts}, " "))
		if err != nil {
			this.logger.Error("Error while scheduling new task")
			this.logger.Error(err.Error())
			time.Sleep(time.Second * 10)
		}
		delay := taskTime.Sub(time.Now())
		this.scheduled = append(this.scheduled, time.AfterFunc(delay, this.applyTasks))
	}
}
