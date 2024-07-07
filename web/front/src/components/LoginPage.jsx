import React, {useState} from 'react';
import {useNavigate} from 'react-router-dom';
import Cookies from 'js-cookie';
import './styles.css'; // Import the provided CSS file
import logout from './globalFunctions.js'

const LoginPage = () => {
    const [isCorrect, setIsCorrect] = useState(false);
    const [key, setKey] = useState(Cookies.get('userKey') || '');
    const [secret, setSecret] = useState(Cookies.get('userSecret') || '');
    const [comment, setComment] = useState('');
    const navigate = useNavigate();
    const [id, setId] = useState('');
    const [status, setStatus] = useState('')

    console.log(`UserKey: ${key}, secret: ${secret}`)
    const handleSubmit = async (e) => {

        e.preventDefault();

        const queryParams = new URLSearchParams({
            key,
            secret,
        });

        const response = await fetch(`/api/setAdmin?${queryParams}`, {
            method: 'GET'
        });
        let resp_json = await response.json();
        console.log(resp_json);
        setStatus(resp_json.status);
        if (resp_json.status === 'OK') {
            setId(resp_json.id);
            console.log(`** id is ${id}`)
            Cookies.set('userKey', key);
            // console.log(`key: ${key}, key from cookies: ${Cookies.get('userKey')}`)

            Cookies.set('userSecret', secret);

            // console.log(`secret: ${secret}, secret from cookies: ${Cookies.get('userSecret')}`)
            localStorage.setItem('isAuthorized', 'true'); // Store the authorization status in local storage
            localStorage.setItem('userId', resp_json.id); // Store the user ID in local storage
            console.log(`** is user auth ${localStorage.getItem('isAuthorized')}`)

            console.log(`** userId is ${localStorage.getItem('userId')}`)
            navigate('/link');
        } else if (resp_json.status === 'FAILED') {
            setComment(resp_json.comment)
            alert(resp_json.comment);
        }

    };

    return (
        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Login to CodeForces</h1>
                    </div>
                    <div className="right-part">
                        <form onSubmit={handleSubmit} autoComplete='on'>
                            <label htmlFor="key">Key:</label>
                            <input type="password"
                                   id="key"
                                   value={key}
                                   onChange={(e) => setKey(e.target.value)}
                                   required className={isCorrect ? 'correct' : 'incorrect'}/><br/><br/>

                            <label htmlFor="secret">Secret:</label>
                            <input type="password"
                                   id="secret"
                                   value={secret}
                                   onChange={(e) => setSecret(e.target.value)}
                                   required className={isCorrect ? 'correct' : 'incorrect'}/><br/><br/>
                            <button type="submit">Submit</button>
                        </form>
                    </div>
                </div>
            </div>
            <div className="navigation">
                <div className="left-navigation-part">

                </div>
                <p>{comment}</p>
                <div className="right-navigation-part">
                    <a href="/">
                        <button className={'logout'} onClick={() => logout()}>Logout
                        </button>
                    </a>
                </div>
            </div>
        </div>
    );
};

export default LoginPage;
