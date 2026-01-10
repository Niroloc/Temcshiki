package data

type ReviewCategory string

const INTERIOR ReviewCategory = "interior"
const SERVICE ReviewCategory = "service"
const FOODR ReviewCategory = "food"
const PRICES ReviewCategory = "prices"

type Review struct {
	Id           int
	UserId       int
	RestorauntId int
	Category     ReviewCategory
	Rate         int
}
