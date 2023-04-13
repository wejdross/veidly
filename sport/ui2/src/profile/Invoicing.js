import { Print } from '@mui/icons-material';
import {
    IconButton, Table, TableBody, TableCell, TableHead, TableRow, Typography,
    useTheme
} from '@mui/material';
import { makeStyles } from "@mui/styles";
import React, { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { getInvoices } from '../apicalls/instructor.api';
import CardWithBg from '../card/cardWithBg';
import { API_URL } from '../conf';
import { prettyPrintDate } from '../harmonogram/trainingDetails';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';
import { getErrorDialog, getNullDialog, StatusDialog } from '../StatusDialog';

export default function Invoicing(props) {


    const theme = useTheme()

    const useStyles = makeStyles({
        divider: {
            marginTop: 10,
            marginBottom: 10,
        },
        widthSettings: {
            [theme.breakpoints.up('sm')]: {
                minWidth: 600,
            }
        },
        purple: {
            backgroundColor: "purple"
        }
    })

    const classes = useStyles()
    const [invoices, setInvoices] = useState([])
    const [info, setInfo] = useState(getNullDialog())

    async function refresh() {
        try {
            let i = JSON.parse(await getInvoices())
            setInvoices(i || [])
            console.log(i)
        } catch (ex) {
            setInfo(getErrorDialog(locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }

    async function print(invoiceID) {
        let url = API_URL + "/api/invoice/print?id=" + invoiceID
        document.getElementById('invoiceframe').src = url
    }

    function printRef(type, id) {
        switch(type) {
        case "RSV": 
            return (<Link style={{
                textDecoration:"none"
            }} to={"/rsv_details?id=" + id}>
                Rezerwacja
            </Link>)
        case "SUB":
            return (<Link style={{
                textDecoration:"none"
            }} to={"/sub_details?id=" + id}>
                Karnet
            </Link>)
        default:
            return ""
        }
    }

    useEffect(() => {
        refresh()
    }, [props.instructor])

    if (!props.instructor)
        return null

    return (<React.Fragment>
        <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
        <iframe id="invoiceframe" style={{
            display:"none"
        }}></iframe>
        <CardWithBg img="/static/form-backgrounds/moto.webp">
            <div className={classes.widthSettings}>

                <Typography>
                    Twoje faktury
                </Typography>

                <Table>
                    <TableHead>
                        <TableRow>
                            <TableCell>
                                Numer
                            </TableCell>
                            <TableCell>
                                Dotyczy
                            </TableCell>
                            <TableCell>
                                Data wystawienia
                            </TableCell>
                            <TableCell>
                                Kwota
                            </TableCell>
                            <TableCell>
                                Drukuj
                            </TableCell>
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {invoices && invoices.map((invoice, i) => (<TableRow key={i}>
                            <TableCell>
                                {invoice.Number}
                            </TableCell>
                            <TableCell>
                                {printRef(invoice.ObjType, invoice.ObjID)}
                            </TableCell>
                            <TableCell>
                                {prettyPrintDate(new Date(invoice.DateOfIssue))}
                            </TableCell>
                            <TableCell>
                                PLN {invoice.Paid / 100}
                            </TableCell>
                            <TableCell align='center'>
                                <IconButton onClick={() => print(invoice.ID)}>
                                    <Print style={{
                                        color:MulwiColors.blueDark
                                    }}/>
                                </IconButton>
                            </TableCell>
                        </TableRow>))}
                    </TableBody>
                </Table>

            </div>
        </CardWithBg>
    </React.Fragment>)
}