package chat_test

import (
	"encoding/json"
	"net/url"
	"sport/api"
	"sport/chat"
	"sport/chat_integrator"
	"sport/helpers"
	"sport/user"
	"testing"
	"time"
)

func TestCreateAnonChatromIntegrator(t *testing.T) {

	ur2 := user.ApiTestUserRequest()
	ur2.Name = helpers.CRNG_stringPanic(20)
	user2VeidlyToken, err := userCtx.ApiCreateAndLoginUser(&ur2)
	if err != nil {
		t.Fatal(err)
	}

	user2ID, err := userCtx.Api.AuthorizeUserFromToken(user2VeidlyToken)
	if err != nil {
		t.Fatal(err)
	}

	invalidReq := chat_integrator.CreateUserChatroomReq{
		PeerUserID: user2ID,
	}

	req := chat_integrator.CreateUserChatroomReq{
		PeerUserID: user2ID,
		AnonData: &chat_integrator.AnonData{
			DisplayName: helpers.CRNG_stringPanic(20),
			Email:       helpers.CRNG_EmailOrPanic(),
		},
	}

	m := "POST"
	u := "/api/chat_integrator/room"
	tm := "GET"
	tu := "/api/chat_integrator/token"

	var res chat_integrator.CreateUserChatroomRes

	var user2ChatTokenRes chat.TokenResponse

	userCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			RequestReader:      helpers.JsonMustSerializeReader(invalidReq),
			ExpectedStatusCode: 400,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			RequestReader:      helpers.JsonMustSerializeReader(invalidReq),
			ExpectedStatusCode: 400,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			RequestReader:      helpers.JsonMustSerializeReader(req),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &res)
			},
		}, {
			RequestMethod:  tm,
			RequestUrl:     tu,
			RequestHeaders: user.GetAuthHeader(user2VeidlyToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &user2ChatTokenRes)
			},
		},
	})

	userID, err := chatIntegratorCtx.Chat.AuthorizeToken(res.AccessToken)
	if err != nil {
		t.Fatal(err)
	}

	userCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:  "GET",
			RequestUrl:     "/api/chat/room",
			RequestHeaders: user.GetAuthHeader(res.AccessToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var members []chat.ChatRoomMember
				if err := json.Unmarshal(b, &members); err != nil {
					return err
				}
				if len(members) != 1 {
					return err
				}
				e := members[0]
				return helpers.AssertJsonErr("chatoom member for user 1", e, chat.ChatRoomMember{
					ChatRoomID:   res.ChatRoomID,
					ServerID:     "",
					UserID:       userID,
					LastNotified: time.Time{},
					Data: chat.MemberData{
						Email:        req.AnonData.Email,
						ChatRoomName: ur2.Name,
						MemberPubData: chat.MemberPubData{
							DisplayName: req.AnonData.DisplayName,
							IconRelpath: "",
						},
					},
				})
			},
		}, {
			RequestMethod:  "GET",
			RequestUrl:     "/api/chat/room",
			RequestHeaders: user.GetAuthHeader(user2ChatTokenRes.Token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var members []chat.ChatRoomMember
				if err := json.Unmarshal(b, &members); err != nil {
					return err
				}
				if len(members) != 1 {
					return err
				}
				e := members[0]
				return helpers.AssertJsonErr("chatoom member for user 2", e, chat.ChatRoomMember{
					ChatRoomID:   res.ChatRoomID,
					ServerID:     "",
					UserID:       user2ID,
					LastNotified: time.Time{},
					Data: chat.MemberData{
						Email:        ur2.Email,
						ChatRoomName: req.AnonData.DisplayName,
						MemberPubData: chat.MemberPubData{
							DisplayName: ur2.Name,
							IconRelpath: "",
						},
					},
				})
			},
		},
	})

}

