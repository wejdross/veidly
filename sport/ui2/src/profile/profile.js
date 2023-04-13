import {
  Avatar, Divider,
  Grid,
  Typography,
  useTheme
} from '@mui/material';
import { makeStyles } from "@mui/styles";
import React, { useEffect } from 'react';
import { Link } from 'react-router-dom';
import CardWithBg from '../card/cardWithBg';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';
import AboutMeEdit from "./AboutMeEdit";
import AvatarEdit from './avatarEdit';
import InvoiceEdit from './invoiceEdit';
import NameEdit from './nameEdit';
import PassEdit from './passEdit';
import UrlsEdit from './UrlsEdit';

export function AvatarContainer(props) {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  const useStyles = makeStyles({
    purple: {
      backgroundImage: "static/no-profile-photo.jpg"
    }
  })

  const classes = useStyles()

  if (!props.user)
    return null

  return (
    <div style={{
      width: "100%",
    }}>
      <center>

      {(props.user.AvatarUrl && (
        <Avatar style={{
          width: props.large ? 160 : 80,
          height: props.large ? 160 : 80
        }} src={props.user.AvatarUrl} alt="-" />
        )) || ((
          (
            <Avatar className={classes.purple} style={{
              width: props.large ? 160 : 80,
              height: props.large ? 160 : 80
            }}>
          </Avatar>
        )
        )) || (
          <Avatar className={classes.purple} style={{
            width: props.large ? 160 : 80,
            height: props.large ? 160 : 80
          }}>
          </Avatar>
        )}
        </center>
      </div>
  )
}

export default function Profile(props) {

  const theme = useTheme()

  const useStyles = makeStyles({
    divider: {
      marginTop: 10,
      marginBottom: 10,
    },
    widthSettings: {
      [theme.breakpoints.up('sm')]: {
        minWidth: 600,
      }
    },
    purple: {
      backgroundColor: "purple"
    }
  })

  const classes = useStyles()

  if (!props.user) return null

  return (
      <CardWithBg
      header={locale2.MY_DATA[props.lang]}
      
      >
          {props.instructor && <React.Fragment>
            <center><Typography
              variant="body2"
              component="p"
              align="center"
              style={{
                maxWidth: 400
              }} color="primary">
              {locale2.MY_DATA_INSTR[props.lang]} <Link to="/instr_profile" style={{
                textDecoration: "none",
                color: MulwiColors.blueDark
              }}>
                <strong>{locale2.MY_DATA_INSTR_PROFILE[props.lang]}</strong>
              </Link></Typography> </center>
            <br />
          </React.Fragment>}

          <Grid direction={'row'} alignItems={'center'} justifyItems={'center'} style={{width: "100%"}}>
            <AvatarContainer user={props.user} />
          </Grid>


          <NameEdit
            lang={props.lang} user={props.user} main={props.main} />
          <Divider className={classes.divider} />

          <PassEdit lang={props.lang} user={props.user} main={props.main} />
          <Divider className={classes.divider} />

          <AvatarEdit lang={props.lang} main={props.main} refreshInstr />


          {/* </Grid></Grid> */}

          <Divider className={classes.divider} />

          <Grid container spacing={3}>
            <Grid item xs={3}>
              <Typography variant="body2" component={'span'}
                style={{ color: MulwiColors.subtitleTypography }}>
                {locale2.EMAIL[props.lang]}
              </Typography>
            </Grid>
            <Grid item xs={9}>
              <Typography variant="body2" component={'span'} noWrap>
                {(props.user && props.user.Email) || ""}
              </Typography>
            </Grid>
          </Grid>
          <Divider className={classes.divider} />


          {/* <KnownLangEdit lang={props.lang} user={props.user} main={props.main} /> */}

          <UrlsEdit lang={props.lang} urls={props.Urls} user={props.user} main={props.main} />

          {
            // ABOUT ME
          }
          <Divider className={classes.divider} />
          <AboutMeEdit
            lang={props.lang}
            user={props.user}
            main={props.main} />
          <Divider className={classes.divider} />

          {props.instructor && (<React.Fragment>
            <InvoiceEdit
              lang={props.lang}
              instr={props.instructor}
              main={props.main} />
            <Divider className={classes.divider} />
          </React.Fragment>)}
      </CardWithBg>
  );
}