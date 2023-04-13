import {
    useMediaQuery,
    useTheme
} from '@mui/material';
import React from 'react';
import { DrawerBig } from './DrawerBig';
import DrawerSmall from './DrawerSmall';


export function DrawerResponsive(props) {
    const theme = useTheme()
    const isLowRes = useMediaQuery(theme.breakpoints.down('md'))
    //const isXLowRes = useMediaQuery(theme.breakpoints.down('xs'))
    
    return isLowRes ? (
        <DrawerSmall {...props}/>
    ) : (
        <DrawerBig
            {...props} >

        </DrawerBig>
    )
}