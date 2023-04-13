import {
    Avatar, Divider, Grid,
    Typography,
    useMediaQuery,
    useTheme
} from '@mui/material';
import { Twitter, YouTube } from '@mui/icons-material';
import FacebookIcon from "@mui/icons-material/Facebook";
import InstagramIcon from "@mui/icons-material/Instagram";
import React from 'react';
import { MulwiColors } from '../mulwiColors';
import { locale2 } from '../locale';

export function urlNameToIcon(n) {
    switch (n) {
        case "facebook":
            return <FacebookIcon />
        case "youtube":
            return <YouTube />
        case "instagram":
            return <InstagramIcon />
        case "twitter":
            return <Twitter />
        default:
            return null
    }
}

export function UserInfo(props) {

    const t = useTheme()
    const isLowRes = useMediaQuery(t.breakpoints.down('sm'))

    if (!props.userInfo || !props.contactData) 
        return null

    return (<React.Fragment>
        <Grid container direction="column" justify="center"
            alignItems="center">
            <Grid item>
                <Avatar alt="asd" src={props.userInfo.AvatarUrl 
                || "placeholder.png"}
                    style={{
                        height: 150,
                        width: 150
                    }} />
            </Grid>
            <Grid item>
                <br />
                <Typography variant="h6">{props.userInfo.Name}</Typography>
            </Grid>
            <Grid item>
                <Grid
                    container
                    direction="row"
                    justify="space-between"
                    alignItems="center">
                    {props.userInfo.Urls && props.userInfo.Urls.map((u, i) => {
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
                    {props.userInfo.Urls && props.userInfo.Urls.map((u, i) => {
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

            {props.contactData.Email && (
                <Grid item>
                    {props.contactData.Email}
                </Grid>
            )}

            {props.contactData.Phone && (
                <Grid item>
                    {props.contactData.Phone}
                </Grid>
            )}


            {props.userInfo.AboutMe && (
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
                            overflowWrap:"break-word",
                            maxWidth: 320,
                            textAlign: "center"
                        }}>{props.instructor && props.instructor.UserInfo 
                            && props.instructor.UserInfo.AboutMe}</Typography>
                    </Grid>
                    </Grid>
                </Grid>
            )}

        </Grid>
    </React.Fragment>)
}