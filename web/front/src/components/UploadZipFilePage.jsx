import React, {useState} from 'react';
import {useNavigate, useParams} from "react-router-dom";
import './styles.css';
import logout, {show404page} from "./globalFunctions.jsx";
import Cookies from "js-cookie";

const UploadZipFilePage = () => {
    const {groupCode, contestId} = useParams();
    const navigate = useNavigate();
    const [comment, setComment] = useState('It is not required step, you can skip this page. Press "Submit"');
    const [isCorrect, setIsCorrect] = useState(true);
    const [zipFile, setZipFile] = useState(null); // State for the ZIP file

    if (!localStorage.getItem('isAuthorized')) {
        return show404page();
    }

    const fileSubmit = async (e) => {
        e.preventDefault();
        console.log(' processing...')

        if (!zipFile) {
            navigate(`/contest-details/${groupCode}/${contestId}`);
            return;
        }

        // Convert the date to a string in ISO format for sending to the API
        const formData = new FormData();
        formData.append('zipFile', zipFile);
        Cookies.set('zip-file', zipFile);

        navigate(`/contest-details/${groupCode}/${contestId}`);

    };

    const handleZipChange = (e) => {
        setZipFile(e.target.files[0]);
    };

    return (
        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Upload Submissions</h1>
                    </div>
                    <div className="right-part">
                        <form onSubmit={fileSubmit} autoComplete='on'>
                            <label htmlFor="zipFile">Choose ZIP file:</label>
                            <input
                                type="file"
                                id="zipFile"
                                accept=".zip"
                                onChange={handleZipChange}
                                required
                            />

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

export default UploadZipFilePage;

