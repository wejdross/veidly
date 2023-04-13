import React from 'react';
import { Typography } from '@mui/material';
import { Link } from "react-router-dom";
import { MulwiColors } from './mulwiColors';

export function LinkWithTypo(props) {

    const [hovered, setHovered] = React.useState(false)

    return (
        
        <Link to={props.to} style={{
            textDecoration: "none",
            color: hovered ? MulwiColors.greenDark : MulwiColors.blueDark,
        }}
        onMouseEnter={() => setHovered(true)}
        onMouseLeave={() => setHovered(false)}
        >
            <Typography variant='h6' style={{
                fontWeight: 500,
            }}>
                <strong>{props.text}</strong>
            </Typography>
        </Link>
    )
}