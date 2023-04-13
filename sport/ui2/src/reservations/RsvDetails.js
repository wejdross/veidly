import {
    Button, Card, CardContent,
    Container, Grid, Typography, useMediaQuery,
    useTheme
} from '@mui/material';
import { Alert, AlertTitle } from '@mui/lab';
import React, { useEffect, useState } from 'react';
import { useLocation } from 'react-router';
import { getInstrRsvByID, getUserRsvByAccessToken, getUserRsvByID } from '../apicalls/instructor.api';
import { getUserReview } from '../apicalls/review';
import { RsvInteractMenu } from '../harmonogram/rsvInteract';
import { MulwiColors } from '../mulwiColors';
import { CreateReview } from '../review/createReview';
import { UserReview } from '../review/userReview';
import { getNullDialog, StatusDialog } from '../StatusDialog';
import { InstructorInfo } from './instructorInfo';
import { RsvInfo } from './RsvInfo';
import { TrainingSummary } from './trainingSummary';
import { UserInfo } from './UserInfo';
import { locale2 } from '../locale';
import Donate from '../donations/donate';
import { removeFromQs } from '../helpers';

export function RsvDetails(props) {

    const [rsv, setRsv] = useState(null)
    const [sch, setSch] = useState(null)

    const [isOwner, setIsOwner] = useState(false)

    const [reviewState, setReviewState] = useState(0)
    const [accessToken, setAccessToken] = useState("")
    const [content, setContent] = useState(null)
    const [donateOpen, setDonateOpen] = useState(false)

    const [at, setat] = useState(null)

    async function setReviewInfo(rsvID) {
        try {
            let ur = await getUserReview(rsvID)
            ur = JSON.parse(ur)
            switch (ur.Type) {
                case "token":
                    setReviewState(1)
                    setAccessToken(ur.Token.AccessToken)
                    break
                case "content":
                    setReviewState(2)
                    setContent(ur.Content)
                    break
                default:
                    throw "unkown review state"
            }
        } catch (ex) {
            setReviewState(0)
            // no review available yet
            if (ex != 404) {
                console.log(ex)
            }
        }
    }

    async function refreshFromUrl(l) {
        let query = new URLSearchParams(l.search)
        let id = query.get("id")
        if (!id) return
        setReviewInfo(id)
        // can be either 
        //     - id
        //     - token
        let tp = query.get("type")
        if (!tp)
            tp = "id"
        let _rsv = null
        let instr = query.get("instr")
        if (instr) {
            setIsOwner(true)
        }
        let isnew = query.get("new")
        if(isnew) {
            setDonateOpen(true)
            removeFromQs("new")
        }
        try {
            switch (tp) {
                case "id":
                    if (isOwner || instr) {
                        _rsv = await getInstrRsvByID(id)
                    } else {
                        _rsv = await getUserRsvByID(id)
                    }
                    break
                case "token":
                    _rsv = await getUserRsvByAccessToken(id)
                    setat(id)
                    break
                default:
                    return
            }
            _rsv = JSON.parse(_rsv)
        } catch (ex) {
            console.log(ex)
        }
        if (_rsv && _rsv.Rsv && _rsv.Rsv[0]) {
            setRsv(_rsv.Rsv[0])
            let s = {
                Start: new Date(_rsv.Rsv[0].DateStart),
                End: new Date(_rsv.Rsv[0].DateEnd),
                Occ: _rsv.Rsv[0].Occ
            }
            setSch(s)
            //if(_rsv.Rsv[0].InstructorID)
        }
    }

    useEffect(() => {
        if (!props.instructor || !rsv) return
        if (props.instructor.id === rsv.InstructorID) {
            setIsOwner(true)
        }
    }, [props.instructor, rsv])

    function refresh() {
        refreshFromUrl(window.location)
    }

    const location = useLocation();
    useEffect(() => {
        refresh()
    }, [location])

    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))

    const [info, setInfo] = useState(getNullDialog())

    const [ropen, setRopen] = useState(false)

    if (!rsv) return null

    return (
        <React.Fragment>
            <Donate open={donateOpen} lang={props.lang} />
            {rsv && reviewState === 1 && (
                <React.Fragment>
                    <CreateReview lang={props.lang}
                        accessToken={accessToken}
                        onChange={() => {
                            setRopen(false)
                            setReviewInfo(rsv.ID)
                        }}
                        open={ropen}
                        setOpen={setRopen}
                        training={rsv.Training.Title}
                    />
                    <Alert severity="info" action={<Button onClick={() => setRopen(true)} style={{
                        marginLeft: 10
                    }}>{locale2.ADD_REVIEW[props.lang]}</Button>}>
                        <AlertTitle>
                        {locale2.HOW_DID_YOU_LIKE_TRAINING[props.lang]}
                        </AlertTitle>
                    </Alert>
                </React.Fragment>
            )}
            <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
            <Grid direction={"row"}
                style={{
                    marginBottom: isLowRes ? 0 : 10,
                    backgroundColor: "white",
                    paddingTop: isLowRes ? 70 : 10,
                    paddingLeft: isLowRes ? 0 : 10,
                    paddingRight: isLowRes ? 10 : 0,
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
                        {locale2.RESERVATION[props.lang]}
                    </Typography>
                </Grid>
                {rsv && (
                    <RsvInteractMenu
                        lang={props.lang}
                        grid={!isLowRes}
                        rsv={rsv}
                        noDetails
                        at={at}
                        onChange={refresh}
                        setInfo={setInfo}
                        instructor={isOwner} />
                )}
            </Grid>
            {isLowRes ? (
                <Grid spacing={2}
                    container
                    direction="row"
                    justify="center"
                    alignItems="stretch">

                    <Grid item>
                        <Card>
                            <CardContent>
                                <RsvInfo lang={props.lang} rsv={rsv} />
                            </CardContent>
                        </Card>
                    </Grid>

                    <Grid item>
                        <Card>
                            <CardContent>
                                <InstructorInfo lang={props.lang}
                                    user={props.user}
                                    setInfo={setInfo} 
                                    instructor={rsv.Instructor} />
                            </CardContent>
                        </Card>
                    </Grid>

                    <Grid item>
                        <Card>
                            <CardContent>
                                <TrainingSummary lang={props.lang}
                                    sch={sch}
                                    training={rsv.Training}
                                />
                            </CardContent>
                        </Card>
                        {reviewState === 2 && (<Card>
                            <CardContent>
                                <UserReview lang={props.lang} content={content} />
                            </CardContent>
                        </Card>)}
                    </Grid>
                </Grid>
            ) : (
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
                                <Grid item lg={8} >
                                    <Card style={{
                                        marginBottom: 10
                                    }}>
                                        <CardContent>
                                            <TrainingSummary lang={props.lang}
                                                sch={sch}
                                                training={rsv.Training}
                                            />
                                        </CardContent>
                                    </Card>
                                    {reviewState === 2 && (<Card>
                                        <CardContent>
                                            <UserReview lang={props.lang} content={content} />
                                        </CardContent>
                                    </Card>)}
                                </Grid>
                                <Grid item lg={4} style={{
                                    position: "relative"
                                }}>
                                    <Grid
                                        style={{
                                            position: "sticky",
                                            top: 100
                                        }}
                                        container
                                        direction="column">
                                        <Grid item style={{
                                            marginBottom: 10
                                        }}>
                                            <Card>
                                                <CardContent>
                                                    <RsvInfo lang={props.lang} rsv={rsv} />
                                                </CardContent>
                                            </Card>
                                        </Grid>
                                        <Grid item>
                                            <Card>
                                                <CardContent>
                                                    {isOwner ? (
                                                        <UserInfo lang={props.lang}
                                                            user={props.user}
                                                            setInfo={setInfo} 
                                                            userInfo={rsv.UserInfo}
                                                            contactData={rsv.UserContactData} />
                                                    ) : (
                                                        <InstructorInfo lang={props.lang}
                                                            user={props.user}
                                                            setInfo={setInfo} 
                                                            instructor={rsv.Instructor} />
                                                    )}
                                                </CardContent>
                                            </Card>
                                        </Grid>
                                    </Grid>
                                </Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </Container>
            )}
        </React.Fragment>
    )
}
