import { Typography } from '@mui/material'
import React from 'react'
import ModalEdit from '../profile/ModalEdit'
import { locale2 } from '../locale'

export function DcEdit(props) {
    if(!props.dc) return null
    
    return (<ModalEdit 
        lang={props.lang}
        hideSaveButton nocontent
        title={locale2.DCS[props.lang]}
        label={<Typography style={{color:"gray"}}>{locale2.DCS[props.lang]}</Typography>}
        value={props.dc.map((g,i) => 
            <React.Fragment key={g.Name}>
                <Typography variant="body2">
                    <strong>{g.Name}</strong>, {g.Discount}%, {g.RedeemedQuantity}/{g.Quantity}
                </Typography>
            </React.Fragment>)}
        content={null}
    />)
}