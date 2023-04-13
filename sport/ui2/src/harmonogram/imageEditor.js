import {
    Button,
    Dialog,
    DialogActions,
    DialogContent,
    DialogTitle, Grid, Typography
} from '@mui/material';
import { Add, Delete } from '@mui/icons-material';
import EditIcon from '@mui/icons-material/Edit';
import React, { useState } from 'react';
import { deleteImg as apiDeleteImg, postImg } from '../apicalls/instructor.api';
import { randomString, sprintf } from '../helpers';
import { MulwiColors } from '../mulwiColors';
import { getErrorDialog } from '../StatusDialog';
import { locale2 } from '../locale';

// api will also validate this number
const maxImgs = 6

  export default function TrainingImageEditor(props) {

    async function onImgSelected(e, ismain) {
        if(!e.target.files || !e.target.files[0]) return
        // setOpen(true)
        // setLoading(true)
        let formData = new FormData()
        formData.append("image", e.target.files[0])
        formData.append("training_id", props.training.ID)
        if(ismain)
            formData.append("main", "1")
        try {
            await postImg(formData)
            props.onChange()
        } catch(ex) {
            props.setInfo(getErrorDialog(locale2.ERROR[props.lang], ex))
        }

        // setFd(formData)
        // // load img to preview
        // var fr = new FileReader()
        // fr.onload = function () {
        //     //document.getElementById("prevavatar").src = fr.result
        //     //setLoading(false)
        //     props.onChange()
        // }
        // fr.readAsDataURL(e.target.files[0])

    }

    async function deleteImg(id) {
        try {
            await apiDeleteImg(id, props.training.ID)
            props.onChange()
        } catch(ex) {
            props.setInfo(getErrorDialog(
                locale2.ERROR[props.lang], 
                ex))
        }
    }

    const [fnOpen, setFnOpen] = useState(false)

    function img(src, id, large, isproto) {
        if(isproto) {
            let id = randomString(10)
            return (<React.Fragment>
                <input 
                    onChange={e => onImgSelected(e, large)}
                    accept="image/*" 
                    style={{display:"none"}} 
                    id={id} type="file" />

                    <label htmlFor={id}>
                        <Button style={{
                                width: large ? 470 : 230,
                                height: large ? 310 : 150,
                            }} aria-label="upload picture" component="span">
                                <Add/>
                        </Button>
                    </label>
                </React.Fragment>)
        }
        if(!src) return null
        return <React.Fragment>
            <div style={{
                position:"relative"
            }}>
                <img style={{
                maxHeight: large ? 320 : 160,
                maxWidth: large ? 480 : 240,
                marginBottom: -5,
                paddingBottom: 5
            }} src={src} />
            <Button style={{
                position:"absolute",
                top: 0,
                right: 0,
                color: MulwiColors.redError
            }} onClick={() => deleteImg(id)}>
                <Delete/>
            </Button>
            </div>
        </React.Fragment>
    }

    function renderSecondaryImgs(urls, ids) {
        let rows = []
        if(urls) {
            for(let i = 0; i < urls.length; i+=2) {
                rows.push(
                    <Grid item key={i}>
                <Grid container direction="column" alignItems="stretch">
                    <Grid item>
                        {img(urls[i], ids[i], 0, 0)}
                    </Grid>
                    <Grid item>
                        {urls.length > i+1 ? img(urls[i+1], ids[i+1], 0,0) : img("", "", 0, 1)}
                    </Grid>
                </Grid>
                        </Grid>)
            }
        }
        if(!urls || ((urls.length % 2) === 0) && urls.length < maxImgs) {
            rows.push(
                <Grid key={999} item><Grid container direction="column" justify="center" alignItems="stretch">
            <Grid item>
                {img("", "", 0, 1)}
            </Grid>
        </Grid></Grid>)
        }
        return rows
    }

    function modalContent() {
        return (<Dialog open={fnOpen}
                        maxWidth={false}
                        onClose={() => setFnOpen(false)} >
            <div>
                <DialogTitle>{locale2.TRAINING_PHOTOS[props.lang]}</DialogTitle>
                <DialogContent >
                    <Typography>{sprintf(locale2.PHOTO_NUM_FMT[props.lang], maxImgs)}</Typography>
                    <Typography>{ locale2.PHOTO_DISPLAY[props.lang]}</Typography>
                    <br/>
                    <Grid container direction="row" justify="center">
                        <Grid item style={{
                            marginRight: 5
                        }}>
                            {props.training.MainImgUrl ? 
                                img(props.training.MainImgUrl, props.training.MainImgID, 1) 
                                : 
                                img("", "", 1, 1)}
                        </Grid>
                            {renderSecondaryImgs(
                                props.training.SecondaryImgUrls, 
                                props.training.SecondaryImgIDs, 
                                Boolean(props.training.MainImgUrl))}
                    </Grid>
                </DialogContent>
                <DialogActions>
                    <Button onClick={() => setFnOpen(false)}>{ locale2.CLOSE[props.lang] }</Button>
                </DialogActions>
            </div>
        </Dialog>)
    }

    return (
        <React.Fragment>
            <Grid container spacing={3}>
                <Grid item xs={4}>
                    <Typography
                        style={{color:"gray"}}
                        component={'span'}>{ locale2.TRAINING_PHOTOS[props.lang] }</Typography>
                </Grid>
                <Grid container item xs={6}>
                    {props.training && props.training.MainImgUrl ? 1 : 0} + {props.training && props.training.SecondaryImgUrls ? props.training.SecondaryImgUrls.length : 0}
                </Grid>
                <Grid container item xs={2}>
                    <Button onClick={() => setFnOpen(true)} 
                            color="primary" size='small' aria-label="edit">
                    <EditIcon/>
                    </Button>
                </Grid>
            </Grid>
            {props.training && modalContent()}
        </React.Fragment>
    )
  }