import React, {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import './styles.css'; // Import the provided CSS file

const ContestDetails = () => {
    const {groupCode, contestId} = useParams(); // Extracting groupCode and contestId from URL parameters
    const [googleSheetLink, setGoogleSheetLink] = useState('');
    const [csvData, setCsvData] = useState('');
    const [submissionsData, setSubmissionsData] = useState('');
    const [loading, setLoading] = useState(true); // Add a loading state

    useEffect(() => {
        const fetchContestDetails = async () => {
            try {
                const queryParams = new URLSearchParams({
                    groupCode,
                    contestId,
                });

                const response = await fetch(`/api/proceed?${queryParams}`);
                const data = await response.json();

                console.log(data);

                setGoogleSheetLink(data.googleSheets);
                setCsvData(data.csv);
                // setSubmissionsData(data.Submissions);
                setLoading(false); // Data fetching complete, set loading to false
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
        </div>
    );
};

export default ContestDetails;
