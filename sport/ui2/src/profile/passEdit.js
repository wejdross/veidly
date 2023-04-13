import {
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  Grid,
  TextField
} from '@mui/material'
import React, {useEffect, useState} from 'react'
import EditIcon from '@mui/icons-material/Edit';
import {patchUserPassword} from '../apicalls/user.api';
import {Check, Close} from '@mui/icons-material';
import {Link} from 'react-router-dom'
import { MulwiColors } from '../mulwiColors';
import { locale2 } from '../locale';
import { makeStyles } from "@mui/styles";


export default function PassEdit(props) {

  const useStyles = makeStyles({
    padded: {
      marginLeft: 100, marginRight: 100, marginTop: 50, marginBottom: 50
    }
  })
  const classes = useStyles()

  const [dialogOpen, setDialogOpen] = useState(false)

  function openPassEdit() {
    setPassEdit({oldPass: "", confPass: "", pass: ""})
    setPassEditSt({st: "", msg: ""})
    setDialogOpen(true)
  }

  const [passEdit, setPassEdit] = useState({oldPass: "", confPass: "", pass: ""})
  const [passEditSt, setPassEditSt] = useState({st: "", msg: ""})

  const [err, setErr] = useState({oldPassErr: false, passErr: false})

  function validate(forceErr) {
    let e = {oldPassErr: false, passErr: false}
    let res = true
    if (!passEdit.oldPass) {
      if (forceErr)
        e.oldPassErr = true
      res = false
    }
    if (!passEdit.pass) {
      if (forceErr)
        e.passErr = true
      res = false
    }
    if (passEdit.confPass !== passEdit.pass) {
      e.passErr = true
      res = false
    }
    if (err.oldPassErr !== e.oldPassErr || err.passErr !== e.passErr)
      setErr(e)
    return res
  }

  useEffect(() => {
    validate(false)
  })

  function passEditForm() {
    return passEditSt.st === "" && (
      <React.Fragment>
        <DialogContent>
          <DialogContentText>
            {locale2.SET_NEW_PASSWORD[props.lang]}
          </DialogContentText>
          <TextField value={passEdit.oldPass}
                     onChange={(e) => {
                       setPassEdit(c => ({...c, oldPass: e.target.value}))
                     }}
                     autoFocus label={locale2.OLD_PASSWORD[props.lang]}
                     error={err.oldPassErr}
                     type="password" fullWidth/>
          <Link to="/forgot_password" style={{textDecoration: "none"}}>
            {locale2.FORGOT_PASSWORD[props.lang]}
          </Link>
          <TextField value={passEdit.pass}
                     onChange={(e) => {
                       setPassEdit(c => ({...c, pass: e.target.value}))
                     }}
                     error={err.passErr}
                     label={locale2.NEW_PASSWORD[props.lang]} type="password" fullWidth/>
          <TextField value={passEdit.confPass}
                     onChange={(e) => {
                       setPassEdit(c => ({...c, confPass: e.target.value}))
                     }}
                     error={err.passErr}
                     label={locale2.CONFIRM_NEW_PASSWORD[props.lang]} type="password" fullWidth/>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setDialogOpen(false)} color="secondary">
            {locale2.CANCEL[props.lang]}
          </Button>
          <Button onClick={savePass} color="primary">
            {locale2.SAVE[props.lang]}
          </Button>
        </DialogActions>
      </React.Fragment>
    )
  }

  function passEditWaiter() {
    return passEditSt.st === "wip" && (
      <CircularProgress className={classes.padded}/>
    )
  }

  function passEditOK() {
    return passEditSt.st === "ok" && (
      <Check color="primary" fontSize="large" className={classes.padded}/>
    )
  }

  function passEditEX() {
    return passEditSt.st === "ex" && (
      <React.Fragment>
        <DialogContent>
          <DialogContentText color="secondary" align="center">
            {locale2.SOMETHING_WENT_WRONG[props.lang]}
          </DialogContentText>
          <DialogContentText color="secondary" align="center">
            {passEditSt.msg}
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button color="primary" onClick={() => setPassEditSt({st: "", msg: ""})}>
            {locale2.ONCE_AGAIN[props.lang]}
          </Button>
          <Button onClick={() => setDialogOpen(false)} color="secondary">
            <Close/>
          </Button>
        </DialogActions>
      </React.Fragment>
    )
  }

  async function savePass() {
    if (!validate(true)) return
    try {
      setPassEditSt({st: "wip", msg: ""})
      await patchUserPassword({
        OldPassword: passEdit.oldPass,
        NewPassword: passEdit.pass
      })
      setPassEditSt({st: "ok"})
      setTimeout(async () => {
        setDialogOpen(false)
      }, 100)
    } catch (ex) {
      setPassEditSt({
        st: "ex",
        msg: ex
      })
    }
  }

  return (
    <React.Fragment>
      <Grid
        container
        spacing={3}
      >
        <Grid item xs={4} style={{ color: MulwiColors.subtitleTypography }}>
          {locale2.PASSWORD[props.lang]}
        </Grid>
        {(props.user && props.user.OauthProvider && (
          <>
            <Grid item xs={8}>
              {locale2.PASSWORD_MANAGED_BY_GOOGLE[props.lang]}
            </Grid>
          </>
        )) || (

            <React.Fragment>
                <Grid item xs={6}>
                    <Dialog open={dialogOpen} onClose={() => setDialogOpen(false)} 
                            aria-labelledby="form-dialog-title">
                        <DialogTitle id="form-dialog-title">
                          {locale2.PASSWORD[props.lang]}
                        </DialogTitle>
                        {passEditForm()}
                        {passEditWaiter()}
                        {passEditOK()}
                        {passEditEX()}
                    </Dialog>
                    ********
                </Grid>
                <Grid item xs={2}>
                    <Button onClick={openPassEdit} color="primary" size='small' 
                              style={{marginLeft: 5}} aria-label="edit">
                        <EditIcon/>
                    </Button>
                </Grid>
            </React.Fragment>
        )}
      </Grid>

    </React.Fragment>
  )
}