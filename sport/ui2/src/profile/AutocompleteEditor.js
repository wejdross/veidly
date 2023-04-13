import {
  Button, Chip, CircularProgress,
  Grid, InputAdornment,
  TextField, Typography
} from '@mui/material'
import React, { useState } from 'react'
import EditIcon from '@mui/icons-material/Edit';
import { Cancel, Check, Save } from '@mui/icons-material';
import Autocomplete, { createFilterOptions } from '@mui/material/Autocomplete';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';
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

export default function AutocompleteEdit(props) {

  const classes = useStyles()
  const [rd, setRd] = useState(true)
  //const [ev, setEv] = useState("")
  const [st, setSt] = useState({ id: "", msg: "" })

  function beginEdit() {
    //setEv(props.value)
    setValue(props.value)
    setSt({ id: "", msg: "" })
    setRd(false)
  }

  const [value, setValue] = useState([])

  async function save() {
    try {
      setSt({ id: "wip", msg: "" })
      await props.onChange(value)
      setRd(true)
      //setSt({id: "ok", msg: ""})
      // setTimeout(async () => {
      //   setRd(true)
      // }, 200)
    } catch (ex) {
      setSt({ id: "ex", msg: ex })
    }
  }

  function okFrag() {
    if (st.id !== "ok") {
      return
    }
    return (
      <React.Fragment>
        <Check color="primary" />
      </React.Fragment>
    )
  }

  function exFrag() {
    if (st.id !== "ex") {
      return
    }
    return (
      <React.Fragment>
        {`${locale2.FAILED_REASON[props.lang]} st.msg`}
        <Button
          onClick={() => beginEdit()} color="primary" aria-label="edit">
          {locale2.ONCE_AGAIN[props.lang]}
        </Button>
        <Button
          onClick={() => setRd(true)} color="secondary" aria-label="edit">
          {locale2.CANCEL[props.lang]}
        </Button>
      </React.Fragment>
    )
  }

  function wipFrag() {
    if (st.id !== "wip") {
      return
    }
    return (
      <React.Fragment>
        <CircularProgress size={24} />
      </React.Fragment>
    )
  }

  function editFrag() {
    if (st.id !== "") {
      return
    }

    return (
      <Autocomplete multiple fullWidth
        value={value}
        onChange={(event, newValue) => {
          let cpy = []
          if (newValue)
            for (let i = 0; i < newValue.length; i++) {
              let x = newValue[i]
              if (cpy.indexOf(x) < 0)
                cpy.push(x)
            }
          setValue(cpy)
          //props.onChange(cpy)
        }}
        selectOnFocus
        clearOnBlur
        handleHomeEndKeys
        filterSelectedOptions
        options={props.options}
        renderInput={(params) => {
          let c = params.InputProps.startAdornment
          params.InputProps.startAdornment = (
            <React.Fragment>
              <InputAdornment position="start">
                <Button
                  style={{
                    padding: 0,
                    marginTop: 7,
                    minHeight: 0,
                    minWidth: 0,
                    display: 'inline-block'
                  }}
                  onClick={save} example color="primary" aria-label="edit">
                  <Save />
                </Button>
                <Button
                  style={{
                    padding: 0,
                    marginTop: 7,
                    minHeight: 0,
                    minWidth: 0,
                    display: 'inline-block'
                  }}
                  onClick={() => setRd(true)} color="secondary" aria-label="edit">
                  <Cancel />
                </Button>
              </InputAdornment>
              {c}
            </React.Fragment>
          )
          return (<TextField
            onKeyPress={(e) => {
              if (e.key === "Enter") {
                save()
              }
            }}
            variant="outlined" {...params}
            label={props.title} />)
        }}
      />
    )

    return (
      <React.Fragment>
        <Autocomplete
          multiple fullWidth
          value={value}
          onChange={(event, newValue) => {
            let cpy = []
            if (newValue)
              for (let i = 0; i < newValue.length; i++) {
                let x = newValue[i]
                if (cpy.indexOf(x) < 0)
                  cpy.push(x)
              }
            setValue(cpy)
            //props.onChange(cpy)
          }}
          //defaultValue={tags}
          filterOptions={(options, params) => {
            // let o = []
            // for(let i = 0; i < options.length; i++) {
            //   if(value.indexOf())
            // }
            const filtered = filter(options, params)
            return filtered
          }}
          classes={{
            option: classes.option
          }}
          selectOnFocus
          clearOnBlur
          handleHomeEndKeys
          filterSelectedOptions
          options={props.options}
          getOptionLabel={o => o.val}
          renderOption={(_, o) => o.val}
          renderInput={(params) => {
            let c = params.InputProps.startAdornment
            params.InputProps.startAdornment = (
              <React.Fragment>
                <InputAdornment position="start">
                  <Button
                    style={{
                      padding: 0,
                      marginTop: 7,
                      minHeight: 0,
                      minWidth: 0,
                      display: 'inline-block'
                    }}
                    onClick={save} example color="primary" aria-label="edit">
                    <Save />
                  </Button>
                  <Button
                    style={{
                      padding: 0,
                      marginTop: 7,
                      minHeight: 0,
                      minWidth: 0,
                      display: 'inline-block'
                    }}
                    onClick={() => setRd(true)} color="secondary" aria-label="edit">
                    <Cancel />
                  </Button>
                </InputAdornment>
                {c}
              </React.Fragment>
            )
            return (<TextField
              onKeyPress={(e) => {
                if (e.key === "Enter") {
                  save()
                }
              }}
              variant="outlined" {...params}
              label={props.title} />)
          }}
        />
      </React.Fragment>
    )
  }

  return (
    <React.Fragment>
      {rd && (
        <Grid container spacing={3}>
          <Grid container item xs={4}>
            <Typography
              style={{ color: "gray" }}
              component={'span'} className={classes.typographyLineHeight}>
              {props.label}
            </Typography>
          </Grid>
          {
            (props.value && props.value() && (
              <>
                <Grid container item xs={6}>
                  {props.value().map((c, i) => (
                    <Chip key={i} label={props.valueLabel(c)} style={{ marginRight: 4, marginBottom: 4 }} />
                  ))}
                  {props.placeholder && props.value().length === 0 && <Typography style={{
                    marginTop: 5,
                    color: MulwiColors.subtitleTypography
                  }} variant="body2">
                    {props.placeholder}
                  </Typography>}
                </Grid>
                <Grid container item xs={2}>
                  <Button onClick={beginEdit} color="primary"
                    size='small' aria-label="edit"
                    className={classes.moveEditRight}>
                    <EditIcon className={classes.moveEditRight} />
                  </Button>
                </Grid>
              </>
            )) || (
              <>
                {/* <Grid container item xs={6}>
                  {props.placeholder}
                </Grid> */}
                <Grid container item xs={2}>
                  <Button onClick={beginEdit} color="primary" size='small' aria-label="edit">
                    <EditIcon />
                  </Button>
                </Grid>
              </>
            )
          }
        </Grid>)}
      {rd || (
        <React.Fragment>
          {editFrag()}
          {okFrag()}
          {exFrag()}
          {wipFrag()}
        </React.Fragment>)}
    </React.Fragment>
  )
}