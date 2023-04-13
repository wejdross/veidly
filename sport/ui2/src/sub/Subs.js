import { Avatar, Card, Grid, CardHeader, 
        Tab, Tabs, Typography, CardContent, 
        useTheme, useMediaQuery } from '@mui/material'
import React, { useEffect, useState } from 'react'
import { getInstrSub, getUserSub } from '../apicalls/sm';
import { getRsvStatus, prettyPrintDay } from '../harmonogram/trainingDetails';
import { MulwiColors } from '../mulwiColors';
import { getErrorDialog, getNullDialog, StatusDialog } from '../StatusDialog';
import { KeyVal, SmContent } from './SubCard';
import { useHistory, useLocation } from 'react-router'
import { locale2 } from '../locale';

export function Subs(props) {
    
    const [sub, setSub] = useState(null)
    const [csub, setCsub] = useState(null)

    const history = useHistory()

    const [value, setValue] = React.useState(0)
    const handleChange = (_, newValue) => {
        setValue(newValue);
    }
    
    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))

    const [info, setInfo] = useState(getNullDialog())

    async function setSubFromApi() {
        try {
            let s = []
            if(props.instr) {
                s = await getInstrSub()
            } else {
                s = await getUserSub()
            }
            s = JSON.parse(s)
            let os = []
            let cs = []
            if(!s) s = []
            for(let i = 0; i < s.length; i++) {
                if(s[i].IsActive) {
                    os.push(s[i])
                }else {
                    cs.push(s[i])
                }
            }
            setSub(os)
            setCsub(cs)
        } catch(ex) {
            setInfo(getErrorDialog(
                locale2.FAILED_TO_FETCH_CARNETS[props.lang], ex))
        } 
    }


    const location = useLocation();
    useEffect(() => {
        setSubFromApi()
    }, [location])

    function displaySubs(_subs, lang) {
        return (<Grid container direction="row">
            {_subs && _subs.map((s, i) => (
                <Grid item xs={12} sm={6} md={4} lg={4}>
                    <Card style={{
                        margin: 10,
                        padding: 10,
                        cursor: "pointer"
                    }} onClick={() => {
                        history.push("/sub_details?id=" + s.ID 
                            + (props.instr  ? "&instr=1" : ""))
                    }}>
                        <CardHeader
                            avatar={
                                <Avatar aria-label="recipe"
                                    src={s.Instructor.UserInfo.AvatarUrl 
                                        || "static/empty_avatar.png"}>
                                    R
                                </Avatar>
                            }
                            title={s.SubModel.Name}
                            subheader={<React.Fragment>
                                <Typography>{s.Instructor.UserInfo.Name}</Typography>
                                <Typography>{getRsvStatus(s)}</Typography>
                            </React.Fragment>}
                        />
                        <CardContent>
                            <SmContent lang={props.lang} noname sm={s.SubModel} />
                            <KeyVal k={locale2.REMAINING_ENTRIES[props.lang]} 
                                v={s.RemainingEntries === -1 
                                    ? locale2.UNLIMITED[props.lang] 
                                    : s.RemainingEntries} />
                            <KeyVal k={locale2.CARNET_VALID_UNTIL[props.lang]} 
                                v={prettyPrintDay(new Date(s.DateEnd))} />
                        </CardContent>
                    </Card>
                </Grid>
            ))}
        </Grid>)
    }

    return (<React.Fragment>
        <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
        <Grid direction={"row"}
            style={{
                marginBottom: 10,
                backgroundColor: props.embedded ? null : "white",
                paddingTop: isLowRes ? 0 : 10,
                paddingLeft: isLowRes ? 0 : 10,
            }}
            spacing={2}
            justify={(isLowRes || props.embedded) ? "center" : "flex-start"}
            alignItems="center"
            container>
            <Grid item>
                <Typography style={{
                    paddingLeft: 5,
                    color: MulwiColors.greenDark,
                }} variant="h4">
                    {props.instr 
                        ? locale2.CARNETS_BOUGHT_FOR_TRAININGS[props.lang] 
                        : locale2.YOUR_CARNETS[props.lang]}
                </Typography>
            </Grid>
        </Grid>
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

                    <Tab label={locale2.ACTIVE[props.lang]} style={{
                        fontSize: 12
                    }} />
                    <Tab label={locale2.CANCELLED[props.lang]}
                        id="tt1" aria-controls="stt1" style={{
                            fontSize: 12
                        }} />
                </Tabs>

            </Grid>
        </Grid>
        <div>
            {(value === 0) && <React.Fragment >
                {displaySubs(sub)}
            </React.Fragment>}
            {(value === 1) && <React.Fragment >
                {displaySubs(csub)}
            </React.Fragment>}
        </div>
    </React.Fragment>)
}