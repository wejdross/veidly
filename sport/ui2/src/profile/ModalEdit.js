import {
    Button,
    CircularProgress,
    Dialog,
    DialogActions,
    DialogContent, DialogTitle,
    Grid,
    Typography
} from '@mui/material';
import { Check, Close } from '@mui/icons-material';
import EditIcon from '@mui/icons-material/Edit';
import React, { useEffect, useState } from 'react';
import { MulwiColors } from '../mulwiColors';
import { errToStr } from '../StatusDialog';
import { locale2 } from '../locale';
import { makeStyles } from "@mui/styles";

export function emptyInfo() {
    return {
        st: "",
        msg: ""
    }
}

export default function ModalEdit(props) {

    const useStyles = makeStyles({
        padded: {
            marginLeft: 100, marginRight: 100, marginTop: 50, marginBottom: 50
        }
    })
    const classes = useStyles()

    const [fnOpen, _setFnOpen] = useState(false)
    function setFnOpen(x) {
        setInfo({st: "", msg: ""})
        _setFnOpen(x)
    }
    function openEditor() {
        props.onOpen && props.onOpen()
        setFnOpen(true)
    }
    const [info,setInfo] = useState({
        st: "",
        msg: ""
    })

    useEffect(() => {
        setInfo(props.info)
    }, [props.info])

    useEffect(() => {
        setFnOpen(props.open)
    }, [props.open])

    useEffect(() => {
        if(props.open === true || props.open === false)
            setFnOpen(props.open)
    }, [props.open])

    function editForm() {
        return info.st === "" && (
            <React.Fragment>
                {(props.nocontent && (
                    props.content
                )) || (
                    <DialogContent>
                        {props.content}
                    </DialogContent>
                )}
                {!props.hideSaveButton && (
                    <DialogActions>
                        <Button onClick={() => setFnOpen(false)}>
                        {locale2.CANCEL[props.lang]}
                        </Button>
                        {props.customActions && props.customActions(save)}
                        <Button onClick={() => save()} color="primary" variant="contained" style={{
                            color: "white",
                            backgroundColor: props.btnColor || MulwiColors.greenDark,
                        }}>
                            {props.btnLabel || locale2.SAVE[props.lang]}
                        </Button>
                    </DialogActions>
                )}
                {props.hideSaveButton && (
                    <DialogActions>
                        {props.customActions && props.customActions(save)}
                        <Button onClick={() => setFnOpen(false)}>
                            {locale2.CLOSE[props.lang]}
                        </Button>
                    </DialogActions>
                )}
            </React.Fragment>
        )
    }

    function editWaiter() {
        return info.st === "wip" && (
            <CircularProgress className={classes.padded}/>
        )
    }

    function exitOk() {
        return info.st === "ok" && (
            <Check color="primary" fontSize="large" className={classes.padded}/>
        )
    }

    function editEx() {
        return info.st === "ex" && (
            <React.Fragment>
                <DialogContent>
                    {errToStr(info.msg)}
                </DialogContent>
                <DialogActions>
                    <Button color="primary" onClick={() => 
                            setInfo({st: "", msg: ""})
                        }>
                        {locale2.ONCE_AGAIN[props.lang]}
                    </Button>
                    <Button onClick={() => setFnOpen(false)} color="secondary">
                        <Close/>
                    </Button>
                </DialogActions>
            </React.Fragment>
        )
    }

    async function save(args) {
        try {
            setInfo({st: "wip", msg: ""})
            await props.onSave(args)
            setTimeout(async () => {
                setFnOpen(false)
            }, 100)
        } catch(ex) {
            setInfo({
                st: "ex",
                msg: ex
            })
        }
    }

    return (
        <React.Fragment>
            {(props.onlyButton && (
                <Button onClick={openEditor} {...props.buttonProps}>{props.label}</Button>
            )) || (
                <Grid
                    container
                    spacing={3}
                    alignItems="center"
                    >
                    {props.label && <Grid item xs={4} style={props.labelStyle}>
                        {props.label}
                    </Grid>}
                    
                    <Grid item xs={props.label ? 6 : 10}>
                        {(props.custom ? (
                                props.value
                            ) : (
                                (props.multiline ? (
                                        <pre style={{ fontFamily: 'inherit' }}>
                                        {props.value}
                                    </pre>
                                ) : (
                                    <Typography>
                                        {props.value}
                                    </Typography>
                                ))
                            ))}
                    </Grid>
                    {props.content && (
                        <Grid item xs={2}>
                            <Button onClick={openEditor} color="primary" size='small' aria-label="edit">
                                <EditIcon />
                            </Button>
                        </Grid>
                    )}
                </Grid>
            )}
            <Dialog open={Boolean(fnOpen)} onClose={() => setFnOpen(false)} aria-labelledby="form-dialog-title">
                <DialogTitle id="form-dialog-title">{props.title}</DialogTitle>
                {editForm()}
                {editWaiter()}
                {exitOk()}
                {editEx()}
            </Dialog>
        </React.Fragment>
    )
}