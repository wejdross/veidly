import {
	Button, Container, Dialog,
	DialogActions,
	DialogContent, DialogTitle, Fab, Grid,
	ImageList,
	ImageListItem,
	TextField,
	Typography,
} from '@mui/material';
import { Add, Delete, Edit } from '@mui/icons-material';
import React, { useEffect, useState } from 'react';
import { locale2 } from '../locale';
import { MulwiColors } from '../mulwiColors';
import AvatarEdit from '../profile/avatarEdit';
import { AvatarContainer } from '../profile/profile';
import {
	deleteProfileImg, getInstructor, PATCHInstructor,
	postProfileImg
} from '../apicalls/instructor.api';
import StickyFooter from '../Footer';
import { errToStr, getErrorDialog, getNullDialog, StatusDialog } from '../StatusDialog';
import { Harmonogram } from '../harmonogram/harmonogram';
import AboutMeEdit from '../profile/AboutMeEdit';
import ModalEdit from '../profile/ModalEdit';
import { useHistory, useLocation } from 'react-router-dom/cjs/react-router-dom.min';
import { InstructorShedule } from '../reservations/InstructorShedule';
import NameEdit from '../profile/nameEdit';
import { defaultLogoPath } from '../helpers';
import { makeStyles } from "@mui/styles";
import useMediaQuery from '@mui/material/useMediaQuery';
import { useTheme } from "@mui/material/styles";
import { Box } from '@mui/system';

const ImgHeight = "500"
const ImgWidth = "100%"

const btnBg = MulwiColors.blueDark


function ProfileSecondaryImg(props) {
	const theme = useTheme();
	const belowSMSize = useMediaQuery(theme.breakpoints.down("sm"));

	const [open, _setOpen] = useState(false)

	function setOpen(v) {
		if (v)
			setError(null)
		_setOpen(v)
	}

	const [err, setError] = useState(null)

	return (
			<Grid container style={{
				height: "100%",
				maxWidth: belowSMSize ? 300 : "auto"
				}}>
			<Grid item xs={12}>
			{!props.readonly && (
				<React.Fragment>
				<Dialog open={open} onClose={() => setOpen(false)}>
					<DialogTitle>
						{locale2.DELETING_IMG[props.lang]}
					</DialogTitle>
					<DialogContent>
						{err ?
							err :
							<Typography>
								{locale2.ARE_YOU_SURE[props.lang]}
							</Typography>}
					</DialogContent>
					<DialogActions>
						<Button style={{
							color: "white",
							backgroundColor: MulwiColors.blueDark
						}} onClick={async () => {
							try {
								await deleteProfileImg(props.path, 0)
								props.main.refreshInstructor()
								setOpen(false)
							} catch (ex) {
								setError(errToStr(ex))
							}
						}}>
							{locale2.YES[props.lang]}
						</Button>
						<Button onClick={() => {
							setOpen(false)
						}}>
							{locale2.CANCEL[props.lang]}
						</Button>
					</DialogActions>
				</Dialog>
				<center>

				<Fab
					onClick={() => {
						setOpen(true)
					}}
					style={{
						// position: "absolute",
						// right: belowSMSize ? 0 : "40%",
						// top: belowSMSize ? 10 : "90%" /*45*/,
						color: "white",
						backgroundColor: MulwiColors.redError,
						
					}}
					size={"small"}
					>
					<Delete />
				</Fab>
						</center>
			</React.Fragment>)}
			</Grid>
			<Grid item xs={12} style={{height: "100%"}}>

			<img style={{
				// maxWidth: belowSMSize ? "75vw" : "30vw",
				// maxHeight: belowSMSize ? "75vh" : "50vh",
				width: "100%",
				objectFit: "contain",
				height: "100%",
				maxWidth: "100%",
				maxHeight: "100%",
				
			}} src={props.url} alt="secondary logo" />
			</Grid>
			</Grid>
	)
}

function AddProfileSecondaryImg(props) {
	return (<React.Fragment>
		<Grid spacing={1}>
			<Grid item xs={12}>
				<AvatarEdit
					refreshInstr
					putImg={fd => postProfileImg(fd, 0)}
					fab={id => (<label htmlFor={id}>
						<Box textAlign={"center"}>

						<Button aria-label="upload picture"
							component="span" style={{
								color: "white",
								backgroundColor: btnBg
							}}>
							<Add />
						</Button>
								</Box>
						</label>)}
					{...props} />
			</Grid>
		</Grid>
	</React.Fragment>)
}

