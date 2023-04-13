import { Wrapper } from '@googlemaps/react-wrapper';
import { Backdrop, CircularProgress, Grid, Typography } from '@mui/material';
import React from 'react';
import CookieConsent from "react-cookie-consent";
import { Route, Switch, withRouter } from 'react-router-dom';
import { apiValidateToken, removeChatToken, setChatToken } from './apicalls/chat';
import { getInstructor } from './apicalls/instructor.api';
import { getUserInfo } from './apicalls/user.api';
import ForgotPassword from './auth/ForgotPassword';
import Login from './auth/login';
import OauthFinish from './auth/OauthFinish';
import Register from './auth/register';
import RegisterFinish from './auth/RegisterFinish';
import ResetPassword from './auth/ResetPassword';
import { BecomeTrainer } from './becomeTrainer';
import { ChatIndex } from './chat';
import { ChatNotifications } from './chat/chatNotification';
import { endNotify, startNotify } from './chat/commons';
import { MiniChatWindow } from './chat/miniChatWindow';
import { G_API_KEY } from './conf';
import Configure from './configure/configure';
import { DcCard } from './dc/dcCard';
import { GroupCard } from './group/GroupCard';
import { Harmonogram } from './harmonogram/harmonogram';
import { gettoken, isUuid, rmtoken } from './helpers';
import InstructorBenefits from './InstructorBenefits';
import InstrProfile from './instr_profile/InstrProfile';
import { getSupportedLanguage, lsKey } from './locale';
import { MulwiColors } from './mulwiColors';
import Navbar from './navbar/Navbar';
import InstructorPayments from './profile/InstructorPayments';
import Invoicing from './profile/Invoicing';
import ManageCard from './profile/manageCard';
import Profile from './profile/profile';
import { QrEval } from './qr/qrEval';
import { InstructorShedule } from './reservations/InstructorShedule';
import { Rsv } from './reservations/rsv';
import { RsvDetails } from './reservations/RsvDetails';
import { SubPurch } from './reservations/subPurch';
import SearchResults from './search/results';
import { SmCard } from './sub/smCard';
import { SubDetails } from './sub/subDetails';
import { Subs } from './sub/Subs';
import SupportPage from './Support';
import UserBenefits from './UserBenefits';
import InstrVacationCard from './vacation/InstrVacationCard';
import Welcome from './welcome';
import StickyFooter from './Footer';
class Main extends React.Component {

  state = {
    userInfo: null,
    instructor: null,
    lang: getSupportedLanguage(),
    loading: false,
    chatToken: null,
    nots: null
  }

  setLang(v) {
    this.setState({ lang: v })
    window.location.reload()
  }


  async resetChatToken(force, persist) {
    if(!force && await this.validateToken())
      return

      try {
          await setChatToken(persist)
          this.setState({chatToken: true})
          console.log("set chat token")
      } catch (ex) {
          console.log(ex)
      }
  }

  logout() {
    rmtoken()
    this.setState({ userInfo: null })
    // im completely reloading application to fix positioning of certain elements
    window.location.href = "/"
    removeChatToken()
    //this.props.history.push("/")
  }

  async refreshUser() {
    if (!gettoken()) {
      await this.resetChatToken(true, false)
      this.setState({ userInfo: null })
      return
    }
    try {
      let u = JSON.parse(await getUserInfo())
      if (!localStorage.getItem(lsKey)) {
        localStorage.setItem(lsKey, u.Language)
        this.setState({ lang: u.Language })
      }
      this.setState({ userInfo: u })
      await this.resetChatToken(true, true)
    } catch (ex) {
      await this.resetChatToken(true, false)
      if (ex == 401) {
        // 
        this.logout()
      } else {
        rmtoken()
        this.setState({ userInfo: null })
        console.error(ex)
      }
    }
  }

  async refreshInstructor() {
    if (!gettoken()) {
      this.setState({ instructor: null })
      return
    }
    try {
      this.setState({ instructor: JSON.parse(await getInstructor()) })
    } catch (ex) {
      this.setState({ instructor: null })
    }
  }

  async validateToken() {
    try {
      await apiValidateToken()
      this.setState({chatToken: true})
      console.log("chat token is still valid")
      return true
    } catch(ex) {
      return false
    }
  }

  async refresh() {
    await this.refreshUser()
    await this.refreshInstructor()
    // add here more deps
  }

  loadImg(picture) {
    return new Promise((res, rej) => {
      const img = new Image()
      img.onload = res
      img.onerror = rej
      img.src = picture
    })
  }

  async componentDidMount() {
    try {
      this.setState({ loading: true })
      await this.refresh()
    } finally {
      this.setState({ loading: false })
    }
  }

