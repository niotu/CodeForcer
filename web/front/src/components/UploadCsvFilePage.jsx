import React, {useState} from 'react';
import {useNavigate, useParams} from "react-router-dom";
import './styles.css';
import logout, {show404page} from "./globalFunctions.jsx";

const UploadCsvFilePage = () => {
    const {groupCode, contestId} = useParams();
    const navigate = useNavigate();
    const [comment, setComment] = useState('It is not required step, you can skip this page. Press "Submit"');
    const [isCorrect, setIsCorrect] = useState(true);
    const [csvFile, setCsvFile] = useState(null); // State for the CSV file

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

        try {
            const response = await fetch(`/api/uploadUsers`, {
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
        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Set up the handle-email mapping</h1>
                    </div>
                    <div className="right-part">
                        <form onSubmit={fileSubmit} autoComplete='on'>
                            <label htmlFor="csvFile">Choose CSV file:</label>
                            <input
                                type="file"
                                id="csvFile"
                                accept=".csv"
                                onChange={handleCsvChange}
                            /><br/><br/>

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

export default UploadCsvFilePage;