import {
  CircularProgress, Grid,
  TextField, Autocomplete, Button
} from "@mui/material";
import { useEffect, useState } from "react";
import { getTags } from "../apicalls/user.api";
import { getCategoryLabel, getTagLabel } from "../harmonogram/trainingDetails";
import { locale2 } from "../locale"
import React from "react";
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
    searchBarElements: {
      [t.breakpoints.down("sm")]: {
        marginTop: 20,
        minWidth: 250,
      },
      // for column alignement
      // [t.breakpoints.up("sm")]: {
      //   marginLeft: 40,
      // },
      // for row alignement
      [t.breakpoints.up("lg")]: {
        minWidth: 350,
        marginLeft: 40,
      },
      backgroundColor: "white",
    }
  }
))

export default function TagEditor(props) {

  const classes = useStyles()

  const [cv, setCv] = useState("")
  const [th, setTh] = useState(null)

  const [options, setOptions] = useState([])
  function updateOptions(value) {
    if (th) window.clearTimeout(th)
    setTh(setTimeout(async () => {
      try {
        let c = await getTags(value, 1)
        c = JSON.parse(c)
        setOptions(c.slice(0, 8))
      } catch (ex) {
        console.log(ex)
      } finally {
        setTh(null)
      }
    }, 100))
  }

  useEffect(() => {
    setCv(props.val)
    setVal({
      Tag: {
        Name: props.val
      }
    })
  }, [props.val])


  const [val, setVal] = useState(null)

  return (
    <Autocomplete
      value={val}
      onChange={(event, newValue) => {
        if (!newValue) return
        setVal(newValue)
        props.setVal && props.setVal(getTagLabel(newValue))
      }}
      filterOptions={(options, params) => {
        return options;
      }}
      classes={{
        option: classes.option,
      }}
      selectOnFocus
      clearOnBlur
      size={props.size}
      handleHomeEndKeys
      filterSelectedOptions
      options={options}
      //getOptionSelected={(o,v) => o && o.Tag && v && v.Tag && o.Tag.Name===v.Tag.Name}
      getOptionLabel={getTagLabel}
      //valueLabel={getTagLabel}
      // renderOption={t =>
      //     <Grid container direction="row" justifyContent="flex-start" alignContent={"flex-start"}>
      //     <Grid item>
      //       <Button fullWidth>{t.key}</Button>
      //       <br />
      //     </Grid>
      //   </Grid>
      // }

      renderInput={(params) => {
        if (th)
          params.InputProps.endAdornment = (
            <CircularProgress style={{ width: 30, height: 30 }} />
          )
        return (
          <TextField
            className={props.class}
            style={{
              marginLeft: `${props.noMargin && 0}`
            }}
            value={cv}
            onChange={async e => {
              setCv(e.target.value)
              props.setVal && props.setVal(e.target.value)
              //_cv = e.target.value
              if (e.target.value) {
                updateOptions(e.target.value)
              }
            }}
            variant="outlined"
            {...params}
            InputLabelProps={props.noshrink ? { shrink: false } : null}
            placeholder={props.noshrink ? (cv ? null
              : locale2.DISCIPLINE[props.lang])
              : locale2.DISCIPLINE[props.lang]}
          />)
      }}
      freeSolo
    />)
}