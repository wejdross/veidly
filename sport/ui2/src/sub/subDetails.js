import { Button, Card, CardContent, Container, Grid, 
        Typography, useMediaQuery, useTheme } from '@mui/material'
import React, { useEffect, useState } from 'react'
import { useLocation } from 'react-router'
import { getInstrSub, getTrainingsForSm, getUserSub } from '../apicalls/sm'
import { SubInteractMenu } from '../harmonogram/subInteract'
import { daySuffix as dl, getRsvStatus, 
        prettyPrintDay } from '../harmonogram/trainingDetails'
import { MulwiColors } from '../mulwiColors'
import { InstructorInfo } from '../reservations/instructorInfo'
import { UserInfo } from '../reservations/UserInfo'
import { getNullDialog, StatusDialog } from '../StatusDialog'
import { SubSchedule } from './SubSchedule'
import { locale2 } from '../locale';

function smTerm(sub) {
    let now = new Date()
    now.setDate(now.getDate() + sub.SubModel.Duration)
    return now
}

function kv(k, v) {
    return (
        <Grid item>
            <Grid container direction="row" spacing={2} alignItems="center" justify="space-between">
                <Grid item>
                    <Typography>{k}</Typography>
                </Grid>
                <Grid item>
                    <strong>
                        {v}
                    </strong>
                </Grid>
            </Grid>
        </Grid>)
}

export function SubInfo(props) {

    let sub = props.sub

    function dd(sub) {
        let sd = smTerm(sub)
        let now = new Date()
        if(sd < now) {
            return <Typography style={{
                color: MulwiColors.redError
            }}>
                {prettyPrintDay(sd)}
            </Typography>
        }
        return <Typography style={{
            color: MulwiColors.greenDark
        }}>
            {prettyPrintDay(sd)}
        </Typography>
    }

    return (<React.Fragment>
        <Grid item>
                <Typography variant="h6">{locale2.CARNET[props.lang]}</Typography>
            </Grid>
            {kv(locale2.CARNET_PAID[props.lang], 
                    sub.IsConfirmed ? <Typography style={{
                color: MulwiColors.greenDark
            }}>
                {locale2.YES[props.lang]}
            </Typography> : <Typography style={{
                color: MulwiColors.redError
            }}>
                {locale2.NO[props.lang]}
            </Typography>)}
            {!props.instr && !sub.IsConfirmed && <Typography variant="body2" style={{
                color: MulwiColors.subtitleTypography
            }}>
                {locale2.CARNET_MUST_BE_PAID_TO_BE_VALIDATED[props.lang]}
            </Typography>}
            {!props.instr && (<center>
                {sub.State === "link_express" && <Button onClick={() => {
                    window.open(sub.LinkUrl)
                }}>
                    {locale2.LINK_TO_PAYMENT[props.lang]}
                </Button>}
            </center>)}
            {kv(locale2.PAYMENT_STATUS[props.lang], getRsvStatus(sub))}
            {!props.instr && (
                kv(locale2.PERIOD_OF_VALIDITY[props.lang], <React.Fragment>
                {sub.SubModel.Duration} {dl(sub.SubModel.Duration, props.lang)}
            </React.Fragment>)
            )}
            {kv(locale2.CARNET_VALID_UNTIL[props.lang], dd(sub))}
            {kv(locale2.MAX_NUMBER_OF_ENTRIES[props.lang], 
                    sub.SubModel.MaxEntrances === -1 
                        ? locale2.UNLIMITED[props.lang] 
                        : sub.SubModel.MaxEntrances)}
            {kv(locale2.REMAINING_ENTRIES[props.lang], sub.RemainingEntries === -1 
                ? locale2.UNLIMITED[props.lang] 
                : sub.RemainingEntries)}
        </React.Fragment>)
}

