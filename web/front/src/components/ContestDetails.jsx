import React, {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import './styles.css';
import logout, {show404page} from "./globalFunctions.jsx";
import localForage from "localforage";
import logo from "../assets/logo.svg";
import logoutIcon from "../assets/logout.png";
import loadingGif from "../assets/loading.gif"

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
    console.log(`splitParts: ${splitParts}`);

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
    const [headers, setHeaders] = useState(sessionStorage.getItem('headers').replaceAll('\n', '-') || []);
    const [isCorrect, setIsCorrect] = useState(true)
    const [isAuth, setIsAuth] = useState(localStorage.getItem('isAuth') || true)

    // const crypto = require('crypto');

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

    // const Blob = require('blob');

    let jsonPart = null;
    let zipFilePart = null;

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
                    headers: headers,
                    mode: mode
                });

                const formData = new FormData();
                const fileInput = await localForage.getItem('zipFile');
                formData.append('file', fileInput);
                console.log(fileInput);

                let url =
                    process.env.REACT_APP_BACKEND_URL +
                    '/api/proceed?' + queryParams;

                const response = await fetch(url, {
                    method: 'POST',
                    body: formData,
                })

                console.log(response);

                if (response.headers.get('Content-Type').includes('multipart/mixed')) {
                    const responseBlob = await response.blob();
                    console.log(`headers: ${response.headers.get('Content-Type')}`);
                    // Parse the multipart response
                    const boundary = getBoundary(response.headers.get('Content-Type'));
                    const parts = await parseMultipart(responseBlob, boundary);
                    console.log(parts);
                    // Process the parts

                    parts.forEach(part => {
                        const contentType = part.headers['Content-Type'];
                        // const contentDisposition = part.headers['Content-Disposition'];
                        if (contentType) {
                            if (contentType.includes('application/json')) {
                                jsonPart = JSON.parse(part.body);
                            } else if (contentType.includes('application/zip')) {
                                part.responseType = "arraybuffer";
                                // console.log(part.body.length);
                                // // const base64 = atob(part.body);
                                // const conv = Iconv('windows-1251', 'utf8');
                                // const text = conv.convert(part.body).toString();
                                zipFilePart = new Blob([part.body],
                                    {
                                        type: "application/zip", endings: 'native'
                                    }
                                );
                            }
                        }
                    });
                    if (jsonPart && zipFilePart) {
                        console.log('JSON part:', jsonPart);
                        console.log('ZIP part:', zipFilePart);
                        if (jsonPart.status === 'OK') {
                            const result = jsonPart.result;
                            setGoogleSheetLink(result.googleSheets);
                            setCsvData(result.csv);
                            setSubmissionsData(zipFilePart);

                            setLoading(false);
                        } else if (jsonPart.status === 'FAILED') {
                            setComment(jsonPart.comment);
                            setIsCorrect(false);
                            alert(jsonPart.comment);
                        }
                    } else {
                        setComment('Some error caught while processing response');
                        setIsCorrect(false);
                        alert(comment);
                        return show404page();
                    }
                } else {
                    let jsonPart = await response.json();
                    console.log(jsonPart);
                    if (jsonPart.status === 'OK') {
                        const result = jsonPart.result;
                        setGoogleSheetLink(result.googleSheets);
                        setCsvData(result.csv);

                        setLoading(false);
                    } else if (jsonPart.status === 'FAILED') {
                        setComment(jsonPart.comment);
                        setIsCorrect(false);
                        alert(jsonPart.comment);
                    }
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
        // submissionsData.lastModifiedDate = new Date();
        // submissionsData.name = 'submissions'
        // console.log(submissionsData);
        const f = new File([submissionsData], 'submissions.zip');
        // f.type = 'zip';
        console.log(f);
        const url = window.URL.createObjectURL(f);
        const a = document.createElement('a');
        a.href = url;
        a.download = 'submissions.zip';
        a.click();
        window.URL.revokeObjectURL(url);
    };

    if (loading) {
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
                        <div className="loading-spinner">
                            <h1>Loading contest details...</h1>
                            <img src={loadingGif} width={200} height={200} alt='loading'/>
                        </div>
                    </div>
                </div>
            </div>); // Render a loading indicator while data is being fetched
    }

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
                    <p className={isCorrect ? 'correct-comment' : 'incorrect-comment'}>{comment}</p>
                    <div className="right-navigation-part">
                        <a href="/">
                            <button className={'logout'} onClick={() => logout()}>Logout
                            </button>
                        </a>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default ContestDetails;
