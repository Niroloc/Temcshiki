package data

type Alias string

const (
	USER    Alias = "user"
	REST    Alias = "rest"
	DATE    Alias = "date"
	APPROVE Alias = "approve"
	REVIEW  Alias = "review"
	RATING  Alias = "rating"
)

type FactoryRights struct {
	roleRights map[UserRights]map[Alias]struct{}
}

func CreateFactoryRights() *FactoryRights {
	return &FactoryRights{
		map[UserRights]map[Alias]struct{}{
			ADMIN: {
				USER:    struct{}{},
				REST:    struct{}{},
				DATE:    struct{}{},
				APPROVE: struct{}{},
				REVIEW:  struct{}{},
				RATING:  struct{}{},
			},
			RESERVATOR: {
				REST:    struct{}{},
				DATE:    struct{}{},
				APPROVE: struct{}{},
				REVIEW:  struct{}{},
				RATING:  struct{}{},
			},
			VISITOR: {
				REST:    struct{}{},
				DATE:    struct{}{},
				APPROVE: struct{}{},
				REVIEW:  struct{}{},
				RATING:  struct{}{},
			},
			SPECTATOR: {
				RATING: struct{}{},
			},
		},
	}
}

func (this *FactoryRights) CheckUserRights(alias Alias, user *User) bool {
	var aliases map[Alias]struct{}
	var exists bool
	if aliases, exists = this.roleRights[user.Rights]; !exists {
		return false
	}
	if _, exists = aliases[alias]; !exists {
		return false
	}
	return true
}
