import React, { useEffect, useState } from 'react'
import CardWithBg from '../card/cardWithBg';
import { makeStyles } from "@mui/styles";
import Stepper from '@mui/material/Stepper';
import Step from '@mui/material/Step';
import StepLabel from '@mui/material/StepLabel';
import Button from '@mui/material/Button';
import Typography from '@mui/material/Typography';
import { StepContent } from '@mui/material';
import { MulwiColors } from '../mulwiColors';
import NameEdit from '../profile/nameEdit';
import { Link } from 'react-router-dom';
import ContactEdit from '../profile/ContactEdit';
import AboutMeEdit from '../profile/AboutMeEdit';
import { locale2 } from '../locale';
import InvoiceEdit from '../profile/invoiceEdit';

export default function Configure(props) {

  const cl = makeStyles({
    step: {
      "&$completed": {
        color: MulwiColors.greenLight
      },
      "&$active": {
        color: MulwiColors.blueDark
      },
      "&$error": {
        color: MulwiColors.redError
      },
    },
    active: {},
    completed: {},
    error: {},
  })()

  let c = (props.instructor && props.instructor.Config) || 0

  function getPayStep() {
    let ok = (c & 2) == 0
    return (<Step completed={ok} expanded={!ok}>
      <StepLabel StepIconProps={{
        classes: {
          root: cl.step,
          active: cl.active,
          completed: cl.completed,
          error: cl.error
        }
      }} error={!ok}>{locale2.ADD_PAYOUT_DATA[props.lang]}</StepLabel>
      {!ok && (
        <StepContent>
          <Typography>{locale2.WITHOUT_PAYOUT_DETAILS[props.lang]}</Typography>
          <Link to="/payments?return_to=configure" style={{
            textDecoration: "none"
          }}>
            <Button size="small" variant="contained" style={{
              color: "white",
              backgroundColor: MulwiColors.blueDark
            }}>
              {locale2.GOTO_PAYMENT[props.lang]}
              </Button>
          </Link>
        </StepContent>
      )}
    </Step>)
  }

  function getInvoiceStep() {
    let ok = (c & 32) == 0
    return (<Step completed={ok} expanded={!ok}>
      <StepLabel StepIconProps={{
        classes: {
          root: cl.step,
          active: cl.active,
          completed: cl.completed,
          error: cl.error
        }
      }} error={!ok}>{locale2.PROVIDE_INVOICE_DATA[props.lang]}</StepLabel>
      {!ok && (
        <StepContent>
          <Typography>{locale2.WITHOUT_INVOICE_DATA[props.lang]}</Typography>
          <InvoiceEdit lang={props.lang} onlyButton
            label={locale2.EDIT[props.lang]}
            buttonProps={{
              variant: "contained",
              style: {
                color: "white",
                backgroundColor: MulwiColors.blueDark
              },
              size: "small"
            }}
            main={props.main}
            instr={props.instructor} />
        </StepContent>
      )}
    </Step>)
  }

  // function getContactStep() {
  //   let ok = (c & 8) == 0
  //   return (<Step completed={ok} expanded={!ok}>
  //     <StepLabel StepIconProps={{
  //       classes: {
  //         root: cl.step,
  //         active: cl.active,
  //         completed: cl.completed,
  //         error: cl.error
  //       }
  //     }} error={!ok}>{locale2.ENTER_CONTACT[props.lang]}</StepLabel>
  //     {!ok && (
  //       <StepContent>
  //         <Typography>{locale2.WITHOUT_CONTACT_DATA[props.lang]}</Typography>
  //         <ContactEdit onlyButton
  //           lang={props.lang}
  //           label={locale2.EDIT[props.lang]}
  //           buttonProps={{
  //             variant: "contained",
  //             style: {
  //               color: "white",
  //               backgroundColor: MulwiColors.blueDark
  //             },
  //             size: "small"
  //           }}
  //           main={props.main}
  //           user={props.user} />
  //       </StepContent>
  //     )}
  //   </Step>)
  // }

  function getAboutMeStep() {
    let ok = (c & 16) == 0
    return (<Step completed={ok} expanded={!ok}>
      <StepLabel StepIconProps={{
        classes: {
          root: cl.step,
          active: cl.active,
          completed: cl.completed,
          error: cl.error
        }
      }} error={!ok}>{locale2.WRITE_ABOUT_YOURSELF[props.lang]}</StepLabel>
      {!ok && (
        <StepContent>
          <Typography>{locale2.GIVE_INFO_TO_CLIENTS[props.lang]}</Typography>
          <AboutMeEdit lang={props.lang} onlyButton
            label={locale2.EDIT[props.lang]}
            buttonProps={{
              variant: "contained",
              style: {
                color: "white",
                backgroundColor: MulwiColors.blueDark
              },
              size: "small"
            }}
            main={props.main}
            user={props.user} />
        </StepContent>
      )}
    </Step>)
  }

  const [open, setOpen] = useState(false)

  function getNameStep() {
    let ok = (c & 1) == 0
    return (<Step completed={ok} expanded={!ok}>
      <StepLabel StepIconProps={{
        classes: {
          root: cl.step,
          active: cl.active,
          completed: cl.completed,
          error: cl.error
        }
      }} error={!ok}>{locale2.ENTER_YOUR_NAME[props.lang]}</StepLabel>
      {!ok && (
        <StepContent>
          <Button size="small" variant="contained"
            onClick={() => setOpen(true)}
            style={{
              backgroundColor: MulwiColors.blueDark,
              color: "white"
            }}>{locale2.EDIT[props.lang]}</Button>
          <NameEdit external
            lang={props.lang}
            main={props.main}
            user={props.user}
            setOpen={setOpen}
            open={open} />
        </StepContent>
      )}
    </Step>)
  }

  return (
    <React.Fragment>
      <CardWithBg header={locale2.CONFIGURE_INSTR_ACC[props.lang]} 
          img="/static/form-backgrounds/moto.webp">
        <Stepper activeStep={-1} style={{
          maxWidth: 500
        }} variant="outlined" orientation="vertical">
          {getNameStep()}
          {/* {getPayStep()} */}
          {/* {getContactStep()} */}
          {getAboutMeStep()}
          {/* {getInvoiceStep()} */}
        </Stepper>
        {c == 0 && (
          <React.Fragment>
            <br />
            <Typography style={{
              color: MulwiColors.greenDark
            }}>
              <strong>{locale2.EVERYTHING_IS_READY[props.lang]}</strong>
            </Typography>
            <br />
            <Typography>
              {locale2.REFER_CREATE_YOUR_OFFER[props.lang]}
            </Typography>
            <br />
            <Link to="/harmonogram?open_add_train=1" style={{
              textDecoration: "none"
            }}>
              <Button variant="contained" style={{
                color: "white",
                backgroundColor: MulwiColors.greenDark
              }}>{locale2.CREATE_YOUR_OFFER[props.lang]}</Button>
            </Link>
          </React.Fragment>
        )}
      </CardWithBg>
    </React.Fragment>
  )
}