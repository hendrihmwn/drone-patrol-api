// This file contains types that are used in the repository layer.
package repository

import "time"

type GetEstateByIdInput struct {
	Id string
}

type Estate struct {
	Id        string    `json:"id" db:"id"`
	Length    int       `json:"length" db:"length"`
	Width     int       `json:"width" db:"width"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type GetTreeByPlot struct {
	X        int
	Y        int
	EstateId string
}

type ListTreesByEstateIdInput struct {
	EstateId string
}

type Tree struct {
	Id        string    `json:"id" db:"id"`
	EstateId  string    `json:"estate_id" db:"estate_id"`
	X         int       `json:"x" db:"x"`
	Y         int       `json:"y" db:"y"`
	Height    int       `json:"height" db:"height"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}
