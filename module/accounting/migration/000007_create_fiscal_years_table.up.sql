CREATE TABLE IF NOT EXISTS fiscal_years
(
    id         SERIAL PRIMARY KEY,
    start_date TIMESTAMP WITH TIME ZONE NOT NULL,
    end_date   TIMESTAMP WITH TIME ZONE NOT NULL,
    closed     BOOLEAN DEFAULT FALSE,

    UNIQUE (start_date, end_date)
);
