import {
    Button
} from '@mui/material'
import { NavigateBefore, NavigateNext } from '@mui/icons-material'
import React from 'react'
import { lMonths } from '../helpers'

export function MonthLabel(props) {
    let months = []
    for(let i =0; i < 12; i++) {
        months.push(lMonths[i][props.lang])
    }

    return (
        <React.Fragment>
            <center>
                <div style={{
                    marginBottom: 10
                }}>
                    <Button
                        onClick={() => {
                            let c = new Date(props.monthDate)
                            c.setMonth(c.getMonth() - 1)
                            c = new Date(c)
                            c.setDate(1)
                            props.setMonthDate(c)
                            //setWkFromMonth(c)
                        }}
                        variant="text" color="primary"><NavigateBefore /></Button>
                    <div style={{ display: "inline-block", minWidth: 120, textAlign: "center" }}>
                        {months[props.monthDate.getMonth()]} {props.monthDate.getFullYear()}
                    </div>
                    <Button
                        onClick={() => {
                            let c = new Date(props.monthDate)
                            c.setMonth(c.getMonth() + 1)
                            c = new Date(c)
                            c.setDate(1)
                            props.setMonthDate(c)
                            //setWkFromMonth(c)
                        }}
                        variant="text" color="primary"><NavigateNext /></Button>
                </div></center>
        </React.Fragment>
    )
}
