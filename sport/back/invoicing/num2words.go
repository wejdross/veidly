package invoicing

import (
	"fmt"
	"strings"
)

var polishMegas = [][]string{
	{"", "", ""},
	{"tysiąc", "tysiące", "tysięcy"},
	{"milion", "miliony", "milionów"},
	{"miliard", "miliardy", "miliardów"},
	{"bilion", "biliony", "bilionów"},
	{"biliard", "biliardy", "biliardów"},
	{"trylion", "tryliony", "trylionów"},
	{"tryliard", "tryliardy", "tryliardów"},
	{"kwadrylion", "kwadryliony", "kwadrylionów"},
	{"kwintylion", "kwintyliony", "kwintylionów"},
	{"sekstylion", "sekstyliony", "sekstylionów"},
	{"septylion", "septyliony", "septylionów"},
	{"oktylion", "oktyliony", "oktylionów"},
	{"nonylion", "nonyliony", "nonylionów"},
	{"decylion", "decyliony", "decylionów"}}
var polishUnits = []string{"", "jeden", "dwa", "trzy", "cztery", "pięć", "sześć", "siedem", "osiem", "dziewięć"}
var polishTens = []string{"", "dziesięć", "dwadzieścia", "trzydzieści", "czterdzieści", "pięćdziesiąt", "sześćdziesiąt", "siedemdziesiąt", "osiemdziesiąt", "dziewięćdziesiąt"}
var polishTeens = []string{"dziesięć", "jedenaście", "dwanaście", "trzynaście", "czternaście", "piętnaście", "szesnaście", "siedemnaście", "osiemnaście", "dziewiętnaście"}
var polishHundreds = []string{"", "sto", "dwieście", "trzysta", "czterysta", "pięćset", "sześćset", "siedemset", "osiemset", "dziewięćset"}

func IntToTriplets(number int) []int {
	triplets := []int{}
	for number > 0 {
		triplets = append(triplets, number%1000)
		number = number / 1000
	}
	return triplets
}

func PLNToWords(amount int) string {
	rem := amount % 10
	if rem == 0 || rem >= 5 {
		return "złotych"
	}
	if rem == 1 {
		return "złoty"
	}
	return "złote"
}

func IntToPlWords(val int) string {

	words := []string{}

	triplets := IntToTriplets(val)

	if len(triplets) == 0 {
		return "zero"
	}

	for i := len(triplets) - 1; i >= 0; i-- {
		t := triplets[i]
		if t == 0 {
			continue
		}
		hundreds := t / 100 % 10
		tens := t / 10 % 10
		units := t % 10

		if hundreds > 0 {
			words = append(words, polishHundreds[hundreds])
		}

		if tens == 0 && units == 0 {
			goto tripletEnd
		}

		switch tens {
		case 0:
			words = append(words, polishUnits[units])
		case 1:
			words = append(words, polishTeens[units])
		default:
			if units > 0 {
				word := fmt.Sprintf("%s %s", polishTens[tens], polishUnits[units])
				words = append(words, word)
			} else {
				words = append(words, polishTens[tens])
			}
		}

	tripletEnd:
		if t == 1 {
			words = append(words, polishMegas[i][0])
			continue
		}

		megaIndex := 2
		if units >= 2 && units <= 4 {
			megaIndex = 1
		}

		if mega := polishMegas[i][megaIndex]; mega != "" {
			words = append(words, mega)
		}
	}

	return strings.TrimSpace(strings.Join(words, " "))
}
