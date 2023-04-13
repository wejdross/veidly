import React from 'react'
import Backdrop from '@mui/material/Backdrop';
import ImageList from '@mui/material/ImageList';
import ImageListItem from '@mui/material/ImageListItem';

export function ImgOverlay(props) {
    return (
        <>
            {
                props.MainImgUrl &&
                <Backdrop style={{
                    zIndex: 9999
                }} open={props.open} onClick={() => {
                    props.setOpen(false)
                }}>
                    <ImageList cellHeight={200} style={{ maxWidth: 600 }} cols={2}>
                        <ImageListItem cols={2} rows={2}>
                            <img alt="training img" src={props.MainImgUrl} />
                        </ImageListItem>
                        {props.SecondaryImgUrls && props.SecondaryImgUrls.map(d => (
                            <ImageListItem key={d}>
                                <img alt="training img" src={d} />
                            </ImageListItem>
                        ))}
                    </ImageList>
                </Backdrop>
            }
        </>
    )
}