import React, {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import './styles.css';
import Cookies from "js-cookie"; // Import the provided CSS file

const ContestDetails = () => {
    let comment = 'Congratulations!'
    const {groupCode, contestId} = useParams(); // Extracting groupCode and contestId from URL parameters
    const [googleSheetLink, setGoogleSheetLink] = useState('');
    const [csvData, setCsvData] = useState('');
    const [submissionsData, setSubmissionsData] = useState('');
    const [loading, setLoading] = useState(true); // Add a loading state
    const [taskWeights, setWeignts] = sessionStorage.getItem('weights')
    const [mode, setMode] = sessionStorage.getItem('mode')

    useEffect(() => {
        const fetchContestDetails = async () => {
            try {
                const queryParams = new URLSearchParams({
                    groupCode: groupCode,
                    contestId: contestId,
                    userID: localStorage.getItem('userId'),
                    weights: taskWeights.replace(',', '-'),
                    mode: mode
                });

                const response = await fetch(`/api/proceed?${queryParams}`);
                const data = await response.json();

                console.log(data);

                setGoogleSheetLink(data.googleSheets);
                setCsvData(data.csv);
                setLoading(false); // Data fetching complete, set loading to false

                if (data.status === 'OK') {
                    // setSubmissionsData(data.Submissions);
                } else if (data.status === 'FAILED') {
                    comment = data.comment;
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
        return <div>Loading...</div>; // Render a loading indicator while data is being fetched
    }

    history.goBack = () => {
        history.go(-1)
    };
    return (
        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Contest Details</h1>
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
                        <button className={'logout'} onClick={() => {
                            localStorage.clear();
                            sessionStorage.clear();
                            Cookies.clear()
                        }}>Logout
                        </button>
                    </a>
                </div>
            </div>
        </div>
    );
};

export default ContestDetails;
