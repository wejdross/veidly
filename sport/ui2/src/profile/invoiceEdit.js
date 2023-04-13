import { TextField, useTheme, Typography } from '@mui/material'
import { makeStyles } from "@mui/styles";
import React, { useEffect, useState } from 'react'
import { patchUserData } from '../apicalls/user.api'
import ModalEdit from './ModalEdit'
import { locale2 } from '../locale'
import { MulwiColors } from '../mulwiColors'
import { PATCHInstructor } from '../apicalls/instructor.api'
import InstructorPayments from './InstructorPayments'

export function linesArrToStr(arr) {
    let lines = ""
    if (!arr)
        return lines
    for (let i = 0; i < arr.length; i++) {
        lines += arr[i] + "\n"
    }
    return lines
}

export default function InvoiceEdit(props) {

    const [invoiceLines, setInvoiceLines] = useState("")
    const theme = useTheme()

    const useStyles = makeStyles({
        dialog: {
            [theme.breakpoints.only('xs')]: {
                width: 290,
            },
            [theme.breakpoints.only('sm')]: {
                width: 500,
            },
            [theme.breakpoints.up('md')]: {
                width: 400,
            },
        }
    })

    const styles = useStyles()

    useEffect(() => {
        if (props.instr && props.instr.InvoiceLines) {
            setInvoiceLines(linesArrToStr(props.instr.InvoiceLines))
        }
    }, [props.instr])

    async function save(del) {
        if (!invoiceLines)
            return
        let instr = { ...props.instr }
        let lines = invoiceLines.split("\n")
        let safeLines = []
        for (let i = 0; i < lines.length; i++) {
            if (lines[i]) {
                safeLines.push(lines[i])
            }
        }
        lines = safeLines
        instr.InvoiceLines = lines
        await PATCHInstructor(instr)
        props.main.refresh()
    }

    if (!props.instr) return null

    return (<React.Fragment>
        <ModalEdit
            lang={props.lang}

            buttonProps={props.buttonProps}
            onlyButton={props.onlyButton}

            title={`${locale2.INVOICE_DATA[props.lang]}`}
            label={props.label || `${locale2.INVOICE_DATA[props.lang]}`}
            labelStyle={{ color: MulwiColors.subtitleTypography }}

            custom
            onSave={save}
            value={
                <Typography style={{ overflowWrap: "break-word", whiteSpace: "pre-wrap", maxWidth: 320 }} variant="body2">
                    {linesArrToStr(props.instr.InvoiceLines)}
                </Typography>}
            content={(<div>
                <Typography
                        variant="body2">
                    {locale2.INVOICE_DISCLAIMER[props.lang]}
                </Typography>
                <TextField
                    className={styles.dialog}
                    multiline
                    inputProps={
                        { maxLength: 250 }
                    }
                    rows={5}
                    variant={"outlined"}
                    helperText={((invoiceLines && invoiceLines.length) || 0) + "/250"}
                    fullWidth={true}
                    value={invoiceLines}
                    onChange={(event => (setInvoiceLines(event.target.value)))}
                />
            </div>)} />
    </React.Fragment>)
}

