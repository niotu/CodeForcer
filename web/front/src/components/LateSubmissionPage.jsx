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
    const [date, setDate] = useState(new Date());
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

        // ... (Your logic to get groupCode, contestId, and userID) ...

        // Convert the date to a string in ISO format for sending to the API
        const formattedDate = dateToUnix(date); // Get only the date part
        console.log(formattedDate);
        sessionStorage.setItem('date', formattedDate);
        sessionStorage.setItem('penalty', penalty);

        navigate(`/contest-details/${groupCode}/${contestId}`);
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
                            <label htmlFor="date">Choose the date:</label>
                            <input
                                type="datetime-local" // Change to "date"
                                id="date"
                                value={date.toISOString().slice(0, 16)} // Format for date input
                                onChange={(e) => setDate(new Date(e.target.value))}
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
