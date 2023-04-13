import React from 'react';
import Card from '@mui/material/Card';
import {
    Grid,
    CardContent,
    useTheme, CardHeader, Typography, Divider
} from '@mui/material'
import { DrawerBig } from './DrawerBig';
import DrawerSmall from './DrawerSmall';
import { makeStyles } from "@mui/styles";
import useMediaQuery from '@mui/material/useMediaQuery';
import '../veidly-styles.css';

const useStyles = makeStyles((theme) => ({
    outer: props => ({

        backgroundImage: "url(" + props.img + ")",
        backgroundSize: "cover",
        borderRadius: 10,
        marginTop: 20,
        [theme.breakpoints.down('xs')]: {
            backgroundImage: "none",
            width: "100%",
            marginTop: 60,
        }
    }),
    card: props => ({
        margin: 70,
        backgroundColor: props.bg || "rgba(250,250,250,0.85)",
        [theme.breakpoints.down('xs')]: {
            margin: 0,
            backgroundColor: props.bg || "rgba(250,250,250,1)",
            boxShadow: "none",
        },
    })
}))

export default function CardWithBg(props) {
    const classes = useStyles({ img: props.img, bg: props.bg });
    const theme = useTheme()

    const isLowRes = useMediaQuery(theme.breakpoints.down('sm'))

    function getpad() {
        if (props.nocontent) return 0
        return isLowRes ? 10 : 40
    }

    return (

        <Grid container
            direction="column"
            justify="center"
            alignItems="center">
            {
                /*
                    All generic styles should be assigned in veidly-styles.css
                    but for some component we need custom sizes
                */
            }
            <div className={"veidlyDataSheet"} style={props.style || {}}>
                <Grid container
                    direction="column"
                    justify="center"
                    alignItems="center">

                    {props.header && (
                        <>
                            <Typography className='cardWithBgHeader' variant="h6" align={"center"} style={{
                                marginTop: 20,
                                marginBottom: 5,
                            }}>
                                {props.header}
                            </Typography>
                            <Divider className='cardWithBgDivider' style={{marginBottom: 10}}/>
                        </>
                    )
                    }
                            <Typography className='cardWithBgSubHeader' variant="h5" align='left'>
                                {props.subheader}
                            </Typography>
                </Grid>
                {(props.nocontent && (
                    props.children
                )) || (
                        props.children
                    )}
            </div>
        </Grid>
    )

}