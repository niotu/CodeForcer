import React, {useState} from 'react';
import {useNavigate} from 'react-router-dom';
import './styles.css'; // Import the provided CSS file

const LoginPage = () => {
    localStorage.setItem('isAuthorized', false);
    const [handle, setHandle] = useState('');
    const [password, setPassword] = useState('');
    const [key, setKey] = useState('');
    const [secret, setSecret] = useState('');
    const navigate = useNavigate();
    let comment = '';

    const handleSubmit = async (e) => {
        e.preventDefault();

        const queryParams = new URLSearchParams({
            handle,
            password,
            key,
            secret,
        });

        const response = await fetch(`/api/setAdmin?${queryParams}`, {
            method: 'GET'
        });
        let id = 0;
        let status = await ((response.json()).then(r => r))['status'];
        if (status === 'OK') {
            id = await ((response.json()).then(r => r))[id];
            localStorage.setItem('isAuthorized', 'true'); // Store the authorization status in local storage
            localStorage.setItem('userId', id); // Store the user ID in local storage
            navigate('/groups');
        } else if (status === 'FAILED') {
            alert('Login failed');
            comment = await ((response.json()).then(r => r))['comment'];
        }
    };

    const [isCorrect, setIsCorrect] = useState(false)

    function logout() {
        localStorage.setItem('isAuthorized', 'false');
        localStorage.setItem('userId', null);
        navigate('/');
    }

    return (
        <body>
        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Login to CodeForces</h1>
                    </div>
                    <div className="right-part">
                        <form onSubmit={handleSubmit}>
                            <label htmlFor="handle">Handle:</label>
                            <input type="text" id="handle" value={handle} onChange={(e) => setHandle(e.target.value)}
                                   required className={isCorrect ? 'correct' : 'incorrect'}/><br/><br/>

                            <label htmlFor="password">Password:</label>
                            <input type="password" id="password" value={password}
                                   onChange={(e) => setPassword(e.target.value)} required
                                   className={isCorrect ? 'correct' : 'incorrect'}/><br/><br/>

                            <label htmlFor="key">Key:</label>
                            <input type="password" id="key" value={key} onChange={(e) => setKey(e.target.value)}
                                   required className={isCorrect ? 'correct' : 'incorrect'}/><br/><br/>

                            <label htmlFor="secret">Secret:</label>
                            <input type="password" id="secret" value={secret}
                                   onChange={(e) => setSecret(e.target.value)}
                                   required className={isCorrect ? 'correct' : 'incorrect'}/><br/><br/>

                            <button type="submit">Submit</button>
                        </form>
                    </div>
                </div>
            </div>
            <button className={'logout'} onSubmit={(e) => logout()}>Logout</button>
        </div>
        </body>
    );
};

export default LoginPage;
