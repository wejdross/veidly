import { Avatar, Grid, Typography } from '@mui/material'
import { Rating } from '@mui/lab'
import React from 'react'
import { locale2 } from '../locale'

export function ReviewContent(props) {
    if(!props.content) return null
    return (<Grid container direction="row" alignItems="center">
    {props.content.UserInfo.AvatarUrl && (
        <Grid item style={{
            marginRight: 10
        }}>
            <Avatar src={props.content.UserInfo.AvatarUrl} />
        </Grid>
    )}
    <Grid item>
        <Typography><strong>{props.content.UserInfo.Name}</strong></Typography>
        {props.content.Mark && (
            <Rating size="small" readOnly value={props.content.Mark} max={6} />
        )}
        {props.content.Review && <Typography>
            {props.content.Review}
        </Typography>}
    </Grid>
</Grid>)
}

export function UserReview(props) {
    
    return (<React.Fragment>
        <Typography variant="h5" style={{
            marginBottom: 10
        }}>{locale2.YOUR_REVIEW[props.lang]}</Typography>
        <ReviewContent content={props.content} />
    </React.Fragment>)
}
