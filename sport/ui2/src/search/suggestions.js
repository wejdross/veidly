// import React  from "react";
// import SingleTraining from "./singleTraining";
// import { Grid, Typography, 
//         useMediaQuery, useTheme } from "@mui/material";
// import { locale2 } from "../locale";
// import { activateMapItem, deactivateMapItem } from "./results";

// export default function Suggestions(props) {

//     const data = props.data
//     const theme = useTheme()
//     const isSmall = useMediaQuery(theme.breakpoints.down('sm'))

//     if (!data || data.length < 1) return null

//     return (
//         <React.Fragment>
//             <div style={{
//                 marginTop: "auto"
//             }}>
//                 <center>
//                     <Typography variant="h5" style={{
//                         marginTop: 20,
//                         marginBottom: 20,
//                     }}>
//                         <strong>{locale2.SIMILAR_OFFERS[props.lang]}</strong>
//                     </Typography>
//                 </center>
//                 <React.Fragment>
//                     <Grid
//                         justify="center"
//                         alignItems="center"
//                         // spacing={isSmall ? 0 : 3}
//                         style={{ marginTop: 20 }}
//                         container direction="row" item xs={12}>
//                         {data && data.map((v, i) => {
//                             return (
//                                 <Grid item key={i} lg={!isSmall ? 12 : 4}
//                                     style={{
//                                         marginBottom: isSmall ? 20 : null,
//                                         maxWidth: !isSmall ? 900 : 390,
//                                         marginBottom: 10
//                                     }}>
//                                     <SingleTraining lang={props.lang}
//                                         user={props.user}
//                                         onMouseEnter={() => {
//                                             activateMapItem(v)
//                                         }}
//                                         onMouseLeave={() => {
//                                             deactivateMapItem(v)
//                                         }}
//                                         dr={{
//                                             DateStart: props.apiRequest.DateStart,
//                                             DateEnd: props.apiRequest.DateEnd
//                                         }}
//                                         list={!isSmall}
//                                         d={v}
//                                         key={i} />
//                                 </Grid>
//                             )
//                         })}
//                     </Grid>
//                 </React.Fragment>
//             </div>
//         </React.Fragment>
//     )
// }
