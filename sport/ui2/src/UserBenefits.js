import { FitnessCenter, Payment } from "@mui/icons-material";
import CalendarMonthIcon from "@mui/icons-material/CalendarMonth";
import DoneAllIcon from "@mui/icons-material/DoneAll";
import FlagIcon from "@mui/icons-material/Flag";
import PersonIcon from "@mui/icons-material/Person";
import PersonSearchIcon from '@mui/icons-material/PersonSearch';
import QrCode2Icon from '@mui/icons-material/QrCode2';
import QrCodeScannerIcon from '@mui/icons-material/QrCodeScanner';
import { Button, Card, Divider, Grid, Typography } from "@mui/material";
import Container from "@mui/material/Container";
import { useTheme } from "@mui/material/styles";
import useMediaQuery from '@mui/material/useMediaQuery';
import { makeStyles } from "@mui/styles";
import React, { useEffect, useState } from "react";
import { Parallax } from "react-parallax";
import { Link } from "react-router-dom";
import StickyFooter from "./Footer";
import { locale2 } from "./locale";
import { MulwiColors } from "./mulwiColors";

export default function InstructorBenefits(props) {
  useEffect(() => {
    window.scrollTo(0, 0);
  }, []);

  const theme = useTheme();
  const belowSMSize = useMediaQuery(theme.breakpoints.down("sm"));
  return (
    <>
      <Container maxWidth="xl" style={{ marginTop: belowSMSize ? 0 : 80, marginBottom: 100 }}>
        {
          ///////////////////////////////////////////////////////
        }
        <Container maxWidth={"xl"} style={{ marginTop: 80 }}>
          <Typography
            style={{ marginTop: 25, marginBottom: 60, fontWeight: 300 }}
            align='center'
            variant={belowSMSize ? "h4" : "h3"}
            component="div"
            gutterBottom
          >
            {locale2.WHAT_ARE_BENEFITS[props.lang]}
          </Typography>
          <Grid spacing={3} container>
            <Grid item xs={12} md={4}>
              <Grid
                container
                direction="column"
                alignContent="flex-start"
                justifyContent="flex-start"
              >
                <FlagIcon
                  
                  style={{
                    textAlign: "center",
                    width: "100%",
                    fontSize: 40,
                  }}
                />
                <Typography
                  variant="h5"
                  style={{ marginTop: 15, marginBottom: 15 }}
                  
                  align={"center"}
                >
                  {locale2.PROF_PARTNERS[props.lang]}
                </Typography>
                <Typography variant="subtitle1" align={"center"}>
                    {locale2.PROFPARTNERS_SUB_TEXT[props.lang]}
                 </Typography>
              </Grid>
            </Grid>
            <Grid item xs={12} md={4}>
              <Grid
                container
                direction="column"
                alignContent="flex-start"
                justifyContent="flex-start"
              >
                <DoneAllIcon
                  style={{
                    textAlign: "center",
                    width: "100%",
                    fontSize: 40,
                  }}
                />
                <Typography
                  variant="h5"
                  style={{ marginTop: 15, marginBottom: 15 }}
                  
                  align={"center"}
                >
                  {locale2.SIMPLICITY[props.lang]}
                </Typography>
                <Typography variant="subtitle1" align={"center"}>
                  {locale2.SIMPLICITY_SUB_TEXT[props.lang]}
                </Typography>
              </Grid>
            </Grid>
            <Grid item xs={12} md={4}>
              <Grid
                container
                direction="column"
                alignContent="flex-start"
                justifyContent="flex-start"
              >
                <CalendarMonthIcon
                  style={{ textAlign: "center", width: "100%", fontSize: 40 }}
                />
                <Typography
                  variant="h5"
                  style={{ marginTop: 15, marginBottom: 15 }}
                  
                  align={"center"}
                >
                  {locale2.TIME[props.lang]}
                </Typography>
                <Typography variant="subtitle1" align={"center"}>
                  {locale2.TIME_SUB_TEXT[props.lang]}
                </Typography>
              </Grid>
            </Grid>
          </Grid>
        </Container>
        <Parallax
          blur={0}
          strength={300}
          bgImage={"static/boxing3.webp"}
          bgImageAlt="the cat"
          style={{ height: belowSMSize ? "30vh" : "55vh", marginTop: 80, marginBottom: 50 }}
        />
        {
          ////////////////////////////////////////////////////////////////////////////////
        }
        <Typography
          variant="h3"
          align="center"
          style={{ marginTop: 25, marginBottom: 60, fontWeight: 300 }}
          component="div"
          gutterBottom
        >
          {locale2.HOW_TO_START_TRAININGS[props.lang]}
        </Typography>
        <Grid container spacing={3} style={{ marginTop: 50 }}>
          {
           belowSMSize || <Grid item xs={0} md={4}></Grid>
          }
          <SingleFunctionality
            header={"1. " + locale2.FIND_ACTIVITY_FOR_YOU[props.lang]}
            icon={<PersonSearchIcon fontSize="large" />}
            />
          {
           belowSMSize || <Grid item xs={0} md={4}></Grid>
          }
          {
            // separator try
          }
          <VerticalSeparator />
            {
             belowSMSize || <Grid item xs={0} md={4}></Grid>
            }
            <SingleFunctionality

              header={"2. " + locale2.PICK_DATE_AND_PLACE[props.lang]}
              icon={<FitnessCenter fontSize="large" />}
                />
            {
             belowSMSize || <Grid item xs={0} md={4}></Grid>
            }

          <VerticalSeparator />

            {
             belowSMSize || <Grid item xs={0} md={4}></Grid>
            }
            <SingleFunctionality

              header={"3. " + locale2.USER_BENEFITS_CONTACT_WITH_INTSRUCTOR[props.lang]}
              icon={<Payment fontSize="large" />}
                />
            {
             belowSMSize || <Grid item xs={0} md={4}></Grid>
            }

          <VerticalSeparator />

            {
             belowSMSize || <Grid item xs={0} md={4}></Grid>
            }
            <SingleFunctionality

              header={"4. " + locale2.AFTER_PAYMENT_QR[props.lang]}
              icon={<QrCode2Icon fontSize="large" />}
                />
            {
             belowSMSize || <Grid item xs={0} md={4}></Grid>
            }

          <VerticalSeparator />

          {
           belowSMSize || <Grid item xs={0} md={4}></Grid>
          }
          <SingleFunctionality
            header={"5. " + locale2.QR_SHOW_IT_BRO[props.lang]}
            icon={<QrCodeScannerIcon fontSize="large" />}
            />
          {
           belowSMSize || <Grid item xs={0} md={4}></Grid>
          }
        </Grid>
        <RedirectToSearch lang={props.lang} />
      </Container>
    </>
  );
}

