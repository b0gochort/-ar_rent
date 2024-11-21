package model

import "time"

type Car struct {
	Id         int64  `json:"id" db:"id"`
	Make       string `json:"make" db:"make"`
	Model      string `json:"model" db:"model"`
	Year       int64  `json:"year" db:"year"`
	FreeStatus bool   `json:"free_status" db:"free_status"`
}

type RentalSession struct {
	ID         int64     `json:"id" db:"id"`
	CarID      int64     `json:"car_id" db:"car_id"`
	UserID     int64     `json:"user_id" db:"user_id"`
	StartDate  time.Time `json:"start_date" db:"start_date"`
	EndDate    time.Time `json:"end_date" db:"end_date"`
	TotalPrice int       `json:"total_price" db:"total_price"`
	Status     bool      `json:"status" db:"status"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type CreateSessionReq struct {
	Price    int   `json:"price"`
	UserId   int64 `json:"user_id"`
	CarId    int64 `json:"car_id"`
	RentDays int   `json:"rent_days"`
}

type UtilizationReport struct {
	CarID                 *int    `json:"car_id" db:"id"`
	Make                  string  `json:"make" db:"make"`
	Model                 string  `json:"model" db:"model"`
	DaysRented            int     `json:"days_rented" db:"days_rented"`
	TotalDays             int     `json:"total_days" db:"total_days"`
	UtilizationPercentage float64 `json:"utilization_percentage" db:"utilization_percentage"`
}

type GetPriceResponse struct {
	Price int `json:"price"`
}

type GetCarsResponse struct {
	Cars []Car `json:"cars"`
}
