// import AdyenCheckout from '@adyen/adyen-web';
// import '@adyen/adyen-web/dist/adyen.css';
// import {
//     Backdrop,
//     Button, CircularProgress, Container, Grid, Paper, TextField, Typography
// } from '@mui/material';
// import React, { useEffect, useState } from "react";
// import { useHistory } from "react-router-dom";
// import {
//     deletePayoutData, getPayoutData,
//     updatePayoutData
// } from "../apicalls/instructor.api";
// import CardWithBg from "../card/cardWithBg";
// import { MulwiColors } from "../mulwiColors";
// import {
//     getDialogWithOptions, getErrorDialog,
//     getNullDialog, StatusDialog
// } from "../StatusDialog";
// import { getSupportedLanguage, getSupportedLocale, locale2 } from '../locale';


// export default function InstructorPayments(props) {

//     const [info, setInfo] = useState(getNullDialog())
//     const [load, setLoad] = useState(false)

//     const [pd, setPd] = useState(null)

//     const h = useHistory()

//     async function refresh() {
//         try {
//             let pd = await getPayoutData()
//             pd = JSON.parse(pd)
//             setPd(pd)
//             let query = new URLSearchParams(window.location.search)
//             switch (query.get("return_to")) {
//                 case "configure":
//                     h.push("/configure?payout=1")
//                     break
//                 default:
//                     break
//             }
//         } catch (ex) {
//             setPd(null)
//             if (ex != 404) {
//                 setInfo(getErrorDialog(`${locale2.SOMETHING_WENT_WRONG[props.lang]}`, ex))
//             }
//         }
//     }

//     useEffect(() => {
//         refresh()
//     }, [])

//     const configuration = {
//         paymentMethodsResponse: {
//             "paymentMethods": [
//                 {
//                     "brands": [
//                         "amex",
//                         "bcmc",
//                         "cup",
//                         "diners",
//                         "discover",
//                         "jcb",
//                         "maestro",
//                         "mc",
//                         "visa"
//                     ],
//                     "details": [
//                         {
//                             "key": "number",
//                             "type": "text"
//                         },
//                         {
//                             "key": "expiryMonth",
//                             "type": "text"
//                         },
//                         {
//                             "key": "expiryYear",
//                             "type": "text"
//                         },
//                         {
//                             "key": "cvc",
//                             "type": "text"
//                         },
//                         {
//                             "key": "holderName",
//                             "optional": false,
//                             "type": "text"
//                         }
//                     ],
//                     "name": "Credit Card",
//                     "type": "scheme"
//                 },
//             ]
//         }, // The `/paymentMethods` response from the server.
//         // Web Drop-in versions before 3.10.1 use originKey instead of clientKey.
//         clientKey: "test_SZSPQFLFEVDSFMJ2FWDIQCLBBEOFWQBE",
//         locale: getSupportedLanguage(),
//         environment: "test",
//         translations: {
//             "pl": {
//                 payButton: `${locale2.SAVE[props.lang]}`
//             },
//             "en": {
//                 payButton: `${locale2.SAVE[props.lang]}`
//             },
//             "de": {
//                 payButton: `${locale2.SAVE[props.lang]}`
//             }
//         },
//         onSubmit: async (state, dropin) => {
//             if (!state ||
//                 !state.data ||
//                 !state.data.paymentMethod) return
//             if (!state.isValid) return

//             let pm = state.data.paymentMethod
//             if (pm.type !== "scheme") return

//             let PayoutData = {
//                 CardNumber: pm.encryptedCardNumber,
//                 ExpiryMonth: pm.encryptedExpiryMonth,
//                 ExpiryYear: pm.encryptedExpiryYear,
//                 Cvc: pm.encryptedSecurityCode,
//                 //
//                 HolderName: pm.holderName,
//                 Brand: pm.brand,
//             }

//             setLoad(true)

//             try {
//                 await updatePayoutData(PayoutData)
//                 props.main.refreshInstructor()
//                 refresh()
//             } catch (ex) {
//                 setInfo(getErrorDialog(`${locale2.SOMETHING_WENT_WRONG[props.lang]}`, ex))
//             } finally {
//                 setLoad(false)
//             }

//         },
//         paymentMethodsConfiguration: {
//             card: { // Example optional configuration for Cards
//                 hasHolderName: true,
//                 holderNameRequired: true,
//                 hideCVC: false, // Change this to true to hide the CVC field for stored cards
//                 name: `${locale2.CREDIT_CARD_OR_DEBIT[props.lang]}`,
//             }
//         }
//     };

//     useEffect(() => {
//         const checkout = new AdyenCheckout(configuration);

//         const dropin = checkout
//             .create('dropin', {
//                 // Starting from version 4.0.0, 
//                 // Drop-in configuration only accepts props related to itself 
//                 //     and cannot contain generic configuration like the onSubmit event.
//                 openFirstPaymentMethod: true,
//                 showPayButton: true,
//             })
//             .mount('#dropin-container')
//     }, [])

