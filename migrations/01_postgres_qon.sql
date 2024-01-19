-- +goose Up

-- +goose StatementBegin

BEGIN;
-- Компании
CREATE TABLE IF NOT EXISTS companies (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(256),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT chck_company_name_not_empty CHECK ( NOT (name IS NULL OR name = '') )
);

-- Должности
CREATE TABLE IF NOT EXISTS positions (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    company_id INTEGER,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(256),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT chck_position_company_not_empty CHECK ( company_id IS NOT NULL ),
    CONSTRAINT fk_position_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT chck_position_name_not_empty CHECK ( NOT (name IS NULL OR name = '') )
);

-- Пользователи
CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    company_id INTEGER,
    position_id INTEGER,
    active BOOLEAN NOT NULL DEFAULT FALSE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    admin BOOLEAN NOT NULL DEFAULT FALSE,
    email VARCHAR(256),
    enc_password VARCHAR(256),
    name VARCHAR(128) NOT NULL DEFAULT '',
    patronymic VARCHAR(128) NOT NULL DEFAULT '',
    surname VARCHAR(128) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT chck_user_company_not_empty CHECK ( company_id IS NOT NULL ),
    CONSTRAINT fk_user_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE,
    CONSTRAINT chck_user_position_not_empty CHECK ( position_id IS NOT NULL ),
    CONSTRAINT fk_user_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE CASCADE,
    CONSTRAINT chck_user_email_not_empty CHECK ( NOT (email IS NULL OR email = '') ),
    CONSTRAINT unq_user_email UNIQUE (email),
    CONSTRAINT chck_user_encpassword_not_empty CHECK ( NOT (enc_password IS NULL OR enc_password = '') )
);

-- Курсы
CREATE TABLE IF NOT EXISTS courses (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_by INTEGER,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    name VARCHAR(256),
    description VARCHAR(512) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT chck_course_creater_not_empty CHECK ( created_by IS NOT NULL ),
    CONSTRAINT fk_course_creater FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chck_course_name_not_empty CHECK ( NOT (name IS NULL OR name = '') )
);

-- Уроки
CREATE TABLE IF NOT EXISTS lessons (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    course_id INTEGER,
    created_by INTEGER,
    active BOOLEAN NOT NULL DEFAULT TRUE,
    archived BOOLEAN NOT NULL DEFAULT FALSE,
    number INTEGER,
    name VARCHAR(256),
    description VARCHAR(512) NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT chck_lesson_course_not_empty CHECK ( course_id IS NOT NULL AND course_id <> 0 ),
    CONSTRAINT fk_lesson_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT chck_lesson_creater_not_empty CHECK ( created_by IS NOT NULL AND created_by <> 0 ),
    CONSTRAINT fk_lesson_creater FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chck_lesson_name_not_empty CHECK ( NOT (name IS NULL OR name = '') )
);

-- Тексты
CREATE TABLE IF NOT EXISTS texts (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    lesson_id INTEGER,
    created_by INTEGER,
    number INTEGER,
    header VARCHAR(512) DEFAULT '',
    content TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT chck_text_lesson_not_empty CHECK ( lesson_id IS NOT NULL ),
    CONSTRAINT fk_text_lesson FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
    CONSTRAINT chck_text_creater_not_empty CHECK ( created_by IS NOT NULL ),
    CONSTRAINT fk_text_creater FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chck_text_content_not_empty CHECK ( NOT (content IS NULL OR content = '') )
);

-- Картинки
CREATE TABLE IF NOT EXISTS pictures (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    lesson_id INTEGER,
    created_by INTEGER,
    number INTEGER,
    name VARCHAR(256) DEFAULT '',
    url_picture VARCHAR(1024),
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT chck_picture_lesson_not_empty CHECK ( lesson_id IS NOT NULL ),
    CONSTRAINT fk_picture_lesson FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
    CONSTRAINT chck_picture_creater_not_empty CHECK ( created_by IS NOT NULL ),
    CONSTRAINT fk_picture_creater FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chck_url_picture_not_empty CHECK ( NOT (url_picture IS NULL OR url_picture = '') )
);

-- Должности-Курсы (назначенные на должность курсы)
CREATE TABLE IF NOT EXISTS position_course (
    id INTEGER PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    position_id INTEGER,
    course_id INTEGER,
    CONSTRAINT chck_positioncourse_position_not_empty CHECK ( position_id IS NOT NULL ),
    CONSTRAINT fk_positioncourse_position FOREIGN KEY (position_id) REFERENCES positions(id) ON DELETE CASCADE,
    CONSTRAINT chck_positioncourse_course_not_empty CHECK ( course_id IS NOT NULL ),
    CONSTRAINT fk_positioncourse_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT unq_positioncourse UNIQUE (position_id, course_id)
);

-- Назначенные сотрудникам курсы (или прогресс сотрудников по курсам)
CREATE TABLE IF NOT EXISTS course_assign (
    user_id INTEGER,
    course_id INTEGER,
    status VARCHAR(16) NOT NULL DEFAULT 'not-started',
    started_at TIMESTAMP NOT NULL DEFAULT now(),
    finished_at TIMESTAMP NOT NULL DEFAULT now(),
    CONSTRAINT chck_courseassign_user_not_empty CHECK ( user_id IS NOT NULL ),
    CONSTRAINT fk_courseassign_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT chck_courseassign_course_not_empty CHECK ( course_id IS NOT NULL ),
    CONSTRAINT fk_courseassign_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT unq_assigncourse UNIQUE (course_id, user_id),
    CONSTRAINT chck_course_status_type CHECK (
        status IN ('not-started', 'in-process', 'done')
    )
);

-- Прогресс сотрудников по урокам
CREATE TABLE IF NOT EXISTS lesson_results (
    user_id   INTEGER,
    course_id INTEGER,
    lesson_id INTEGER,
    status VARCHAR(16) NOT NULL DEFAULT 'not-started',
    CONSTRAINT chck_lessonresult_course_not_empty CHECK ( course_id IS NOT NULL ),
    CONSTRAINT fk_lessonresult_course FOREIGN KEY (course_id) REFERENCES courses(id) ON DELETE CASCADE,
    CONSTRAINT chck_lessonresult_lesson_not_empty CHECK ( lesson_id IS NOT NULL ),
    CONSTRAINT fk_lessonresult_lesson FOREIGN KEY (lesson_id) REFERENCES lessons(id) ON DELETE CASCADE,
    CONSTRAINT chck_lessonresult_user_not_empty CHECK ( user_id IS NOT NULL ),
    CONSTRAINT fk_lessonresult_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT unq_assignlesson UNIQUE (course_id, lesson_id, user_id),
    CONSTRAINT chck_lesson_status_type CHECK (
        status IN ('not-started', 'in-process', 'done')
    )
);

CREATE INDEX IF NOT EXISTS position_course_idx ON position_course (position_id, course_id);
CREATE INDEX IF NOT EXISTS course_assign_idx ON course_assign (user_id, course_id);
CREATE INDEX IF NOT EXISTS lesson_result_idx ON lesson_results (user_id, course_id, lesson_id);

COMMIT;
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

DROP INDEX IF EXISTS position_course_idx;

DROP INDEX IF EXISTS course_assign_idx;

DROP INDEX IF EXISTS lesson_result_idx;
-- +goose StatementEnd
