package data

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/Niroloc/Temcshiki/v2/src/logger"
	_ "github.com/mattn/go-sqlite3"
)

type Db struct {
	connection *sql.DB
	logger     *logger.Logger
}

type ExportedUser struct {
	Id       int
	TgId     int
	Username string
	Rights   string
}

type ExportedEvent struct {
	id        int
	visitDate sql.NullString
	restId    sql.NullInt64
	isVisited int
}

type ExportedDate struct {
	id        int
	candidate string
	eventId   int
}

type ExportedVote struct {
	id      int
	userId  int
	eventId int
	restId  sql.NullInt64
	dateId  sql.NullInt64
}

func GetDb(file string) *Db {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		panic(err)
	}
	return &Db{db, logger.GetLogger(reflect.TypeFor[Db]())}
}

func (this *Db) InitDb(migration_file string) {
	query, err := os.ReadFile(migration_file)
	if err != nil {
		this.logger.Error("No initial sql script")
		panic(err)
	}
	if _, err := this.connection.Exec(string(query)); err != nil {
		this.logger.Error("Error occured while execute initial script")
		panic(err)
	}
	res, err := this.connection.Query("SELECT current_stage FROM stage")
	if err != nil {
		this.logger.Error("Stage check failed")
		panic(err)
	}
	if !res.Next() {
		this.logger.Warn("No stage in db, inserting stage CHOOSING")
		this.connection.Exec("INSERT INTO stage (current_stage) VALUES (0)")
	}
	res, err = this.connection.Query("SELECT * FROM users WHERE rights = 'admin'")
	if err != nil {
		this.logger.Error("Error while checking admin user")
		panic(err)
	}
	if !res.Next() {
		this.logger.Warn("No admin user, default admin set")
		defaultAdminId, err := strconv.Atoi(os.Getenv("DEFAULT_ADMIN_ID"))
		if err != nil {
			this.logger.Error("No valid default admin id in envs")
			panic(err)
		}
		defaultAdminName := os.Getenv("DEFAULT_ADMIN_NAME")
		if _, err := this.connection.Exec(
			fmt.Sprintf("INSERT INTO users (tgid, username, rights) VALUES (%d, '%s', '%s')",
				defaultAdminId, defaultAdminName, ADMIN)); err != nil {
			this.logger.Error("Error while inserting default admin")
			panic(err)
		}
	}
	res, err = this.connection.Query("SELECT dt FROM next_task WHERE id = 0")
	if err != nil || !res.Next() {
		this.logger.Warn("No next task date scheduled, inserting next Monday")
		nextTaskDay := getNextWeekday(time.Now(), time.Monday)
		this.logger.Info(nextTaskDay.Format(time.DateOnly))
		if _, err := this.connection.Exec(
			fmt.Sprintf(
				"INSERT INTO next_task (id, dt) VALUES (0, '%s')",
				nextTaskDay.Format(time.DateOnly),
			),
		); err != nil {
			this.logger.Error("Error while setting new task date")
			panic(err)
		}
	}
}

func (this *Db) GetUsers() []ExportedUser {
	rows, err := this.connection.Query("SELECT * FROM users")
	if err != nil {
		this.logger.Warn("An error occured while getting users from DB")
		this.logger.Warn(err.Error())
		return []ExportedUser{}
	}
	var ans []ExportedUser
	for rows.Next() {
		user := ExportedUser{}
		err := rows.Scan(&user.Id, &user.TgId, &user.Username, &user.Rights)
		if err != nil {
			this.logger.Error("An error occured while parsing user")
			this.logger.Error(err.Error())
			continue
		}
		ans = append(ans, user)
	}
	return ans
}

func (this *Db) GetEvents() []ExportedEvent {
	res := []ExportedEvent{}
	rows, err := this.connection.Query("SELECT id, visit_date, rest_id, is_visited FROM events")
	if err != nil {
		this.logger.Warn("An error occured while getting events from DB")
		this.logger.Warn(err.Error())
		return res
	}
	for rows.Next() {
		event := ExportedEvent{}
		err := rows.Scan(&event.id, &event.visitDate, &event.restId, &event.isVisited)
		if err != nil {
			this.logger.Error("An error occured while parsing event")
			this.logger.Error(err.Error())
			continue
		}
		if event.visitDate.Valid {
			event.visitDate.String = strings.Split(event.visitDate.String, "T")[0]
		}
		res = append(res, event)
	}
	return res
}

func (this *Db) GetStage() Stage {
	rows, err := this.connection.Query("SELECT current_stage FROM stage")
	if err != nil || !rows.Next() {
		this.logger.Warn("Cannot get stage from DB")
		this.logger.Warn(err.Error())
		return CHOOSING
	}
	var stage int
	rows.Scan(&stage)
	return Stage(stage)
}

