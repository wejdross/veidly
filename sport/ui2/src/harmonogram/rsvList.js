import { Avatar, Card, CardHeader, CardMedia, 
        Grid, Tab, Tabs, Typography } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getUserRsvs } from '../apicalls/instructor.api';
import { TabPanel } from '../tabPanel';
import { prettyPrintDate } from './trainingDetails';
import { locale2 } from '../locale';
import { defaultLogoPath } from '../helpers';

function RsvPresentation(props) {
    return (<div style={{marginLeft: 10, marginRight: 10}}>
        <Grid container spacing={3} direction="row">
            {props.data && props.data.map((di, i) => 
                <React.Fragment key={i}>
                    <Grid item lg={4} 
                        style={{
                            marginBottom: null,
                            maxWidth: 390,
                        }}>
                            <Link to={"/rsv_details?id=" + di.ID} style={{
                                textDecoration:"inherit",
                                color: "inherit"
                            }}>
                            <Card style={{
                                margin: 10
                            }}>
                                    <CardHeader
                                        avatar={
                                            <Avatar aria-label="recipe"
                                            src={di.Instructor.AvatarUrl ||  "static/empty_avatar.png"}>
                                            R
                                            </Avatar>
                                        }
                                        title={di.Training.Title}
                                        subheader={prettyPrintDate(new Date(di.DateStart), props.lang)}
                                    />
                                    <CardMedia style={{
                                        height: 0,
                                        paddingTop: "56.25%"
                                    }} image={di.Training.MainImgUrl || defaultLogoPath}
                                    title="training"/>
                            </Card>
                            </Link>
                    </Grid>
                </React.Fragment>
            )}
        </Grid>
    </div>)
}

export function ListRsv(props) { 

    const [dataIncoming, setDataIncoming] = useState([])
    const [dataDone, setDataDone] = useState([])
    const [dataCancelled, setDataCancelled] = useState([])
    const [value, setValue] = React.useState(0)
    const handleChange = (event, newValue) => {
        setValue(newValue);
    }

    useEffect(() => {
        let x = async () => {
            let d = await getUserRsvs()
            d = JSON.parse(d)
            d = d.Rsv
            let di = []
            let dd = []
            let dc = []
            let now = new Date()
            for(let i = 0; i < d.length; i++) {
                if(d[i].IsActive) {
                    di.push(d[i])
                    continue
                }
                if(d[i].IsConfirmed) {
                    dd.push(d[i])
                    continue
                }
                dc.push(d[i])
            }
            setDataIncoming(di)
            setDataDone(dd)
            setDataCancelled(dc)
            //setData(d)
            //console.log(d)
        }
        x()
    }, [props.refreshToken])

    return (
        <React.Fragment>
            <Grid direction={"column"}
                style={{
                    marginBottom: 30,
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

                        <Tab label={ locale2.INCOMING[props.lang] } style={{
                            fontSize: 12
                        }}/>
                        <Tab label={ locale2.COMPLETED[props.lang] }
                            id="tt2" aria-controls="stt2" style={{
                                fontSize: 12
                            }}/>
                        <Tab label={ locale2.CANCELLED[props.lang] } style={{
                            fontSize: 12
                        }}
                            id="tt3" aria-controls="stt3" />

                    </Tabs>

                </Grid>
            </Grid>
            <TabPanel value={value} index={0}>
                {((!dataIncoming || dataIncoming.length === 0) && (
                        <center>
                            <Typography>
                                { locale2.NO_INCOMING_RSV[props.lang] }
                            </Typography>
                        </center>
                    )) || (
                        <RsvPresentation data={dataIncoming} />
                    )}
            </TabPanel>
            <TabPanel value={value} index={1}>
                <RsvPresentation data={dataDone} />
            </TabPanel>
            <TabPanel value={value} index={2}>
                <RsvPresentation data={dataCancelled} />
            </TabPanel>

        </React.Fragment>
    )
}