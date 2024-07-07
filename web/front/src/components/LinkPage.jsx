import React, {useState} from 'react';
import {useNavigate} from "react-router-dom"; // Import the provided CSS file
import './styles.css';
import Cookies from "js-cookie";
import logout from "./globalFunctions.js";

const LinkPage = () => {
    // const [key, setKey] = useState('');
    // const [secret, setSecret] = useState('');
    const [link, setUrl] = useState('');
    const navigate = useNavigate();
    const [comment, setComment] = useState('');

    const linkSubmit = async (e) => {
        e.preventDefault();
        console.log('** processing...')
        let id = 0,
            status,
            groupCode,
            contestId;

        // Extract groupCode and contestId from the link
        console.log(link);
        const url = new URL(link);
        // const params = new URLSearchParams(url.search);
        const components = url.toString().split('/')

        groupCode = components[4]
        contestId = components[6]

        console.log(groupCode);
        console.log(contestId);


        // Make a request to the backend to proceed with the login
        // Replace '/api/proceed' with the actual API endpoint for the login
        // Make sure to pass the necessary parameters and handle the response appropriately

        const queryParams = new URLSearchParams({
            groupCode: groupCode,
            contestId: contestId,
            userID: localStorage.getItem('userId')
        });

        const response = await fetch(`/api/getTasks?${queryParams}`, {
            method: 'GET',
            mode: 'no-cors'
        });

        let resp_json = await ((response.json()).then(r => r));
        console.log(resp_json);
        status = resp_json.status;
        console.log(status);

        if (status === 'OK') {
            navigate(`/weights-distribution/${groupCode}/${contestId}`);
        } else if (status === 'FAILED') {
            alert(comment);
            setComment(resp_json.comment);
        }
    }

    console.log('** process');

    const [isCorrect, setIsCorrect] = useState(true)

    function logout() {
        localStorage.setItem('isAuthorized', 'false');
        localStorage.setItem('userId', null);
    }

    return (
        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Enter the link to a contest</h1>
                    </div>
                    <div className="right-part">
                        <form onSubmit={linkSubmit} autoComplete='on'>

                            <label htmlFor="link">Paste the link:</label>
                            <input type="url" id="link" value={link}
                                   onChange={(e) => setUrl(e.target.value)}
                                   required className={isCorrect ? 'correct' : 'incorrect'}/><br/><br/>
                            <button type="submit">Submit</button>
                        </form>
                    </div>
                </div>
            </div>
            <div className="navigation">
                <div className="left-navigation-part">
                    <a href="">
                        <button className="previous-page" onClick={(e) => {
                            e.preventDefault();
                            history.go(-1);
                        }}>Back
                        </button>
                    </a>
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

export default LinkPage;
