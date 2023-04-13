import { TextField } from "@mui/material";
import Autocomplete, { createFilterOptions } from "@mui/lab/Autocomplete";
import { useState } from "react";
import { makeStyles } from "@mui/styles";


const filter = createFilterOptions()

const useStyles = makeStyles((t) => (
    {
        option: {
            minHeight: 'auto',
            alignItems: 'flex-start',
            padding: 8,
            '&[aria-selected="true"]': {
                backgroundColor: 'transparent',
            },
            '&[data-focus="true"]': {
                backgroundColor: 'rgba(0, 0, 0, 0.15)',
            },
        },
        typographyLineHeight: {
            lineHeight: 2,
        },
    }
))

export function ArrEdit(props) {
    const classes = useStyles()
    const [options, setOptions] = useState([])

    return (
        <Autocomplete 
          multiple
            value={props.value}
            fullWidth
            onChange={(event, newValue) => {
              let cpy = []
              if(newValue)
                for(let i = 0; i < newValue.length; i++) {
                  let x = newValue[i].replace("Dodaj \"", "")
                  x = x.replace("\"", "")
                  if(cpy.indexOf(x) < 0)
                    cpy.push(x)
                }
              props.setValue(cpy)
            }}
            filterOptions={(options, params) => {
              let o = []
              for(let i = 0; i < options.length; i++) {
                if(props.value.indexOf(options[i]) < 0)
                  o.push(options[i])
              }
  
              const filtered = filter(o, params);
          
              // Suggest the creation of a new value
              if (params.inputValue !== '') {
                filtered.push( `Dodaj "${params.inputValue}"`);
              }
          
              return filtered;
            }}
            classes={{
              option: classes.option
            }}
            selectOnFocus
            clearOnBlur
            handleHomeEndKeys
            filterSelectedOptions
            id="tagatc"
            options={options} 
            getOptionLabel={o => o} 
            renderOption={(o) => o}
            renderInput={(params) => {
                return (
              <TextField
                variant="outlined" 
                {...params}
                label={props.label} />
            )}}
            freeSolo
          />)
}