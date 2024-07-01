import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import './styles.css'; // Import the provided CSS file

const LoginPage = () => {
    const [handle, setHandle] = useState('');
    const [password, setPassword] = useState('');
    const [key, setKey] = useState('');
    const [secret, setSecret] = useState('');
    const navigate = useNavigate();

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

        if (response.ok) {
            navigate('/groups');
        } else {
            alert('Login failed');
        }
    };

    return (
        <div className="wizard">
            <div className="panel">
                <div className="left-part">
                    <h1>Login to CodeForces</h1>
                </div>
                <div className="right-part">
                    <form onSubmit={handleSubmit}>
                        <label htmlFor="handle">Handle:</label>
                        <input type="text" id="handle" value={handle} onChange={(e) => setHandle(e.target.value)} required /><br /><br />

                        <label htmlFor="password">Password:</label>
                        <input type="password" id="password" value={password} onChange={(e) => setPassword(e.target.value)} required /><br /><br />

                        <label htmlFor="key">Key:</label>
                        <input type="password" id="key" value={key} onChange={(e) => setKey(e.target.value)} required /><br /><br />

                        <label htmlFor="secret">Secret:</label>
                        <input type="password" id="secret" value={secret} onChange={(e) => setSecret(e.target.value)} required /><br /><br />

                        <button type="submit">Submit</button>
                    </form>
                </div>
            </div>
        </div>
    );
};

export default LoginPage;