function ModifyProfileSection(props) {

	const [v, setv] = useState("")
	const [t, sett] = useState("")

	useEffect(() => {
		if (!props.edit)
			return
		setv(props.instructor.ProfileSections[props.index].Content)
		sett(props.instructor.ProfileSections[props.index].Title)
	}, [props.index, props.instructor])

	async function save(isDelete) {
		if (isDelete) {
			props.instructor.ProfileSections.splice(props.index, 1)
		} else if (props.edit) {
			props.instructor.ProfileSections[props.index] = {
				Content: v,
				Title: t
			}
		} else {
			if (!props.instructor.ProfileSections)
				props.instructor.ProfileSections = []
			props.instructor.ProfileSections.push({
				Title: t,
				Content: v
			})
		}
		await PATCHInstructor(props.instructor)
		props.main.refresh()
	}

	return (<React.Fragment>
		<Grid container direction="row">

			<ModalEdit
				customActions={props.edit && (_s => (
					<Button onClick={() => _s(1)}>
						{locale2.DELETE[props.lang]}
					</Button>
				))}
				lang={props.lang}
				buttonProps={{
					variant: "contained",
					style: {
						color: "white",
						backgroundColor: MulwiColors.blueDark
					},
					size: "small"
				}}
				onlyButton={true}
				label={props.edit
					? locale2.EDIT[props.lang] : <Add />}
				title={props.edit ?
					locale2.EDIT_SECTION[props.lang] :
					locale2.ADD_SECTION[props.lang]}
				custom
				onSave={save}
				content={(<React.Fragment>
					<TextField
						style={{
							marginBottom: 10
						}}
						fullWidth
						value={t}
						onChange={(event => (sett(event.target.value)))}
						placeholder={locale2.TITLE[props.lang]}
						variant="outlined"
					/>
					<TextField
						multiline
						placeholder={locale2.CONTENT[props.lang]}
						inputProps={{ maxLength: 250 }}
						rows={5}
						variant={"outlined"}
						helperText={((v && v.length) || 0) + "/250"}
						fullWidth={true}
						value={v}
						onChange={(event => (setv(event.target.value)))}
					/>
				</React.Fragment>)}
			/>

		</Grid>
	</React.Fragment>)
}

function ProfileSection(props) {
	const theme = useTheme();
	const belowSMSize = useMediaQuery(theme.breakpoints.down("sm"));
	let align = props.right ? "right" : "left"

	function buttons() {
		return (<React.Fragment>
			{props.sp === 1 && (
				<AboutMeEdit
					label={locale2.EDIT[props.lang]}
					lang={props.lang} onlyButton
					buttonProps={{
						variant: "contained",
						style: {
							color: "white",
							backgroundColor: MulwiColors.blueDark
						},
						size: "small"
					}}
					main={props.main}
					user={props.user} />
			)}
			{!props.sp && (<React.Fragment>
				<ModifyProfileSection
					{...props}
					edit />
				{/* <Button size="small" style={{
					color: "white",
					backgroundColor: MulwiColors.blueDark
				}}>
					{locale2.EDIT[props.lang]}
				</Button> */}
			</React.Fragment>)}
		</React.Fragment>)
	}

	return (<React.Fragment>
		<div style={{ position: "relative", marginBottom: 20 }}>
			<Grid container direction="row" spacing={2}>
				{!props.readonly && <Grid item>
					{buttons()}
				</Grid>}
				<Grid item>
					<Typography align={align} variant="h6">
						{props.title}
					</Typography>
				</Grid>
			</Grid>
			<br />
			<Typography style={{
				whiteSpace: "pre-wrap",
				overflowWrap: "break-word",
				maxWidth: "100%",
				textAlign: align
			}}>
				{props.content}
			</Typography>
		</div>
	</React.Fragment>)
}

const logoPath = defaultLogoPath

