package data

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/Niroloc/Temcshiki/v2/src/db"
	"github.com/Niroloc/Temcshiki/v2/src/logger"
)

type Data struct {
	tgIdToUser map[int]User
	stage      db.Stage
	tasks      []Task
	bot        *Bot
	db         *db.Db
	logger     *logger.Logger
}

func CreateData(dbWrapper *db.Db, bot *Bot, tasks []Task) *Data {
	tgIdToUser := map[int]User{}
	stage := dbWrapper.GetStage()
	for _, user := range dbWrapper.GetUsers() {
		tgIdToUser[user.TgId] = User{user.Id, user.Username, db.UserRights(user.Rights)}
	}
	return &Data{
		tgIdToUser: tgIdToUser,
		stage:      stage,
		tasks:      tasks,
		db:         dbWrapper,
		logger:     logger.GetLogger(reflect.TypeFor[Data]())}
}

func (this *Data) GetTasks() []Task {
	return this.tasks
}

func (this *Data) GetStage() db.Stage {
	return this.stage
}

func (this *Data) GetUser(tgId int) (User, error) {
	if user, exists := this.tgIdToUser[tgId]; exists {
		return user, nil
	}
	this.logger.Warn("Unknown user has been searched")
	return User{0, "", ""}, errors.New("Unknown user")
}

func (this *Data) GetUsers() map[int]User {
	return this.tgIdToUser
}

func (this *Data) GetBot() *Bot {
	return this.bot
}

func (this *Data) NextStage() {
	this.stage.Next()
}

func (this *Data) InfinitePolling() error {
	bot := this.bot.bot
	updates, err := bot.UpdatesViaLongPolling(context.Background(), nil)
	if err != nil {
		this.logger.Error("Cannot start polling")
		panic(err)
	}
	this.logger.Info("Starting polling")
	for update := range updates {
		if update.Message != nil {
			tgId := int(update.Message.Chat.ID)
			this.bot.SendMessage(tgId, fmt.Sprint(this.tgIdToUser[tgId].Username))
		} else if update.CallbackQuery != nil {
			continue
		}
	}
	return errors.New("Polling closed")
}
