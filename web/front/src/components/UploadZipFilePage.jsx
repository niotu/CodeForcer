import React, {useState} from 'react';
import {useNavigate, useParams} from "react-router-dom";
import './styles.css';
import logout, {show404page} from "./globalFunctions.jsx";
import localForage from "localforage";
import logo from "../assets/logo.svg";
import logoutIcon from "../assets/logout.png";
import InfoComponent from "./InfoComponent.jsx";
import data from "./infoDistr.json";

const UploadZipFilePage = () => {
    const {group, contestId} = useParams();
    const data = require('./infoDistr.json');
    const infoData = {
        content: data.UploadZipFilePage
    };
    const navigate = useNavigate();
    const [comment, setComment] = useState('It is a not required step, you can skip it. Just click "Submit"');
    const [isCorrect, setIsCorrect] = useState(true);
    const [zipFile, setZipFile] = useState(null);
    const [isAuth, setIsAuth] = useState(localStorage.getItem('isAuthorized') || true) // State for the ZIP file


    console.log(`group : ${group}`);
    console.log(`contest : ${contestId}`);
    console.log(useParams());

    if (!localStorage.getItem('isAuthorized')) {
        return show404page();
    }

    const fileSubmit = async (e) => {
        e.preventDefault();
        console.log(' processing...')

        // Convert the date to a string in ISO format for sending to the API
        const formData = new FormData();
        formData.append('zipFile', zipFile);
        await localForage.setItem('zipFile', zipFile);

        console.log(`/contest-details/${group}/${contestId}`);

        navigate(`/contest-details/${group}/${contestId}`);

    };

    const handleZipChange = (e) => {
        setZipFile(e.target.files[0]);
        console.log('* changed file ', e.target.files[0]);
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
            <div className="page-active">
                <div className="wizard">
                    <div className="filler">
                        <InfoComponent infoData={infoData}/>
                    </div>
                    <div className="panel">
                        <div className="left-part">
                            <h1>Upload Submissions</h1>
                            <p className={isCorrect ? 'correct-comment' : 'incorrect-comment'}>{comment}</p>
                        </div>
                        <div className="right-part">
                            <form onSubmit={fileSubmit} autoComplete='on'>
                                <label htmlFor="zipFile">Choose ZIP file:</label>
                                <input
                                    type="file"
                                    id="zipFile"
                                    accept=".zip"
                                    onChange={handleZipChange}
                                />
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
                                <button type="submit" onClick={fileSubmit}>Next</button>
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default UploadZipFilePage;

