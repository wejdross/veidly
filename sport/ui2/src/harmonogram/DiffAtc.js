// import { makeStyles, TextField } from "@mui/material";
// import Autocomplete, { createFilterOptions } from "@mui/lab/Autocomplete";
// import { useEffect, useState } from "react";
// import { diffs, getDiffsByLang } from "../diffs";
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

// export function DiffAtc(props) {
    
//     const classes = useStyles()

//     const [value, setValue] = useState([])

//     useEffect(() => {
//       let x = props.diff
//       let y = []
//       for(let i = 0; i < x.length; i++) {
//         let v = x[i]
//         y.push({id: v, val: diffs[v]["pl"]})
//       }
//       setValue(y)
//     }, [props.diff])

//     return (<Autocomplete multiple fullWidth
//         value={value}
//         onChange={(event, newValue) => {
//           let cpy = []
//           if(newValue)
//             for(let i = 0; i < newValue.length; i++) {
//               let x = newValue[i]
//               if(cpy.indexOf(x) < 0)
//                 cpy.push(x.id)
//             }
//             props.setDiff(cpy)
//           //props.onChange(cpy)
//         }}
//         //defaultValue={tags}
//         filterOptions={(options, params) => {
//           let o = []
//           for(let i = 0; i < options.length; i++) {
//             if(props.diff.indexOf(options[i].id) < 0)
//               o.push(options[i])
//           }
//           const filtered = filter(o, params)
//           return filtered
//         }}
//         classes={{
//           option: classes.option
//         }}
//         selectOnFocus
//         clearOnBlur
//         handleHomeEndKeys
//         filterSelectedOptions
//         options={getDiffsByLang("pl")} 
//         getOptionLabel={o => o.val} 
//         renderOption={(o) => o.val}
//         renderInput={(params) => {
//             return (<TextField
//                 variant="outlined" {...params}  
//                 label={locale2.SELECT_DIFF[props.lang]} />)
//         }}/>)
// }