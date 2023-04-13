package rsv

// import (
// 	"sport/adyen"
// 	"sport/api"
// 	"sport/helpers"

// 	"github.com/google/uuid"
// )

// var adyenSuccessRes = "[accepted]"

// func boolToAdyenStr(v bool) string {
// 	ss := ""
// 	if v {
// 		ss = "true"
// 	} else {
// 		ss = "false"
// 	}
// 	return ss
// }

// func (ctx *Ctx) ApiSimulateAdyenWHCall(n *adyen.Notification) error {
// 	return ctx.Api.TestAssertCaseErr(&api.TestCase{
// 		RequestMethod:      "POST",
// 		RequestUrl:         adyen.ApiWhPath,
// 		ExpectedStatusCode: 200,
// 		RequestHeaders:     ctx.Adyen.CorrectAdyenAuthHdr(),
// 		ExpectedBody:       &adyenSuccessRes,
// 		RequestReader:      helpers.JsonMustSerializeReader(n),
// 	})
// }

// type AdyenNtf_AuthOpts struct {
// 	Success bool
// }

// func AdyenNtf_Auth(r *DDLRsvWithInstr, o AdyenNtf_AuthOpts) *adyen.Notification {
// 	s := boolToAdyenStr(o.Success)
// 	return &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					Amount: adyen.Amount{
// 						Value:    int64(r.Training.Price),
// 						Currency: r.Training.Currency,
// 					},
// 					Success:           s,
// 					EventCode:         "AUTHORISATION",
// 					PspReference:      uuid.New().String(),
// 					MerchantReference: adyenWhPrefix + r.RefID,
// 				},
// 			},
// 		},
// 	}
// }

// func AdyenNtf_PositiveAuth(r *DDLRsvWithInstr) *adyen.Notification {
// 	return AdyenNtf_Auth(r, AdyenNtf_AuthOpts{
// 		Success: true,
// 	})
// }

// func AdyenNtf_Capture(r *DDLRsvWithInstr, success bool) *adyen.Notification {
// 	ss := "true"
// 	if !success {
// 		ss = "false"
// 	}
// 	return &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					Amount: adyen.Amount{
// 						Value:    int64(r.Training.Price),
// 						Currency: r.Training.Currency,
// 					},
// 					Success:           ss,
// 					EventCode:         "CAPTURE",
// 					MerchantReference: adyenWhPrefix + r.RefID,
// 				},
// 			},
// 		},
// 	}
// }

// func AdyenNtf_PositiveCapture(r *DDLRsvWithInstr) *adyen.Notification {
// 	return AdyenNtf_Capture(r, true)
// }

// func AdyenNtf_Payout(r *DDLRsvWithInstr, success bool) *adyen.Notification {
// 	ss := boolToAdyenStr(success)
// 	return &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					Amount: adyen.Amount{
// 						Value:    int64(r.SplitPayout),
// 						Currency: r.Training.Currency,
// 					},
// 					Success:           ss,
// 					EventCode:         "PAYOUT_THIRDPARTY",
// 					MerchantReference: adyenWhPrefix + r.RefID,
// 				},
// 			},
// 		},
// 	}
// }

// func AdyenNtf_PositivePayout(r *DDLRsvWithInstr) *adyen.Notification {
// 	return AdyenNtf_Payout(r, true)
// }

// func AdyenNtf_Refund(r *DDLRsvWithInstr, success bool) *adyen.Notification {
// 	ss := boolToAdyenStr(success)
// 	return &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					Amount: adyen.Amount{
// 						Value:    int64(r.Training.Price),
// 						Currency: r.Training.Currency,
// 					},
// 					Success:           ss,
// 					EventCode:         "REFUND",
// 					MerchantReference: adyenWhPrefix + r.RefID,
// 				},
// 			},
// 		},
// 	}
// }

// func AdyenNtf_CancelOrRefund(r *DDLRsvWithInstr, success bool) *adyen.Notification {
// 	ss := boolToAdyenStr(success)
// 	return &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					Amount: adyen.Amount{
// 						Value:    int64(r.Training.Price),
// 						Currency: r.Training.Currency,
// 					},
// 					Success:           ss,
// 					EventCode:         "CANCEL_OR_REFUND",
// 					MerchantReference: adyenWhPrefix + r.RefID,
// 				},
// 			},
// 		},
// 	}
// }
