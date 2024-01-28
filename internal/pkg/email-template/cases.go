package emailtemplate

import (
	"log"
	"strings"

	"github.com/integralist/go-findroot/find"
)

var templatePath = getTemplatePath()

// templates - шаблоны писем.
var templates = map[MailTemplate][]string{
	Verification:     preparePaths("base.layout.tmpl", "footer.partial.tmpl", "verification.mail.tmpl"),
	InvitationLink:   preparePaths("base.layout.tmpl", "footer.partial.tmpl", "invitation-link.mail.tmpl"),
	PasswordRecovery: preparePaths("base.layout.tmpl", "footer.partial.tmpl", "password-recovery.mail.tmpl"),
	NewCourse:        preparePaths("base.layout.tmpl", "footer.partial.tmpl", "new-course.mail.tmpl"),
}

// getTemplatePath возвращает путь к html-шаблонам.
func getTemplatePath() string {
	// Путь к html-шаблонам по умолчанию
	var defaultMigrationPath = "/templates/email"

	// получить получить путь к корню
	rep, err := find.Repo()
	if err != nil {
		log.Printf("cannot get root dir: %s", err.Error())
		log.Println("use default path to mail html-template")
		return defaultMigrationPath
	}

	path := strings.Join([]string{rep.Path, "templates/email"}, "/")

	if strings.EqualFold(rep.Path, "./") || strings.EqualFold(rep.Path, "/") {
		path = strings.Join([]string{rep.Path, "templates/email"}, "")
	}

	return path
}

// preparePaths - формирование путей файлов с шаблонами, необходимых для создания письма.
func preparePaths(files ...string) []string {
	var paths []string
	for _, f := range files {
		paths = append(paths, templatePath+"/"+f)
	}

	return paths
}
