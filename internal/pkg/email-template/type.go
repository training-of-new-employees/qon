package emailtemplate

// MailTemplate - тип для определения html-шаблона.
type MailTemplate string

const (
	// Verification - шаблон письма для верификации
	Verification MailTemplate = "verification"
	// Verification - шаблон письма для пригласительной ссылки
	InvitationLink MailTemplate = "invitation-link"
	// PasswordRecovery - шаблон письма для восстановления пароля
	PasswordRecovery MailTemplate = "password-recovery"
	// PasswordRecovery - шаблон письма для уведомления пользователя о назначенном курсе
	NewCourse MailTemplate = "new-course"
)
