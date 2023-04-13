package rsv

const DefaultLang = "en"

var Locale map[string][]string = map[string][]string{
	"en": {
		"Reservation created",
		"Funds on hold",
		"Reservation confirmed",
		"Reservation dispute",
		"Reservation cancelled",
		"Problem with transaction",
		"Problem with payout",
	},
	"pl": {
		"Nowa rezerwacja",
		"Oczekiwanie na pobranie zapłaty",
		"Rezerwacja potwiedzona",
		"Problem z rezerwacją",
		"Rezerwacja anulowana",
		"Problem z płatnością",
		"Problem z wypłatą",
	},
}

type LocIx int

const NewRsv LocIx = 0
const RsvHold LocIx = 1
const RsvConfirm LocIx = 2
const RsvDispute LocIx = 3
const RsvCancelled LocIx = 4
const RsvFailCapture LocIx = 5
const RsvFailPayout LocIx = 6