func (this *Db) GetRests() []Rest {
	res := []Rest{}
	rows, err := this.connection.Query("SELECT id, rest_name, map_url, reference_by, closest_metro FROM restoraunts")
	if err != nil {
		this.logger.Error("An error occured while getting restorants")
		this.logger.Error(err.Error())
		return res
	}
	for rows.Next() {
		rest := Rest{}
		err := rows.Scan(&rest.Id, &rest.RestName, &rest.MapUrl, &rest.ReferenceBy, &rest.ClosestMetro)
		if err != nil {
			this.logger.Error("An error occured while parsing restorant")
			this.logger.Error(err.Error())
			continue
		}
		res = append(res, rest)
	}
	return res
}

func (this *Db) GetReviews() []Review {
	res := []Review{}
	rows, err := this.connection.Query(
		"SELECT id, user_id, restoraunt_id, category, rate FROM reviews",
	)
	if err != nil {
		this.logger.Error("An error occured while getting reviews")
		this.logger.Error(err.Error())
		return res
	}
	for rows.Next() {
		review := Review{}
		err := rows.Scan(&review.Id, &review.UserId, &review.RestorauntId, &review.Category, &review.Rate)
		if err != nil {
			this.logger.Error("An error occured while parsing review")
			this.logger.Error(err.Error())
			continue
		}
		res = append(res, review)
	}
	return res
}

func (this *Db) GetDates() []ExportedDate {
	res := []ExportedDate{}
	rows, err := this.connection.Query(
		"SELECT id, candidate, event_id FROM dates",
	)
	if err != nil {
		this.logger.Error("An error occured while getting dates from db")
		this.logger.Error(err.Error())
		return res
	}
	for rows.Next() {
		date := ExportedDate{}
		err := rows.Scan(&date.id, &date.candidate, &date.eventId)
		if err != nil {
			this.logger.Error("An error occured while parsing date")
			this.logger.Error(err.Error())
			continue
		}
		date.candidate = strings.Split(date.candidate, "T")[0]
		res = append(res, date)
	}
	return res
}

func (this *Db) GetVotes() []ExportedVote {
	res := []ExportedVote{}
	rows, err := this.connection.Query(
		"SELECT id, user_id, event_id, rest_id, date_id FROM votes",
	)
	if err != nil {
		this.logger.Error("An error occured while getting votes from db")
		this.logger.Error(err.Error())
		return res
	}
	for rows.Next() {
		vote := ExportedVote{}
		err := rows.Scan(&vote.id, &vote.userId, &vote.eventId, &vote.restId, &vote.dateId)
		if err != nil {
			this.logger.Error("An error occured while parsing vote")
			this.logger.Error(err.Error())
			continue
		}
		res = append(res, vote)
	}
	return res
}

func (this *Db) UpdateStage(prevStage, curStage Stage) {
	_, err := this.connection.Exec(
		fmt.Sprintf("UPDATE stage SET dt = %d WHERE dt = %d",
			curStage,
			prevStage,
		),
	)
	if err != nil {
		this.logger.Warn("Error while updating stage in db, rollback to default")
		this.GetStage()
	}
}

func (this *Db) UpdateNextTaskDate(updatedDate time.Time) {
	_, err := this.connection.Exec(
		fmt.Sprintf(
			"UPDATE next_task SET dt = '%s' WHERE id = 0",
			updatedDate.Format(time.DateOnly),
		),
	)
	if err != nil {
		this.logger.Error("Next task date was not updated")
	}
}

func (this *Db) GetNextTaskDate() string {
	res, err := this.connection.Query(
		"SELECT dt FROM next_task WHERE id = 0",
	)
	if err != nil || !res.Next() {
		this.logger.Error("Cannot get next task date from db")
	}
	result := ""
	res.Scan(&result)
	result = strings.Split(result, "T")[0]
	return result
}

func (this *Db) GetQueuedRests() []Rest {
	rows, err := this.connection.Query(
		"SELECT id, rest_name, map_url, reference_by, closest_metro FROM restoraunts r left join events e on r.id = e.rest_id WHERE e.id IS NULL",
	)
	result := []Rest{}
	if err != nil {
		this.logger.Error("Error while getting rests for voting")
		this.logger.Error(err.Error())
		return result
	}
	for rows.Next() {
		rest := Rest{}
		rows.Scan(&rest.Id, &rest.RestName, &rest.ReferenceBy, &rest.MapUrl, &rest.ClosestMetro)
		result = append(result, rest)
	}
	return result
}

func (this *Db) CreateNewEvent() (Event, error) {
	rows, err := this.connection.Query("INSERT INTO events (visit_date, rest_id, is_visited) VALUES (NULL, NULL, 0) RETURNS id")
	if err != nil {
		this.logger.Error("Error while adding new event into db")
		this.logger.Error(err.Error())
		return Event{}, err
	}
	id := 0
	if rows.Next() {
		rows.Scan(&id)
	} else {
		return Event{}, errors.New("No return value")
	}
	return Event{id, nil, -1, false}, nil
}
