// import {
//     Grid
// } from '@mui/material';
// import { KeyboardDatePicker } from '@mui/lab/';
// import React, { useState } from 'react';
// import { dateIsNotZero } from '../helpers';
// import ModalEdit from '../profile/ModalEdit';
// import { locale2, returnLocaleString } from '../locale';

// export function DateRangeEditor(props) {

//     const [dateStart, setDateStart] = useState(null)
//     const [dateEnd, setDateEnd] = useState(null)

//     return (
//         <ModalEdit
//             title={props.title + ' - ' + locale2.DURATION[props.lang]}
//             value={
//                 dateIsNotZero(props.dateStart) &&
//                 dateIsNotZero(props.dateEnd) && (
//                     props.dateStart.toDateString() +
//                     " - " +
//                     props.dateEnd.toDateString()
//                 )}
//             label={locale2.DURATION[props.lang]}
//             onSave={async () => {
//                 try {
//                     await props.onChange(dateStart, dateEnd)
//                 } catch (ex) {
//                     console.log(ex)
//                     throw ex
//                 }
//             }}
//             content={
//                 <React.Fragment>
//                     <Grid container spacing={2} style={{ marginBottom: 20 }}>
//                         <Grid item sm={12} md={6}>
//                             <KeyboardDatePicker fullWidth
//                                 margin="normal"
//                                 id="date-picker-dialog"
//                                 label={locale2.START_DATE[props.lang]}
//                                 format="MM/dd/yyyy"
//                                 value={dateStart}
//                                 onChange={e => {
//                                     let x = new Date(e)
//                                     x.setHours(0, 0, 0, 0)
//                                     setDateStart(x)
//                                 }}
//                                 KeyboardButtonProps={{
//                                     'aria-label': 'change date',
//                                 }} />
//                         </Grid>
//                         <Grid item sm={12} md={6}>
//                             <KeyboardDatePicker fullWidth
//                                 margin="normal"
//                                 id="date-picker-dialog"
//                                 label={locale2.END_DATE[props.lang]}
//                                 format="MM/dd/yyyy"
//                                 value={dateEnd}
//                                 onChange={e => {
//                                     let x = new Date(e)
//                                     x.setHours(23, 59, 59, 0)
//                                     setDateEnd(x)
//                                 }}
//                                 KeyboardButtonProps={{
//                                     'aria-label': 'change date',
//                                 }} />
//                         </Grid>
//                     </Grid>
//                 </React.Fragment>}
//             onOpen={() => {
//                 if (dateIsNotZero(props.dateStart))
//                     setDateStart(props.dateStart)
//                 if (dateIsNotZero(props.dateEnd))
//                     setDateEnd(props.dateEnd)
//             }}
//         />
//     )
// }
