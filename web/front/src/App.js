import React from 'react';
import {BrowserRouter as Router, Navigate, Route, Routes} from 'react-router-dom';
import LoginPage from './components/LoginPage.jsx';
import GroupsPage from './components/GroupsPage.jsx';
import ContestsPage from './components/ContestsPage.jsx';
import ContestDetails from './components/ContestDetails.jsx';
import WeightsDistrPage from "./components/WeightsDistrPage.jsx";

const App = () => {
    return (
        <Router>
            <Routes>
                {/* Default Route for Unauthorized Users */}
                <Route path="/" element={<LoginPage/>}/>

                {/* Protected Routes (with Redirection) */}
                {localStorage.getItem('isAuthorized') === false ? (
                    <>
                        <Route path="/groups" element={<GroupsPage/>}/>
                        <Route path="/contests/:groupCode" element={<ContestsPage/>}/>
                        <Route
                            path="/contest-details/:groupCode/:contestId"
                            element={<ContestDetails/>}
                        />
                        <Route
                            path="/weightsDistr/:GroupCode/:contestId"
                            element={<WeightsDistrPage/>}
                        />
                    </>
                ) : (
                    <Route
                        path="*" // Catch-all route for unauthorized access
                        element={<Navigate to="/" replace={true}/>} // Redirect to the default route
                    />
                )}
            </Routes>
        </Router>
    );
};

export default App;
