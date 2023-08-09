CREATE TYPE grade AS ENUM ('trainee', 'junior', 'middle', 'senior');

CREATE TABLE IF NOT EXISTS usr (
                                id          SERIAL PRIMARY KEY,
                                name        TEXT NOT NULL
                                surname     TEXT NOT NULL
                                position    grade NOT NULL
                                project     TEXT
);