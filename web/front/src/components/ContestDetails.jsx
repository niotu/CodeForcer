import React, {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import './styles.css';
import logout, {show404page} from "./globalFunctions.jsx";
import localForage from "localforage";

function getBoundary(contentType) {
    const boundaryPrefix = 'boundary=';
    const start = contentType.indexOf(boundaryPrefix) + boundaryPrefix.length;
    let end = contentType.indexOf(';', start);
    if (end === -1) end = contentType.length;
    return contentType.substring(start, end).trim();
}

async function parseMultipart(blob, boundary) {
    const text = await blob.text();
    const parts = [];
    const delimiter = `--${boundary}`;
    const closeDelimiter = `--${boundary}--`;
    const splitParts = text.split(delimiter);

    for (let part of splitParts) {
        if (part === '' || part === closeDelimiter || part === '--') continue;

        const headersEndIndex = part.indexOf('\r\n\r\n');
        const headersText = part.slice(0, headersEndIndex);
        const bodyText = part.slice(headersEndIndex + 4);

        const headers = {};
        headersText.split('\r\n').forEach(header => {
            const [key, value] = header.split(': ');
            headers[key] = value;
        });

        parts.push({
            headers,
            body: bodyText,
        });
    }

    return parts;
}


const ContestDetails = () => {
    const [comment, setComment] = useState('Congratulations!');

    const {groupCode, contestId} = useParams(); // Extracting groupCode and contestId from URL parameters

    const [googleSheetLink, setGoogleSheetLink] = useState('');

    const [csvData, setCsvData] = useState('');
    const [submissionsData, setSubmissionsData] = useState(null);
    const [loading, setLoading] = useState(true); // Add a loading state
    const [result, setResult] = useState(null);

    if (!localStorage.getItem('isAuthorized')) {
        return show404page();
    }

    const [taskWeights, setTaskWeights] = useState(sessionStorage.getItem('weights').replaceAll(',', '-') || '')
    const [userId, setUserId] = useState(localStorage.getItem('userId') || '');
    const [late, setLate] = useState(sessionStorage.getItem('lateHours'));
    const [penalty, setPenalty] = useState(sessionStorage.getItem('penalty') || '');
    const [mode, setMode] = useState(sessionStorage.getItem('mode') || '');
    //
    // console.log(`
    //     groupCode ${groupCode}
    //     contestId ${contestId}
    //     userId ${userId}
    //     taskWeights ${taskWeights}
    //     late ${late}
    //     penalty ${penalty}
    //     mode ${mode}`
    // );
    //

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
                const fileInput = await localForage.getItem('zipFile');
                formData.append('file', fileInput);
                console.log(fileInput);

                let url = process.env.REACT_APP_BACKEND_URL +
                    '/api/proceed?' + queryParams;

                const response = await fetch(url, {
                    method: 'POST',
                    body: formData,
                })
                // .then(response => response.text())
                // .then(data => {
                //     console.log(data);
                //     setResult(data.result)
                // })
                // .catch(error => console.error('Error:', error));
                // console.log(await response.text());
                const responseBlob = await response.blob();

                // Parse the multipart response
                const boundary = getBoundary(response.headers.get('Content-Type'));
                const parts = await parseMultipart(responseBlob, boundary);

                console.log(parts);

                // Process the parts
                let jsonPart = null;
                let zipFilePart = null;

                parts.forEach(part => {
                    const contentType = part.headers['Content-Type'];
                    const contentDisposition = part.headers['Content-Disposition'];
                    if (contentType) {
                        if (contentType.includes('application/json')) {
                            jsonPart = JSON.parse(part.body);
                        } else if (contentType.includes('application/zip')) {
                            zipFilePart = new Blob([part.body], {type: 'application/zip'});
                        }
                    }
                });

                if (jsonPart && zipFilePart) {
                    console.log('JSON part:', jsonPart);
                    if (jsonPart.status === 'OK') {
                        const result = jsonPart.result;
                        setGoogleSheetLink(result.googleSheets);
                        setCsvData(result.csv);
                        setSubmissionsData(zipFilePart);

                        setLoading(false);
                    } else if (jsonPart.status === 'FAILED') {
                        setComment(jsonPart.comment);
                        alert(jsonPart.comment);
                    }
                } else {
                    setComment('Some error caught while processing response');
                    alert(comment);
                    return show404page();
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
        const url = window.URL.createObjectURL(submissionsData);
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
                    <img src="/web/front/assets/loading.gif" width={200} height={200} alt='loading'/>
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
