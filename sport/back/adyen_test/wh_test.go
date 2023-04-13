package adyen_test

// func TestWH(t *testing.T) {

// 	m := "POST"
// 	p := adyen.ApiWhPath

// 	k := "FOO"

// 	cr := &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					MerchantReference: k + "_" + uuid.NewString(),
// 				},
// 			},
// 		},
// 	}
// 	cr2 := &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					MerchantReference: "BAR_" + uuid.NewString(),
// 				},
// 			},
// 		},
// 	}
// 	apiCtx.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         p,
// 			ExpectedStatusCode: 401,
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         p,
// 			ExpectedStatusCode: 400,
// 			RequestHeaders:     adyenCtx.CorrectAdyenAuthHdr(),
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         p,
// 			ExpectedStatusCode: 200,
// 			RequestHeaders:     adyenCtx.CorrectAdyenAuthHdr(),
// 			RequestReader: helpers.JsonMustSerializeReader(&adyen.Notification{
// 				Live: "false",
// 				NotificationItems: &[]adyen.NotificationItem{
// 					{
// 						NotificationRequestItem: adyen.NotificationRequestItem{},
// 					},
// 				},
// 			}),
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         p,
// 			ExpectedStatusCode: 422,
// 			RequestHeaders:     adyenCtx.CorrectAdyenAuthHdr(),
// 			RequestReader: helpers.JsonMustSerializeReader(&adyen.Notification{
// 				Live: "false",
// 				NotificationItems: &[]adyen.NotificationItem{
// 					{
// 						NotificationRequestItem: adyen.NotificationRequestItem{
// 							MerchantReference: "_",
// 						},
// 					},
// 				},
// 			}),
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         p,
// 			ExpectedStatusCode: 422,
// 			RequestHeaders:     adyenCtx.CorrectAdyenAuthHdr(),
// 			RequestReader:      helpers.JsonMustSerializeReader(cr),
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         p,
// 			ExpectedStatusCode: 422,
// 			RequestHeaders:     adyenCtx.CorrectAdyenAuthHdr(),
// 			RequestReader: helpers.JsonMustSerializeReader(&adyen.Notification{
// 				Live: "false",
// 				NotificationItems: &[]adyen.NotificationItem{
// 					{
// 						NotificationRequestItem: adyen.NotificationRequestItem{
// 							MerchantReference: k + "_BAR",
// 						},
// 					},
// 				},
// 			}),
// 		},
// 	})

// 	ok := "[accepted]"

// 	adyenCtx.AddAdyenHandler("FOO", func(ni *adyen.NotificationRequestItem) error {
// 		return nil
// 	})
// 	adyenCtx.AddAdyenHandler("BAR", func(ni *adyen.NotificationRequestItem) error {
// 		return fmt.Errorf("bar")
// 	})

// 	apiCtx.TestAssertCases(t, []api.TestCase{
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         p,
// 			ExpectedStatusCode: 422,
// 			RequestHeaders:     adyenCtx.CorrectAdyenAuthHdr(),
// 			RequestReader: helpers.JsonMustSerializeReader(&adyen.Notification{
// 				Live: "false",
// 				NotificationItems: &[]adyen.NotificationItem{
// 					{
// 						NotificationRequestItem: adyen.NotificationRequestItem{
// 							MerchantReference: k + "_BAR",
// 						},
// 					},
// 				},
// 			}),
// 		}, {
// 			RequestMethod:      m,
// 			RequestUrl:         p,
// 			ExpectedStatusCode: 200,
// 			ExpectedBody:       &ok,
// 			RequestHeaders:     adyenCtx.CorrectAdyenAuthHdr(),
// 			RequestReader:      helpers.JsonMustSerializeReader(cr),
// 		},
// 		{
// 			RequestMethod:      m,
// 			RequestUrl:         p,
// 			ExpectedStatusCode: 422,
// 			RequestHeaders:     adyenCtx.CorrectAdyenAuthHdr(),
// 			RequestReader:      helpers.JsonMustSerializeReader(cr2),
// 		},
// 	})

// }
