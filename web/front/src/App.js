import React from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import LoginPage from './components/LoginPage.jsx';
import ContestDetails from './components/ContestDetails.jsx';
import WeightsDistrPage from "./components/WeightsDistrPage.jsx";
import LinkPage from "./components/LinkPage.jsx";
import Cookies from "js-cookie";

const RequireAuth = ({children}) => {
    /* ... logic to check if user is logged in */
    const isAuthenticated = localStorage.getItem('isAuthorized') === true

    if (!isAuthenticated) {
        // Redirect to login if not authenticated
        return <Navigate to="/"/>;
    }

    return children; // Render the protected component if authenticated
};


const App = () => {
    return (
        <Router>
            <Routes>
                {/* Default Route for Unauthorized Users */}
                <Route path="/" element={<LoginPage/>}/>

                {/* Protected Routes (with Redirection) */}
                <Route path="/link"
                       element={
                           localStorage.getItem('isAuthorized') ?
                               (<LinkPage/>) : (<Navigate to='/'/>)
                       }
                />
                <Route path="/contest-details/:groupCode/:contestId"
                       element={
                           localStorage.getItem('isAuthorized') ?
                               (<ContestDetails/>) : (<Navigate to='/'/>)
                       }
                />
                <Route path="/weights-distribution/:groupCode/:contestId"
                       element={
                           localStorage.getItem('isAuthorized') ?
                               (<WeightsDistrPage/>) : (<Navigate to='/'/>)
                       }
                />
                {/*<Route path="*" element={<Navigate to="/" replace/>}/>*/}
            </Routes>
        </Router>
    );
};

export default App;
