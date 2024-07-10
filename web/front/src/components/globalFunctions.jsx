import Cookies from 'js-cookie'
import React from "react";

export default function logout() {
    localStorage.clear();
    sessionStorage.clear();
    Cookies.remove('userKey');
    Cookies.remove('userSecret')
    // history.clear();
}

export function show404page() {
    return (<div className='page-active'>
        <div className='wizard'>
            <h1>Code 404. Page not found or you are not authorized</h1>
            <a href='/'>Please, return to home page</a>
        </div>
    </div>);
}
