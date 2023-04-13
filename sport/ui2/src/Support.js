import React from 'react';
import { makeStyles } from '@mui/styles';
import { Grid, Typography } from '@mui/material';
import CardWithBg from './card/cardWithBg';
import { MulwiColors } from './mulwiColors';
import MailOutlineIcon from '@mui/icons-material/MailOutline';
import { locale2 } from './locale';
import { LinkWithTypo } from './Commons';
import { supportEmail } from './StatusDialog';
import Link from '@mui/material/Link';
import './veidly-styles.css'
import StickyFooter from './Footer';
const useStyles = makeStyles((theme) => ({

  TypographyGeneral: {
    // fontSize: "1.8em",
    [theme.breakpoints.down('md')]: {
      fontSize: "1.0em",
    },
  },

}));


export default function SupportPage(props) {
  const classes = useStyles();

  const linkStyle = {
    marginLeft: 20,
    marginTop: 10,
    marginBottom: 10,
    // fontSize: "1.5em",
  }

  return (
    <>
      <CardWithBg header={locale2.SUPPORT[props.lang]} subheader={locale2.SUPPORT_DESC[props.lang]}>

        <Typography variant='body1' className={'typographyBlackTextRegularsize'}>
          {locale2.BENEFITS_INFO[props.lang]}
        </Typography>

        <LinkWithTypo style={linkStyle}
          text={locale2.INSTR_BENEFITS[props.lang]} to="/benefits_instructor" />
        <br />
        <LinkWithTypo style={linkStyle}
          text={locale2.CUSTOMER_BENEFITS[props.lang]} to="/benefits_user" />

        <Typography variant='body1' className={'typographyBlackTextRegularsize'}>
          {locale2.INTERESTED_IN_LEGAL[props.lang]}
        </Typography>

        <Link target={'_blank'} rel="noopener noreferrer" href={"polityka-prywatnosci.pdf"} style={{
          textDecoration: "none",
          color: MulwiColors.blueDark,
          margin: 5
        }}>
          <Typography >
            <strong>{locale2.PRIVACY_POLICY[props.lang]}</strong>
          </Typography>
        </Link>
        <Link target={'_blank'} rel="noopener noreferrer" href={"trener-regulamin.pdf"} style={{
          textDecoration: "none",
          color: MulwiColors.blueDark,
          margin: 5
        }}>
          <Typography >
            <strong>Regulamin trener</strong>
          </Typography>
        </Link>
        <Link target={'_blank'} rel="noopener noreferrer" href={"konsument-regulamin.pdf"} style={{
          textDecoration: "none",
          color: MulwiColors.blueDark,
          margin: 5
        }}>
          <Typography >
            <strong>Regulamin użytkownik</strong>
          </Typography>
        </Link>

        <Typography variant='body1' className={'typographyBlackTextRegularsize'}>
          {locale2.CONTACT_SUPPORT[props.lang]}
        </Typography>

        <a href={"https://m.me/veidly"} target={'_blank'} rel="noopener noreferrer" style={{
          lineHeight: 1,
          textDecoration: "none",
          color: MulwiColors.greenDark
        }}>
          <Grid container spacing={3} style={{
            marginTop: 20
          }}>
            <Grid item xs={3}>

              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="currentColor" d="M12 2C6.36 2 2 6.13 2 11.7c0 2.91 1.19 5.44 3.14 7.17c.16.13.26.35.27.57l.05 1.78c.04.57.61.94 1.13.71l1.98-.87c.17-.06.36-.09.53-.06c.9.27 1.9.4 2.9.4c5.64 0 10-4.13 10-9.7C22 6.13 17.64 2 12 2m6 7.46l-2.93 4.67c-.47.73-1.47.92-2.17.37l-2.34-1.73a.6.6 0 0 0-.72 0l-3.16 2.4c-.42.33-.97-.17-.68-.63l2.93-4.67c.47-.73 1.47-.92 2.17-.4l2.34 1.76a.6.6 0 0 0 .72 0l3.16-2.4c.42-.33.97.17.68.63Z" /></svg>

            </Grid>
            <Grid item xs={9}>

              {"Messenger"}
            </Grid>
          </Grid>
        </a>
        <a href={"mailto: " + supportEmail} style={{
          lineHeight: 1,
          textDecoration: "none",
          color: MulwiColors.greenDark
        }}>
          <Grid container spacing={3} style={{
            marginTop: 20
          }}>
            <Grid item xs={3}>

              <MailOutlineIcon />
            </Grid>
            <Grid item xs={9}>

              {supportEmail}
            </Grid>
          </Grid>
        </a>
      </CardWithBg>
    </>
  );
}

// https://m.me/veidly