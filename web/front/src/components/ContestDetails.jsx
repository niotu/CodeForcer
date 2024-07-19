import React, {useEffect, useState} from 'react';
import {useParams} from 'react-router-dom';
import './styles.css';
import logout, {show404page} from "./globalFunctions.jsx";
import localForage from "localforage";
import logo from "../assets/logo.svg";
import logoutIcon from "../assets/logout.png";
import loadingGif from "../assets/loading.gif"

// import * as url from "node:url";

function getBoundary(contentType) {
    const boundaryPrefix = 'boundary=';
    const start = contentType.indexOf(boundaryPrefix) + boundaryPrefix.length;
    let end = contentType.indexOf(';', start);
    if (end === -1) end = contentType.length;
    return contentType.substring(start, end).trim();
}

async function parseMultipart(blob, boundary) {
    const parts = [];
    const delimiter = `--${boundary}`;
    const closeDelimiter = `--${boundary}--`;

    return new Promise((resolve, reject) => {
        const reader = new FileReader();

        reader.onload = async function (e) {
            try {
                const responseText = e.target.result;
                const splitParts = responseText.split(delimiter);

                for (let i = 1; i < splitParts.length - 1; i++) {
                    let part = splitParts[i].trim();

                    if (part === closeDelimiter) continue;

                    // Handle headersEndIndex correctly
                    const headersEndIndex = part.indexOf('\r\n\r\n');

                    // Handle the case where headersEndIndex is -1
                    const headersText = headersEndIndex !== -1 ? part.slice(0, headersEndIndex) : part;
                    const bodyText = headersEndIndex !== -1 ? part.slice(headersEndIndex + 4) : '';

                    const headers = {};

                    headersText.split('\r\n').forEach(header => {
                        const [key, value] = header.split(': ');
                        headers[key] = value;
                    });

                    if (headers['Content-Type'] && headers['Content-Type'].includes('application/zip')) {
                        const zipFilePart = new Blob([bodyText], {type: "application/zip"});
                        parts.push({headers, body: zipFilePart});
                    } else {
                        parts.push({headers, body: bodyText});
                    }
                }
                resolve(parts);
            } catch (error) {
                reject(error);
            }
        };

        reader.onerror = (error) => {
            reject(error);
        };

        reader.readAsText(blob);
    });
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

                // const request = new Request(url, {
                //     method: 'POST',
                //     body: formData,
                //     headers: {'Content-Type': 'multipart/form-data'}
                // })

                const response = await fetch(url, {
                    method: 'POST',
                    body: formData,
                    // headers: {'Content-Type': 'multipart/form-data'}
                });


                jsonPart = await response.json();

                console.log('JSON part:', jsonPart);
                if (jsonPart.status === 'OK') {
                    const result = jsonPart.result;
                    zipFilePart = result.zipLink;
                    setGoogleSheetLink(result.googleSheets);
                    setCsvData(result.csv);
                    setSubmissionsData(zipFilePart);
                    setIsFinished(true);
                    setLoading(false);
                } else if (jsonPart.status === 'FAILED') {
                    setComment(jsonPart.comment);
                    setIsCorrect(false);
                    setIsFinished(true);
                    // setLoading(false);
                    // alert(jsonPart.comment);
                }
            } catch (error) {
                console.error('Error fetching contest details:', error);
                setLoading(false); // Set loading to false even if there's an error
                return show404page();
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
        const a = document.createElement('a');
        a.href = submissionsData;
        a.download = 'submissions.zip';
        a.click();
        // window.URL.revokeObjectURL(url);
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
                        <div className={'filler'}></div>
                        <div className="loading-spinner">
                            <h1>Loading contest details...</h1>
                            <img src={loadingGif} width={200} height={200} alt='loading'/>
                        </div>
                        <div className={'navigation'}></div>
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
                    <div className={'filler'}>

                    </div>
                    <div className="panel">
                        <div className="left-part">
                            <h1>Contest Details</h1>
                            <h4>
                                weights: {taskWeights},
                                late: {late},
                                penalty: {penalty},
                                mode: {mode}
                            </h4>
                            <p className={isCorrect ? 'correct-comment' : 'incorrect-comment'}>{comment}</p>
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
                            <a href="/">
                                <button className={'home'}>Home
                                </button>
                            </a>
                        </div>
                    </div>
                </div>

            </div>
        </div>
    );
};

export default ContestDetails;
