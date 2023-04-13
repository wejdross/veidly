import {Button, CircularProgress, Grid, 
        InputAdornment, 
        TextField, Typography} from '@mui/material'
import React, {useState} from 'react'
import EditIcon from '@mui/icons-material/Edit';
import {Cancel, Check, Save} from '@mui/icons-material';
import { locale2 } from '../locale';
import { makeStyles } from "@mui/styles";

export default function GenericInlineEdit(props) {

  const [rd, setRd] = useState(true)
  const [ev, setEv] = useState("")
  const [st, setSt] = useState({id: "", msg: ""})

  function beginEdit() {
    setEv(props.value)
    setSt({id: "", msg: ""})
    setRd(false)
  }

  async function save() {
    try {
      setSt({id: "wip", msg: ""})
      await props.onChange(ev)
      setRd(true)
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
      return null
    }
    return (
      <React.Fragment>
        ${locale2.SOMETHING_WENT_WRONG[props.lang]} {String(st.msg)}
        <Button
          onClick={() => beginEdit()} color="primary" aria-label="edit">
          ${locale2.ONCE_AGAIN[props.lang]}
        </Button>
        <Button
          onClick={() => setRd(true)} color="secondary" aria-label="edit">
          ${locale2.CANCEL[props.lang]}
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

  function editFrag() {
    if (st.id !== "") {
      return
    }
    return (
      <React.Fragment>
        <TextField value={ev} variant="outlined" 
                  onChange={(e) => setEv(e.target.value)} 
                  fullWidth
                   size="small"
                   InputProps={{
                     startAdornment: (
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
                     ),
                   }}/>
      </React.Fragment>
    )
  }
  const useStyles = makeStyles({
    typographyLineHeight: {
      lineHeight: 2,
    },
  })
  const classes = useStyles()

  return (
    <React.Fragment>
      {rd && (
        <Grid container spacing={3}>
          <Grid container item xs={4}>
            <Typography 
                style={{color:"gray"}}
                component={'span'} className={classes.typographyLineHeight}>
              {props.label}
            </Typography>
          </Grid>
          {
            (props.value && (
              <>
                <Grid container item xs={6}>
                  <Typography  noWrap className={classes.typographyLineHeight}>
                    {props.value}
                  </Typography>
                </Grid>
                <Grid container item xs={2}>
                  <Button onClick={beginEdit} color="primary" size='small' 
                          aria-label="edit" className={classes.moveEditRight}>
                    <EditIcon className={classes.moveEditRight}/>
                  </Button>
                </Grid>
              </>
            )) || (
              <>
                <Grid container item xs={6} />
                <Grid container item xs={2}>
                  <Button onClick={beginEdit} color="primary" 
                          size='small' aria-label="edit">
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