import {
    BeachAccess, CardMembership, Chat, Edit, GroupTwoTone,
    Person, PersonPin, Receipt, Settings, SportsBaseball, SportsMmaOutlined
} from '@mui/icons-material';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import CreditCardIcon from '@mui/icons-material/CreditCard';
import EventIcon from '@mui/icons-material/Event';
import HealthAndSafetyIcon from '@mui/icons-material/HealthAndSafety';
import HelpIcon from '@mui/icons-material/Help';
import HowToRegIcon from '@mui/icons-material/HowToReg';
import LockOpenIcon from '@mui/icons-material/LockOpen';
import LogoutIcon from '@mui/icons-material/Logout';
import { Typography } from '@mui/material';
import List from '@mui/material/List';
import ListItem from '@mui/material/ListItem';
import ListItemIcon from '@mui/material/ListItemIcon';
import ListItemText from '@mui/material/ListItemText';
import { makeStyles } from "@mui/styles";
import React from 'react';
import { Link } from 'react-router-dom';
import { hasUnreadMessages } from '../chat/chatNotification';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';
const useStyles = makeStyles((theme) => ({
    root: {
        display: 'flex',
    },
    appBar: {
        zIndex: theme.zIndex.drawer + 1,
        transition: theme.transitions.create(['width', 'margin'], {
            easing: theme.transitions.easing.sharp,
            duration: theme.transitions.duration.leavingScreen,
        }),
    },
    menuButton: {
        marginRight: 36,
    },
    hide: {
        display: 'none',
    },
    content: {
        flexGrow: 1
    },
    white: {
        color: "white"
    },
    activeLink: {
        backgroundColor: MulwiColors.black + " !important",
        color: "white"
    },
}));



