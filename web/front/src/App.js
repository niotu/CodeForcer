import React from 'react';
import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
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
                {localStorage.getItem('isAuthorized') ? (
                    <Route path="/groups" element={<GroupsPage/>}/>) : (
                    <Route path="/" element={<LoginPage/>}/> // Redirect to the default route
                )}
                {localStorage.getItem('isAuthorized')? (
                    <Route path="/contests/:groupCode" element={<ContestsPage/>}/>) : (
                    <Route path="/" element={<LoginPage/>}/> // Redirect to the default route
                )}
                {localStorage.getItem('isAuthorized') ? (
                    <Route path="/contests/:groupCode/:contestId" element={<ContestDetails/>}/>) : (
                    <Route path="/" element={<LoginPage/>}/> // Redirect to the default route
                )}
                {localStorage.getItem('isAuthorized') ? (
                    <Route
                        path="/weightsDistr/:GroupCode/:contestId"
                        element={<WeightsDistrPage/>}
                    />) : (
                    <Route path="/" element={<LoginPage/>}/> // Redirect to the default route
                )}
            </Routes>
        </Router>
    );
};

export default App;
