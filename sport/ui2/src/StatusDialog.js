import {
  Button, Dialog, DialogActions,
  DialogContent, DialogContentText,
  DialogTitle,
  Typography
} from '@mui/material'
import React from 'react'
import { useHistory } from 'react-router'
import { redirectToLogin, rmtoken } from './helpers'
import { MulwiColors } from './mulwiColors'
import { locale2, getSupportedLanguage } from './locale'

export function getNullDialog() {
  return {
    open: false,
    hdr: "",
    msg: "",
    c: null,
    onClose: null,
    buttons: null
  }
}

export const supportEmail = "support@veidly.com"


const mailpart = (lang) =>  (<Typography>
  {locale2.IN_CASE_OF_QUESTIONS_SUPPORTS[lang]} <strong>{supportEmail}</strong>
</Typography>)

function RedirectToLogin(props) {
  const h = useHistory()
  return (<React.Fragment><center>
    <Typography variant="h6" style={{
      color: MulwiColors.redError
    }}>
      {locale2.IT_SEEMS_YOU_WERE_LOGGED_OUT[props.lang]}
    </Typography>
    <Typography>
      {locale2.IF_YOU_WANT_TO_PERFORM_THIS_ACTION[props.lang]}
    </Typography>
    <Button onClick={() => redirectToLogin(null, h)}>
      {locale2.LOGIN_AGAIN[props.lang]}
    </Button>
    {mailpart(props.lang)}
  </center>
  </React.Fragment>)
}

export function errToStr(err) {

  /*
   err to string is used in far too many cases to pass lang as arg,
   so im just grabbing lang from storage. it shouldnt be issue tho since 
   this function is mostly used in popup
  */
 let lang = getSupportedLanguage()

  let errm = ""

  if (err == null) {
    errm = locale2.UNKOWN_ERROR_OCCURRED[lang]
  } else {
    if (typeof err === 'object') {
      err = JSON.stringify(err)
    }

    console.log(err)

    switch (Number(err)) {
      case 400:
        return (<React.Fragment>
          <center>
            <Typography variant="h6" style={{
              color: MulwiColors.redError
            }}>
              {locale2.INVALID_DATA[lang]}
            </Typography>
            <Typography>
              {locale2.MAKE_SURE_INPUT_IS_OK[lang]}
            </Typography>
            {mailpart(lang)}
          </center>
        </React.Fragment>)
      case 401:
        return <RedirectToLogin lang={lang} />
      case 409:
        return <center>
            <Typography>
                {locale2.CONFLICT[lang]}
            </Typography>
          </center>
      case 413:
        return (
          <>
            <center>
              <Typography>
                {locale2.IMAGE_IS_TOO_BIG[lang]}
              </Typography>
            </center>
          </>
        )

      case 500:
        return (<React.Fragment><center>
          <Typography variant="h6" style={{
            color: MulwiColors.redError
          }}>
            {locale2.SERVER_IS_NOT_ABLE_TO_PROCESS[lang]}
          </Typography>
          <Typography>
            {locale2.PROBLEM_SHOULD_BE_SOLVED_SOON[lang]}
          </Typography>
          {mailpart(lang)}
        </center>
        </React.Fragment>)
      default:
        errm = locale2.UNKOWN_ERROR_OCCURRED[lang] + " (" + err + ")"
        break
    }

  }

  return <Typography>{errm}</Typography>
}

// this adds fatal error button
export function getFatalErrorDialog(hdr, err) {
  let lang = getSupportedLanguage()
  return {
    open: true,
    hdr: hdr,
    msg: errToStr(err),
    onClose: null,
    buttons: (
      <React.Fragment>
        <Button onClick={() => {
          rmtoken()
          window.location = "/"
        }} color="primary">
          {locale2.RESET_APP[lang]}
        </Button>
      </React.Fragment>
    )
  }
}

export function getInfoDialog(hdr, msg, btns, c) {
  return {
    open: true,
    hdr: hdr,
    msg: msg,
    onClose: null,
    buttons: btns,
    c: c
  }
}

export function getErrorDialog(hdr, err, btns, c) {
  return getInfoDialog(hdr, errToStr(err), btns, c)
}

export function getDialogWithOptions(hdr, content, buttons, hideExit) {
  return {
    open: true,
    hdr: hdr,
    c: content,
    onClose: null,
    buttons: buttons,
    hideExit: hideExit
  }
}

export function StatusDialog(props) {

  function resetDialog() {
    props.info.onClose && props.info.onClose()
    props.setInfo(getNullDialog())
  }

  return (
    <Dialog open={props.info.open} onClose={resetDialog}>
      <DialogTitle>{props.info.hdr}</DialogTitle>
      <DialogContent>
        {props.info.c}
        <DialogContentText>
          {props.info.msg}
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        {props.info.buttons}
        {!props.info.hideExit && (
          <Button onClick={resetDialog} color="secondary">
            {locale2.CLOSE[props.lang]}
          </Button>
        )}
      </DialogActions>
    </Dialog>)
}