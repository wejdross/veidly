import { Grid, Typography } from '@mui/material'
import React from 'react'

export function Info(props) {
    return (
        <Grid container
            style={props.style} 
            justify="space-between"
            alignItems="center"
            direction="row">
            
            <Grid item>
                <Typography variant="label">

                </Typography>
            </Grid>

        </Grid>
    )
}