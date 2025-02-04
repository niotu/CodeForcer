import React, {useState} from 'react';
import {useNavigate, useParams} from "react-router-dom";
import './styles/styles.css';
import logout, {show404page} from "./additional/globalFunctions.jsx";
import logo from "../assets/logo.svg";
import logoutIcon from "../assets/logout.png";
import InfoComponent from "./additional/InfoComponent.jsx";

const UploadCsvFilePage = () => {
    const {groupCode, contestId} = useParams();
    const data = require('./additional/infoDistr.json');
    const infoData = {
        content: data.UploadCsvFilePage
    };
    const navigate = useNavigate();
    const [comment, setComment] = useState('It is not required step, you can skip this page. Press "Submit"');
    const [isCorrect, setIsCorrect] = useState(true);
    const [csvFile, setCsvFile] = useState(null);
    const [isAuth, setIsAuth] = useState(localStorage.getItem('isAuthorized') || true) // State for the CSV file

    if (!localStorage.getItem('isAuthorized')) {
        return show404page();
    }

    const fileSubmit = async (e) => {
        e.preventDefault();
        console.log(' processing...')

        if (!csvFile) {
            navigate(`/upload-zip/${groupCode}/${contestId}`);
            return;
        }

        // Convert the date to a string in ISO format for sending to the API
        const formData = new FormData();
        formData.append('file', csvFile);

        let url = process.env.REACT_APP_BACKEND_URL +
            '/api/uploadUsers'

        try {
            const response = await fetch(url, {
                method: 'POST',
                body: formData,
            });

            const data = await response.json();

            if (data.status === 'OK') {
                navigate(`/upload-zip/${groupCode}/${contestId}`);
            } else if (data.status === 'FAILED') {
                setComment(data.comment);
                alert(data.comment);
            }
        } catch (error) {
            console.error('Error submitting late submission:', error);
            alert('An error occurred. Please try again later.');
        }
    };

    const handleCsvChange = (e) => {
        setCsvFile(e.target.files[0]);
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
                    <div className={'filler'}>

                        <InfoComponent infoData={infoData}/>
                    </div>
                    <div className="panel">
                        <div className="left-part">
                            <h1>Set up the handle-email mapping</h1>
                            <p className={isCorrect ? 'correct-comment' : 'incorrect-comment'}>{comment}</p>
                        </div>
                        <div className="right-part">
                            <form onSubmit={fileSubmit} autoComplete='on'>
                                <input type="file" id="file" accept=".csv" onChange={handleCsvChange}/>
                                <label htmlFor="file" className="upload">Upload CSV</label>
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

export default UploadCsvFilePage;