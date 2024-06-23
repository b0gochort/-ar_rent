package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/b0gochort/car-rent/internal/repository/postgres"
	"github.com/b0gochort/car-rent/model"
	"github.com/jmoiron/sqlx"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

const (
	baseRate       = 1000
	maxRentalDays  = 30
	discount5Days  = 0.05
	discount10Days = 0.10
	discount18Days = 0.15
)

func CheckAuto(db *sqlx.DB, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var carId string

		carId = r.URL.Query().Get("car_id")

		// Update car statuses after 3 days from rental
		err := postgres.UpdateCarsStatus(db)
		if err != nil {
			log.Info("postgres.UpdateCarsStatus: ", err.Error())
		}

		status, err := postgres.GetFreeStatusCar(carId, db)
		if err != nil {
			log.Info("postgres.GetFreeStatusCar err: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
		if !status {
			w.WriteHeader(http.StatusOK)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte("false"))

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("true"))
	}
}

func GetPrice(db *sqlx.DB, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			rentDays string
			res      model.GetPriceResponse
		)

		rentDays = r.URL.Query().Get("rent_days")

		price, err := calculateRentalCost(rentDays)
		if err != nil {
			log.Info("calculateRentalCost err: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		res.Price = price

		response, err := json.Marshal(res)
		if err != nil {
			log.Info("failed to marshal car utilization report: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func CreateRentSesion(db *sqlx.DB, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req model.CreateSessionReq

		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			log.Info("createSessionReq err: ", err)
			w.WriteHeader(http.StatusUnprocessableEntity)

			return
		}

		endDate, err := postgres.GetDateLastSession(req.CarId, db)
		if err != nil {
			log.Info("getDateLastSession err: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}
		//if 3 days have not passed then you cannot book
		if !endDate.IsZero() {
			if passed := time.Since(endDate) >= 72*time.Hour; !passed {
				w.WriteHeader(http.StatusUnprocessableEntity)
				w.Write([]byte("false"))

				return
			}
		}

		newRentSession := model.RentalSession{
			CarID:      req.CarId,
			UserID:     req.UserId,
			StartDate:  time.Now(),
			EndDate:    time.Now().AddDate(0, 0, req.RentDays),
			TotalPrice: req.Price,
			Status:     false,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		tx, err := db.Beginx()
		if err != nil {
			log.Info("db.Beginx err: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if err := postgres.CreateRentSession(newRentSession, tx); err != nil {
			RollbackOrCatchError(tx, log)
			log.Info("postgres.CreateRentSession err: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		if err := postgres.SetRentStatus(req.CarId, tx); err != nil {
			RollbackOrCatchError(tx, log)
			log.Info("postgres.SetFreeStatus err: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		CommitOrCatchError(tx, log)

		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("true"))
	}
}

func GetCarUtilizationReport(db *sqlx.DB, log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		monthQuery := r.URL.Query().Get("month")
		yearQuery := r.URL.Query().Get("year")
		month, err := strconv.Atoi(monthQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		year, err := strconv.Atoi(yearQuery)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)

			return
		}

		if month == 0 || year == 0 {
			w.WriteHeader(http.StatusUnprocessableEntity)

			return
		}

		reports, err := postgres.GetUtilizationReport(month, year, db)
		if err != nil {
			log.Info("postgres.GetUtilizationReport err: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		response, err := json.Marshal(reports)
		if err != nil {
			log.Info("failed to marshal car utilization report: ", err)
			w.WriteHeader(http.StatusInternalServerError)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

func calculateRentalCost(rentDays string) (int, error) {

	days, err := strconv.Atoi(rentDays)
	if err != nil {
		return 0, fmt.Errorf("converting %s to int: %w", rentDays, err)
	}

	if days <= 0 {
		return 0, errors.New("must be greater than zero")
	}

	if days > maxRentalDays {
		return 0, errors.New("must be less than or equal to max rental days (30)")
	}

	totalCost := 0

	for i := 1; i <= days; i++ {
		rate := baseRate
		switch {
		case i >= 18:
			rate -= int(float64(baseRate) * discount18Days)
		case i >= 10:
			rate -= int(float64(baseRate) * discount10Days)
		case i >= 5:
			rate -= int(float64(baseRate) * discount5Days)
		}
		totalCost += rate
	}

	return totalCost, nil
}

func RollbackOrCatchError(tx *sqlx.Tx, log *slog.Logger) {
	if err := tx.Rollback(); err != nil {
		log.Info("service - unable to rollback transaction: ", err.Error())
	}

	log.Info("service - rollback transaction ...")
}

func CommitOrCatchError(tx *sqlx.Tx, log *slog.Logger) error {
	if err := tx.Commit(); err != nil {
		log.Info("service - unable to commit transaction:", err.Error())
		return err
	}
	log.Info("service - commit transaction ...")

	return nil
}
