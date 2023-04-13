import { CircularProgress, IconButton, 
      InputAdornment, TextField } from "@mui/material";
import { Close, Loupe } from "@mui/icons-material";
import Autocomplete from "@mui/lab/Autocomplete";
import React, { useState } from "react";
import { searchPubTrainings, searchTrainings } from "../apicalls/instructor.api";
import { locale2 } from "../locale";
import { makeStyles } from "@mui/styles";

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
  let _cv = ""

export default function TrainingAtc(props) {
    const classes = useStyles()

    const [cv, setCv] = useState("")
    const [to, setTo] = useState(null)

    const [options, setOptions] = useState([])

   function updateOptions(input) {
        if(to) return
        setTo(setTimeout(async () => {
          try {
            let c = []
            if(props.forUsr) {
              c = await searchPubTrainings(_cv, props.instructorID)
            } else {
              c = await searchTrainings(_cv)
            }
            setOptions(c)
          } catch(ex) {
            console.log(ex)
          } finally {
            setTo(null)
          }
        }, 500))
      }

    function clearComponent() {
      if(props.setValue) props.setValue(null)
      _cv = ""
      setCv("")
      setOptions([])
    }
      
    return (
      <Autocomplete
          value={props.value}
          onChange={(event, newValue) => {
            if(!newValue) return
            if(props.setValue)
                props.setValue(newValue)
          }}
          filterOptions={(options, params) => {
            return options;
          }}
          classes={{
            option: classes.option,
          }}
          selectOnFocus
          clearOnBlur
          handleHomeEndKeys
          filterSelectedOptions
          options={options} 
          getOptionLabel={o => (o.Training && o.Training.Title) || ""} 
          renderOption={(o) => (o.Training && o.Training.Title) || ""}
          renderInput={(params) => {
            if(to) {
                params.InputProps.endAdornment = (
                    <React.Fragment>
                        <CircularProgress style={{width: 30, height: 30}} />
                        <InputAdornment position="end">
                                <Loupe style={{
                                    color: "gray"
                                }}/>
                            </InputAdornment>
                    </React.Fragment>
                  )
            } else {
                params.InputProps.endAdornment = (<InputAdornment position="end">
                    <Loupe style={{
                        color: "gray"
                    }}/>
                    <IconButton size="small" onClick={clearComponent}>
                        <Close/>
                    </IconButton>
                </InputAdornment>)
            }
            params.InputProps.className = ""
            params.size = "small"
            return (
            <TextField 
                variant="outlined"
                style={{
                  marginLeft: `${props.noMargin && 0}`
                }}
                value={cv}
                onChange={async e => {
                setCv(e.target.value)
                _cv = e.target.value
                if(e.target.value) {
                  updateOptions(e.target.value)
                }
              }}
              {...params}
              label={locale2.SEARCH[props.lang]} />
          )}}
          freeSolo
        />)
}