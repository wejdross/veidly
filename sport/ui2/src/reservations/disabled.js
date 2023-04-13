import {  Grid, Tooltip } from '@mui/material';
import React, { useEffect, useState } from 'react';
import { locale2 } from '../locale';

export function DisabledSupport(props) {

    if (!props.training)
        return null

    return (
        <Grid container direction="row">
            {
                props.training.TrainingSupportsDisabled && (
                    <Tooltip title={locale2.TRAINING_SUPPORTS_DISABLED[props.lang]}>
                        <Grid item>
                            <img height={props.small ? 25 : null} src="/1.svg" alt="" />
                        </Grid>
                    </Tooltip>
                )
            }
            {
                props.training.PlaceSupportsDisabled && (
                    <Tooltip title={locale2.PLACE_SUPPORTS_DISABLED[props.lang]}>
                        <Grid item>
                                <img height={props.small ? 25 : null} src="/2.svg" alt="" />
                        </Grid>
                    </Tooltip>
                )
            }
        </Grid>
    )
}