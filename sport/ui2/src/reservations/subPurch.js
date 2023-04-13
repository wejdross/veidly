import { Card, CardContent, Container, Grid, 
        Typography, useMediaQuery, useTheme } from '@mui/material'
import React, { useEffect, useState } from 'react'
import { useHistory } from 'react-router'
import { getInstructor } from '../apicalls/instructor.api'
import { getTrainingsForSm } from '../apicalls/sm'
import { MulwiColors } from '../mulwiColors'
import { getNullDialog, StatusDialog } from '../StatusDialog'
import { SubSchedule } from '../sub/SubSchedule'
import { CreateSub } from './createSub'
import { InstructorInfo } from './instructorInfo'
import { locale2 } from '../locale'

export function SubPurch(props) {

    const [sm, setsm] = useState(null)
    const [instructor, setInstructor] = useState(null)
    const [trainings, setTrainings] = useState(null)

    const history = useHistory()

    async function setSmAndTrain(smID, instrID) {
        try {
            let t = await getTrainingsForSm(smID, instrID)
            if(t.length < 1) {
                return
            } 
            setsm(t[0].Sms[0])
            setTrainings(t)
        } catch (ex) {
            console.log(ex)
            // TODO: log ex
        }
    }

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

    function setStateFromQuery(l) {
        let query = new URLSearchParams(l.search)

        let _instructorID = query.get("instructorID")
        if (_instructorID) {
            //setInstructorID(_instructorID)
            setInstructorFromApi(_instructorID)
        }

        let _trainingID = query.get("trainingID")
        let smID = query.get("smID")
        if (smID) {
            setSmAndTrain(smID, _instructorID, _trainingID)
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

    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))

    const [info, setInfo] = useState(getNullDialog())

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
                    {locale2.BUY_CARNET[props.lang]}
                </Typography>
            </Grid>
        </Grid>

        {isLowRes ? (<React.Fragment>

        </React.Fragment>) : (<Container>
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
                            {/* column 1 */}
                            <Card style={{
                                height: "100%",
                            }}><CardContent>
                                <SubSchedule lang={props.lang}
                                    setInfo={setInfo} 
                                    smID={sm && sm.ID}
                                    trainings={trainings} 
                                    instructor={instructor} />
                            </CardContent></Card>
                        </Grid>
                        <Grid item lg={4} style={{
                            position: "relative",
                        }}>
                            {/* column 2 */}
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
                                            <CreateSub lang={props.lang}
                                                sm={sm}
                                                setInfo={setInfo}
                                                user={props.user}/>
                                        </CardContent>
                                    </Card>
                                </Grid>
                                <Grid item>
                                    {instructor && (<Card>
                                        <CardContent>
                                            <InstructorInfo lang={props.lang}
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
            </Grid>
        </Container>)}

    </React.Fragment>)
}