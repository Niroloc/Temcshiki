package data

import (
	"fmt"
)

type Rest struct {
	Id           int
	RestName     string
	MapUrl       string
	ReferenceBy  int
	ClosestMetro string
}

func (this Rest) GetDescription(i int) string {
	return fmt.Sprintf(
		"%d) %s: %s, встретимся на станции %s\n",
		i+1,
		this.RestName,
		this.MapUrl,
		this.ClosestMetro,
	)
}

func (this Rest) GetButtonTitle() string {
	return this.RestName
}

func (this Rest) GetCallbackData(eventId int) string {
	return fmt.Sprintf("%s_%d_%d", REST, eventId, this.Id)
}
