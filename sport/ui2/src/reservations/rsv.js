import React, { useEffect, useState } from 'react'
import { 
    getInstructor, getUserSchedule
} from '../apicalls/instructor.api'
import { getErrorDialog, getNullDialog, 
        StatusDialog } from '../StatusDialog'
import {
    Card, CardContent, Container, Grid,
    Typography, useMediaQuery, useTheme
} from '@mui/material'
import { MulwiColors } from '../mulwiColors'
import { TrainingSummary } from './trainingSummary'
import { useHistory } from 'react-router'
import { InstructorInfo } from './instructorInfo'
import { CreateRsv } from './createRsv'
import RsvReviews from './RsvReviews'
import { epochToDate } from '../helpers'
import { locale2 } from '../locale'

export function Rsv(props) {

    const [instructor, setInstructor] = useState(null)
    const [tr, setTraining] = useState(null)
    const [info, setInfo] = useState(getNullDialog())
    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))
    const history = useHistory()
    const [gDateStart, setDateStart] = useState(null)


    async function setInstructorFromApi(iid) {
        try {
            let i = await getInstructor(iid)
            i = JSON.parse(i)
            setInstructor(i)
        } catch (ex) {
            console.log(ex)
            // TODO: log ex
        }
    }

    async function setTrainingFromApi(iid, tid, dateStart) {
        try {
            // let i = await getTrainingByID(tid)
            // setTraining(i[0])

            let sd = await getUserSchedule(dateStart, dateStart, iid, tid)
            if(sd.length != 1 || sd[0].Schedule.length != 1) {
                throw locale2.TRAINING_NOT_FOUND[props.lang]
            }
            //console.log(sd[0])
            setTraining(sd[0])

            //setReady(true)
        } catch (ex) {
            setInfo(getErrorDialog(locale2.PROBLEM_FETCHING_TRAINING[props.lang], ex))
            // TODO: log ex
        }
    }

    function setStateFromQuery(l) {
        let query = new URLSearchParams(l.search)

        let _instructorID = query.get("instructorID")
        if (_instructorID) {
            //setInstructorID(_instructorID)
            setInstructorFromApi(_instructorID)
        }

        let _trainingID = query.get("trainingID")
        let ds = query.get("dateStart")
        let d = epochToDate(ds)
        setDateStart(d)
        if (_trainingID && d) {
            setTrainingFromApi(_instructorID, _trainingID, d)
        }
    }

    useEffect(() => {
        setStateFromQuery(window.location)
    }, [])


    useEffect(() => {
        return history.listen((location) => {
            setStateFromQuery(location)
        })
    }, [history])

    return (<React.Fragment>
        <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />

        <Grid direction={"row"}
            style={{
                marginBottom: isLowRes ? 0 : 10,
                backgroundColor: "white",
                paddingTop: isLowRes ? 60 : 10,
                paddingLeft: isLowRes ? 0 : 10,
            }}
            spacing={2}
            justify="flex-start"
            alignItems="center"
            container>
            <Grid item>
                <Typography style={{
                    paddingLeft: 5,
                    color: MulwiColors.greenDark,
                }} variant="h4">
                    {locale2.SIGN_IN_FOR_THIS_TRAINING[props.lang]}
                </Typography>
            </Grid>
        </Grid>

        {isLowRes ? (
            <Grid container direction="column" spacing={2}>

                {tr && (
                    <Grid item>
                        <Card style={{
                            marginLeft: -5
                        }}>
                            <CardContent>
                                <CreateRsv lang={props.lang}
                                    setInfo={setInfo}
                                    user={props.user}
                                    training={tr.Training} />
                            </CardContent>
                        </Card>
                    </Grid>
                )}

                {tr && (
                    <Grid item>
                        <Card>
                            <CardContent>
                                <TrainingSummary lang={props.lang} 
                                    dateStart={gDateStart} sch={tr.Schedule[0]} 
                                    training={tr.Training} />
                            </CardContent>
                        </Card>
                    </Grid>
                )}

                {instructor && (
                    <Grid item>
                        <Card>
                            <CardContent>
                                <InstructorInfo lang={props.lang}
                                    setInfo={setInfo} 
                                    user={props.user}
                                    instructor={instructor} />
                            </CardContent>
                        </Card>
                    </Grid>
                )}

                {
                    <Grid item>
                        <Card>
                            <CardContent>
                                <RsvReviews lang={props.lang} />
                            </CardContent>
                        </Card>
                    </Grid>
                }
                
            </Grid>
        ) : (
            <Container><Grid container direction="column"
                justify="flex-start"
                spacing={2}
                style={{
                    marginBottom: 30,
                }}
                alignItems="stretch">
                <Grid item>
                    <Grid spacing={2}
                        container
                        direction="row"
                        justify="center"
                        alignItems="stretch">
                        <Grid item md={8}>
                            {tr && (<Card style={{
                                marginBottom: 10
                            }}>
                                <CardContent>
                                    <TrainingSummary lang={props.lang} 
                                        dateStart={gDateStart}
                                        sch={tr.Schedule[0]} training={tr.Training}
                                    />
                                </CardContent>
                            </Card>)}
                            {tr && tr.Training && (<Card>
                                <CardContent>
                                    <RsvReviews lang={props.lang} training={tr.Training} />
                                </CardContent>
                            </Card>)}
                        </Grid>
                        <Grid item md={4} style={{
                            position: "relative",
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
                                    {tr && (<Card>
                                        <CardContent>
                                            <CreateRsv lang={props.lang}
                                                setInfo={setInfo}
                                                user={props.user}
                                                training={tr.Training} />
                                        </CardContent>
                                    </Card>)}
                                </Grid>
                                <Grid item>
                                    {instructor && (<Card>
                                        <CardContent>
                                            <InstructorInfo lang={props.lang}
                                                training={(tr && tr.Training) || null}
                                                setInfo={setInfo} 
                                                user={props.user}
                                                instructor={instructor} />
                                        </CardContent>
                                    </Card>)}
                                </Grid>
                            </Grid>
                        </Grid>
                    </Grid>
                </Grid>
            </Grid></Container>)}
    </React.Fragment>)
}