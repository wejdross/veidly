package user

const DefaultLang = "en"

var Locale map[string][]string = map[string][]string{
	"en": {
		"Veidly - password recovery",
		"Veidly - registration",
	},
	"pl": {
		"Veidly - odzyskiwanie hasła",
		"Veidly - rejestracja",
	},
}

type LocIx int

const PassRecoveryTitle LocIx = 0
const RegisterTitle LocIx = 1
