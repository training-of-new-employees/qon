package main

import (
	"flag"
)

type cfg struct {
	Address   string
	Env       string
	Admins    string
	Users     string
	Positions string
	Courses   string
	Lessons   string
}

func InitFlags() *cfg {
	c := cfg{}
	flag.StringVar(&c.Address, "a", "http://localhost:8080", "Адрес сервер, на котором развернуть тестовые данные")
	flag.StringVar(&c.Env, "env", "test-data", "Папка со всеми json-файлами, в которых описаны желаемые сущности")
	flag.StringVar(&c.Admins, "admins", "admins.json", "Данные администраторов")
	flag.StringVar(&c.Courses, "courses", "courses.json", "Описания курсов")
	flag.StringVar(&c.Lessons, "lessons", "lessons.json", "Описания уроков")
	flag.StringVar(&c.Positions, "pos", "positions.json", "Описания должностей")
	flag.StringVar(&c.Users, "users", "users.json", "Данные пользователей")
	flag.Parse()
	return &c
}
