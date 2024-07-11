import React, {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import './styles.css';
import logout, {show404page} from "./globalFunctions.jsx";


const ContestDetails = () => {
    const [comment, setComment] = useState('Congratulations!');

    const {groupCode, contestId} = useParams(); // Extracting groupCode and contestId from URL parameters

    const [googleSheetLink, setGoogleSheetLink] = useState('');

    const [csvData, setCsvData] = useState('');
    const [submissionsData, setSubmissionsData] = useState('');
    const [loading, setLoading] = useState(true); // Add a loading state

    if (!localStorage.getItem('isAuthorized')) {
        return show404page();
    }

    const [taskWeights, setTaskWeights] = useState(sessionStorage.getItem('weights').replaceAll(',', '-') || '')
    const [userId, setUserId] = useState(localStorage.getItem('userId') || '');
    const [late, setLate] = useState(sessionStorage.getItem('lateHours'));
    const [penalty, setPenalty] = useState(sessionStorage.getItem('penalty') || '');
    const [mode, setMode] = useState(sessionStorage.getItem('mode') || '');

    console.log(`groupCode ${groupCode}\n
    contestId ${contestId}\n
    userId ${userId}\n
    taskWeights ${taskWeights}\n
    late ${late}\n
    penalty ${penalty}\n
    mode ${mode}`);


    useEffect(() => {
        const fetchContestDetails = async () => {
            try {
                const queryParams = new URLSearchParams({
                    groupCode: groupCode,
                    contestId: contestId,
                    userID: userId,
                    weights: taskWeights.replaceAll(',', '-'),
                    late: late,
                    penalty: penalty,
                    mode: mode
                });

                const formData = new FormData();
                const fileInput = document.querySelector('input[type="file"]');
                formData.append('file', fileInput.files[0]);

                fetch(`/api/proceed/${queryParams}`, {
                    method: 'POST',
                    body: formData,
                })
                    .then(response => response.text())
                    .then(data => console.log(data))
                    .catch(error => console.error('Error:', error));
                const data = await response.json();

                console.log(data);

                setGoogleSheetLink(data.googleSheets);
                setCsvData(data.csv);
                setLoading(false); // Data fetching complete, set loading to false

                if (data.status === 'OK') {
                    // setSubmissionsData(data.Submissions);
                } else if (data.status === 'FAILED') {
                    setComment(data.comment);
                    alert(data.comment);
                }
            } catch (error) {
                console.error('Error fetching contest details:', error);
                setLoading(false); // Set loading to false even if there's an error
            }
        };

        fetchContestDetails();
    }, [groupCode, contestId]);

    const downloadCsv = () => {
        const blob = new Blob([atob(csvData)], {type: 'text/csv'});
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'data.csv';
        a.click();
        window.URL.revokeObjectURL(url);
    };

    const downloadSubmissions = () => {
        const blob = new Blob([atob(submissionsData)], {type: 'application/zip'});
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'submissions.zip';
        a.click();
        window.URL.revokeObjectURL(url);
    };

    if (loading) {
        return (<div className="page-active">
            <div className="wizard">
                <div className="loading-spinner">
                    <h1>Loading contest details...</h1>
                    <img src={"/web/front/assets/loading.gif"} width={200} height={200} alt='loading'/>
                </div>
            </div>
        </div>); // Render a loading indicator while data is being fetched
    }

    return (
        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Contest Details</h1>
                        <h4>
                            weights: {taskWeights},
                            late: {late},
                            penalty: {penalty},
                            mode: {mode}
                        </h4>
                    </div>
                    <div className="right-part">
                        <div>
                            <label>Google Sheet: </label>
                            <a href={googleSheetLink} target="_blank" rel="noopener noreferrer">See Google Sheet</a>
                        </div>
                        <div>
                            <button onClick={downloadCsv}>Download CSV</button>
                        </div>
                        <div>
                            <button onClick={downloadSubmissions}>Download Submissions</button>
                        </div>
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

export default ContestDetails;
