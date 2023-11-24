package main

import (
	"github.com/training-of-new-employees/qon/internal/app"
	"github.com/training-of-new-employees/qon/internal/config"
)

func main() {
	cfg := config.New()
	
	a := app.New(app.Config{})
	a.Start()
}
