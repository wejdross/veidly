import {
  Button,
  Divider,
  Grid,
  Typography,
  useTheme
} from '@mui/material';
import { makeStyles } from "@mui/styles";

import React, { useEffect, useState } from 'react';
import { useHistory } from 'react-router';
import { canDeleteInstructor as apiCanDeleteInstructor, 
          deleteInstructor, PATCHInstructor } from '../apicalls/instructor.api';
import { deleteUser } from '../apicalls/user.api';
import CardWithBg from '../card/cardWithBg';
import { MulwiColors } from '../mulwiColors';
import { getDialogWithOptions, getErrorDialog,
        getNullDialog, StatusDialog } from '../StatusDialog';
import { locale2 } from '../locale';

  export default function ManageCard(props) {
  
    const theme = useTheme()
  
    const [canDeleteInstructor, setCanDeleteInstructor] = useState(false)
  
    const useStyles = makeStyles({
      divider: {
        marginTop: 10,
        marginBottom: 10,
      },
      widthSettings: {
        [theme.breakpoints.up('sm')]: {
          minWidth: 500,
        }
      },
      purple: {
        backgroundColor: "purple"
      }
    })
    const classes = useStyles()
  
    useEffect(() => {
      (async function (){
        try {
          await apiCanDeleteInstructor()
          setCanDeleteInstructor(true)
        } catch(ex) {
          setCanDeleteInstructor(false)
        }
      })()
    }, [props.instructor])
  
    const [info, setInfo] = useState(getNullDialog())
  
    async function patchActivateInstructor(disabled) {
      setInfo(getDialogWithOptions((disabled ? 
          locale2.DEACTIVATION[props.lang] 
            : 
          locale2.ACTIVATION[props.lang]) 
            + locale2.OF_ACCOUNT[props.lang], 
      <React.Fragment>
        {disabled ? (
          <React.Fragment>
            <p style={{
          maxWidth: 500
            }}>{locale2.DEACTIVATION_WARN[props.lang]}</p>
          <p style={{
            color: MulwiColors.redError
          }}><strong>{locale2.EXISTING_RSV_WARN[props.lang]}</strong></p>
          </React.Fragment>
        ) : (
          <p style={{
            maxWidth: 500
          }}>{locale2.ACTIVATION_WARN[props.lang]}</p>
        )}
      </React.Fragment>, 
      <React.Fragment>
        <Button variant="contained" style={{
          color: "white",
          backgroundColor: MulwiColors.blueDark
        }} onClick={async () => {if(!props.instructor) return
            try {
              let i = {...props.instructor}
              i.Disabled = disabled
              await PATCHInstructor(i)
              props.main.refreshInstructor()
              setInfo(getNullDialog())
            } catch(ex) {
              setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
            }}}>
          {disabled ? locale2.DEACTIVATE[props.lang] 
              : locale2.ACTIVATE[props.lang]} {locale2.ACCOUNT[props.lang]}
        </Button>
      </React.Fragment>))
    }
  
    function delUserAccWarn() {
      return (<p style={{
        color: MulwiColors.redError,
        maxWidth: 600
      }}><strong>{locale2.DELETE_ACC_WARN[props.lang]}</strong></p>)
    }

    const h = useHistory()
  
    async function deleteUserAccount() {
      setInfo(getDialogWithOptions(locale2.DELETING_ACCOUNT[props.lang], 
      <React.Fragment>
        {delUserAccWarn()}
      </React.Fragment>, 
      <React.Fragment>
        <Button variant="contained" style={{
          color: "white",
          backgroundColor: MulwiColors.redError
        }} onClick={async () => {
            if(!props.user) return
            try {
              await deleteUser()
              props.main.refresh()
              //setPassword("")
              h.push("/")
            } catch(ex) {
              setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
            }
          }}>
           {locale2.DELETE_ACCOUNT[props.lang]}
        </Button>
      </React.Fragment>))
    }
    async function deleteInstructorAccount() {
      setInfo(getDialogWithOptions(locale2.DELETING_ACCOUNT[props.lang], 
      <React.Fragment>
       {delInstrWarnFragment()}
      </React.Fragment>, 
      <React.Fragment>
        <Button variant="contained" style={{
          color: "white",
          backgroundColor: MulwiColors.redError
        }} onClick={async () => {
              if(!props.instructor) return
              try {
                await deleteInstructor()
                props.main.refreshInstructor()
                //setPassword("")
                setInfo(getNullDialog())
              } catch(ex) {
                setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
              }
          }}>
           {locale2.DELETE_ACCOUNT[props.lang]}
        </Button>
      </React.Fragment>))
    }
  
    async function deleteAllAccounts() {
      setInfo(getDialogWithOptions(locale2.DELETING_ACCOUNT[props.lang], 
      <React.Fragment>
       {delInstrWarnFragment()}
      </React.Fragment>, 
      <React.Fragment>
        <Button variant="contained" style={{
          color: "white",
          backgroundColor: MulwiColors.redError
        }} onClick={async () => {
              if(!props.user || !props.instructor) return
              try {
                await deleteInstructor()
                await deleteUser()
                props.main.refresh()
                h.push("/")
                setInfo(getNullDialog())
              } catch(ex) {
                setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
              }
          }}>
           {locale2.DELETE_ACCOUNT[props.lang]}
        </Button>
      </React.Fragment>))
    }
  
    function delInstrWarnFragment() {
      return <React.Fragment>
        <p style={{
          maxWidth: 400
        }}>{locale2.YOU_MAY_DEACITVATE_INSTR_ACC_[props.lang]}</p>
        <p style={{
          maxWidth: 400
        }}>{locale2.DEACTIVATE_INSTR_WARN[props.lang]}</p>
        <p style={{
          maxWidth: 500,
        }}>
          {locale2.IF_YOU_TRY_TO_RECREATE[props.lang]}
        </p>
        <p style={{
          color: MulwiColors.redError,
          maxWidth: 500,
        }}><strong>{locale2.TO_DEACTIVATE_INSTR_ACC_[props.lang]}</strong></p>
      </React.Fragment>
    }
  
    return (
      <React.Fragment>
      <StatusDialog lang={props.lang} info={info} setInfo={setInfo}/>
        <CardWithBg 
        style={{marginBottom: 50}}
        header={locale2.MGMT[props.lang]}
        >            
            {props.instructor && (
            <React.Fragment>
              <Typography variant="h6" align="center">
                {props.instructor.Disabled ? locale2.REACTIVATING[props.lang]
                        : locale2.DEACTIVATING[props.lang]} {locale2.OF_ACCOUNT[props.lang]}
              </Typography>
              <Grid container direction="column" spacing={2}>
              <Grid item style={{marginTop: 10}}>
                {props.instructor.Disabled ? (
                  <Button variant="contained" style={{
                    color: "white",
                    backgroundColor: MulwiColors.blueDark
                  }}  onClick={() => {
                    patchActivateInstructor(false)
                  }}>
                    {locale2.REACTIVATE_ACCOUNT[props.lang]}
                  </Button>
                ) : (<React.Fragment>
                  <p>{locale2.AT_ANY_MOMENT_YOU_CAN_DEACTIVATE[props.lang]}</p>
                  <p style={{
                    maxWidth: 600
                  }}>{locale2.DEACTIVATION_EFFECTS[props.lang]}</p>
                  <p style={{
                    color: MulwiColors.redError
                  }}><strong>{locale2.EXISTING_RSV_WARN[props.lang]}</strong></p>
                  <Button variant="contained" style={{
                  color: "white",
                  backgroundColor: MulwiColors.blueDark
                }} onClick={() => {
                  patchActivateInstructor(true)
                }}>
                  {locale2.TEMPORARILY_DEACTIVATE_INSTR[props.lang]}
                </Button>
                </React.Fragment>)}
              </Grid>
  
              <Divider style={{marginTop: 10}}/>
              <Typography variant="h6" align="center">
                {locale2.DELETING_ACCOUNT[props.lang]}
              </Typography>
  
              {canDeleteInstructor && (
                <React.Fragment>
                  <Grid item>
                    {delInstrWarnFragment()}
                  <Button variant="contained" style={{
                    color: "white",
                    backgroundColor: MulwiColors.redError
                  }} onClick={deleteInstructorAccount}>
                    {locale2.STOP_BEING_INSTR[props.lang]}
                  </Button>
                </Grid><Grid item>
                  <Button variant="contained" style={{
                    color: "white",
                    backgroundColor: MulwiColors.redError
                  }} onClick={deleteAllAccounts}>
                    {locale2.COMPLETELY_DELETE_ACC[props.lang]}
                  </Button>
                </Grid>
                </React.Fragment>
              )}
  
              {!canDeleteInstructor && (
                <label style={{
                  color: MulwiColors.redError
                }}>{locale2.CANNOT_DELETE_ACCOUNT_YET[props.lang]}</label>
              )}
            </Grid>
            </React.Fragment>)}
  
            {!props.instructor && (
              <Grid item>
                  <p>{locale2.AT_ANY_MOMENT_YOU_MAY_DELETE_USERs_ACCOUNT[props.lang]}</p>
                  {delUserAccWarn()}
                <Button variant="contained" style={{
                  color: "white",
                  backgroundColor: MulwiColors.redError
                }} onClick={deleteUserAccount}>
                  {locale2.DELETE_USER_ACC[props.lang]}
                </Button>
              </Grid>
            )}
  
        </CardWithBg>
      </React.Fragment>
  
    );
  }