export default function InstrProfile(props) {
	const theme = useTheme();
	const belowSMSize = useMediaQuery(theme.breakpoints.down("sm"));
	const [info, setInfo] = useState(getNullDialog())

	const location = useLocation()
	const h = useHistory()

	const [instructor, setInstructor] = useState(null)
	const [forcedReadonly, setForcedReadonly] = useState(null)

	const [readonly, setReadonly] = useState(false);

	useEffect(() => {
		if (!readonly)
			setInstructor(props.instructor)
	}, [props.instructor, readonly])

	async function setReadonlyInstr(iid) {
		try {
			let i = await getInstructor(iid)
			i = JSON.parse(i)
			setInstructor(i)
			setReadonly(true)
		} catch (ex) {
			setInfo(getErrorDialog(
				locale2.SOMETHING_WENT_WRONG[props.lang],
				ex
			))
		}
		setReadonly(true)
	}

	useEffect(() => {
		if (!location) {
			setReadonly(false)
			setForcedReadonly(false)
			return
		}
		let x = new URLSearchParams(location.search)
		let iid = x.get("instructorID")
		if (!iid) {
			setReadonly(false)
		} else {
			setReadonlyInstr(iid)
		}
		let forced = x.get("f")
		if (forced) {
			setForcedReadonly(true)
		} else {
			setForcedReadonly(false)
		}
	}, [location])

	function enterReadonly() {
		if (!instructor)
			return
		h.push("/instr_profile?instructorID=" + instructor.id)
	}

	function exitReadonly() {
		h.push("/instr_profile")
	}
	const [open, setOpen] = useState(false)

	if (!instructor)
		return null

	let url = window.location.origin + "/instr_profile?f=1&instructorID=" + instructor.id


	return (<React.Fragment>
		<StatusDialog lang={props.lang} info={info} setInfo={setInfo} />

		<div style={{
			position: "relative",
			width: ImgWidth,
			maxHeight: 500,
			backgroundColor: instructor.BgImgUrl ? "none" :  MulwiColors.blueDark
		}}>
			<img style={{
				objectFit: "cover",
				height: "100%",
				width: ImgWidth,
				maxHeight: 500,
				paddingBottom: 50
			}} src={instructor.BgImgUrl || (logoPath)} alt="instructor logo">
			</img>
			{!readonly && (
				<AvatarEdit
					refreshInstr
					putImg={fd => postProfileImg(fd, 1)}
					fab={id => (<label htmlFor={id}>
						<Fab aria-label="upload picture"
							component="span" style={{
								position: "absolute",
								right: 10,
								top: belowSMSize ? 80 : 10,
								color: MulwiColors.blueDark,
								backgroundColor: MulwiColors.whiteBackground
							}}>
							<Edit />
						</Fab></label>)}
					lang={props.lang} main={props.main} />
			)}
		</div>

		<Container style={{
			marginTop: -80,
			position: "relative"
		}}>
			<div style={{
			}}>
				{!readonly && (<Grid container direction='row' justifyContent='center'
					style={{
						position: "absolute",
						left: 65,
						top: 105,
					}}>
					<AvatarEdit refreshInstr
						fab={id => (
							<label htmlFor={id}>
								<Fab aria-label="upload picture" component="span"
									style={{
										width: 40,
										height: 40,
										// zIndex must be used in desperation and carefulness as it can be visible on popups
										zIndex: 10,
										color: MulwiColors.blueDark,
										backgroundColor: MulwiColors.whiteBackground
									}}>
									<Edit />
								</Fab>
							</label>)}
						lang={props.lang} main={props.main} />
				</Grid>)}

				<AvatarContainer large user={instructor.UserInfo} />

			</div>

			<br />
			<center>
				<Typography variant="h4">
					{instructor.UserInfo.Name || "Anonymous"}
					{!readonly && (<React.Fragment>
						<Button variant="contained" size="small"
							onClick={() => setOpen(true)}
							style={{
								marginLeft: 10,
								backgroundColor: MulwiColors.blueDark,
								color: "white"
							}}>{locale2.EDIT[props.lang]}</Button>
						<NameEdit external
							lang={props.lang}
							main={props.main}
							user={instructor.UserInfo}
							setOpen={setOpen}
							open={open} />
					</React.Fragment>)}
				</Typography>
			</center>

			<center>
				<Typography variant="h6">
					{locale2.INSTRUCTOR_SINCE[props.lang]} {(() => {
						let x = new Date(instructor.CreatedOn);
						return x.getFullYear();
					})()}
					{/* {!props.readonly && (<React.Fragment>
					<InstructorEdit onlyButton
						label={locale2.EDIT[props.lang]}
						lang={props.lang}
						main={props.main}
						instructor={instructor} />
					</React.Fragment>)} */}
				</Typography>
			</center>

			<br />

			{!readonly && (
				<React.Fragment>
					<Grid container direction="row"
						spacing={2}
						justifyContent='center' alignContent='space-between'>
						<Grid item>
							<center>
								<div style={{
									border: "1px solid grey",
									borderRadius: 30,
									padding: 30,
								}}>

									<Typography>
										{locale2.YOU_CAN_SHARE_PROFILE_ANYWHERE[props.lang]}
									</Typography>
									<TextField
										readonly multiline
										variant="outlined"
										style={{
											marginTop: 20,
										}}
										value={url} />
								</div>
							</center>

						</Grid>
					</Grid>
					<br />
				</React.Fragment>
			)}

			<br />

			<Grid container direction="row"
				justifyContent="space-between">

				<Grid item xs={12} md={6}>
					{!readonly && <React.Fragment>
						<Typography variant="h6" align='center'>
							{locale2.TEXT_SECTIONS[props.lang]}
						</Typography>
						<br />
					</React.Fragment>}
					{(!readonly || instructor.UserInfo.AboutMe) && <ProfileSection readonly={readonly} title={locale2.ABOUT_ME[props.lang]}
						content={instructor.UserInfo.AboutMe} sp={1}
						user={instructor.UserInfo}
						main={props.main}
						lang={props.lang} />}
					{instructor.ProfileSections
						&& instructor.ProfileSections.map((c, i) => (
							<ProfileSection
								readonly={readonly}
								index={i}
								title={c.Title}
								content={c.Content}
								user={instructor.UserInfo}
								main={props.main}
								lang={props.lang}
								instructor={instructor}
							/>
						))}

					{!readonly && <ModifyProfileSection
						instructor={instructor}
						main={props.main}
						lang={props.lang} />}
				</Grid>
				<Grid item xs={12}>
					{!readonly && (<Typography align='center' variant="h6">
						{locale2.PHOTO_SECTIONS[props.lang]
							// Your photos starts here		
						}
					</Typography>)}
					{!readonly &&
								<AddProfileSecondaryImg
									lang={props.lang} main={props.main} />
						}
					<Grid container
						direction={belowSMSize ? "column" : "row"}
						// justifyContent="center"
						// alignItems="center"
						justifyContent="center"
						alignItems="center" 
						spacing={3}
						style={{marginTop: 10}}
						>
						{ instructor.ExtraImgUrls &&
						<ImageList cols={belowSMSize ? 1 : 3} rowHeight={belowSMSize ? 350 : 500} style={{marginBottom: 15, marginTop: 10, overflowY: "hidden"}}>
							{instructor.ExtraImgUrls
								&& instructor.ExtraImgUrls.map((c, i) => (
									<ImageListItem key={i} style={{
										border: readonly ? "none" : "1px solid grey",
										borderRadius: 10,
										margin: readonly ? 0 : 10,
										marginLeft: 0,
									}}>
										<ProfileSecondaryImg
											readonly={readonly}
											main={props.main}
											lang={props.lang}
											url={c}
											path={instructor.ExtraImgPaths[i]} />
									</ImageListItem>
								))}
						</ImageList>
						}
					</Grid>
				</Grid>
			</Grid>
		</Container>

		<br />

		{!readonly && (<center><Typography variant="h6">
			{locale2.SCHEDULE[props.lang]}
		</Typography></center>)}

		{readonly ? (
			<InstructorShedule instructor={instructor} onlySchedule
				lang={props.lang} user={props.user} />
		) : (
			<Harmonogram nohdr lang={props.lang} instructor={instructor} />
		)}


		{!forcedReadonly && instructor && (<React.Fragment>
			<Button
				onClick={readonly ? exitReadonly : enterReadonly}
				style={{
					position: "fixed",
					right: 10,
					bottom: 10,
					color: "white",
					backgroundColor: MulwiColors.blueDark
				}} variant='contained'>
				{readonly ? locale2.END_PREVIEW[props.lang] : locale2.PREVIEW[props.lang]}
			</Button>
		</React.Fragment>)}
	</React.Fragment>)
}
