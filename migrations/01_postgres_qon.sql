-- +goose Up

-- +goose StatementBegin

-- Компании
CREATE TABLE IF NOT EXISTS companies (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name VARCHAR(256) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

-- Должности
CREATE TABLE IF NOT EXISTS positions (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    company_id INTEGER NOT NULL,
    name VARCHAR(256) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_position_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE
);

-- Пользователи
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    company_id INTEGER NOT NULL,
    position_id INTEGER NOT NULL,
    email VARCHAR(256) NOT NULL,
    enc_password VARCHAR(256) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    admin BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(128) NOT NULL,
    surname VARCHAR(128) NOT NULL,
    patronymic VARCHAR(128),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_user_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE CASCADE,
    CONSTRAINT unq_user_email UNIQUE (email)
);

-- Курсы
CREATE TABLE IF NOT EXISTS courses (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_by INTEGER NOT NULL,
    name VARCHAR(256) NOT NULL,
    description VARCHAR(256) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_course_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Уроки
CREATE TABLE IF NOT EXISTS lessons (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    course_id INTEGER NOT NULL,
    created_by INT NOT NULL,
    number INTEGER,
    name VARCHAR(256) NOT NULL,
    description VARCHAR(256) NOT NULL,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_lesson_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT fk_lesson_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Тексты
CREATE TABLE IF NOT EXISTS texts (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    lesson_id INTEGER NOT NULL,
    created_by INT NOT NULL,
    number INT,
    header VARCHAR(256) NOT NULL,
    content TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_text_lesson FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
    CONSTRAINT fk_text_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Картинки
CREATE TABLE IF NOT EXISTS pictures (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    lesson_id INTEGER NOT NULL,
    created_by INT NOT NULL,
    number INT,
    name VARCHAR(256) NOT NULL,
    link VARCHAR(1024) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_picture_lesson FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
    CONSTRAINT fk_picture_user FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE
);

-- Должности-Курсы (назначенные на должность курсы)
CREATE TABLE IF NOT EXISTS position_course (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    position_id INTEGER NOT NULL,
    course_id INTEGER NOT NULL,
    CONSTRAINT fk_positioncourse_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE CASCADE,
    CONSTRAINT fk_positioncourse_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT unq_positioncourse UNIQUE (position_id, course_id)
);

-- Назначенные сотрудникам курсы (или прогресс сотрудников по курсам)
CREATE TABLE IF NOT EXISTS course_assign (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id INT NOT NULL,
    course_id INT NOT NULL,
    pass_course BOOLEAN NOT NULL DEFAULT false,
    started_at TIMESTAMP NOT NULL DEFAULT now(),
    finished_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT fk_courseassign_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_courseassign_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT unq_usercourse UNIQUE (user_id, course_id)
);

-- Прогресс сотрудников по урокам
CREATE TABLE IF NOT EXISTS lesson_results (
    assign_id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    lesson_id INTEGER NOT NULL,
    pass_lesson BOOLEAN NOT NULL DEFAULT false,
    CONSTRAINT fk_lessonresult_courseassign FOREIGN KEY (assign_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT fk_lessonresult_lesson FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
    CONSTRAINT unq_assignlesson UNIQUE (assign_id, lesson_id)  
);

-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS lesson_results;

DROP TABLE IF EXISTS course_assign;

DROP TABLE IF EXISTS position_course;

DROP TABLE IF EXISTS pictures;

DROP TABLE IF EXISTS texts;

DROP TABLE IF EXISTS lessons;

DROP TABLE IF EXISTS courses;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS positions;

DROP TABLE IF EXISTS companies;

-- +goose StatementEnd