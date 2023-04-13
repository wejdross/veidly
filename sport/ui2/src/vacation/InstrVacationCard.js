import { ArrowRightAltSharp, Delete } from '@mui/icons-material';
import {
    Button, CircularProgress, Grid, IconButton,
    TablePagination, Typography, useMediaQuery
} from '@mui/material';
import Paper from '@mui/material/Paper';
import { useTheme } from '@mui/material/styles';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TableRow from '@mui/material/TableRow';
import TextField from '@mui/material/TextField';
import { DatePicker } from '@mui/x-date-pickers/DatePicker';

import React, { useEffect, useState } from 'react';
import {
    deleteVacation as apiDeleteVacation,
    getVacations, postVacation
} from '../apicalls/instructor.api';
import CardWithBg from '../card/cardWithBg';
import { prettyPrintDate, prettyPrintDay } from '../harmonogram/trainingDetails';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';
import {
    getDialogWithOptions, getErrorDialog,
    getNullDialog, StatusDialog
} from '../StatusDialog';

export default function InstrVacationCard(props) {

    const [rowsPerPage, setRowsPerPage] = React.useState(5);
    const [page, setPage] = React.useState(0);


    const t = useTheme()
    const isLowRes = useMediaQuery(t.breakpoints.down('sm'))

    const [data, setData] = useState([])

    const handleChangePage = (event, newPage) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event) => {
        setRowsPerPage(parseInt(event.target.value, 10));
        setPage(0);
    };

    const [dateStart, setdateStart] = useState(() => {
        let d = new Date()
        d.setHours(0, 0, 0, 0)
        return d
    })
    const [dateEnd, setdateEnd] = useState(() => {
        let d = new Date()
        d.setHours(23, 59, 0, 0)
        return d
    })

    async function refresh() {
        try {
            let d = await getVacations()
            setData(d)
            setPage(0)
        } catch (ex) {
            setInfo(getErrorDialog(
                locale2.SOMETHING_WENT_WRONG_FETCH_VACATION[props.lang], ex))
        }
    }

    useEffect(() => {
        refresh()
    // eslint-disable-next-line react-hooks/exhaustive-deps
    }, [])

    async function createVacation() {
        setInfo(getDialogWithOptions(
                locale2.WAIT_A_MOMENT[props.lang], <Grid style={{
            minWidth: 200,
        }}>
            <center>
                <CircularProgress />
            </center>
        </Grid>, null, true))
        try {
            await postVacation({
                DateStart: dateStart,
                DateEnd: dateEnd
            })
            await refresh()
            setInfo(getDialogWithOptions("OK", <Typography>
                {locale2.YOUR_VACATION_WILL_START[props.lang]}
                    {prettyPrintDay(dateStart)}
            </Typography>))
        } catch (ex) {
            setInfo(getErrorDialog(
                locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }

    async function deleteVacation(id) {
        setInfo(getDialogWithOptions(
                locale2.WAIT_A_MOMENT[props.lang], <Grid style={{
            minWidth: 200,
        }}>
            <center>
                <CircularProgress />
            </center>
        </Grid>, null, true))
        try {
            await apiDeleteVacation(id)
            await refresh()
            setInfo(getNullDialog())
        } catch (ex) {
            setInfo(getErrorDialog(
                    locale2.SOMETHING_WENT_WRONG[props.lang], ex))
        }
    }

    const [info, setInfo] = useState(getNullDialog())

    return (
        <React.Fragment>
            <StatusDialog lang={props.lang} info={info} setInfo={setInfo} />
            <CardWithBg img="/static/form-backgrounds/kosz.webp">
                <Grid container>
                    <Grid item>
                        <Typography style={{
                            marginBottom: 7
                        }} variant="h5">
                            {locale2.VACATION_WHENEVER_YOU_WANT[props.lang]}
                        </Typography>
                        <Typography>
                            {locale2.NOONE_WILL_BE_ABLE_TO_SIGN_UP_DURING_VACATION[props.lang]}
                        </Typography>
                        <Typography style={{
                            padding: 10,
                            color: MulwiColors.blueDark
                        }}>
                            {locale2.YOU_WILL_HAVE_TO_COMPLETE_RSVS[props.lang]}
                        </Typography>

                        <Grid container
                            direction={isLowRes ? "column" : "row"}
                            justify="space-between"
                            alignItems="center">
                            <Grid item>
                                <Grid container direction="column">
                                    <Grid item>
                                        <DatePicker fullWidth
                                            renderInput={(params) => <TextField {...params} />}
                                            margin="normal"
                                            id="date-picker-dialog"
                                            label={locale2.START_DATE[props.lang]}
                                            inputFormat='"dd/MM/yyyy"'
                                            minDate={new Date()}
                                            minDateMessage={locale2.CANT_ADD_VACATION_IN_THE_PAST[props.lang]}
                                            value={dateStart}
                                            onChange={v => {
                                                if (v instanceof Date && isFinite(v)) {
                                                    let d = new Date(dateStart)
                                                    d.setMonth(v.getMonth())
                                                    d.setDate(v.getDate())
                                                    d.setFullYear(v.getFullYear())
                                                    d.setHours(0)
                                                    d.setMinutes(0)
                                                    d.setSeconds(0)
                                                    d.setMilliseconds(0)
                                                    setdateStart(d)
                                                }
                                            }}
                                            KeyboardButtonProps={{
                                                'aria-label': 'change date',
                                            }} />
                                    </Grid>
                                </Grid>
                            </Grid>
                            <Grid item>
                                <ArrowRightAltSharp fontSize="large" style={{
                                    color: MulwiColors.blueLight
                                }} />
                            </Grid>
                            <Grid item>
                                <Grid container direction="column">
                                    <Grid item>
                                        <DatePicker fullWidth
                                            renderInput={(params) => <TextField {...params} />}
                                            margin="normal"
                                            id="date-picker-dialog2"
                                            label={locale2.END_DATE[props.lang]}
                                            inputFormat="dd/MM/yyyy"
                                            value={dateEnd}
                                            minDate={dateStart}
                                            minDateMessage={locale2.START_SHOULD_BE_AFTER_END[props.lang]}
                                            onChange={v => {
                                                if (v instanceof Date && isFinite(v)) {
                                                    let d = new Date(dateEnd)
                                                    d.setMonth(v.getMonth())
                                                    d.setDate(v.getDate())
                                                    d.setFullYear(v.getFullYear())
                                                    d.setSeconds(0)
                                                    d.setMilliseconds(0)
                                                    setdateEnd(d)
                                                }
                                            }}
                                            KeyboardButtonProps={{
                                                'aria-label': 'change date',
                                            }} />
                                    </Grid>
                                </Grid>
                            </Grid>
                        </Grid>

                        <div style={{
                            marginBottom: 10
                        }}></div>

                        <Button fullWidth variant="contained" style={{
                            color: "white",
                            backgroundColor: MulwiColors.greenDark
                        }} onClick={createVacation}>
                            {locale2.ADD[props.lang]}
                        </Button>

                        <div style={{
                            marginBottom: 20
                        }}></div>

                        <Typography style={{
                            marginBottom: 15
                        }} variant="h5">{locale2.INCOMING_VACATION[props.lang]}</Typography>
                        <TableContainer component={Paper}>
                            <Table aria-label="simple table" >
                                <TableHead>
                                    <TableRow>
                                        <TableCell align="center">
                                            {locale2.FROM[props.lang]}
                                        </TableCell>
                                        <TableCell align="center">
                                            {locale2.TO[props.lang]}
                                        </TableCell>
                                        <TableCell align="center">
                                            {locale2.DELETE[props.lang]}
                                    </TableCell>
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {data && data.map((d, i) => (
                                        <TableRow key={i}>
                                            <TableCell align="center" component="th" scope="row">
                                                {prettyPrintDate(d.DateStart)}
                                            </TableCell>
                                            <TableCell align="center" component="th" scope="row">
                                                {prettyPrintDate(d.DateEnd)}
                                            </TableCell>
                                            <TableCell align="center">
                                                <IconButton onClick={() => deleteVacation(d.ID)}>
                                                    <Delete style={{
                                                        color: MulwiColors.redError
                                                    }} />
                                                </IconButton>
                                            </TableCell>
                                        </TableRow>
                                    ))}
                                </TableBody>
                            </Table>
                        </TableContainer>
                        <TablePagination
                            rowsPerPageOptions={[5, 10, 25]}
                            component="div"
                            count={data.length}
                            rowsPerPage={rowsPerPage}
                            page={page}
                            onPageChange={handleChangePage}
                            onRowsPerPageChange={handleChangeRowsPerPage}
                            labelRowsPerPage={isLowRes ? "w." 
                                : locale2.NUMBER_OF_ROWS_PER_PAGE[props.lang]}
                            labelDisplayedRows={p => p.from + "-" + p.to + " z " + p.count}
                        />
                    </Grid>
                </Grid>
            </CardWithBg>
        </React.Fragment>
    );
}