//     function cardInfo() {
//         if (!pd) return null
//         return (<React.Fragment>
//             <Paper style={{ paddingLeft: 20, paddingRight: 20, padding: 10, maxWidth: 400 }}>
//                 <Grid container
//                     direction="row"
//                     justify="space-between"
//                     style={{
//                         marginBottom: 10
//                     }}
//                     alignItems="center">
//                     <Grid item>
//                         <Typography>{`${locale2.TYPE[props.lang]}`}</Typography>
//                     </Grid>
//                     <Grid item>
//                         <Typography><strong>{pd.CardBrand}</strong></Typography>
//                     </Grid>
//                 </Grid>
//                 <Grid container
//                     style={{
//                         marginBottom: 10
//                     }}
//                     direction="row"
//                     justify="space-between"
//                     alignItems="center">
//                     <Grid item>
//                         <Typography>{`${locale2.ISSUED_FOR[props.lang]}: `}</Typography>
//                     </Grid>
//                     <Grid item>
//                         <Typography> <strong>{pd.CardHolderName}</strong></Typography>
//                     </Grid>
//                 </Grid>
//                 <Grid container
//                     direction="row"
//                     justify="space-between"
//                     alignItems="center">
//                     <Grid item>
//                         <Typography>{`${locale2.LAST_4_DIGITS[props.lang]}`}</Typography>
//                     </Grid>
//                     <Grid item>
//                         <Typography>**** **** **** {pd.CardSummary}</Typography>
//                     </Grid>
//                 </Grid>
//             </Paper>
//         </React.Fragment>)
//     }

//     function deleteCardPopup() {
//         setInfo(getDialogWithOptions(`${locale2.DELETING_CARD[props.lang]}`,
//             <React.Fragment>
//                 <Typography>{`${locale2.ARE_YOU_SURE[props.lang]}`}</Typography>
//                 <br />
//                 {cardInfo()}
//                 <br />
//                 <Typography style={{
//                     maxWidth: 400
//                 }}>{`${locale2.IF_YOU_REMOVE_CARD[props.lang]}`}</Typography>
//             </React.Fragment>,
//             <React.Fragment>
//                 <Button variant="contained" style={{
//                     color: "white",
//                     backgroundColor: MulwiColors.redError
//                 }} onClick={async () => {
//                     try {
//                         await deletePayoutData()
//                         props.main.refreshInstructor()
//                         setInfo(getNullDialog())
//                         refresh()
//                     } catch (ex) {
//                         setInfo(getErrorDialog(`${locale2.SOMETHING_WENT_WRONG[props.lang]}`, ex))
//                     }
//                 }}>
//                     {`${locale2.DELETE[props.lang]}`}
//                 </Button>
//             </React.Fragment>))
//     }

//     return (<React.Fragment>
//         <Backdrop open={load} style={{
//             zIndex: 9999
//         }}>
//             <CircularProgress />
//         </Backdrop>
//         <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
//         <CardWithBg bg="#f7f8f9" img={"/static/form-backgrounds/kayak.webp"}>
//             <div style={{
//                 marginLeft: 15,
//                 marginRight: 15
//             }}>
//                 {pd && ( 
//                     <React.Fragment>
//                         <Typography variant="h6">
//                             {`${locale2.CURRENTLY_FOR_PAYOUTS[props.lang]}`}
//                         </Typography>
//                         <br />
//                         {cardInfo()}
//                         <br />
//                         <Button fullWidth variant="contained" style={{
//                             color: "white",
//                             backgroundColor: MulwiColors.redError
//                         }} onClick={deleteCardPopup}>
//                             {`${locale2.DELETE[props.lang]}`}
//                         </Button>
//                         <br />
//                     </React.Fragment>
//                 )}
//                 <br />
//                 <Typography variant="h6" style={{
//                     maxWidth: 600
//                 }}>
//                     {pd ?
//                         locale2.MODIFY[props.lang] :
//                         locale2.ADD[props.lang]}
//                     {locale2.CARD_TO_WHICH_WE_TRANSFER[props.lang]}
//                 </Typography>
//                 <br />
//                 <Typography>
//                     {locale2.WE_TAKE_SECURITY[props.lang]}
//                 </Typography>
//                 <Typography style={{
//                     margin: 5,
//                     color: MulwiColors.greenDark
//                 }}>
//                     <strong>{locale2.WE_DONT_KNOW_CARD[props.lang]}</strong>
//                 </Typography>
//                 <Typography style={{
//                     maxWidth: 600
//                 }}>
//                     {locale2.ONCE_YOU_ENTER_CARD[props.lang]} <a style={{
//                         color: MulwiColors.blueDark,
//                         textDecoration: "none"
//                     }} href="https://www.adyen.com/" target="_blank" rel="noreferrer">
//                         {locale2.PAYMENT_BROKER[props.lang]}
//                     </a>
//                     {locale2.COMPLIES_WITH_PCI_DSS[props.lang]}
//                 </Typography>
//                 <Typography>
//                     {locale2.DATA_THAT_WE_KNOW_AND_STORE[props.lang]}
//                 </Typography>
//                 <ul>
//                     <li>{locale2.NAME_AND_LAST_NAME_ON_CARD[props.lang]}</li>
//                     <li>{locale2.CARD_TYPE[props.lang]}</li>
//                     <li>{locale2.LAST_4_DIGITS_AKA[props.lang]}</li>
//                 </ul>
//                 {/* <Link style={{
//                     color: MulwiColors.blueDark,
//                     textDecoration: "none"
//                 }}>further read...</Link> */}
//             </div>

//             <div id="dropin-container"></div>
//         </CardWithBg>
//         <br/>
//     </React.Fragment>)
// }

