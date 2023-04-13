import {
    Grid,
    Typography,
} from "@mui/material";
import React, { useState } from "react";
import Button from "@mui/material/Button";
import { MulwiColors } from "./mulwiColors";
import { Link } from "react-router-dom";
import { locale2 } from "./locale";
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';


export function JoinUs(props) {

    const [hovered, setHovered] = useState(false)

    return (<React.Fragment>
        <Grid
            container
            direction={"column"}
            justifyContent={"center"}
            alignItems={"center"}
            style={{
                marginTop: 100,
            }} >
            <Typography variant={"h3"} style={{ marginTop: 10 }} align={"center"}>
                {locale2.ARE_YOU_INSTRUCTOR[props.lang]}
            </Typography>
            <Link to="/register" style={{ textDecoration: "none" }}>
                <Button
                    onMouseEnter={(e) => {setHovered(true)}} onMouseLeave={(e) => {setHovered(false)}}
                    variant="contained"
                    disableElevation
                    style={{
                        backgroundColor: hovered ? MulwiColors.blueDark : MulwiColors.greenDark,
                        color: "white",
                        height: 100,
                        width: 270,
                        borderRadius: 100,
                        marginTop: 20,
                    }}
                >
                    <Grid container alignItems="center" justifyContent="center">
                        <Typography variant="h6" style={{ color: MulwiColors.whiteSurface, fontWeight: 400, fontSize: "1.4em" }}>
                            {locale2.JOIN_US[props.lang]}
                        </Typography>
                        <ArrowForwardIcon style={{ marginTop: 1, marginLeft: 30, fontSize: "2em" }} />
                    </Grid>
                </Button>
            </Link>

        </Grid>
    </React.Fragment>)
}