func TestCreateChatromIntegrator(t *testing.T) {
	ur := user.ApiTestUserRequest()
	ur.Name = helpers.CRNG_stringPanic(20)
	user1VeidlyToken, err := userCtx.ApiCreateAndLoginUser(&ur)
	if err != nil {
		t.Fatal(err)
	}
	user1ID, err := userCtx.Api.AuthorizeUserFromToken(user1VeidlyToken)
	if err != nil {
		t.Fatal(err)
	}

	ur2 := user.ApiTestUserRequest()
	ur2.Name = helpers.CRNG_stringPanic(20)
	user2VeidlyToken, err := userCtx.ApiCreateAndLoginUser(&ur2)
	if err != nil {
		t.Fatal(err)
	}

	user2ID, err := userCtx.Api.AuthorizeUserFromToken(user2VeidlyToken)
	if err != nil {
		t.Fatal(err)
	}

	req := chat_integrator.CreateUserChatroomReq{
		PeerUserID: user2ID,
	}

	m := "POST"
	u := "/api/chat_integrator/room"
	tm := "GET"
	tu := "/api/chat_integrator/token"

	var res chat_integrator.CreateUserChatroomRes

	var user1ChatTokenRes, user2ChatTokenRes chat.TokenResponse

	userCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 400,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			RequestReader:      helpers.JsonMustSerializeReader(req),
			ExpectedStatusCode: 400,
		}, {
			RequestMethod:      m,
			RequestUrl:         u,
			RequestHeaders:     user.GetAuthHeader(user1VeidlyToken),
			RequestReader:      helpers.JsonMustSerializeReader(req),
			ExpectedStatusCode: 200,
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &res)
			},
		}, {
			RequestMethod:  tm,
			RequestUrl:     tu,
			RequestHeaders: user.GetAuthHeader(user1VeidlyToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &user1ChatTokenRes)
			},
		}, {
			RequestMethod:  tm,
			RequestUrl:     tu,
			RequestHeaders: user.GetAuthHeader(user2VeidlyToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				return json.Unmarshal(b, &user2ChatTokenRes)
			},
		},
	})

	userCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:  "GET",
			RequestUrl:     "/api/chat/room",
			RequestHeaders: user.GetAuthHeader(user1ChatTokenRes.Token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var members []chat.ChatRoomMember
				if err := json.Unmarshal(b, &members); err != nil {
					return err
				}
				if len(members) != 1 {
					return err
				}
				e := members[0]
				return helpers.AssertJsonErr("chatoom member for user 1", e, chat.ChatRoomMember{
					ChatRoomID:   res.ChatRoomID,
					ServerID:     "",
					UserID:       user1ID,
					LastNotified: time.Time{},
					Data: chat.MemberData{
						Email:        ur.Email,
						ChatRoomName: ur2.Name,
						MemberPubData: chat.MemberPubData{
							DisplayName: ur.Name,
							IconRelpath: "",
						},
					},
				})
			},
		}, {
			RequestMethod:  "GET",
			RequestUrl:     "/api/chat/room",
			RequestHeaders: user.GetAuthHeader(user2ChatTokenRes.Token),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				var members []chat.ChatRoomMember
				if err := json.Unmarshal(b, &members); err != nil {
					return err
				}
				if len(members) != 1 {
					return err
				}
				e := members[0]
				return helpers.AssertJsonErr("chatoom member for user 2", e, chat.ChatRoomMember{
					ChatRoomID:   res.ChatRoomID,
					ServerID:     "",
					UserID:       user2ID,
					LastNotified: time.Time{},
					Data: chat.MemberData{
						Email:        ur2.Email,
						ChatRoomName: ur.Name,
						MemberPubData: chat.MemberPubData{
							DisplayName: ur2.Name,
							IconRelpath: "",
						},
					},
				})
			},
		},
	})

}

func TestToken(t *testing.T) {

	ur := user.UserRequest{
		Email:    helpers.CRNG_stringPanic(32) + "@test.test",
		Password: "!@213@#!asdASADAS",
		UserData: user.UserData{
			Name:     helpers.CRNG_stringPanic(32),
			Language: "pl",
			Country:  "PL",
		},
	}

	userToken, err := userCtx.ApiCreateAndLoginUser(&ur)

	if err != nil {
		t.Fatal(err)
	}

	req := chat.CreateChatRoomRequest{
		ChatRoomRequest: chat.ChatRoomRequest{
			Flags: chat.ForceRedirectEnabled,
		},
		UserData: chat.MemberData{
			MemberPubData: chat.MemberPubData{
				DisplayName: ur.UserData.Name,
			},
			ChatRoomName: ur.UserData.Name,
			Email:        ur.Email,
		},
	}

	m := "POST"
	u := "/api/chat/room"
	integratorMethod := "GET"
	integratorPath := "/api/chat_integrator/token"
	var ires chat.TokenResponse

	userCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u,
			ExpectedStatusCode: 401,
		}, {
			RequestMethod:      integratorMethod,
			RequestUrl:         integratorPath,
			ExpectedStatusCode: 200,
		}, {
			RequestMethod:      integratorMethod,
			RequestUrl:         integratorPath,
			ExpectedStatusCode: 401,
			RequestHeaders:     user.GetAuthHeader("abc"),
		}, {
			RequestMethod:      integratorMethod,
			RequestUrl:         integratorPath,
			ExpectedStatusCode: 200,
			RequestHeaders:     user.GetAuthHeader(userToken),
			ExpectedBodyVal: func(b []byte, i interface{}) error {
				helpers.JsonMustDeserialize(b, &ires)
				return nil
			},
		},
	})

	userCtx.Api.TestAssertCases(t, []api.TestCase{
		{
			RequestMethod:      m,
			RequestUrl:         u + "?t=" + url.QueryEscape(ires.Token),
			ExpectedStatusCode: 200,
			RequestReader:      helpers.JsonMustSerializeReader(req),
		},
	})
}
