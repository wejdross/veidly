import { ButtonBase, Grid } from '@mui/material'
import React, { useState } from 'react'
import { ImgOverlay } from './ImgOverlay'

export function ImgGrid(props) {
    
    const [open, setOpen] = useState(false)

    if(!props.MainImgUrl || !props.SecondaryImgUrls) return null
    
    let w = 320
    let h = 270
    if(!props.SecondaryImgUrls) {
        w = 480
        h = 240
    }

    // let urls = props.SecondaryImgUrls

    // return (<GridList cellHeight={160} >
    //     <GridListTile cols={2} rows={2}>
    //             <img src={props.MainImgUrl} alt=":(" />
    //     </GridListTile>
    //     {urls.map(t => (
    //         <GridListTile>
    //             <img src={t} alt=":(" />
    //         </GridListTile>
    //     ))}
    // </GridList>)

    const brad = 15

    function renderSecondaryImgs(urls) {
        let rows = []
        let isLast = false
        if(urls) {
            for(let i = 0; i < urls.length; i+=2) {
                if((i+2) >= (urls.length)) {
                    isLast = true
                }
                let hasBottom = false
                if(urls.length > i+1) {
                    hasBottom = true
                }
                rows.push(
                    <Grid item>
                        <Grid container direction="column" alignItems="stretch">
                            <Grid item  >
                                <img alt="training img" style={{
                                    maxHeight: h/2,
                                    maxWidth: w/2,
                                    borderTopRightRadius: isLast ? brad: 0,
                                    borderBottomRightRadius: (!hasBottom && isLast) ? brad: 0,
                                    marginBottom: -5,
                                }} src={urls[i]} />
                            </Grid>
                            {hasBottom && (
                            <Grid item>
                                <img alt="training img" style={{
                                    maxHeight: h/2,
                                    maxWidth: w/2,
                                    borderBottomRightRadius: isLast ? brad: 0,
                                    marginBottom: -5,
                                }} src={urls[i+1]} />
                            </Grid>
                            )}
                        </Grid>
                    </Grid>)
            }
        }
        return rows
    }

    let hasSec = props.SecondaryImgUrls && props.SecondaryImgUrls.length > 0
    

    return (<React.Fragment>
        <ImgOverlay 
            MainImgUrl={props.MainImgUrl} 
            SecondaryImgUrls={props.SecondaryImgUrls}
            open={open} setOpen={setOpen} />
        <ButtonBase style={{
            borderRadius: brad
            }} onClick={() => {
                setOpen(true)
            }}>
            <Grid container direction="row" justify="center"
                        alignItems="center">
                {props.MainImgUrl && <Grid item inline><img alt="training img" style={{
                        maxHeight: h,
                        maxWidth: w,
                        borderTopLeftRadius: brad,
                        borderBottomLeftRadius: brad,
                        borderTopRightRadius: hasSec ? 0 : brad,
                        borderBottomRightRadius: hasSec ? 0 : brad,
                        marginBottom: -5,
                    }} src={props.MainImgUrl} />
                </Grid>
                }
                {renderSecondaryImgs(props.SecondaryImgUrls)}
            </Grid>
        </ButtonBase>
    </React.Fragment>)
}