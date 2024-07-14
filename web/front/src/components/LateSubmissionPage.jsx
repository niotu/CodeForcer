import React, {useState} from 'react';
import {useNavigate, useParams} from "react-router-dom";
import './styles.css';
import logout, {show404page} from "./globalFunctions.jsx";

function dateToUnix(date) {
    // Ensure we're working with a Date object
    if (!(date instanceof Date)) {
        throw new Error("Input must be a valid Date object.");
    }

    // Calculate the Unix timestamp (seconds since the Unix epoch)
    return Math.floor(date.getTime() / 1000);
}


const LateSubmissionPage = () => {
    // ... other state variables ...
    const {groupCode, contestId} = useParams();
    const [lateHours, setLateHours] = useState(12);
    const navigate = useNavigate();

    const [comment, setComment] = useState('');

    const [isCorrect, setIsCorrect] = useState(true);
    const [penalty, setPenalty] = useState(50);

    if (!localStorage.getItem('isAuthorized')) {
        return show404page();
    }

    const lateSubmit = async (e) => {
        e.preventDefault();
        console.log(' processing...')

        sessionStorage.setItem('lateHours', lateHours);
        sessionStorage.setItem('penalty', penalty);

        navigate(`/upload-csv/${groupCode}/${contestId}`);
    };

    // ... (Rest of your component logic) ...

    return (
        <div className="page-active">
            <div className="wizard">
                <div className="panel">
                    <div className="left-part">
                        <h1>Set up late submission date</h1>
                    </div>
                    <div className="right-part">
                        <form onSubmit={lateSubmit} autoComplete='on'>
                            <label htmlFor="number">Enter the late submission hours:</label>
                            <input
                                type="number" // Change to "date"
                                id="hours"
                                value={lateHours} // Format for date input
                                onChange={(e) => setLateHours(e.target.value)}
                                required
                                className={isCorrect ? 'correct' : 'incorrect'}
                            /><br/><br/>
                            <label htmlFor="penalty">Penalty in percents:</label>
                            <input
                                type="number"
                                id="penalty"
                                value={penalty}
                                onChange={(e) => setPenalty(e.target.value)}
                                required
                                className={isCorrect ? 'correct' : 'incorrect'}
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

export default LateSubmissionPage;
