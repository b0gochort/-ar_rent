package postgres

import (
	"database/sql"
	"fmt"
	"github.com/b0gochort/car-rent/model"
	"github.com/jmoiron/sqlx"
	"time"
)

func GetFreeStatusCar(id string, db *sqlx.DB) (bool, error) {
	var status bool

	query := `SELECT free_status FROM cars WHERE id = $1`
	if err := db.QueryRow(query, id).Scan(&status); err != nil {
		return false, fmt.Errorf("db.QueryRow: %w", err)
	}

	return status, nil
}

func CreateRentSession(newSession model.RentalSession, tx *sqlx.Tx) error {
	query := `
		INSERT INTO rental_sessions 
		(car_id, user_id, start_date, end_date, total_price, status, created_at, updated_at) 
		VALUES 
		(:car_id, :user_id, :start_date, :end_date, :total_price, :status, :created_at, :updated_at)
	`
	_, err := tx.NamedExec(query, newSession)
	if err != nil {
		return err
	}
	return nil
}

func GetUtilizationReport(month, year int, db *sqlx.DB) ([]model.UtilizationReport, error) {
	startDate := fmt.Sprintf("%04d-%02d-01", year, month)
	query := fmt.Sprintf(`
	WITH days_in_month AS (
	    SELECT 
	        generate_series(
	            date_trunc('month', DATE '%s'),
	            date_trunc('month', DATE '%s') + interval '1 month' - interval '1 day',
	            interval '1 day'
	        )::date AS day
	),
	rental_days AS (
	    SELECT
	        car_id,
	        COUNT(DISTINCT day) AS days_rented
	    FROM
	        rental_sessions
	    JOIN
	        days_in_month ON day BETWEEN start_date AND end_date
	    WHERE
	        (start_date >= date_trunc('month', DATE '%s') AND start_date < date_trunc('month', DATE '%s') + interval '1 month')
	        OR
	        (end_date >= date_trunc('month', DATE '%s') AND end_date < date_trunc('month', DATE '%s') + interval '1 month')
	    GROUP BY
	        car_id
	),
	total_days_in_month AS (
	    SELECT
	        COUNT(*) AS total_days
	    FROM
	        days_in_month
	)
	SELECT
	    cars.id,
	    cars.make,
	    cars.model,
	    COALESCE(rental_days.days_rented, 0) AS days_rented,
	    total_days_in_month.total_days,
	    (COALESCE(rental_days.days_rented, 0)::decimal / total_days_in_month.total_days) * 100 AS utilization_percentage
	FROM
	    cars
	LEFT JOIN
	    rental_days ON cars.id = rental_days.car_id,
	    total_days_in_month
	UNION ALL
	SELECT
	    NULL AS id,
	    'All' AS make,
	    'Cars' AS model,
	    0 AS days_rented,
	    total_days_in_month.total_days,
	    (SUM(COALESCE(rental_days.days_rented, 0))::decimal / (total_days_in_month.total_days * COUNT(cars.id))) * 100 AS utilization_percentage
	FROM
	    cars
	LEFT JOIN
	    rental_days ON cars.id = rental_days.car_id,
	    total_days_in_month
	GROUP BY
	    total_days_in_month.total_days;
	`, startDate, startDate, startDate, startDate, startDate, startDate)

	rows, err := db.Queryx(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []model.UtilizationReport
	for rows.Next() {
		var report model.UtilizationReport
		if err := rows.StructScan(&report); err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func GetDateLastSession(carId int64, db *sqlx.DB) (time.Time, error) {
	var endDate time.Time
	query := `
        SELECT end_date
        FROM rental_sessions
        WHERE car_id = $1
        ORDER BY end_date DESC
        LIMIT 1;
        `

	err := db.Get(&endDate, query, carId)
	if err != nil {
		if err == sql.ErrNoRows {
			endDate = time.Time{}
		} else {
			return time.Time{}, err
		}
	}

	return endDate, nil
}

func SetRentStatus(carId int64, db *sqlx.Tx) error {
	query := `
        UPDATE cars
		SET free_status = FALSE
		WHERE id = $1;
        `

	_, err := db.Query(query, carId)
	if err != nil {
		return err
	}

	return nil
}

func UpdateCarsStatus(db *sqlx.DB) error {
	query := `
	UPDATE cars
	SET free_status = TRUE
	WHERE id IN (
		SELECT car_id
		FROM rental_sessions
		WHERE end_date = CURRENT_DATE - INTERVAL '3 day'
	)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil
}