export function SubDetails(props) {

    const [sub, setSub] = useState(null)
    const [isOwner, setIsOwner] = useState(false)
    const [trainings, setTrainings] = useState(null)
    //const [instructor, setInstructor] = useState(null)


    async function setSubFromQuery() {
        let query = new URLSearchParams(window.location.search)
        let id = query.get("id")
        let _isowner = Boolean(query.get("instr"))
        try {
            let s = null
            if (_isowner) {
                s = await getInstrSub(id)
            } else {
                s = await getUserSub(id)
            }
            if (!s) return
            s = JSON.parse(s)
            if (s.length !== 1) {
                return
            }
            setSub(s[0])
            let iid = s[0].SubModel.InstructorID
            setSmAndTrain(s[0].SubModel.ID, iid)
            setIsOwner(_isowner)
        } catch (ex) {
            console.log(ex)
        }
    }

    async function setSmAndTrain(smID, instrID) {
        try {
            let t = await getTrainingsForSm(smID, instrID)
            if(t.length < 1) {
                return
            }
            setTrainings(t)
        } catch (ex) {
            console.log(ex)
            // TODO: log ex
        }
    }

    const [info, setInfo] = useState(getNullDialog())

    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))

    const location = useLocation();
    useEffect(() => {
        setSubFromQuery()
    }, [location])



    if (!sub) return null

    return (<React.Fragment>
        <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
        <Grid direction={"row"}
            style={{
                marginBottom: isLowRes ? 0 : 10,
                backgroundColor: "white",
                paddingTop: isLowRes ? 60 : 10,
                paddingLeft: isLowRes ? 0 : 10,
                marginTop: isLowRes ? 5 : 0,
                paddingRight: isLowRes ? 10 : 0
            }}
            spacing={2}
            justify={isLowRes ? "space-between" : "flex-start"}
            alignItems="center"
            container>
            <Grid item>
                <Typography style={{
                    paddingLeft: 5,
                    color: MulwiColors.greenDark,
                }} variant="h4">
                    Karnet
                </Typography>
            </Grid>
            {sub && (
                <SubInteractMenu
                    grid={!isLowRes}
                    sub={sub}
                    noDetails
                    onChange={setSubFromQuery}
                    setInfo={setInfo}
                    instructor={isOwner} />
            )}
        </Grid>
        {isLowRes ? (<React.Fragment>
            <Grid spacing={2}
                    container
                    direction="row"
                    justify="center"
                    alignItems="stretch">
                <Grid item>
                    <Card>
                        <CardContent>
                            <SubInfo lang={props.lang} sub={sub} />
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item  style={{
                        width: "100%"
                    }}>
                    <Card>
                        <CardContent>
                            {isOwner ? (
                                <UserInfo lang={props.lang}
                                    userInfo={sub && sub.UserInfo}
                                    setInfo={setInfo}
                                    contactData={sub && sub.ContactData}
                                />
                            ) : (
                                <InstructorInfo lang={props.lang}
                                    user={props.user}
                                    setInfo={setInfo}
                                    instructor={sub && sub.Instructor}
                                />
                            )}
                        </CardContent>
                    </Card>
                </Grid>
                <Grid item>
                    <Card style={{
                        marginBottom: 10,
                    }}>
                        <CardContent>
                            <SubSchedule  lang={props.lang}
                                setInfo={setInfo} 
                                smID={sub.SubModel && sub.SubModel.ID}
                                trainings={trainings} 
                                instructor={sub && sub.Instructor} 
                                />
                        </CardContent>
                    </Card>
                </Grid>
            </Grid>
        </React.Fragment>) : (
            <Container>
                <Grid container direction="column"
                    justify="flex-start"
                    spacing={2}
                    style={{
                        marginBottom: 30
                    }}
                    alignItems="stretch">
                    <Grid item>
                        <Grid spacing={2}
                            container
                            direction="row"
                            justify="center"
                            alignItems="stretch">
                                
                            <Grid item lg={8}>
                                <Card style={{
                                    marginBottom: 10,
                                    height: "100%"
                                }}>
                                    <CardContent>
                                        <SubSchedule lang={props.lang}
                                            setInfo={setInfo} 
                                            smID={sub.SubModel && sub.SubModel.ID}
                                            //instructorID={sub.SubModel.InstructorID}
                                            trainings={trainings} 
                                            instructor={sub && sub.Instructor} 
                                            />
                                    </CardContent>
                                </Card>
                            </Grid>

                            <Grid item lg={4}>
                                <Grid
                                    container
                                    direction="column">
                                    <Grid item style={{
                                        marginBottom: 10,
                                    }}>
                                        <Card>
                                            <CardContent>
                                                <SubInfo lang={props.lang} sub={sub} />
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                    <Grid item >
                                        <Card>
                                            <CardContent>
                                                {isOwner ? (
                                                    <UserInfo lang={props.lang}
                                                        userInfo={sub && sub.UserInfo}
                                                        setInfo={setInfo}
                                                        contactData={sub && sub.ContactData}
                                                    />
                                                ) : (
                                                    <InstructorInfo lang={props.lang}
                                                        user={props.user}
                                                        setInfo={setInfo}
                                                        instructor={sub && sub.Instructor}
                                                    />
                                                )}
                                            </CardContent>
                                        </Card>
                                    </Grid>
                                </Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </Container>)}
    </React.Fragment>)
}