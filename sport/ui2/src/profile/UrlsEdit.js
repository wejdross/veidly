import { Card, Grid, IconButton, TextField, Typography } from '@mui/material'
import { Add, Delete } from '@mui/icons-material'
import React, { useState } from 'react'
import { patchUserData } from '../apicalls/user.api'
import { MulwiColors } from '../mulwiColors'
import ModalEdit, { emptyInfo } from './ModalEdit'
import InputLabel from '@mui/material/InputLabel';
import FormControl from '@mui/material/FormControl';
import Select from '@mui/material/Select';
import MenuItem from '@mui/material/MenuItem';
import FacebookIcon from '@mui/icons-material/Facebook';
import ListItemIcon from '@mui/material/ListItemIcon';
import InstagramIcon from '@mui/icons-material/Instagram';
import TwitterIcon from '@mui/icons-material/Twitter';
import YouTubeIcon from '@mui/icons-material/YouTube';
import HttpIcon from '@mui/icons-material/Http';
import { locale2 } from '../locale'
import AddCircleIcon from '@mui/icons-material/AddCircle';

export default function UrlsEdit(props) {

    const [name, setName] = useState("")
    const [url, setUrl] = useState("")
    const [icon, setIcon] = useState("default")

    const [info, setInfo] = useState(emptyInfo)

    async function save(del) {
        try {
            // deleting bug with null Urls
            if (props.user.Urls === null) {
                props.user.Urls = []
            }
            let u = { ...props.user }
            let urls = [...props.user.Urls]

            if (del && del >= 1) {
                if (!urls) {
                    return
                }
                urls.splice(del - 1, 1)
            } else {
                let _url = url
                if (!url.startsWith("https://")) {
                    _url = "https://" + _url
                }
                let ur = { Name: name, Url: _url, Avatar: icon }
                if (urls) {
                    urls.push(ur)
                } else {
                    urls = [ur]
                }
            }

            u.Urls = [...urls]
            await patchUserData(u)
            props.main.refreshUser()
        } catch (ex) {
            console.log(ex)
            setInfo({
                st: "ex",
                msg: ex
            })
        }
    }

    return (<React.Fragment>
        <ModalEdit hideSaveButton
            lang={props.lang}
            info={info}
            title={locale2.MY_LINKS[props.lang]}
            label={locale2.MY_LINKS[props.lang]}
            labelStyle={{ color: MulwiColors.subtitleTypography }}
            custom
            value={props.user.Urls && props.user.Urls.map((c, i) => (
                <Typography variant="body2">
                    <a key={i} style={{
                        color: MulwiColors.blueDark,
                        textDecoration: "none"
                    }} href={c.Url} target="_blank" rel="noreferrer" >{c.Name}</a>
                </Typography>
            ))} content={(<div>
                    <center>{locale2.ADD_NEW_LINK[props.lang]}</center>
                <Card>
                <Grid container spacing={4} style={{padding: 10}}>
                    <Grid item xs={5}>
                        <FormControl size="small" style={{ width: "100%" }} variant="outlined">
                            <InputLabel id="icon-select"></InputLabel>
                            <Select value={icon} onChange={(e) => { setIcon(e.target.value) }}>
                                <MenuItem value="default">
                                    <ListItemIcon>
                                        {AvatarStringToJSX("default")}
                                    </ListItemIcon>
                                </MenuItem>
                                <MenuItem value="facebook">
                                    <ListItemIcon>
                                    {AvatarStringToJSX("facebook")}
                                    </ListItemIcon>
                                </MenuItem>
                                <MenuItem value="instagram">
                                    <ListItemIcon>
                                    {AvatarStringToJSX("instagram")}
                                    </ListItemIcon>
                                </MenuItem>
                                <MenuItem value="twitter">
                                    <ListItemIcon>
                                    {AvatarStringToJSX("twitter")}
                                    </ListItemIcon>
                                </MenuItem>
                                <MenuItem value="youtube">
                                    <ListItemIcon>
                                    {AvatarStringToJSX("youtube")}
                                    </ListItemIcon>
                                </MenuItem>
                            </Select>
                        </FormControl>
                    </Grid>
                    <Grid item xs={7}>
                        <TextField
                            style={{width: "100%"}}
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            label={locale2.NAME[props.lang]} />
                    </Grid>
                    <Grid item xs={12} >
                        <TextField
                            style={{width: "100%"}}
                            InputProps={{
                                startAdornment:
                                    "https://",
                            }}
                            value={url}
                            onChange={(e) => setUrl(e.target.value)}
                            label="Url" />
                    </Grid>
                    <Grid item xs={12}>
                        <center>
                        <IconButton onClick={save} size={'large'} style={{color: MulwiColors.greenDark}}>
                            <AddCircleIcon />
                        </IconButton>
                        </center>
                    </Grid>
                </Grid>
                </Card>
                <Grid container spacing={2} direction="column">
                    {props.user.Urls && props.user.Urls.map((c, i) => (
                        <Grid item key={i}>
                            <Grid container spacing={2} >
                                <Grid item xs={2}>
                                    <ListItemIcon style={{ lineHeight: 3}}>
                                    {AvatarStringToJSX(c.Avatar ? c.Avatar : "default", true)}
                                    </ListItemIcon>
                                </Grid>
                                <Grid item xs={3}>                                    
                                    <Typography noWrap style={{ lineHeight: 3}}>
                                        {c.Name}
                                    </Typography>
                                </Grid>
                                <Grid item xs={6}>
                                    <Typography noWrap style={{ lineHeight: 3}}>
                                        {c.Url}
                                    </Typography>
                                </Grid>
                                <Grid item xs={1}>
                                    <IconButton onClick={() => {
                                        save(i + 1)
                                    }} style={{
                                        color: MulwiColors.redError,
                                    }}>
                                        <Delete />
                                    </IconButton>
                                </Grid>
                            </Grid>
                        </Grid>
                    ))}
                </Grid>
            </div>)} />
    </React.Fragment>)
}


function AvatarStringToJSX(avatarString, lower) {

    switch (avatarString) {
        case "default":
            return <HttpIcon style={{marginTop: lower ? 12 : 0, color: MulwiColors.greenDark}}/>
        case "facebook":
            return <FacebookIcon style={{marginTop: lower ? 12 : 0, color: "#3b5998"}} />
        case "instagram": 
            return <InstagramIcon style={{marginTop: lower ? 12 : 0, color: "black"}} />
        case "twitter":
            return <TwitterIcon style={{marginTop: lower ? 12 : 0, color: "#1DA1F2"}} />
        case "youtube":
            return <YouTubeIcon style={{marginTop: lower ? 12 : 0, color: "#FF0000"}} />
        default:
            return <HttpIcon style={{marginTop: lower ? 12 : 0, color: MulwiColors.greenDark}} />
    }
}
