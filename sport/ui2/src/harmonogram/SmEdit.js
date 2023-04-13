import { Typography } from '@mui/material'
import React from 'react'
import ModalEdit from '../profile/ModalEdit'
import { locale2 } from '../locale';

export function SmEdit(props) {
    if(!props.sm) return null
    
    return (<ModalEdit 
        lang={props.lang}
        hideSaveButton nocontent
        title={locale2.CARNETS[props.lang]}
        label={<Typography style={{color:"gray"}}>{locale2.CARNETS[props.lang]}</Typography>}
        value={props.sm.map((g,i) => (i === 0 ? (
            <span key={g.Name}>{g.Name}</span>
        ) : (
            <span key={g.Name}>, {g.Name}</span>
        )))}
        content={null}
    />)
}