export function NavList(props) {

    const classes = useStyles();

    function LinkWithIcon(props) {
        return (<Link style={{ color: 'inherit', textDecoration: 'inherit' }} to={props.url}>
            <ListItem button
                selected={props.location === props.url || (props.urls && props.urls.indexOf(props.location) !== -1)}
                classes={{ selected: classes.activeLink }}>
                <ListItemIcon>
                    {props.children}
                </ListItemIcon>
                <ListItemText primary={props.label} style={{
                    color: props.ovColor || null
                }} />
            </ListItem>
        </Link>
        )
    }

    function linkCl(url, urls, c) {
        let f = props.location === url || (urls && urls.indexOf(props.location) !== -1)
        return {
            className: (f && classes.activeLink) || "",
            style: {
                color: c || (!f && MulwiColors.blueDark) || null
            }
        }
    }

    return (<React.Fragment>
        <List>
            {
                props.user && props.instructor && props.state && (
                    <>
                        <ListItem>
                            <ListItemText>
                                <center>
                                    <Typography>
                                        {locale2.USER_SPACE[props.lang]}
                                    </Typography>
                                </center>
                            </ListItemText>
                        </ListItem>
                    </>
                )
            }
            {
                props.user &&
                <>
                    <LinkWithIcon location={props.location} url="/profile"
                        label={locale2.MY_DATA[props.lang]}>
                        <Person {...linkCl("/profile")} />
                    </LinkWithIcon>
                    <LinkWithIcon location={props.location} url="/chat"
                        label={locale2.CHAT[props.lang]} ovColor={hasUnreadMessages(props.nots) && "orange"}>
                        <Chat {...linkCl("/chat", null, hasUnreadMessages(props.nots) && "orange")} />
                    </LinkWithIcon>
                    <LinkWithIcon location={props.location} url="/manage"
                        label={locale2.MGMT[props.lang]}>
                        <Edit {...linkCl("/manage")} />
                    </LinkWithIcon>
                    <LinkWithIcon location={props.location} url="/sub"
                        label={locale2.MY_CARNETS[props.lang]}>
                        <CardMembership {...linkCl("/sub")} />
                    </LinkWithIcon>
                    <LinkWithIcon location={props.location} url="/trainings/list"
                        label={locale2.MY_TRAININGS[props.lang]}>
                        <SportsBaseball {...linkCl("/trainings/list")} />
                    </LinkWithIcon>
                </>
            }
        </List>
        {(props.instructor && (
            <List>
                {
                    props.state && (
                        <ListItem>
                            <ListItemText>
                                <center>
                                    <Typography>
                                        {locale2.INSTR_SPACE[props.lang]}
                                    </Typography>
                                </center>
                            </ListItemText>
                        </ListItem>
                    )
                }
                <LinkWithIcon location={props.location} url="/instr_profile"
                    label={locale2.INSTR_PROFILE[props.lang]}>
                    <PersonPin {...linkCl("/instr_profile")} />
                </LinkWithIcon>
                <LinkWithIcon location={props.location} url="/configure"
                    label={locale2.CONFIGURATION[props.lang]}>
                    <Settings {...linkCl("/configure")} />
                </LinkWithIcon>
                {/* <LinkWithIcon location={props.location} url="/payments"
                    label={locale2.PAYMENTS[props.lang]}>
                    <CreditCardIcon {...linkCl("/payments")} />
                </LinkWithIcon> */}

                {/* <LinkWithIcon location={props.location} url="/invoice"
                    label={locale2.INVOICES[props.lang]}>
                    <Receipt {...linkCl("/invoice")} />
                </LinkWithIcon> */}

                <LinkWithIcon location={props.location} url="/harmonogram"
                    urls={["/harmonogram/list"]}
                    label={locale2.SCHEDULE[props.lang]}>
                    <EventIcon {...linkCl("/harmonogram", ["/harmonogram/list"])} />
                </LinkWithIcon>

                <LinkWithIcon location={props.location} url="/vacations"
                    label={locale2.HOLIDAYS[props.lang]}>
                    <BeachAccess {...linkCl("/vacations")} />
                </LinkWithIcon>
                <LinkWithIcon location={props.location} url="/group" label={locale2.LIMITS[props.lang]}>
                    <GroupTwoTone {...linkCl("/group")} />
                </LinkWithIcon>
                {/* <LinkWithIcon location={props.location} url="/dc" label={locale2.DCS[props.lang]}>
                    <Typography {...linkCl("/dc")}><strong>30%</strong></Typography>
                </LinkWithIcon> */}
                <LinkWithIcon location={props.location} url="/sm" label={locale2.CARNETS[props.lang]}>
                    <CardMembership {...linkCl("/sm")} />
                </LinkWithIcon>
            </List>
        ))}
        {
            props.user && !props.instructor &&
            <List>
                <LinkWithIcon location={props.location} url="/become_trainer"
                    label={locale2.BECOME_INSTRUCTOR[props.lang]}>
                    <SportsMmaOutlined {...linkCl("/become_trainer")} />
                </LinkWithIcon>
            </List>
        }
        <List>
            {
                props.state && 
                <ListItem>
                    <ListItemText>
                        <center>
                            <Typography>
                                {locale2.MISSION[props.lang]}
                                {//locale2.USER_SPACE[props.lang]
                                }
                            </Typography>
                        </center>
                    </ListItemText>
                </ListItem>
            }
            <LinkWithIcon location={props.location} url="/support"
                label={locale2.HELP[props.lang]}>
                <HelpIcon {...linkCl("/support")} />
            </LinkWithIcon>
            <LinkWithIcon location={props.location} url="/benefits_user"
                label={locale2.FOR_TRAINEE[props.lang]}>
                <HealthAndSafetyIcon {...linkCl("/benefits_user")} />
            </LinkWithIcon>
            <LinkWithIcon location={props.location} url="/benefits_instructor"
                label={locale2.FOR_INSTRUCTOR[props.lang]}>
                <AttachMoneyIcon {...linkCl("/benefits_instructor")} />
            </LinkWithIcon>
            <ListItem>
                {
                    // this is just to visually divide login,register and logout from the rest
                }
            </ListItem>
            {
                // while logged in show logout
                props.user &&
                <div onClick={() => {
                    props.logout()
                }} >

                    <LinkWithIcon location={props.location} url="/"
                        label={locale2.LOGOUT[props.lang]}>
                        <LogoutIcon {...linkCl("/")} />
                    </LinkWithIcon>
                </div>
            }

            {
                // while unlogged show login 
                // I know this is just ugly, but I don't have time to debug
                // why React won't allow me mixing conditions and throws unexpected results
                // future dev reading that - don't be mad
                !props.user &&

                <>
                    <LinkWithIcon location={props.location} url="/login"
                        label={locale2.LOGIN[props.lang]}>
                        <LockOpenIcon {...linkCl("/login")} />
                    </LinkWithIcon>
                    <LinkWithIcon location={props.location} url="/register"
                        label={locale2.REGISTER[props.lang]}>
                        <HowToRegIcon {...linkCl("/register")} />
                    </LinkWithIcon>
                </>
            }
        </List>
    </React.Fragment>)
}