  render() {
    return (
      <div style={{ overflowX: "hidden" }} onClick={() => {
        endNotify()
      }}>
        <Wrapper apiKey={G_API_KEY} 
                 libraries={["places", "geometry"]} 
                 language={this.state.lang}>
          <Backdrop
            style={{
              zIndex: 9999999,
              // backgroundColor: MulwiColors.blueLight,
              // opacity: 0.6
            }}
            open={this.state.loading}>
            <Grid container direction="column" alignItems="center">
              <CircularProgress style={{
                color: MulwiColors.blueDark
              }} />
              <Typography paragraph variant="h4" style={{
                color: "white"
              }}>
                Loading...
              </Typography>
            </Grid>
          </Backdrop>

          <ChatNotifications
                onNotification={n => {
                  for(let k in n) {
                    if(isUuid(k) && n[k]) {
                      startNotify()
                      break
                    }
                  }
                  this.setState({nots: n})
                }}
                chatToken={this.state.chatToken} 
                lang={this.state.lang} 
             />
          <MiniChatWindow
              chatToken={this.state.chatToken} 
              lang={this.state.lang} 
            />

          <Switch>
            <Route path="/search"  >
              {/* <SearchResults query={this.state.userQuery} liftQuery={this.setUserQuery}/> */}
              <SearchResults
                main={this}
                lang={this.state.lang} setLang={this.setLang.bind(this)}
                user={this.state.userInfo}
                instructor={this.state.instructor} />
            </Route>
            <Navbar main={this}
              nots={this.state.nots}
              user={this.state.userInfo}
              lang={this.state.lang} setLang={this.setLang.bind(this)}
              instructor={this.state.instructor} >
              <Switch>
                {/* <Route path="/dc">
                  <DcCard lang={this.state.lang} />
                </Route> */}
                <Route path="/qr/eval">
                  <QrEval lang={this.state.lang} />
                </Route>
                <Route path="/login">
                  <Login
                    lang={this.state.lang} main={this} />
                </Route>
                <Route path="/forgot_password">
                  <ForgotPassword
                    lang={this.state.lang} />
                </Route>
                <Route path="/reset_password">
                  <ResetPassword
                    lang={this.state.lang} />
                </Route>
                <Route path="/oauth/finish">
                  <OauthFinish
                    lang={this.state.lang} main={this} />
                </Route>
                <Route path="/register/finish">
                  <RegisterFinish
                    lang={this.state.lang} />
                </Route>
                <Route path="/register">
                  <Register
                    lang={this.state.lang} />
                </Route>
                {/* <Route path="/invoice">
                  <Invoicing 
                    main={this}
                    user={this.state.userInfo}
                    instructor={this.state.instructor}
                    lang={this.state.lang} />
                </Route> */}
                <Route path="/trainings">
                  <Harmonogram usrRsv
                    lang={this.state.lang}
                    instructor={this.state.instructor} basePath="/trainings" />
                </Route>
                <Route path="/configure">
                  <Configure
                    lang={this.state.lang}
                    main={this}
                    user={this.state.userInfo}
                    instructor={this.state.instructor} />
                </Route>
                <Route path="/become_trainer">
                  <BecomeTrainer lang={this.state.lang} main={this} />
                </Route>
                <Route path="/profile">
                  <Profile
                    instructor={this.state.instructor}
                    main={this} lang={this.state.lang}
                    user={this.state.userInfo} />
                </Route>
                <Route path="/vacations">
                  <InstrVacationCard lang={this.state.lang} />
                </Route>
                <Route path="/manage">
                  <ManageCard instructor={this.state.instructor}
                    main={this} lang={this.state.lang}
                    user={this.state.userInfo} />
                </Route>
                <Route path="/harmonogram*">
                  <Harmonogram lang={this.state.lang} instructor={this.state.instructor} />
                </Route>
                <Route path="/rsv_details">
                  <RsvDetails lang={this.state.lang}
                    instructor={this.state.instructor} />
                </Route>
                <Route path="/sub_details">
                  <SubDetails lang={this.state.lang} user={this.state.userInfo} />
                </Route>
                <Route path="/rsv">
                  <Rsv lang={this.state.lang} user={this.state.userInfo} />
                </Route>
                <Route path="/sub_purch">
                  <SubPurch lang={this.state.lang} user={this.state.userInfo} />
                </Route>
                <Route path="/instr/sched">
                  <InstructorShedule lang={this.state.lang} user={this.state.userInfo} />
                </Route>
                {/* <Route path="/payments">
                  {/* <InstructorPayments main={this} lang={this.state.lang} />
                </Route> 
                */}
                <Route path="/group">
                  <GroupCard lang={this.state.lang} />
                </Route>
                <Route path="/support">
                  <SupportPage lang={this.state.lang} />
                  <StickyFooter lang={this.state.lang}/>
                </Route>
                <Route path="/benefits_user">
                  <UserBenefits lang={this.state.lang} />
                  <StickyFooter lang={this.state.lang}/>
                </Route>
                <Route path="/benefits_instructor">
                  <InstructorBenefits lang={this.state.lang} />
                  <StickyFooter lang={this.state.lang}/>
                </Route>
                <Route path="/sm">
                  <SmCard lang={this.state.lang} />
                </Route>
                <Route path="/sub/instr">
                  <Subs lang={this.state.lang} instr user={this.state.userInfo} />
                </Route>
                <Route path="/sub">
                  <Subs lang={this.state.lang} user={this.state.userInfo} />
                </Route>
                <Route path="/instr_profile">
                  <InstrProfile main={this} lang={this.state.lang}
                    user={this.state.userInfo}
                    instructor={this.state.instructor} />
                </Route>
                <Route path="/chat">
                  <ChatIndex 
                    nots={this.state.nots}
                    chatToken={this.state.chatToken} 
                    lang={this.state.lang} 
                    instructor={this.state.instructor}
                    user={this.state.userInfo} />
                </Route>
                <Route path="/">
                  <Welcome lang={this.state.lang} />
                  <StickyFooter lang={this.state.lang}/>
                </Route>
              </Switch>
            </Navbar>
          <CookieConsent
            location="bottom"
            buttonText="OK!"
            cookieName="veidly_cookie_accept"
            style={{ background: "#2B373B", zIndex: "99999" }}
            buttonStyle={{ backgroundColor: MulwiColors.greenDark, fontSize: "15px", color: "white" }}
            expires={150}
          >
            <center>
              
            <Typography variant='h7' style={{textAlign: "center"}}>
              This website uses cookies to enhance the user experience.{" "}
            </Typography>
            </center>
          </CookieConsent>
          </Switch>
        </Wrapper>
      </div>
    )
  }
}

export default withRouter(Main)
