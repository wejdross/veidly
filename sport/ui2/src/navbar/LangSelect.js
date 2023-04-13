import { MenuItem, Select } from '@mui/material';
import Button from '@mui/material/Button';
import Menu from '@mui/material/Menu';
import React from 'react';
import { lsKey } from '../locale';

const flagCss = {
    marginRight: 5,
    width: 16,
    height: 8,
}

const minstyles = {
    maxWidth: 30,
    width: "30 !important",
    border: "1px solid green"
}


export function GenericLangSelect(props) {
    return (
        <React.Fragment>
            <Select
                onChange={(e) => props.onChange(e.target.value)} value={props.value}>
                <MenuItem value="pl" >
                    <img alt="pl" src="/pl.png" style={flagCss} />
                </MenuItem>
                <MenuItem value="en">
                    <img alt="en" src="/uk.png" style={flagCss} />
                </MenuItem>
            </Select>
        </React.Fragment>)
}
export function LangSelect(props) {
    const [anchorEl, setAnchorEl] = React.useState(null);
    const open = Boolean(anchorEl);
    const handleClick = (event) => {
        setAnchorEl(event.currentTarget);
    };
    const handleClose = (e) => {
        setAnchorEl(null)
        if (e.target.id === "pl" || e.target.id === "en") {
            localStorage.setItem(lsKey, e.target.id)
            props.setLang(e.target.id)
        }
    };
    return (
        <>
            <Button
                size='small'

                id="lang-choice-menu"
                aria-controls={open ? 'basic-menu' : undefined}
                aria-haspopup="true"
                aria-expanded={open ? 'true' : undefined}
                onClick={handleClick}
            >
                {
                    props.lang === "pl" ?
                        <img alt="pl" src="/pl.png" style={flagCss} />
                        :
                        <img alt="en" src="/uk.png" style={flagCss} />
                }
            </Button>
            <Menu

                id="basic-menu"
                anchorEl={anchorEl}
                open={open}
                onClose={handleClose}
                MenuListProps={{
                    'aria-labelledby': 'lang-choice-menu',
                }}
            >
                <MenuItem onClick={handleClose} id="pl">
                    <img alt="pl" src="/pl.png" style={flagCss} id="pl" />
                </MenuItem>
                <MenuItem onClick={handleClose} id="en">
                    <img alt="en" src="/uk.png" style={flagCss} id="en" />
                </MenuItem>
            </Menu>
        </>
    )
}