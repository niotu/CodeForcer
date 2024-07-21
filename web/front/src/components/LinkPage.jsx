import React, {useState} from 'react';
import {useNavigate} from "react-router-dom"; // Import the provided CSS file
import './styles.css';
import './inputs.scss'
import logout, {show404page} from "./globalFunctions.jsx";
import logo from "../assets/logo.svg";
import logoutIcon from "../assets/logout.png";
import InfoComponent from "./InfoComponent.jsx";

const LinkPage = () => {
    // const [key, setKey] = useState('');
    // const [secret, setSecret] = useState('');
    const [link, setUrl] = useState('');
    const data = require('./infoDistr.json');
    const infoData = {
        content: data.LinkPage
    };
    const navigate = useNavigate();
    const [comment, setComment] = useState('');
    const [isCorrect, setIsCorrect] = useState(true)
    const [isAuth, setIsAuth] = useState(localStorage.getItem('isAuthorized') || true);

    console.log(localStorage.getItem('isAuthorized'));

    if (!localStorage.getItem('isAuthorized')) {
        return show404page();
    }

    const linkSubmit = async (e) => {
        e.preventDefault();

        if (link === '') {
            setIsCorrect(false);
            setComment('Link is required');
            return;
        } else {
            setIsCorrect(true);
        }
        console.log('** clicked...')

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

        let url0 = process.env.REACT_APP_BACKEND_URL +
            '/api/getTasks?' + queryParams;
        try {
            const response = await fetch(url0, {
                method: 'GET'
            });
            console.log(response);
            let resp_json = await response.json();
            // console.log(resp_json);
            status = resp_json.status;

            console.log(status);
            if (status === 'OK') {
                navigate(`/weights-distribution/${groupCode}/${contestId}`);
            } else if (status === 'FAILED') {
                setComment(resp_json.comment);
                setIsCorrect(false);
            }
        } catch (e) {
            setComment('Error in connection');
            setIsCorrect(false);
        }
    }

    console.log('** process');


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
            <div className="page-active">

                <div className="wizard">
                    <div className={'filler'}>
                        <InfoComponent infoData={infoData}/>
                    </div>
                    <div className="panel">
                        <div className="left-part">
                            <h1>Enter the link to a contest</h1>
                            <p className={isCorrect ? 'correct-comment' : 'incorrect-comment'}>{comment}</p>
                        </div>
                        <div className="right-part">
                            <form onSubmit={linkSubmit} autoComplete='on'>
                                <div className="form__group field">
                                    <input type="url"
                                           className="form__field"
                                           placeholder="link"
                                           name="link"
                                           id='link'
                                           required
                                           onChange={(e) => setUrl(e.target.value)}

                                    />
                                    <label htmlFor="url" className="form__label">Link</label>
                                </div>
                            </form>
                        </div>

                    </div>
                    <div className="navigation">
                        <div className="left-navigation-part">

                        </div>
                        <div className="right-navigation-part">
                            <a href="">
                                <button className="previous-page" onClick={(e) => {
                                    e.preventDefault();
                                    history.go(-1);
                                }}>Back
                                </button>
                            </a>
                            <a>
                                <button type="submit" onClick={linkSubmit}>Next</button>
                            </a>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    );
};

export default LinkPage;
