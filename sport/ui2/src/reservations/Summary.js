// import { Button, Divider, Grid, Link, TextField, Typography } from '@mui/material'
// import React, { useEffect, useState } from 'react'
// import CardWithBg from '../card/cardWithBg'
// import { getSupportedLocale } from '../locale'
// import { MulwiColors } from '../mulwiColors'
// import { returnLocaleString } from '../locale'

// /*
// Summary: {
//             pl: [
//                 "Przekaż instruktorowi coś o sobie",
//                 "Imię nazwisko lub nick",
//                 "Rok urodzenia", // 2
//                 "Rok urodzenia może być potrzebny instruktorowi by oszacować twój wiek\ni określić czy się nadajesz na ten trening",
//                 "Numer telefonu", // 4
//                 "Dane kontaktowe które tu opcjonalnie możesz sprecyzować\nzostaną przekazane tylko instruktorowi\nPo to by się mógł z Tobą skontaktować w sprawie nieścisłości lub pytań.",
//                 "OK, rezerwuję i płacę", // 6
//             ],
//             en: [
//                 "Tell instructor something about Yourself",
//                 "Name or nick",
//                 "Year of birth",
//                 "The year of birth may be needed by an instructor to estimate your age \nand determine whether you are suitable for this training",
//                 "Phone number",
//                 "Contact details will be sent to the instructor, so that they can contact You in case of any problems or questions.",
//                 "OK, make reservation and pay.",
//             ]
//         },
// */

// export function Summary(props) {

//     const [name, setName] = useState("")
//     const [yearOfBirth, setYearOfBirth] = useState(0)
//     const [email, setEmail] = useState("")
//     const [phoneNumber, setPhoneNumber] = useState("")

//     return (
//         <CardWithBg img="/static/form-backgrounds/surfer.webp">
//             <div>
//             <Typography variant="h6" style={{
//                 marginBottom: 10
//             }}>
//                 {returnLocaleString(['reservations', 'Summary'])[0]}
//             </Typography>
//             <Grid container direction="column" spacing={3}>
//                 <Grid item>
//                     <TextField type="text" size="small"
//                         value={name}
//                         fullWidth
//                         onChange={(e) => setName(e.target.value)}
//                         variant="outlined" 
//                         label={returnLocaleString(['reservations', 'Summary'])[1]}/>
//                 </Grid>
//                 <Grid item>
//                     <Divider/>
//                 </Grid>
//                 <Grid item>
//                     <TextField type="text" size="small"
//                         value={yearOfBirth <= 0 ? "" : String(yearOfBirth)}
//                         onChange={(e) => {
//                             let x = Number(e.target.value)
//                             setYearOfBirth(x)
//                         }}
//                         variant="outlined" 
//                         label={returnLocaleString(['reservations', 'Summary'])[2]}/>
//                 </Grid>
//                 <Grid item>
//                     <div style={{whiteSpace:"pre-wrap"}}>
//                         {returnLocaleString(['reservations', 'Summary'])[3]}
//                     </div>
//                 </Grid>
//                 <Grid item>
//                     <Divider/>
//                 </Grid>
//                 <Grid item>
//                     <Grid container direction="row" spacing={3}>
//                         <Grid item>
//                             <TextField type="text" size="small"
//                                 value={email}
//                                 onChange={(e) => setEmail(e.target.value)}
//                                 variant="outlined" label="Email"/>
//                         </Grid>
//                         <Grid item>
//                             <TextField type="text" size="small"
//                                 value={phoneNumber}
//                                 onChange={(e) => setPhoneNumber(e.target.value)}
//                                 variant="outlined" 
//                                 label={returnLocaleString(['reservations', 'Summary'])[4]}/>
//                         </Grid>
//                     </Grid>
//                 </Grid>
//                 <Grid item>
//                     <div style={{whiteSpace:"pre-wrap"}}>
//                         {returnLocaleString(['reservations', 'Summary'])[5]}
//                     </div>
//                 </Grid>
//                 <Grid item>
//                     <Button variant="contained" 
//                         onClick={() => {
//                             const [l, c] = getSupportedLocale()
//                             props.onConfirm({
//                                 Name: name,
//                                 YearOfBirth: yearOfBirth,
//                                 Language: l,
//                                 Country: c
//                             }, {
//                                 PhoneNumber: phoneNumber,
//                                 Email: email
//                             })
//                         }}
//                         style={{
//                             color: "white",
//                             backgroundColor: MulwiColors.greenDark
//                         }}>
//                         {returnLocaleString(['reservations', 'Summary'])[6]}
//                     </Button>
//                 </Grid>
//             </Grid>
//             </div>
//         </CardWithBg>)
// }