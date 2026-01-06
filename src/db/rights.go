package db

type UserRights string

const ADMIN UserRights = "admin"
const RESERVATOR UserRights = "reservator"
const VISITOR UserRights = "visitor"
const SPECTATOR UserRights = "spectator"
