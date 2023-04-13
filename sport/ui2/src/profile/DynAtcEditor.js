import {Button, Chip, CircularProgress, 
        Grid, InputAdornment,
        TextField, Typography} from '@mui/material'
import React, {useState} from 'react'
import EditIcon from '@mui/icons-material/Edit';
import {Cancel, Check, Save} from '@mui/icons-material';
import Autocomplete from '@mui/lab/Autocomplete';
import { MulwiColors } from '../mulwiColors';
import { locale2 } from '../locale';
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
  
export default function DynAtcEdit(props) {

const classes = useStyles()
  const [rd, setRd] = useState(true)
 // const [ev, setEv] = useState("")
  const [st, setSt] = useState({id: "", msg: ""})

  function beginEdit() {
    //setEv(props.value)
    setValue(props.value)
    setSt({id: "", msg: ""})
    setRd(false)
  }

  const [value, setValue] = useState([])

  async function save() {
    try {
      setSt({id: "wip", msg: ""})
      await props.onChange(value)
      setRd(true)
      //setSt({id: "ok", msg: ""})
      // setTimeout(async () => {
      //   setRd(true)
      // }, 200)
    } catch (ex) {
      setSt({id: "ex", msg: ex})
    }
  }

  function okFrag() {
    if (st.id !== "ok") {
      return
    }
    return (
      <React.Fragment>
        <Check color="primary"/>
      </React.Fragment>
    )
  }

  function exFrag() {
    if (st.id !== "ex") {
      return
    }
    return (
      <React.Fragment>
       {locale2.SOMETHING_WENT_WRONG[props.lang]} {st.msg}
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
        <CircularProgress size={24}/>
      </React.Fragment>
    )
  }

  const [cv, setCv] = useState("")
    const [to, setTo] = useState(null)

   function updateOptions(input) {
        if(to) return
        setTo(setTimeout(async () => {
          //let opts = []
          try {
            // opts = await getAllTags(_cv)
            // opts = JSON.parse(opts)
            await props.updateOptions(_cv)
          } catch(ex) {
            console.log(ex)
          } finally {
            setTo(null)
          }
          //opts.push( `Add "${cv}"`)
        }, 1000))
      }

  function editFrag() {
    if (st.id !== "") {
      return
    }
    return (
      <React.Fragment>
        <Autocomplete 
              multiple
                value={value}
                fullWidth
                onChange={(event, newValue) => {
                  if(props.noFreeSolo) {
                    setValue(newValue)
                    return
                  }
                  let cpy = []
                  if(newValue)
                    for(let i = 0; i < newValue.length; i++) {
                      let x = newValue[i].replace(`${locale2.ADD[props.lang]} "`, "")
                      x = x.replace("\"", "")
                      if(cpy.indexOf(x) < 0)
                        cpy.push(x)
                    }
                  setValue(cpy)
                }}
                defaultValue={value}
                filterOptions={(options, params) => {
                  const filtered = options // filter(options, params);
            
                  // Suggest the creation of a new value
                  if (params.inputValue !== '' && !props.noFreeSolo) {
                    filtered.push( `${locale2.ADD[props.lang]} "${params.inputValue}"`);
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
                options={props.options} 
                getOptionSelected={props.equals}
                freeSolo={props.noFreeSolo ? false : true}
                getOptionLabel={props.optionLabel} 
                renderOption={props.renderOption || null}
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
                          <Save/>
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
                          <Cancel/>
                        </Button>
                      </InputAdornment>
                      {c}
                      </React.Fragment>
                  )
                  if(to)
                    params.InputProps.endAdornment = (
                      <CircularProgress style={{width: 30, height: 30}} />
                    )
                  return (
                  <TextField 
                    onKeyPress={(e) => {
                      if(e.key === "Enter") {
                        save()
                      }
                    }}
                    value={cv}
                    onChange={async e => {
                      setCv(e.target.value)
                      _cv = e.target.value
                      if(e.target.value && props.updateOptions) {
                        updateOptions(e.target.value)
                      }
                    }}
                    variant="outlined" 
                    {...params}
                    label={props.atclabel || locale2.ADD_OR_SELECT_CAT[props.lang]} />
                )}}
              />
      </React.Fragment>
    )
  }

  return (
    <React.Fragment>
      {rd && (
        <Grid container spacing={3} alignItems="center">
          <Grid container item xs={4}>
            {props.noLabelTypo && (props.label) || (
              <Typography 
                  style={{color:"gray"}}
                  component={'span'} className={classes.typographyLineHeight}>
                {props.label}
              </Typography>
            )}
          </Grid>
          {
            (props.value && props.value() && (
              <>
                <Grid container item xs={6}>
                    {props.value().map((c, i) => (
                        <Chip key={i} label={props.valueLabel(c)} style={{marginRight: 4}} />
                    ))}
                    {props.placeholder && props.value().length === 0 && <Typography style={{
                      marginTop: 7,
                      color: MulwiColors.subtitleTypography
                    }} variant="body2">
                      {props.placeholder}
                    </Typography>}
                </Grid>
                <Grid container item xs={2}>
                  <Button onClick={beginEdit} color="primary" 
                          size='small' aria-label="edit" 
                          className={classes.moveEditRight}>
                    <EditIcon className={classes.moveEditRight}/>
                  </Button>
                </Grid>
              </>
            )) || (
              <>
                <Grid container item xs={6} />
                <Grid container item xs={2}>
                  <Button onClick={beginEdit} color="primary" size='small' aria-label="edit">
                    <EditIcon/>
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