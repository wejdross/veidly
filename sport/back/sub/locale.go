package sub

const DefaultLang = "en"

var Locale map[string][]string = map[string][]string{
	"en": {
		"Subscription created",
		"Subscription confirmed",
		"Subscription dispute",
		"Subscription cancelled",
		"Problem with transaction",
		"Problem with payout",
	},
	"pl": {
		"Nowy karnet",
		"Karnet potwiedzony",
		"Problem z karnetem",
		"Karnet anulowany",
		"Problem z płatnością",
		"Problem z wypłatą",
	},
}

type LocIx int

const NewSub LocIx = 0
const SubConfirm LocIx = 1
const SubDispute LocIx = 2
const SubCancelled LocIx = 3
const SubFailCapture LocIx = 4
const SubFailPayout LocIx = 5
