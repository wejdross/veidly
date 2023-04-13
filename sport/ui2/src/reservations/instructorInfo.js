import {
    Avatar, Button, DialogActions, DialogContent, DialogTitle, Divider, Grid, TextField, Typography,
} from '@mui/material';
import { ContactMail } from '@mui/icons-material';
import React, { useState } from 'react';
import { useHistory } from 'react-router-dom';
import { MulwiColors } from '../mulwiColors';
import { getErrorDialog } from '../StatusDialog';
import { urlNameToIcon } from './UserInfo';
import { locale2 } from '../locale';
import { getPathFromSearchParams, gettoken } from '../helpers';
import { postUserRoom, storeChatToken } from '../apicalls/chat';
import { Dialog } from '@mui/material';
import { addChatroomToMiniWindow } from '../chat/miniChatWindow';

export function InstructorInfo(props) {

    const history = useHistory()

    const [anonData, setAnonData] = useState({
        DisplayName: "",
        Email: ""
    })

    const [question, setQuestion] = useState("")

    const [anonDialogOpen, setAnonDialogOpen] = useState(false)

    function isAnonDataValid() {
        if (!anonData.DisplayName || !anonData.Email)
            return false
        return true
    }

    async function contactInstructor() {
        try {
            let res;
            let isLoggedIn = gettoken()
            if (isLoggedIn) {
                res = await postUserRoom({
                    PeerUserID: props.instructor.UserID
                })
            } else {
                if (!isAnonDataValid()) {
                    setAnonDialogOpen(true)
                    return
                }
                res = await postUserRoom({
                    PeerUserID: props.instructor.UserID,
                    AnonData: anonData,
                    InitContent: question
                })
            }
            res = JSON.parse(res)

            if (!isLoggedIn) {
                storeChatToken(res.AccessToken)
            }

            if (!res.ChatRoomID)
                throw "invalid ChatRoomID"

            history.push(getPathFromSearchParams(addChatroomToMiniWindow(res.ChatRoomID)))

            setAnonDialogOpen(false)
            
            //history.push("/chat?open=" + res.ChatRoomID)

        } catch (ex) {
            props.setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG, ex))
        }

    }

    if (!props.instructor || !props.instructor.UserInfo) return null

    return (<React.Fragment>

        <Dialog open={anonDialogOpen} onClose={() => setAnonDialogOpen(false)}>
            <DialogTitle>
                {locale2.CONTACT_INSTRUCTOR[props.lang]}
            </DialogTitle>
            <DialogContent>
                <TextField
                    label={locale2.YOUR_EMAIL[props.lang]}
                    value={anonData.Email}
                    onChange={(e) => {
                        if(!anonData.DisplayName || anonData.DisplayName == anonData.Email) {
                            setAnonData(a => ({...a, 
                                DisplayName: e.target.value, Email: e.target.value}))
                        } else {
                            setAnonData(a => ({...a, Email: e.target.value}))
                        }
                    }}
                    margin="dense" fullWidth required
                />
                <TextField
                    label={locale2.YOUR_NAME[props.lang]}
                    value={anonData.DisplayName}
                    onChange={(e) => setAnonData(a => ({...a, DisplayName: e.target.value}))}
                    margin="dense" fullWidth required
                />
                <TextField style={{
                    marginTop: 20
                }} variant='outlined' 
                label={locale2.YOUR_QUESTION[props.lang]}
                        minRows={3} maxRows={4} fullWidth multiline
                        onChange={e => setQuestion(e.target.value)}
                        value={question}
                    />
            </DialogContent>
            <DialogActions>
                <Button variant="contained" style={{
                    color: "white",
                    backgroundColor: MulwiColors.greenDark
                }} onClick={contactInstructor}>
                    {locale2.NEXT[props.lang]}
                </Button>
                <Button onClick={() => setAnonDialogOpen(false)}>
                    {locale2.CLOSE[props.lang]}
                </Button>
            </DialogActions>
        </Dialog>

        <Grid container direction="column" justify="center"
            alignItems="center">
            <Grid item>
                <Avatar src={props.instructor.UserInfo.AvatarUrl || "/placeholder.png"}
                    style={{
                        height: 150,
                        width: 150
                    }} />
            </Grid>
            <Grid item>
                <br />
                <Typography variant="h6">{props.instructor.UserInfo.Name}</Typography>
            </Grid>
            <Grid item>
                <Typography>{locale2.INSTRUCTOR_SINCE[props.lang]}
                    {" " + new Date(props.instructor.CreatedOn).getFullYear()}</Typography>
            </Grid>
            <Grid item>
                <Grid
                    container
                    direction="row"
                    justify="space-between"
                    alignItems="center">
                    {props.instructor.UserInfo.Urls && props.instructor.UserInfo.Urls.map((u, i) => {
                        let nn = urlNameToIcon(u.Name)
                        if (!nn) return null
                        return (<React.Fragment key={i}>
                            <a href={u.Url}>
                                <Avatar style={{
                                    margin: 10,
                                    color: MulwiColors.greenDark,
                                    backgroundColor: MulwiColors.lightGreyAddedByLukasz
                                }}>
                                    {nn}
                                </Avatar>
                            </a>
                        </React.Fragment>)
                    })}
                </Grid>
            </Grid>

            <Grid item>
                <Grid
                    container
                    direction="column"
                    justify="space-between"
                    alignItems="center">
                    {props.instructor.UserInfo.Urls && props.instructor.UserInfo.Urls.map((u, i) => {
                        let nn = urlNameToIcon(u.Name)
                        if (nn) return null
                        return (<React.Fragment key={i}>
                            <a href={u.Url} style={{
                                textDecoration: "none",
                                color: MulwiColors.blueDark
                            }}>
                                <Typography>
                                    {u.Name}
                                </Typography>
                            </a>
                        </React.Fragment>)
                    })}
                </Grid>
            </Grid>

            <Grid item>
                <Button onClick={contactInstructor}><ContactMail style={{
                    color: MulwiColors.blueLight,
                    marginRight: 10
                }} />{locale2.CONTACT_INSTRUCTOR[props.lang]}</Button>
            </Grid>

            <Grid item>
                <Button onClick={() => {
                    let p = "/instr_profile" +
                        "?instructorID=" + props.instructor.id +
                        "&f=1"
                    history.push(p)
                }}>{locale2.INSTRUCTORS_PROFILE[props.lang]}</Button>
            </Grid>

            <Grid item style={{ backgroundColor: "white", padding: 20 }}>
                <Grid
                    container
                    direction="column"
                    justify="center"
                    alignItems="center">
                    <Grid item>
                        <Typography>{locale2.ABOUT_ME[props.lang]}</Typography>
                    </Grid>
                    <Grid item style={{ width: "100%" }}>
                        <Divider />
                    </Grid>
                    <Grid item style={{ marginTop: 10 }}>
                        <Typography style={{
                            whiteSpace: "pre-wrap",
                            overflowWrap: "break-word",
                            maxWidth: 320,
                            textAlign: "center"
                        }}>{props.instructor.UserInfo.AboutMe}</Typography>
                    </Grid>
                </Grid>
            </Grid>
        </Grid>
    </React.Fragment>)
}