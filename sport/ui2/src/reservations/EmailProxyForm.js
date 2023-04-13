// import { Button, Dialog, DialogActions, DialogContent, DialogTitle, 
//         Grid, TextField, useMediaQuery, useTheme } from '@mui/material'
// import React, { useEffect, useState } from 'react'
// import { MulwiColors } from '../mulwiColors'
// import { returnLocaleString } from '../locale'

// /*
// EmailProxyForm: {
//             pl: [
//                 "pytanie w sprawie treningu: ",
//                 "Skontaktuj się z instruktorem",
//                 "Twój email",
//                 "Tytuł",
//                 "Treść",
//                 "Wyślij wiadomość",
//             ],
//             en: [
//                 "Question about training: ",
//                 "Contact with the instructor",
//                 "Your e-mail",
//                 "Title",
//                 "Contents",
//                 "Send a message",
//             ]
//         },
// */

// export function EmailProxyForm(props) {

//     const [sender, setSender] = useState("")
//     const [title, setTitle] = useState("")
//     const [content, setContent] = useState("")
//     const theme = useTheme()
//     const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))

//     useEffect(() => {
//         if(props.user) {
//             if(props.user.ContactData.Email) {
//                 setSender(props.user.ContactData.Email)
//             }
//         }
//     }, [props.user])

//     useEffect(() => {
//         if(props.training) {
//             setTitle(returnLocaleString(['reservations', 'EmailProxyForm'])[0] + props.training.Title)
//         }
//     }, [props.training])

//     return (<Dialog
//                 open={props.open} onClose={() => props.setOpen(false)}>
//         <DialogTitle>
//             {returnLocaleString(['reservations', 'EmailProxyForm'])[1]}
//         </DialogTitle>
//         <DialogContent style={{
//             minWidth: isLowRes ? null : 600
//         }}>
//             <Grid container spacing={1} direction="column">
//                 <Grid item>
//                     <TextField
//                         label={returnLocaleString(['reservations', 'EmailProxyForm'])[2]}
//                         variant="outlined"
//                         size="small"
//                         value={sender}
//                         onChange={e => setSender(e.target.value)}
//                         fullWidth />
//                 </Grid>
//                 <Grid item>
//                     <TextField
//                         fullWidth
//                         value={title}
//                         onChange={e => setTitle(e.target.value)}
//                         variant="outlined"
//                         size="small"
//                         label={returnLocaleString(['reservations', 'EmailProxyForm'])[3]} />
//                 </Grid>
//                 <Grid item>
//                     <TextField
//                         id="emc"
//                         multiline
//                         rows={10}
//                         fullWidth
//                         value={content}
//                         onChange={e => setContent(e.target.value)}
//                         variant="outlined"
//                         size="small"
//                         label={returnLocaleString(['reservations', 'EmailProxyForm'])[4]} />
//                 </Grid>
//             </Grid>
//         </DialogContent>
//         <DialogActions>
//             <Button
//                 fullWidth
//                 style={{
//                     color: "white",
//                     backgroundColor: MulwiColors.blueDark
//                 }}>{returnLocaleString(['reservations', 'EmailProxyForm'])[5]}</Button>
//         </DialogActions>
//     </Dialog>)
// }
