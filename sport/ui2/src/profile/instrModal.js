import { TextField, Typography } from '@mui/material';
import React from 'react';
import { locale2 } from '../locale';

export function InstrModal(props) {
    return (
        <React.Fragment>
            <Typography>
                {locale2.SINCE_WHEN_INSTR[props.lang]}
            </Typography>

            <TextField
                variant="outlined"
                size="small"
                type="number"
                value={props.instrInfo.yearExp || 2021}
                onChange={(e) => {
                    let v = Number(e.target.value)
                    if (v > new Date().getFullYear()) return
                    props.setInstrInfo(c => ({ ...c, yearExp: v }))
                }}
                InputProps={{
                    max: 100
                }}
            />
        </React.Fragment>
    )
}
