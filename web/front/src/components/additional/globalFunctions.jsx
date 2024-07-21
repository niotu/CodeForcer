import Cookies from 'js-cookie'
import React from "react";
import logo from "../../assets/logo.svg";
import logoutIcon from "../../assets/logout.png";

export default function logout() {
    localStorage.clear();
    sessionStorage.clear();
    Cookies.remove('userKey');
    Cookies.remove('userSecret');
    // Cookies.clear();
    // history.clear();
}



export function show404page() {
    return (<div className="content">
            <div className="header">
                <img src={logo} height={50} alt={'logo'}/>
                <a href="/web/front/public">
                    <button className={'logout'} onClick={() => logout()}>
                        <img src={logoutIcon} height={25}
                             alt='logout icon'/>
                    </button>
                </a>
            </div>
            <div className='page-active'>
                <div className='wizard'>
                    <h1>Code 404. Page not found or you are not authorized</h1>
                    <a href='/web/front/public'>Please, return to home page</a>
                </div>
            </div>
        </div>
    );
}
