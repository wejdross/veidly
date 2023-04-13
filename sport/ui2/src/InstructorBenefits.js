import CalendarMonthIcon from "@mui/icons-material/CalendarMonth";
import ChatIcon from "@mui/icons-material/Chat";
import DoneAllIcon from "@mui/icons-material/DoneAll";
import FlagIcon from "@mui/icons-material/Flag";
import GroupsIcon from "@mui/icons-material/Groups";
import HouseboatIcon from "@mui/icons-material/Houseboat";
import PersonIcon from "@mui/icons-material/Person";
import PriceCheckIcon from "@mui/icons-material/PriceCheck";
import SportsMartialArtsIcon from "@mui/icons-material/SportsMartialArts";
import { Card, Divider, Grid, Typography } from "@mui/material";
import Container from "@mui/material/Container";
import { useTheme } from "@mui/material/styles";
import useMediaQuery from '@mui/material/useMediaQuery';
import { makeStyles } from "@mui/styles";
import React, { useEffect, useState } from "react";
import { Parallax } from "react-parallax";
import StickyFooter from "./Footer";
import { JoinUs } from "./JoinUs";
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
            {locale2.YOUR_MANAGEMENT_CENTER[props.lang]}
          </Typography>
          <Grid spacing={3} container>
            <Grid item xs={12} sm={4}>
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
                  {locale2.MISSION[props.lang]}
                </Typography>
                <Typography variant="subtitle1" align={"center"}>
                  {locale2.MISSION_DESC[props.lang]}
                </Typography>
              </Grid>
            </Grid>
            <Grid item xs={12} sm={4}>
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
                  {locale2.PARTNERSHIP[props.lang]}
                </Typography>
                <Typography variant="subtitle1" align={"center"}>
                {locale2.PARTNERSHIP_DESC[props.lang]}
                </Typography>
              </Grid>
            </Grid>
            <Grid item xs={12} sm={4}>
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
                  {locale2.CHOICE[props.lang]}
                </Typography>
                <Typography variant="subtitle1" align={"center"}>
                  {locale2.CHOICE_DESC[props.lang]}
                </Typography>
              </Grid>
            </Grid>
          </Grid>
        </Container>
        <Parallax
          blur={0}
          strength={300}
          bgImage={"static/sunset.webp"}
          bgImageAlt="the cat"
          style={{ height: belowSMSize ? "30vh" : "75vh", marginTop: 80, marginBottom: 50 }}
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
            {locale2.FEATURES[props.lang]}
        </Typography>
        <Grid container spacing={3} style={{ marginTop: 50 }}>
          <SingleFunctionality
            belowSMSize
            header={locale2.PERSONAL_TRAININGS[props.lang]}
            body={locale2.CUSTOMER_IN_CENTER[props.lang]}
          />
          <SingleFunctionality
            belowSMSize
            header={locale2.GROUP_TRAININGS[props.lang]}
            icon={<GroupsIcon fontSize="large" />}
            body={locale2.GROUP_TRAININGS_DESC[props.lang]}
          />
          <SingleFunctionality
            belowSMSize
            header={locale2.CARNETS[props.lang]}
            icon={<SportsMartialArtsIcon fontSize="large" />}
            body={locale2.CARNETS_MARKETING[props.lang]}
          />

          <SingleFunctionality
            belowSMSize
            header={locale2.DCS[props.lang]}
            icon={<PriceCheckIcon fontSize="large" />}
            body={locale2.DCS_MARKETING[props.lang]}
          />
          <SingleFunctionality
            belowSMSize
            header={locale2.ADVANCED_CALENDAR[props.lang]}
            icon={<CalendarMonthIcon fontSize="large" />}
            body={locale2.ADVANCED_CALENDAR_MARKETING[props.lang]}
          />
          <SingleFunctionality
            belowSMSize
            header={"Chat"}
            icon={<ChatIcon fontSize="large" />}
            body={locale2.CHAT_MARKETING[props.lang]}
          />

          <SingleFunctionality
            belowSMSize
            header={locale2.HOLIDAYS[props.lang]}
            icon={<HouseboatIcon fontSize="large" />}
            body={locale2.HOLIDAYS_DESC[props.lang]}
            width={12}
          />
        </Grid>
        <JoinUs lang={props.lang} />
      </Container>
    </>
  );
}

function SingleFunctionality(props) {
  const [hovered, setHovered] = useState(false)
  return (
    <Grid item xs={12} md={props.width | 4}>
      <Card
      elevation={3}
      style={{
        textAlign: 'center',
        borderRadius: 15,
        height: "auto",
        padding: 20,
        color: hovered ? MulwiColors.greenDark : "",
        cursor:  "default",
        userSelect: "none"
      }}
      onMouseEnter={()=> {
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
        <Typography
          variant="h5"
          align={"center"}
          
          >
          {props.header}
        </Typography>
        <Divider style={{width: "85%", marginBottom: 15}}/>
        <Typography variant="p">
          {props.body}
        </Typography>
      </Grid>
          </Card>
    </Grid>
  );
}
