CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password BYTEA NOT NULL,
    role TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS apps(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    secret TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS tataisk (
    id SERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    year_of_prod INTEGER NOT NULL CHECK(year_of_prod >= 1900),
    imdb FLOAT NOT NULL CHECK(imdb > 0 AND imdb <= 10),
    description TEXT NOT NULL,
    country TEXT[] NOT NULL,
    genre TEXT[] NOT NULL,
    film_director TEXT NOT NULL,
    screenwriter TEXT NOT NULL,
    budget INTEGER NOT NULL,
    collection INTEGER NOT NULL
    );
