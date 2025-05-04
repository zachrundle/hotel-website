package models

import "time"


// Users is the user model
type User struct {
	Id int
	FirstName string
	LastName string
	Email string
	Password string
	AccessLevel int
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Rooms is the room model
type Room struct {
	Id int
	RoomName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Restrictions is the restriction model
type Restriction struct {
	Id int
	RestrictionName string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Reservations is the reservation model
type Reservation struct {
	Id int
	FirstName string
	LastName string
	Email string
	Phone string
	StartDate time.Time
	EndDate time.Time
	RoomID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
}

// RoomRestrictions is the room restriction model
type RoomRestriction struct {
	Id int
	StartDate time.Time
	EndDate time.Time
	RoomID int
	ReservationID int
	RestrictionID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
	Reservation Reservation
	Restriction Restriction
}