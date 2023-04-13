package main

import (
	"math/rand"
	"sport/helpers"
	"sport/instr"
	"sport/user"

	"github.com/google/uuid"
)

var n = []string{
	"Isaac",
	"Richard",
	"Andrzej",
	"Max",
	"Jeanette",
	"V",
	"Enrico",
	"Mike",
	"Roger",
	"Łukasz",
	"Kamil",
	"Katherine",
	"Sheldon",
	"James",
	"Clerk",
}

var ln = []string{
	"Maxwell",
	"Dirac",
	"Fermi",
	"Faraday",
	"Planck",
	"Gauss",
	"Euler",
	"Schrodinger",
	"Bohr",
	"Rutheford",
	"Feynman",
	"Newton",
	"Bohr",
	"Rutheford",
	"Maxwell",
	"Dirac",
	"Fermi",
	"Faraday",
	"Planck",
	"Gauss",
	"Euler",
	"Schrodinger",
	"Gołota",
	"Tyson",
	"Federer",
	"Widera",
	"Zagórski",
	"Maulwurf",
	"Cooper",
}

var d = []string{
	"gmail.com",
	"protonmail.onion",
	"interia.pl",
	"wp.pl",
	"bazinga.com",
}

type TokenWithInstrID struct {
	Token string
	IID   uuid.UUID
}

func InstrGen(instrCtx *instr.Ctx, maxTh, instrNo int) ([]TokenWithInstrID, error) {

	iids := make([]TokenWithInstrID, instrNo)

	err := helpers.Sem(maxTh, instrNo, []func(int) error{
		func(i int) error {
			name := n[rand.Int()%len(n)]
			name2 := ln[rand.Int()%len(ln)]
			domain := d[rand.Int()%len(d)]
			randomizer := uuid.New().String()
			password := "password"
			token, err := instrCtx.User.ApiCreateAndLoginUser(&user.UserRequest{
				Email:    name + name2 + randomizer + "@" + domain,
				Password: password,
				UserData: user.UserData{
					Name:     name + " " + name2 + randomizer,
					Language: "pl",
					Country:  "PL",
					AboutMe:  "long text that describes user, just to catch more border cases, long text that describes user",
				},
			})
			if err != nil {
				return err
			}

			userID, err := instrCtx.Api.AuthorizeUserFromToken(token)
			if err != nil {
				return err
			}

			iid, err := instrCtx.CreateTestInstructor(userID, nil)
			if err != nil {
				return err
			}

			iids[i] = TokenWithInstrID{
				Token: token,
				IID:   iid,
			}

			return err
		},
	})

	return iids, err
}
