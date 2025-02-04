import React, {useEffect, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import Cookies from 'js-cookie';
import './styles/styles.css'; // Import the provided CSS file
import logout from './additional/globalFunctions.jsx'
import logo from '../assets/logo.svg'
import logoutIcon from '../assets/logout.png'
import InfoComponent from "./additional/InfoComponent.jsx";
import data from "./additional/infoDistr.json";

const LoginPage = () => {
    const [isCorrect, setIsCorrect] = useState(true);
    const data = require('./additional/infoDistr.json');
    const infoData = {
        content: data.LoginPage
    };
    const [key, setKey] = useState(Cookies.get('userKey') || '');
    const [secret, setSecret] = useState(Cookies.get('userSecret') || '');
    const [comment, setComment] = useState('');
    const navigate = useNavigate();
    let id;
    let status;
    let isAuth = localStorage.getItem('isAuthorized');

    console.log(`UserKey: ${key}, secret: ${secret}`)

    useEffect(() => {
        if (isAuth) {
            navigate('/link');
        }
    }, []);

    const handleSubmit = async (e) => {

        e.preventDefault();

        const queryParams = new URLSearchParams({
            key,
            secret,
        });

        console.log(process.env.REACT_APP_BACKEND_URL)

        const url =
            process.env.REACT_APP_BACKEND_URL +
            '/api/setAdmin?' + queryParams;

        console.log(url);

        const response = await fetch(url, {
            method: 'GET'
        });

        console.log(response);
        // console.log(await response.text());
        try {
            // Handle non-200 status codes
            if (!response.ok) {
                throw new Error(`Request failed: ${response.ok}`);
            }
            const resp_json = await response.json();
            console.log(resp_json);
            status = resp_json.status;
            if (resp_json.status === 'OK') {
                id = resp_json.id;
                console.log(`** id is ${id}`)
                // Cookies.set('userKey', key);
                // console.log(`key: ${key}, key from cookies: ${Cookies.get('userKey')}`)

                // Cookies.set('userSecret', secret);

                // console.log(`secret: ${secret}, secret from cookies: ${Cookies.get('userSecret')}`)
                localStorage.setItem('isAuthorized', true); // Store the authorization status in local storage
                localStorage.setItem('userId', id); // Store the user ID in local storage
                // console.log(`** is user auth ${localStorage.getItem('isAuthorized')}`)
                // console.log(`** userId is ${localStorage.getItem('userId')}`)

                navigate('/link')
            } else if (resp_json.status === 'FAILED') {
                setComment(resp_json.comment)
                alert(resp_json.comment);
            }
        } catch (error) {
            console.error('Error fetching data: ', error);
            setComment('Error fetching data');
            setIsCorrect(false);
            // Display an error message to the user
        }

    };

    return (
        <div className="content">

            <div className="header">
                <img src={logo} height={50} alt={'logo'}/>
                {isAuth ? (<a href="/" className={isAuth ? 'authorized' : 'non-authorized'}>
                    <button className={'logout'} onClick={() => logout()}>
                        <img src={logoutIcon} height={25}
                             alt='logout icon'/>
                    </button>
                </a>) : (<a></a>)}
            </div>
            <div className='page-active'>
                <div className="wizard">
                    <div className={'filler'}>
                        <InfoComponent infoData={infoData}/>
                    </div>
                    <div className="panel">
                        <div className="left-part">
                            <h1>Login to CodeForces</h1>
                            <p className={isCorrect ? 'correct-comment' : 'incorrect-comment'}>{comment}</p>

                        </div>
                        <div className="right-part">
                            <form onSubmit={handleSubmit} autoComplete='on'>
                                <div className="form__group field">
                                    <input type="password"
                                           className="form__field"
                                           placeholder="key"
                                           name="key"
                                           id='key'
                                           value={key}
                                           required
                                           onChange={(e) => setKey(e.target.value)}

                                    />
                                    <label htmlFor="password" className="form__label">Key</label>
                                </div>
                                <div className="form__group field">
                                    <input type="password"
                                           className="form__field"
                                           placeholder="secret"
                                           name="secret"
                                           id='secret'
                                           required
                                           value={secret}
                                           onChange={(e) => setSecret(e.target.value)}
                                    />
                                    <label htmlFor="password" className="form__label">Secret</label>
                                </div>
                                <button type="submit">Log In</button>
                            </form>
                        </div>
                    </div>
                    <div className={'navigation'}>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default LoginPage;
