import {
    Avatar, Card, CardHeader, CardMedia,
    Grid, Tab, Tabs, Typography,
    useMediaQuery, useTheme
} from '@mui/material';
import React, { useEffect, useState } from 'react';
import DrawerSmall from '../card/DrawerSmall';
import { HarmonogramDay } from '../harmonogram/day';
import { getWkFromMonth, WeekSwitch } from '../harmonogram/harmonogram';
import { HarmonogramMonth } from '../harmonogram/month';
import { MonthLabel } from '../harmonogram/MonthLabel';
import { HarmonogramWeek } from '../harmonogram/weekBigRes';
import { defaultLogoPath, trainingResToDrawerData } from '../helpers';
import { locale2 } from '../locale';
import { TrainingSummary } from '../reservations/trainingSummary';

export function SubSchedule(props) {

    const [value, setValue] = React.useState(0)
    const handleChange = (event, newValue) => {
        setValue(newValue);
    }

    const [data, setData] = useState(null)

    const [drawerOpen, setDrawerOpen] = useState(false)
    const [drawerData, setDrawerData] = useState({})

    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))

    const [week, setWeek] = useState(() => getWkFromMonth(new Date()))

    const [day, setDay] = useState(new Date())
    const [wk, setWk] = useState(getWkFromMonth(new Date()))
    const [monthDate, _setMonthDate] = useState(new Date())

    function setMonthDate(v) {
        _setMonthDate(v)
        setWk(getWkFromMonth(v))
    }

    function setWeekLevelWithDate(start, end) {
        setWk({ e: true, start: start, end: end })
    }

    useEffect(() => {
        setData(props.trainings)
    }, [props.trainings])

    return (<React.Fragment>
        <DrawerSmall padding={7}
                content={<TrainingSummary lang={props.lang}
                    setInfo={props.setInfo}
                    sch={drawerData && drawerData.sch}
                    training={drawerData && drawerData.training}
                    setDrawerOpen={setDrawerOpen} />}

                open={drawerOpen}
                width={isLowRes ? "100vw" : 650}
                onClose={() => setDrawerOpen(false)}
                onOpen={() => setDrawerOpen(true)}>

            <Typography variant="h5" style={{
                marginLeft: isLowRes ? 20 : 0,
                marginRight: isLowRes ? 20 : 0
            }}><center>
                {locale2.CARNET_IS_VALID_FOR[props.lang]}
            </center></Typography>

            <Grid direction={"column"}
                style={{
                    marginBottom: 10,
                    paddingLeft: 10,
                }}
                justify="center"
                alignItems="center"
                spacing={2}
                container >
                <Grid item>
                    <Tabs value={value}
                        variant="fullWidth"
                        onChange={handleChange}
                        indicatorColor="primary"
                        scrollButtons="auto"
                        textColor="primary">

                        <Tab label={isLowRes ? locale2.LIST[props.lang]
                            : locale2.GRID[props.lang]}
                            style={{
                                fontSize: 12
                            }} />
                        <Tab label={locale2.SCHEDULE[props.lang]}
                            id="tt2" aria-controls="stt2" style={{
                                fontSize: 12
                            }} />
                    </Tabs>

                </Grid>
            </Grid>


            {(value === 0) && <React.Fragment>
                <Grid container direction="row"
                        alignItems="center"
                    >
                        {data && data.map((di, i) => (
                            <Card style={{
                                margin: isLowRes ? 0 : 10,
                                padding: isLowRes ? 0 : 10,
                                width: isLowRes ? "100%" : 350,
                                cursor: "pointer"
                            }} onClick={() => {
                                if(drawerOpen) {
                                    if(drawerData && drawerData.training 
                                            && di.Training 
                                            && drawerData.training.ID === di.Training.ID) {
                                        // same training - close trawer
                                        setDrawerOpen(false)
                                    } else {
                                        // if drawer was already open, 
                                        //  but was showing different training, then just update content 
                                        setDrawerData(trainingResToDrawerData(di, null, null))
                                    }
                                } else {
                                    setDrawerOpen(true)
                                    setDrawerData(trainingResToDrawerData(di, null, null))
                                }
                            }}>
                                <CardHeader
                                    avatar={props.instructor && (
                                        <Avatar aria-label="recipe"
                                            src={(props.instructor && props.instructor.UserInfo.AvatarUrl) 
                                                ||  "static/empty_avatar.png"}>
                                            R
                                        </Avatar>)
                                    }
                                    title={di.Training.Title}
                                    subheader={di.Training.LocationText}
                                />
                                <CardMedia style={{
                                    height: 0,
                                    paddingTop: "56.25%"
                                }} image={di.Training.MainImgUrl || defaultLogoPath}
                                    title="training" />
                            </Card>
                        ))}
                    </Grid>
            </React.Fragment>}


            {(value === 1) && <React.Fragment>
                <Grid justify={"space-evenly"}
                    alignItems={"center"}
                    direction="row"
                    container
                    spacing={0}>
                    <Grid item>
                        {isLowRes ? (<React.Fragment>
                            <Grid justify={"space-evenly"}
                                    alignItems={"center"}
                                    direction="row"
                                    container
                                    spacing={0}>
                                    
                                    <Grid item>
                                        <Grid>
                                            <MonthLabel monthDate={monthDate} setMonthDate={setMonthDate} />
                                            <HarmonogramMonth
                                                lang={props.lang}
                                                
                                                user
                                                instructorID={(props.instructor && props.instructor.id) 
                                                    || props.instructorID}

                                                setInfo={props.setInfo}

                                                date={monthDate}
                                                week={wk}
                                                switchToWeek={setWeekLevelWithDate}

                                                day={day}
                                                setDay={setDay}

                                                smID={props.smID}
                                            />
                                        </Grid>
                                    </Grid>

                                    <Grid item style={{
                                        width: "95%"
                                    }}>
                                            <HarmonogramDay
                                                lang={props.lang}
                                                
                                                user

                                                day={day}
                                                instructorID={(props.instructor && props.instructor.id) 
                                                        || props.instructorID}

                                                setDrawerData={setDrawerData}
                                                setDrawerOpen={setDrawerOpen}

                                                setInfo={props.setInfo}

                                                smID={props.smID}

                                            />
                                    </Grid>
                                </Grid>
                        </React.Fragment>) : (<React.Fragment>
                            <WeekSwitch week={week} setWeek={setWeek} />
                            <HarmonogramWeek sm user
                                setInfo={props.setInfo}
                                instructorID={(props.instructor && props.instructor.id) 
                                        || props.instructorID}
                                setDrawerData={setDrawerData}
                                setDrawerOpen={setDrawerOpen}
                                week={week}
                                lang={props.lang}
                                smID={props.smID}/>
                        </React.Fragment>)}
                    </Grid>
                </Grid>
            </React.Fragment>}

        </DrawerSmall>

    </React.Fragment>)
}