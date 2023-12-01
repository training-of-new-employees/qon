-- +goose Up

-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS companies (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(256) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS positions (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    company_id INTEGER NOT NULL,
    name VARCHAR(256) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_company_position FOREIGN KEY(company_id) REFERENCES companies(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    company_id INTEGER NOT NULL,
    position_id INTEGER NOT NULL,
    email VARCHAR(256) NOT NULL,
    enc_password VARCHAR(256) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    admin BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(128) NOT NULL,
    surname VARCHAR(128) NOT NULL,
    patronymic VARCHAR(128),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_company_user FOREIGN KEY(company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_position_user FOREIGN KEY(position_id) REFERENCES positions(id) ON DELETE CASCADE
);

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS positions;

DROP TABLE IF EXISTS companies;

-- +goose StatementEnd