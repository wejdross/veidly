// import { CircularProgress, makeStyles, TextField } from "@mui/material";
// import Autocomplete, { createFilterOptions } from "@mui/lab/Autocomplete";
// import { useState } from "react";
// import { getAllTags } from "../apicalls/tags";
// import { locale2, returnLocaleString } from "../locale";

// const filter = createFilterOptions()

// const useStyles = makeStyles((t) => (
//     {
//       option: {
//         minHeight: 'auto',
//         alignItems: 'flex-start',
//         padding: 8,
//         '&[aria-selected="true"]': {
//           backgroundColor: 'transparent',
//         },
//         '&[data-focus="true"]': {
//           backgroundColor: 'rgba(0, 0, 0, 0.15)',
//         },
//       },
//       typographyLineHeight: {
//         lineHeight: 2,
//       },
//     }
//   ))
//   let _cv = ""

// export function TagAtc(props) {
    
//     const classes = useStyles()

//     const [cv, setCv] = useState("")
//     const [to, setTo] = useState(null)

//     const [options, setOptions] = useState([])

//    function updateOptions(input) {
//         if(to) return
//         setTo(setTimeout(async () => {
//           try {
//             // opts = await getAllTags(_cv)
//             // opts = JSON.parse(opts)
//             let c = await getAllTags(_cv)
//             setOptions(JSON.parse(c))
//           } catch(ex) {
//             console.log(ex)
//           } finally {
//             setTo(null)
//           }
//           //opts.push( `Add "${cv}"`)
//         }, 1000))
//       }
      
//     return (
//       <Autocomplete 
//         multiple
//           value={props.tags}
//           fullWidth
//           onChange={(event, newValue) => {
//             let cpy = []
//             if(newValue)
//               for(let i = 0; i < newValue.length; i++) {
//                 let x = newValue[i].replace(`${locale2.ADD[props.lang]} "`, "")
//                 x = x.replace("\"", "")
//                 if(cpy.indexOf(x) < 0)
//                   cpy.push(x)
//               }
//             props.setTags(cpy)
//           }}
//           filterOptions={(options, params) => {
//             let o = []
//             for(let i = 0; i < options.length; i++) {
//               if(props.tags.indexOf(options[i]) < 0)
//                 o.push(options[i])
//             }

//             const filtered = filter(o, params);
        
//             // Suggest the creation of a new value
//             if (params.inputValue !== '') {
//               filtered.push( `${locale2.ADD[props.lang]} "${params.inputValue}"`);
//             }
        
//             return filtered;
//           }}
//           classes={{
//             option: classes.option
//           }}
//           selectOnFocus
//           clearOnBlur
//           handleHomeEndKeys
//           filterSelectedOptions
//           id="tagatc"
//           options={options} 
//           getOptionLabel={o => o} 
//           renderOption={(o) => o}
//           renderInput={(params) => {
//             if(to)
//               params.InputProps.endAdornment = (
//                 <CircularProgress style={{width: 30, height: 30}} />
//               )
//             return (
//             <TextField 
//               value={cv}
//               onChange={async e => {
//                 setCv(e.target.value)
//                 _cv = e.target.value
//                 if(e.target.value) {
//                   updateOptions(e.target.value)
//                 }
//               }}
//               variant="outlined" 
//               {...params}
//               label={locale2.ADD_OR_SELECT_CAT[props.lang]} />
//           )}}
//           freeSolo
//         />)
// }