import { Container, Typography } from "@mui/material";
import React from "react";
import { locale2 } from "../locale";

export default function NoResults(props) {

    return (<React.Fragment>
        <Container>
            <center>
                <Typography variant="h5" style={{marginTop: 20}}>
                    {locale2.COULDNT_FIND_TRAININGS[props.lang]}
                </Typography>
            </center>
        </Container>
    </React.Fragment>)
}