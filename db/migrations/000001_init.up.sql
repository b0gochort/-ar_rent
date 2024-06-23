CREATE TABLE cars (
    id SERIAL PRIMARY KEY ,
    make TEXT NOT NULL,
    model TEXT NOT NULL,
    year INTEGER NOT NULL,
    free_status boolean NOT NULL
);

CREATE TABLE rental_sessions (
    id SERIAL PRIMARY KEY ,
    car_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    total_price INTEGER NOT NULL,
    status BOOLEAN NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (car_id) REFERENCES cars (id)

);

INSERT INTO cars (make, model, year,free_status) VALUES
    ('Lada', 'Granta', 2020,true),
    ('Lada', 'Vesta', 2019,true),
    ('GAZ', 'Volga', 2005,true),
    ('UAZ', 'Patriot', 2021,true),
    ('Lada', 'Xray', 2022,true);
