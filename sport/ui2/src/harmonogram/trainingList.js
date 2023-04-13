import {
    Avatar, Card, Grid, CardHeader,
    CardMedia, Tab, Tabs, Typography,
    CardContent
} from '@mui/material'
import React, { useEffect, useState } from 'react'
import { Link } from 'react-router-dom';
import { getInstructorRsvs, getTrainings } from '../apicalls/instructor.api'
import { getRsvStatus, prettyPrintDate } from './trainingDetails';
import { getUserInfo } from '../apicalls/user.api';
import { defaultLogoPath, trainingResToDrawerData } from '../helpers';
import { locale2 } from '../locale';

export function ListTrainings(props) {

    const [data, setData] = useState([])
    const [rsvSortedByDate, setrsvSortedByDate] = useState([])
    const [endedReservations, setEndedReservations] = useState([])
    const [value, setValue] = React.useState(0)
    const [userAvatar, setUserAvatar] = useState("")
    const handleChange = (event, newValue) => {
        setValue(newValue);
    }
    useEffect(() => {
        let x = async () => {
            let d = await getTrainings(null)
            setData(d)
            let d2 = await getInstructorRsvs()
            let reservations = JSON.parse(d2).Rsv

            let user = await getUserInfo()
            setUserAvatar(JSON.parse(user).AvatarUrl)

            let di = []
            let dd = []
            reservations.map((val, key) => {
                val.Training.DateStart = new Date(val.Training.DateStart)
                val.Training.DateEnd = new Date(val.Training.DateEnd)

                if (!val.IsActive) {
                    dd.push(val)
                } else {
                    di.push(val)
                }
                return null
            })
            setEndedReservations(dd)

            setrsvSortedByDate(di.sort((d1, d2) => {
                let x = new Date(d1.DateStart)
                let y = new Date(d2.DateStart)
                return x - y
            }))
        }
        x()
    }, [props.refreshToken])

    return (
        <React.Fragment>
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

                        <Tab label={locale2.TRAININGS[props.lang]} style={{
                            fontSize: 12
                        }} />
                        <Tab label={locale2.RESERVATIONS[props.lang]}
                            id="tt2" aria-controls="stt2" style={{
                                fontSize: 12
                            }} />
                        <Tab label={locale2.FINISHED[props.lang]} style={{
                            fontSize: 12
                        }}
                            id="tt3" aria-controls="stt3" />

                    </Tabs>

                </Grid>
            </Grid>

            {(value === 0) && <React.Fragment >
                <Grid container direction="row">
                    {data && data.map((di, i) => (
                        <Card style={{
                            margin: 10,
                            padding: 10,
                            width: 350,
                            cursor: "pointer"
                        }} onClick={() => {
                            props.setDrawerOpen(!props.drawerOpen)
                            props.setDrawerData(trainingResToDrawerData(di, null, null))
                        }}>
                            <CardHeader
                                avatar={
                                    <Avatar aria-label="recipe"
                                        src={userAvatar || "/static/empty_avatar.png"}>
                                        R
                                    </Avatar>
                                }
                                title={di.Training.Title}
                                subheader={di.Training.LocationText}
                            />
                            <CardMedia
                                component="img"
                                style={{
                                    maxHeight: 280
                                }}
                                alt="img"
                                image={di.Training.MainImgUrl || "/" + defaultLogoPath}
                                title="training" />
                        </Card>
                    ))}
                </Grid>
            </React.Fragment>
            }

            {(value === 1) && <React.Fragment >
                <Grid container direction="row">
                    {rsvSortedByDate && rsvSortedByDate.map((di, i) => (
                        <Link to={"/rsv_details?id=" + di.ID + "&instr=1"} style={{
                            textDecoration: "none",
                        }}
                        >
                            <Card style={{
                                margin: 10,
                                padding: 10,
                                width: 350,
                                cursor: "pointer"
                            }}
                                key={i}>
                                <CardHeader
                                    avatar={
                                        <Avatar aria-label="recipe"
                                            src={di.UserInfo.AvatarUrl || "/static/empty_avatar.png"}>
                                            R
                                        </Avatar>
                                    }
                                    title={di.Training.Title}
                                    subheader={di.Training.LocationText}
                                />
                                <CardContent>
                                    <EvenlySpacedCardInfo
                                        question={locale2.WITH_WHOM[props.lang]} answer={di.UserInfo.Name} />
                                    <EvenlySpacedCardInfo
                                        question={locale2.WHEN[props.lang]} answer={prettyPrintDate(new Date(di.DateStart), props.lang)} />
                                    <EvenlySpacedCardInfo
                                        question={locale2.STATUS[props.lang]} answer={getRsvStatus(di)} />
                                </CardContent>
                                <CardMedia
                                    component="img"
                                    style={{
                                        maxHeight: 280
                                    }}
                                    alt="img"
                                    image={di.Training.MainImgUrl || "/" + defaultLogoPath}
                                    title="training" />
                            </Card>
                        </Link>
                    ))}
                </Grid>
            </React.Fragment>
            }

            {(value === 2) && <React.Fragment >
                <Grid container direction="row">
                    {endedReservations && Array.from(endedReservations).map((di, i) => (
                        <Link to={"/rsv_details?id=" + di.ID + "&instr=1"} style={{ textDecoration: "none" }}>
                            <Card style={{
                                margin: 10,
                                padding: 10,
                                width: 350,
                                cursor: "pointer"
                            }}
                            >
                                <CardHeader
                                    avatar={
                                        <Avatar aria-label="recipe"
                                            src={di.UserInfo.AvatarUrl || "/static/empty_avatar.png"}>
                                            R
                                        </Avatar>
                                    }
                                    title={di.Training.Title}
                                    subheader={di.Training.LocationText}
                                />
                                <CardContent>
                                    <EvenlySpacedCardInfo
                                        question={locale2.WITH_WHOM[props.lang]} answer={di.UserInfo.Name} />
                                    <EvenlySpacedCardInfo
                                        question={locale2.WHEN[props.lang]} answer={prettyPrintDate(new Date(di.DateStart), props.lang)} />
                                    <EvenlySpacedCardInfo
                                        question={locale2.STATUS[props.lang]} answer={getRsvStatus(di)} />
                                </CardContent>
                                <CardMedia
                                    component="img"
                                    style={{
                                        maxHeight: 280
                                    }}
                                    alt="img"
                                    image={di.Training.MainImgUrl || "/" + defaultLogoPath}
                                    title="training" />
                            </Card>
                        </Link>
                    ))}
                </Grid>
            </React.Fragment>
            }
        </React.Fragment>
    )
}

function EvenlySpacedCardInfo(props) {
    if (props.map) {
        if (!props.LocationLat) {
            throw new Error("lng is broken")
        }
        if (!props.LocationLng) {
            throw new Error("lng is broken")
        }
        if (!props.LocationText) {
            throw new Error("lng is broken")
        }
    }

    return (
        <React.Fragment>
            <Grid container direction="row"
                justify="space-between"
                alignItems="stretch">
                {
                    props.map ?
                        <Link href={`https://www.google.com/maps/search/?api=1&query=${props.LocationLat}%2C${props.LocationLng}`} target={"_blank"} rel={"noreferrer"}>
                            <Typography>
                                {props.LocationText}
                            </Typography>
                        </Link>
                        :
                        <Typography>
                            {props.answer}
                        </Typography>
                }
            </Grid>
        </React.Fragment>
    )
}
