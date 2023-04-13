package sub

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

// func AdyenNtf_Auth(sb *Sub, o AdyenNtf_AuthOpts) *adyen.Notification {
// 	s := boolToAdyenStr(o.Success)
// 	return &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					Amount: adyen.Amount{
// 						Value:    int64(sb.SubModel.Total()),
// 						Currency: sb.SubModel.Currency,
// 					},
// 					Success:           s,
// 					EventCode:         "AUTHORISATION",
// 					PspReference:      uuid.New().String(),
// 					MerchantReference: AdyenSmPrefix + sb.RefID.String(),
// 				},
// 			},
// 		},
// 	}
// }

// func AdyenNtf_PositiveAuth(sub *Sub) *adyen.Notification {
// 	return AdyenNtf_Auth(sub, AdyenNtf_AuthOpts{
// 		Success: true,
// 	})
// }

// func AdyenNtf_Capture(sub *Sub, success bool) *adyen.Notification {
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
// 						Value:    int64(sub.SubModel.Total()),
// 						Currency: sub.SubModel.Currency,
// 					},
// 					Success:           ss,
// 					EventCode:         "CAPTURE",
// 					MerchantReference: AdyenSmPrefix + sub.RefID.String(),
// 				},
// 			},
// 		},
// 	}
// }

// func AdyenNtf_PositiveCapture(sub *Sub) *adyen.Notification {
// 	return AdyenNtf_Capture(sub, true)
// }

// func AdyenNtf_Payout(sub *Sub, success bool) *adyen.Notification {
// 	ss := boolToAdyenStr(success)
// 	return &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					Amount: adyen.Amount{
// 						Value:    int64(sub.SubModel.Total()),
// 						Currency: sub.SubModel.Currency,
// 					},
// 					Success:           ss,
// 					EventCode:         "PAYOUT_THIRDPARTY",
// 					MerchantReference: AdyenSmPrefix + sub.RefID.String(),
// 				},
// 			},
// 		},
// 	}
// }

// func AdyenNtf_PositivePayout(sub *Sub) *adyen.Notification {
// 	return AdyenNtf_Payout(sub, true)
// }

// func AdyenNtf_Refund(sub *Sub, success bool) *adyen.Notification {
// 	ss := boolToAdyenStr(success)
// 	return &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					Amount: adyen.Amount{
// 						Value:    int64(sub.SubModel.Total()),
// 						Currency: sub.SubModel.Currency,
// 					},
// 					Success:           ss,
// 					EventCode:         "REFUND",
// 					MerchantReference: AdyenSmPrefix + sub.RefID.String(),
// 				},
// 			},
// 		},
// 	}
// }

// func AdyenNtf_CancelOrRefund(sub *Sub, success bool) *adyen.Notification {
// 	ss := boolToAdyenStr(success)
// 	return &adyen.Notification{
// 		Live: "false",
// 		NotificationItems: &[]adyen.NotificationItem{
// 			{
// 				NotificationRequestItem: adyen.NotificationRequestItem{
// 					Amount: adyen.Amount{
// 						Value:    int64(sub.SubModel.Total()),
// 						Currency: sub.SubModel.Currency,
// 					},
// 					Success:           ss,
// 					EventCode:         "CANCEL_OR_REFUND",
// 					MerchantReference: AdyenSmPrefix + sub.RefID.String(),
// 				},
// 			},
// 		},
// 	}
// }