function VerticalSeparator(props) {
  return (
    <>
      <Grid item xs={6} style={{
        borderRight: "2px dashed grey",
        height: 40,
        margin: 0,
      }}></Grid>
      <Grid item xs={6}></Grid>
    </>
  )
}

function SingleFunctionality(props) {
  /* 
    this version of Single functionality is slightly different than the one used in InstructorBenefits
    therefore please edit with caution
    #refactor_me 
  */
  const [hovered, setHovered] = useState(false)
  return (
    <Grid item xs={12} md={props.width | 4}>
      <Card
        elevation={3}
        style={{
          textAlign: 'center',
          borderRadius: 15,
          //height: 170,
          padding: 20,
          color: hovered ? MulwiColors.greenDark : "",
          cursor: "default",
          userSelect: "none"
        }}
        onMouseEnter={() => {
          setHovered(true)
        }}
        onMouseLeave={() => {
          setHovered(false)
        }}
      >

        <Grid
          container
          direction="column"
          justifyContent="center"
          alignItems="center"
        >
          {props.icon || <PersonIcon fontSize="large" />}
          <Divider style={{ width: "85%", marginBottom: 15 }} />
          <Typography
            variant="h5"
            align={"center"}
            
          >
            {props.header}
          </Typography>
        </Grid>
      </Card>
    </Grid>
  );
}

function RedirectToSearch(props) {
  const [hovered, setHovered] = useState(false)

  return (
    <center id="#redirect">
      <Link
        to={"/"}
        style={{
          textDecoration: "none",
        }}
        >
        <Button
          onMouseEnter={() => {
            setHovered(true)
          }}
          onMouseLeave={() => {
            setHovered(false)
          }}
          variant="contained"
          style={{
            backgroundColor: hovered ? MulwiColors.blueDark : MulwiColors.greenDark,
            borderRadius: 50,
            width: 270,
            height: 100,
            marginTop: 50,
            marginBottom: 50,
            marginLeft: 0, 
            merginRight: 0,
          }}
          >
          <Typography style={{ color: MulwiColors.whiteSurface, fontWeight: 800, fontSize: "1.2em" }}>
            {locale2.BOOK_TRAINING[props.lang]}
          </Typography>
        </Button>
      </Link>
    </center>
  );
}
