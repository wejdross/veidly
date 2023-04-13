import {
  Avatar, Button, CircularProgress,
  FormControl, Grid,
  TextField,
  Typography
} from '@mui/material';
import Divider from '@mui/material/Divider';
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import validator from 'validator';
import { getOauthGoogleUrl, resendRegisterEmail, userRegister, validatPass } from '../apicalls/user.api';
import StickyFooter from '../Footer';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';
import { defaultRedirect } from './login';
import '../veidly-styles.css'


export default function Register(props) {

  const [withInstr, setWithInstr] = useState(props.withInstructor || false)
  const [emailErr, setEmailErr] = useState(false)
  const [passErr, setPassErr] = useState(false)
  const [workInProgress, setWorkInProgress] = useState(false)
  const [registerRes, setRegisterRes] = useState(null)
  const [cPass, setCPass] = useState("")
  const [resendEmailMsg, setResendEmailMsg] = useState("")

  const [passErrText, setPassErrText] = useState(null)

  const [userRequest, setUserRequest] = useState({
    email: "",
    password: "",
    language: props.lang
  })

  async function vp(p) {
    try {
      await validatPass(p)
      return true
    } catch {
      return false
    }
  }
  useEffect(() => {
    window.scrollTo(0, 0)
  }, [])
  // if force error is defined, all controls will be invalidated even if they are empty
  async function validateForm(forceErr) {
    let e = true
    if (!validator.isEmail(userRequest.email)) {
      if (userRequest.email || forceErr)
        setEmailErr(true)
      e = false
    } else {
      setEmailErr(false)
    }

    if (userRequest.password && !await vp(userRequest.password)) {
      setPassErr(true)
      setPassErrText(locale2.INVALID_PASS_ERR[props.lang])
      return false
    } else {
      setPassErrText(null)
    }

    if (!userRequest.password || userRequest.password !== cPass) {
      if (userRequest.password || cPass || forceErr)
        setPassErr(true)
      e = false
    } else {
      setPassErr(false)
    }

    return e
  }

  async function resendEmail() {
    setWorkInProgress(true)
    try {
      let returl = "";
      if (withInstr) {
        returl = "/become_trainer";
      }
      await resendRegisterEmail(userRequest.email, returl);
      setResendEmailMsg(locale2.RESENT_EMAIL[props.lang])
    } catch (ex) {
      setResendEmailMsg(locale2.COULDNT_RESEND_EMAIL[props.lang])
    } finally {
      setWorkInProgress(false)
    }
  }

  useEffect(() => {
    validateForm()
  })

  async function register(e) {
    if (!await validateForm(true)) return
    e.preventDefault()
    setWorkInProgress(true)
    try {
      let returl = ""
      if (withInstr) {
        returl = "/become_trainer"
      }
      userRequest.language = props.lang
      let res = await userRegister(userRequest, returl)
      setRegisterRes({
        res: JSON.parse(res),
        msg: ""
      })
    } catch (ex) {
      setRegisterRes({
        res: null,
        msg: locale2.COULDNT_REGISTER[props.lang]
      })
    } finally {
      setWorkInProgress(false)
    }
  }


  async function registerGoogle() {
    try {
      let returl = defaultRedirect
      if (withInstr) {
        returl = "/become_trainer"
      }
      let url = await getOauthGoogleUrl(returl);
      window.location.replace(url);
    } catch (ex) {
      console.log(ex)
    }
  }

  function confirmEmailFragment() {
    if (!registerRes || !registerRes.res || !registerRes.res.mfa) {
      return null
    }
    return (
      <React.Fragment>
        <Typography variant="h4" style={{
          marginTop: 30,
          color: MulwiColors.greenDark,
        }}
          align={"center"}
        >
          {locale2.SUCCESS[props.lang]}
        </Typography>
        <br />
        <Typography>
          {locale2.SENT_EMAIL_TO[props.lang]} <strong>{userRequest.email}</strong>
        </Typography>
        <Typography>{locale2.WILL_BE_VALID_FOR[props.lang]} {registerRes.res.ttl_seconds / 60} {locale2.MINUTES[props.lang]}</Typography>
        <br />
        <Typography>{locale2.IF_DIDNT_FIND_EMAIL[props.lang]}</Typography>
        <br />
        <div style={{ marginTop: 20, marginBottom: 20 }}>
          <Divider style={{ fontSize: 12, fontWeight: 400 }}>{locale2.OR[props.lang]}</Divider>
        </div>

        <Button variant="contained" style={{
          color: "white",
          backgroundColor: MulwiColors.blueDark
        }}
          onClick={resendEmail}>
          {locale2.SEND_EMAIL_AGAIN[props.lang]}
        </Button>
        <br />
        <Typography>
          {resendEmailMsg}
        </Typography>
        <br />
        <Typography>
          {locale2.YOU_MAY_CLOSE_THIS_WINDOW[props.lang]}
        </Typography>
      </React.Fragment>
    )
  }

  function registerConfirmFragment() {
    if (!registerRes || !registerRes.res || registerRes.res.mfa) {
      return null
    }
    return (
      <React.Fragment>
        <Typography variant="h4" style={{
          marginTop: 30,
          color: MulwiColors.greenDark
        }}>
          {locale2.SUCCESS[props.lang]}
        </Typography>
        <p>Register was successful and you may proceed to <a href="#/login">login</a> now</p>
        <p>You are on development version of the application so no account confirmation is required</p>
      </React.Fragment>
    )
  }

  function formFragment() {
    if (registerRes && registerRes.res) {
      return null
    }
    return (
      <FormControl id="registerForm">
        <Typography variant="h7" style={{ marginTop: 20, fontWeight: 400 }} align={"center"}>
          {locale2.CREATE_ACCOUNT[props.lang]}
        </Typography>
        <Divider style={{
          marginTop: 10,
          marginBottom: 30,
        }} />
        <Typography variant="h6" style={{
          fontSize: 22,
          marginBottom: 20,
          fontWeight: 600
        }}>
          {locale2.WELCOME_TO_VEIDLY[props.lang]}
        </Typography>
        <Grid container direction="row" spacing={2}>
          <Grid item xs={12}>
            <TextField
              variant='outlined'
              InputLabelProps={{
                classes: {
                  focused: "focused",
                }
              }}
              error={emailErr}
              style={{ marginBottom: 20, width: "100%" }}
              onChange={(e) => {
                setUserRequest(c => ({ ...c, email: e.target.value }))
              }}
              type="email"
              value={userRequest.email}
              label="Email" />
          </Grid>
          <Grid item md={6} sm={12} xs={12}>
            <TextField type="password" style={{ margin: 0 }}
              InputLabelProps={{
                classes: {
                  focused: "focused",
                }
              }}
              variant='outlined'
              value={userRequest.password}
              fullWidth
              onChange={(e) => {
                setUserRequest(c => ({ ...c, password: e.target.value }))
              }}
              error={passErr}
              label={locale2.PASSWORD[props.lang]} />
          </Grid>
          <Grid item md={6} sm={12} xs={12}>
            <TextField type="password" style={{ margin: 0 }}
              InputLabelProps={{
                classes: {
                  focused: "focused",
                }
              }}
              fullWidth
              variant='outlined'
              value={cPass}
              onChange={(e) => {
                setCPass(e.target.value)
              }}
              error={passErr}
              label={locale2.CONFIRM_PASSWORD[props.lang]} />
          </Grid>
        </Grid>
        {passErrText && (
          <Typography variant='body2' style={{
            marginTop: 15,
            marginBottom: 15,
            color: MulwiColors.redError
          }}>
            {passErrText}
          </Typography>
        ) || (
            <Typography variant='body2'>
              &nbsp;&nbsp;
            </Typography>
          )}
        <Typography>
          {locale2.CREAT_ACC_AGREEMENT[props.lang]} <Link to="/support" style={{ textDecoration: "none", color: MulwiColors.blueDark }}>
            {locale2.TERMS_OF_SERVICE[props.lang]}
          </Link>
        </Typography>
        <Typography>
          {locale2.AND_I_ACK[props.lang]}
          <a href="/polityka-prywatnosci.pdf" target='_blank' style={{ textDecoration: "none", color: MulwiColors.blueDark }}>
            {" "}{locale2.PRIVACY_POLICY[props.lang]}
          </a>
        </Typography>

        <Button variant="contained" color="primary" onClick={register} disabled={passErr}
          style={{
            marginTop: 30,
            backgroundColor: passErr ? "grey" : MulwiColors.blueDark
          }}>
          {(workInProgress && (
            <CircularProgress style={{ padding: 5, width: 30, height: 30 }} color="secondary" />
          )) || (
              locale2.REGISTER[props.lang]
            )}
        </Button>

        <Typography color="error">
          {registerRes && registerRes.msg}
        </Typography>

        <Divider style={{ marginTop: 40, fontSize: 12, fontWeight: 400 }}>{locale2.OR[props.lang]}</Divider>

        <Button
          onClick={registerGoogle}
          style={{ marginTop: 10, backgroundColor: "white" }}
          variant="contained" color="secondary">
          <Grid container direction={"row"} alignContent={"flex-start"} justify={"flex-start"}>
            <Grid item xs={5}>

              <Avatar src={"google.png"} style={{ height: 18, width: 18 }} />
            </Grid>
            <Typography style={{ marginLeft: 5, color: "Black", lineHeight: "1.2" }} align={'center'}>
              Google
            </Typography>
          </Grid>
        </Button>
      </FormControl>
    )
  }

  return (
    <>
      <Grid container
        direction="column"
        alignContent='center'
        justifyContent='center'
      >
        <div className='veidlyDataSheet'>
          <Grid container direction="column" alignContent='center' justifyContent='center'>
            {confirmEmailFragment()}
            {formFragment()}
            {registerConfirmFragment()}
          </Grid>
        </div>
      </Grid>
    </>
  )
}