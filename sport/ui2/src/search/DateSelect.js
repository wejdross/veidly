// import { Divider, FormControl, makeStyles, 
//         MenuItem, Select, useTheme } from "@mui/material";
// import React, { useEffect, useState } from "react";
// import { dfInHours, months } from "../helpers";
// import { returnLocaleString } from "../locale";

// const seasonOffset = 20

// const useStyles = makeStyles((theme) => ({
//   formControl: {
//     [theme.breakpoints.down("lg")]: {
//       width: 80,
//     },
//     [theme.breakpoints.between("lg", "xl")]: {
//       width: 140,
//       //  marginLeft: 40,
//     },
//     [theme.breakpoints.up("xl")]: {
//       width: 150,
//     },
//     backgroundColor: "white",
//   },
// }));

// export function DateSelect(props) {

//   const styles = useStyles();
//   const [value, setValue] = useState(seasonOffset);

//   useEffect(() => {
//     if (props.searchRequest && props.searchRequest.DateStart 
//                 && props.searchRequest.DateEnd) {
//       let start = new Date(props.searchRequest.DateStart)
//       let end = new Date(props.searchRequest.DateEnd)
//       if(dfInHours(start, end) > (24*32)) {
//         setValue(seasonOffset + getSeason(start))
//       } else {
//         setValue(start.getMonth())
//       }
//     }
//   }, [props.searchRequest])
  
//   const SEASON_WINTER = 3
//   const SEASON_AUTUMN = 2
//   const SEASON_SUMMER = 1
//   const SEASON_SPRING = 0

//   // months are 0-based
//   const seasons = {
//     [SEASON_SPRING]: {
//       label: returnLocaleString(['search', 'DateSelect'])[0],
//       //desc: "(Marzec - Maj)",
//       months: [2, 3, 4]
//     },
//     [SEASON_SUMMER]: {
//       label: returnLocaleString(['search', 'DateSelect'])[1],
//       //desc: "(Czerwiec - Sierpień)",
//       months: [5, 6, 7]
//     },
//     [SEASON_AUTUMN]: {
//       label: returnLocaleString(['search', 'DateSelect'])[2],
//       //desc: "(Wrzesień - Listopad)",
//       months: [8, 9, 10]
//     },
//     [SEASON_WINTER]: {
//       label: returnLocaleString(['search', 'DateSelect'])[3],
//       //desc: "(Grudzień - Luty)",
//       months: [11, 0, 1]
//     },
//   }

//   function mergeRequest(_value) {
//     let r = props.searchRequest
//     // month
//     if(_value < seasonOffset) {
//       let s = new Date()
//       s.setHours(0, 0, 0, 0)
//       s.setDate(1)
//       s.setMonth(_value)
//       r.DateStart = s
//       let e = new Date(s)
//       e.setHours(23, 59, 59, 0)
//       e.setMonth(e.getMonth() + 1)
//       e.setDate(e.getDate() - 1)
//       r.DateEnd = e
//     // season
//     } else {
//       let s = new Date()
//       s.setHours(0, 0, 0, 0)
//       s.setDate(1)
//       s.setMonth(seasons[_value - seasonOffset].months[0])
//       r.DateStart = s
//       let e = new Date(s)
//       e.setHours(23, 59, 59, 0)
//       e.setMonth(e.getMonth() + 3)
//       e.setDate(e.getDate() - 1)
//       r.DateEnd = e
//     }
//     return r
//   }


//   function getSeason(d) {
//     let month = d.getMonth()
//     if(month < 2) {
//       return SEASON_WINTER
//     }
//     if(month < 5) {
//       return SEASON_SPRING
//     }
//     if(month < 8) {
//       return SEASON_SUMMER
//     }
//     if(month < 11) {
//       return SEASON_AUTUMN
//     }
//     return SEASON_WINTER
//   } 

//   function generateSeasons() {
//     let now = new Date()
//     let season = getSeason(now)
//     let ret = []
//     let yr = now.getFullYear()
//     for(let i = 0; i < 4; i++) {
//       let s = (season + i) % 4
//       if((season + i) == 4) yr++
//       ret.push({
//         season: seasons[s],
//         year: yr,
//         ix: s
//       })
//     }
//     return ret
//   }

//   function generateMonths() {
//     let now = new Date()
//     let mon = now.getMonth()
//     let yr = now.getFullYear()
//     let ret = []
//     for(let i = 0; i < 4; i++) {
//       let m = (mon + i) % 12
//       if((mon + i) == 12) yr++
//       ret.push({
//         ix: m,
//         label: months[m] + " " + yr
//       })
//     }
//     return ret
//   }

//   return (
//     <React.Fragment>

//       <FormControl size="small" className={styles.formControl}>
//         <Select
//           variant="outlined"
//           value={value}
//           onChange={(e) => {
//             props.onChange(mergeRequest(e.target.value))
//           }}
//         >
//           {generateMonths().map((m, i) => (<MenuItem key={m.ix} value={m.ix}>
//               {i == 0 ? <strong>{m.label}</strong> : m.label}
//             </MenuItem>))}
//           <Divider style={{width:'100%'}}/>

//           {generateSeasons().map((s, i) => (
//             <MenuItem key={s.ix+seasonOffset} value={s.ix+seasonOffset}>
//                   {i == 0 ? (<strong>{s.season.label} {s.year}</strong>) : (<span>{s.season.label} {s.year}</span>)}
//                 </MenuItem>
//           ))}
         
//         </Select>
//       </FormControl>
//     </React.Fragment>
//   );
// }